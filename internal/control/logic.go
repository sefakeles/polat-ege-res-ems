package control

import (
	"fmt"
	"sync"

	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/pkg/logger"
)

type ActivePowerControl struct {
	Enabled bool    `json:"enabled"`
	Power   float32 `json:"power"`
}

// Logic handles control logic and automation
type Logic struct {
	bmsManager *bms.Manager
	pcsManager *pcs.Manager
	config     config.EMSConfig
	mode       string // "AUTO", "MANUAL", "MAINTENANCE"
	mutex      sync.RWMutex
	log        logger.Logger

	activePowerControl ActivePowerControl // Active power control state
}

const (
	ModeAutomatic       = "AUTO"
	ModeManual          = "MANUAL"
	ModeMaintenance     = "MAINTENANCE"
	ModeSelfConsumption = "SELF_CONSUMPTION"
)

// NewLogic creates a new control logic instance
func NewLogic(bmsManager *bms.Manager, pcsManager *pcs.Manager, config config.EMSConfig) *Logic {
	// Create component-specific logger
	controlLogger := logger.With(
		logger.String("component", "control_logic"),
	)

	return &Logic{
		bmsManager: bmsManager,
		pcsManager: pcsManager,
		config:     config,
		mode:       ModeManual,
		log:        controlLogger,
	}
}

// SetMode sets the control mode
func (l *Logic) SetMode(mode string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	oldMode := l.mode
	l.mode = mode

	l.log.Info("Control mode changed",
		logger.String("old_mode", oldMode),
		logger.String("new_mode", mode))
}

// GetMode returns the current control mode
func (l *Logic) GetMode() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.mode
}

func (l *Logic) SetActivePowerControl(power float32) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.activePowerControl.Enabled = power != 0
	l.activePowerControl.Power = power
}

func (l *Logic) GetActivePowerControl() ActivePowerControl {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.activePowerControl
}

// ExecuteControl executes the control logic immediately based on fresh data
func (l *Logic) ExecuteControl() {
	l.mutex.RLock()
	mode := l.mode
	l.mutex.RUnlock()

	bms1Service, _ := l.bmsManager.GetService(1)
	pcs1Service, _ := l.pcsManager.GetService(1)

	bmsData := bms1Service.GetLatestBMSData()
	bmsStatusData := bms1Service.GetLatestBMSStatusData()

	// Safety checks
	if bms.IsFaultState(bmsStatusData.SystemStatus) {
		l.log.Warn("BMS in fault state, stopping PCS",
			logger.Uint16("state", bmsStatusData.SystemStatus),
			logger.String("state_description", bms.GetStateDescription(bmsStatusData.SystemStatus)))

		if err := pcs1Service.SetActivePowerCommand(0); err != nil {
			l.log.Error("Failed to set active power to zero during fault state", logger.Err(err))
		}
		return
	}

	// Prevent charging if SOC too high
	if bms.IsChargingState(bmsStatusData.SystemStatus) && float32(bmsData.SOC) >= l.config.MaxSOC {
		l.log.Warn("BMS SOC too high for charging, stopping charge",
			logger.Float32("soc", float32(bmsData.SOC)),
			logger.Float32("max_soc", l.config.MaxSOC))

		if err := pcs1Service.SetActivePowerCommand(0); err != nil {
			l.log.Error("Failed to set active power to zero during high SOC", logger.Err(err))
		}
	}

	// Prevent discharging if SOC too low
	if bms.IsDischargingState(bmsStatusData.SystemStatus) && float32(bmsData.SOC) <= l.config.MinSOC {
		l.log.Warn("BMS SOC too low for discharging, stopping discharge",
			logger.Float32("soc", float32(bmsData.SOC)),
			logger.Float32("max_soc", l.config.MaxSOC))

		if err := pcs1Service.SetActivePowerCommand(0); err != nil {
			l.log.Error("Failed to set active power to zero during low SOC", logger.Err(err))
		}
	}

	if mode != "AUTO" {
		return // Skip automatic control in manual or maintenance mode
	}
}

// GetBESSUpdateChannel returns the BESS data update channel for reactive control
func (l *Logic) GetBESSUpdateChannel() <-chan struct{} {
	bms1Service, _ := l.bmsManager.GetService(1)
	return bms1Service.GetBaseDataUpdateChannel()
}

func (l *Logic) calculateChargePower(bmsData database.BMSData) float32 {
	maxPower := min(float32(bmsData.MaxChargePower), l.config.MaxChargePower)

	// Apply SOC-based ramping
	soc := float32(bmsData.SOC)
	rampStartSOC := l.config.MaxSOC - 5.0 // Start ramping 5% below MaxSOC

	if soc > rampStartSOC {
		// Reduce charge power as SOC approaches MaxSOC
		rampFactor := (l.config.MaxSOC - soc) / 5.0
		if rampFactor < 0 {
			rampFactor = 0
		}
		maxPower *= rampFactor
	}

	return maxPower
}

func (l *Logic) calculateDischargePower(bmsData database.BMSData) float32 {
	maxPower := min(float32(bmsData.MaxDischargePower), l.config.MaxDischargePower)

	// Apply SOC-based ramping
	soc := float32(bmsData.SOC)
	rampStartSOC := l.config.MinSOC + 5.0 // Start ramping 5% above MinSOC

	if soc < rampStartSOC {
		// Reduce discharge power as SOC approaches MinSOC
		rampFactor := (soc - l.config.MinSOC) / 5.0
		if rampFactor < 0 {
			rampFactor = 0
		}
		maxPower *= rampFactor
	}

	return maxPower
}

// ManualPowerCommand handles manual power command
func (l *Logic) ManualPowerCommand(power float32) error {
	if l.GetMode() != "MANUAL" {
		l.log.Warn("Manual power command rejected - not in manual mode",
			logger.String("current_mode", l.GetMode()),
			logger.Float32("requested_power", power))
		return fmt.Errorf("manual power command only allowed in MANUAL mode")
	}

	bms1Service, _ := l.bmsManager.GetService(1)
	pcs1Service, _ := l.pcsManager.GetService(1)

	bmsData := bms1Service.GetLatestBMSData()
	bmsStatusData := bms1Service.GetLatestBMSStatusData()

	// Safety checks even in manual mode
	if bms.IsFaultState(bmsStatusData.SystemStatus) {
		l.log.Error("Manual power command rejected - BMS in fault state",
			logger.Uint16("bms_state", bmsStatusData.SystemStatus),
			logger.Float32("requested_power", power))
		return fmt.Errorf("BMS in fault state, command rejected")
	}

	originalPower := power

	// Check power limits
	if power < 0 { // Charging (negative power)
		maxCharge := l.calculateChargePower(bmsData)
		if -power > maxCharge {
			power = -maxCharge
			l.log.Warn("Manual charge power limited",
				logger.Float32("requested_power", originalPower),
				logger.Float32("limited_power", power),
				logger.Float32("max_charge", maxCharge))
		}
	} else if power > 0 { // Discharging (positive power)
		maxDischarge := l.calculateDischargePower(bmsData)
		if power > maxDischarge {
			power = maxDischarge
			l.log.Warn("Manual discharge power limited",
				logger.Float32("requested_power", originalPower),
				logger.Float32("limited_power", power),
				logger.Float32("max_discharge", maxDischarge))
		}
	}

	l.log.Info("Executing manual power command",
		logger.Float32("requested_power", originalPower),
		logger.Float32("final_power", power),
		logger.Float32("current_soc", float32(bmsData.SOC)))

	err := pcs1Service.SetActivePowerCommand(power)
	if err != nil {
		l.log.Error("Manual power command failed",
			logger.Err(err),
			logger.Float32("power", power))
		return err
	}

	l.SetActivePowerControl(power)

	l.log.Info("Manual power command executed successfully",
		logger.Float32("power", power))
	return nil
}

// ManualReactivePowerCommand handles manual reactive power command
func (l *Logic) ManualReactivePowerCommand(power float32) error {
	if l.GetMode() != "MANUAL" {
		l.log.Warn("Manual reactive power command rejected - not in manual mode",
			logger.String("current_mode", l.GetMode()),
			logger.Float32("requested_power", power))
		return fmt.Errorf("manual reactive power command only allowed in MANUAL mode")
	}

	bms1Service, _ := l.bmsManager.GetService(1)
	pcs1Service, _ := l.pcsManager.GetService(1)

	bmsData := bms1Service.GetLatestBMSData()
	bmsStatusData := bms1Service.GetLatestBMSStatusData()

	// Safety checks even in manual mode
	if bms.IsFaultState(bmsStatusData.SystemStatus) {
		l.log.Error("Manual reactive power command rejected - BMS in fault state",
			logger.Uint16("bms_state", bmsStatusData.SystemStatus),
			logger.Float32("requested_power", power))
		return fmt.Errorf("BMS in fault state, command rejected")
	}

	originalPower := power

	l.log.Info("Executing manual reactive power command",
		logger.Float32("requested_power", originalPower),
		logger.Float32("final_power", power),
		logger.Float32("current_soc", float32(bmsData.SOC)))

	err := pcs1Service.SetReactivePowerCommand(power)
	if err != nil {
		l.log.Error("Manual reactive power command failed",
			logger.Err(err),
			logger.Float32("power", power))
		return err
	}

	l.log.Info("Manual reactive power command executed successfully",
		logger.Float32("power", power))
	return nil
}
