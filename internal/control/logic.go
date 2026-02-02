package control

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/internal/pcs"
)

type ActivePowerControl struct {
	Enabled bool    `json:"enabled"`
	Power   float32 `json:"power"`
}

// Logic handles control logic and automation
type Logic struct {
	config     config.EMSConfig
	bmsManager *bms.Manager
	pcsManager *pcs.Manager
	log        *zap.Logger

	mutex              sync.RWMutex
	mode               string             // "AUTO", "MANUAL", "MAINTENANCE"
	activePowerControl ActivePowerControl // Active power control state
}

const (
	ModeAutomatic       = "AUTO"
	ModeManual          = "MANUAL"
	ModeMaintenance     = "MAINTENANCE"
	ModeSelfConsumption = "SELF_CONSUMPTION"
)

// NewLogic creates a new control logic instance
func NewLogic(config config.EMSConfig, bmsManager *bms.Manager, pcsManager *pcs.Manager, logger *zap.Logger) *Logic {
	// Create component-specific logger
	controlLogger := logger.With(
		zap.String("component", "control_logic"),
	)

	return &Logic{
		config:     config,
		bmsManager: bmsManager,
		pcsManager: pcsManager,
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
		zap.String("old_mode", oldMode),
		zap.String("new_mode", mode))
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

	// Check all BMS-PCS pairs
	l.checkBMSPCSPairs()

	if mode != "AUTO" {
		return // Skip automatic control in manual or maintenance mode
	}
}

// checkBMSPCSPairs checks SOC limits for each BMS-PCS pair and stops PCS if needed
func (l *Logic) checkBMSPCSPairs() {
	// Each PCS is connected to 2 BMS units
	// PCS1 -> BMS1, BMS2
	// PCS2 -> BMS3, BMS4
	// PCS3 -> BMS5, BMS6
	// PCS4 -> BMS7, BMS8

	pcsCount := 4

	for pcsID := 1; pcsID <= pcsCount; pcsID++ {
		bms1ID := (pcsID-1)*2 + 1
		bms2ID := (pcsID-1)*2 + 2

		shouldStopPCS := false
		reason := ""

		// Get PCS data to check power direction
		pcsService, err := l.pcsManager.GetService(pcsID)
		if err != nil {
			l.log.Error("Failed to get PCS service",
				zap.Error(err),
				zap.Int("pcs_id", pcsID))
			continue
		}
		pcsCommandState := pcsService.GetCommandState()
		pcsPower := pcsCommandState.ActivePowerCommand

		// Check BMS1 for this PCS
		bms1Service, err := l.bmsManager.GetService(bms1ID)
		if err == nil {
			bmsData := bms1Service.GetLatestBMSData()
			bmsStatusData := bms1Service.GetLatestBMSStatusData()

			// Check for fault state
			if bms.IsFaultState(bmsStatusData.SystemStatus) {
				shouldStopPCS = true
				reason = fmt.Sprintf("BMS%d in fault state", bms1ID)
			}

			// Check for high SOC during charging (negative power)
			if pcsPower < 0 && (bms.IsFullChargeState(bmsStatusData.SystemStatus) || float32(bmsData.SOC) >= l.config.MaxSOC) {
				shouldStopPCS = true
				reason = fmt.Sprintf("BMS%d SOC at MaxSOC during charging", bms1ID)
			}

			// Check for low SOC during discharging (positive power)
			if pcsPower > 0 && (bms.IsFullDischargeState(bmsStatusData.SystemStatus) || float32(bmsData.SOC) <= l.config.MinSOC) {
				shouldStopPCS = true
				reason = fmt.Sprintf("BMS%d SOC at MinSOC during discharging", bms1ID)
			}
		}

		// Check BMS2 for this PCS (if it exists)
		bms2Service, err := l.bmsManager.GetService(bms2ID)
		if err == nil {
			bmsData := bms2Service.GetLatestBMSData()
			bmsStatusData := bms2Service.GetLatestBMSStatusData()

			// Check for fault state
			if bms.IsFaultState(bmsStatusData.SystemStatus) {
				shouldStopPCS = true
				if reason != "" {
					reason += fmt.Sprintf(", BMS%d in fault state", bms2ID)
				} else {
					reason = fmt.Sprintf("BMS%d in fault state", bms2ID)
				}
			}

			// Check for high SOC during charging (negative power)
			if pcsPower < 0 && (bms.IsFullChargeState(bmsStatusData.SystemStatus) || float32(bmsData.SOC) >= l.config.MaxSOC) {
				shouldStopPCS = true
				if reason != "" {
					reason += fmt.Sprintf(", BMS%d SOC at MaxSOC during charging", bms2ID)
				} else {
					reason = fmt.Sprintf("BMS%d SOC at MaxSOC during charging", bms2ID)
				}
			}

			// Check for low SOC during discharging (positive power)
			if pcsPower > 0 && (bms.IsFullDischargeState(bmsStatusData.SystemStatus) || float32(bmsData.SOC) <= l.config.MinSOC) {
				shouldStopPCS = true
				if reason != "" {
					reason += fmt.Sprintf(", BMS%d SOC at MinSOC during discharging", bms2ID)
				} else {
					reason = fmt.Sprintf("BMS%d SOC at MinSOC during discharging", bms2ID)
				}
			}
		}

		// Stop PCS if needed
		if shouldStopPCS {
			l.log.Warn("Stopping PCS due to BMS condition",
				zap.Int("pcs_id", pcsID),
				zap.String("reason", reason))

			// Set active power to zero
			if err := pcsService.SetActivePowerCommand(0); err != nil {
				l.log.Error("Failed to set active power to zero",
					zap.Error(err),
					zap.Int("pcs_id", pcsID))
			}

			// Optionally stop the PCS completely
			// if err := pcsService.StartStopCommand(false); err != nil {
			// 	l.log.Error("Failed to stop PCS",
			// 		logger.Err(err),
			// 		logger.Int("pcs_id", pcsID))
			// }
		}
	}
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

// GetBESSUpdateChannel returns the BESS data update channel for reactive control
func (l *Logic) GetBESSUpdateChannel() <-chan struct{} {
	bms1Service, _ := l.bmsManager.GetService(1)
	return bms1Service.GetSystemDataUpdateChannel()
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
			zap.String("current_mode", l.GetMode()),
			zap.Float32("requested_power", power))
		return fmt.Errorf("manual power command only allowed in MANUAL mode")
	}

	bms1Service, _ := l.bmsManager.GetService(1)
	pcs1Service, _ := l.pcsManager.GetService(1)

	bmsData := bms1Service.GetLatestBMSData()
	bmsStatusData := bms1Service.GetLatestBMSStatusData()

	// Safety checks even in manual mode
	if bms.IsFaultState(bmsStatusData.SystemStatus) {
		l.log.Error("Manual power command rejected - BMS in fault state",
			zap.Uint16("bms_state", bmsStatusData.SystemStatus),
			zap.Float32("requested_power", power))
		return fmt.Errorf("BMS in fault state, command rejected")
	}

	originalPower := power

	// Check power limits
	if power < 0 { // Charging (negative power)
		maxCharge := l.calculateChargePower(bmsData)
		if -power > maxCharge {
			power = -maxCharge
			l.log.Warn("Manual charge power limited",
				zap.Float32("requested_power", originalPower),
				zap.Float32("limited_power", power),
				zap.Float32("max_charge", maxCharge))
		}
	} else if power > 0 { // Discharging (positive power)
		maxDischarge := l.calculateDischargePower(bmsData)
		if power > maxDischarge {
			power = maxDischarge
			l.log.Warn("Manual discharge power limited",
				zap.Float32("requested_power", originalPower),
				zap.Float32("limited_power", power),
				zap.Float32("max_discharge", maxDischarge))
		}
	}

	l.log.Info("Executing manual power command",
		zap.Float32("requested_power", originalPower),
		zap.Float32("final_power", power),
		zap.Float32("current_soc", float32(bmsData.SOC)))

	err := pcs1Service.SetActivePowerCommand(power)
	if err != nil {
		l.log.Error("Manual power command failed",
			zap.Error(err),
			zap.Float32("power", power))
		return err
	}

	l.SetActivePowerControl(power)

	l.log.Info("Manual power command executed successfully",
		zap.Float32("power", power))
	return nil
}

// ManualReactivePowerCommand handles manual reactive power command
func (l *Logic) ManualReactivePowerCommand(power float32) error {
	if l.GetMode() != "MANUAL" {
		l.log.Warn("Manual reactive power command rejected - not in manual mode",
			zap.String("current_mode", l.GetMode()),
			zap.Float32("requested_power", power))
		return fmt.Errorf("manual reactive power command only allowed in MANUAL mode")
	}

	bms1Service, _ := l.bmsManager.GetService(1)
	pcs1Service, _ := l.pcsManager.GetService(1)

	bmsData := bms1Service.GetLatestBMSData()
	bmsStatusData := bms1Service.GetLatestBMSStatusData()

	// Safety checks even in manual mode
	if bms.IsFaultState(bmsStatusData.SystemStatus) {
		l.log.Error("Manual reactive power command rejected - BMS in fault state",
			zap.Uint16("bms_state", bmsStatusData.SystemStatus),
			zap.Float32("requested_power", power))
		return fmt.Errorf("BMS in fault state, command rejected")
	}

	originalPower := power

	l.log.Info("Executing manual reactive power command",
		zap.Float32("requested_power", originalPower),
		zap.Float32("final_power", power),
		zap.Float32("current_soc", float32(bmsData.SOC)))

	err := pcs1Service.SetReactivePowerCommand(power)
	if err != nil {
		l.log.Error("Manual reactive power command failed",
			zap.Error(err),
			zap.Float32("power", power))
		return err
	}

	l.log.Info("Manual reactive power command executed successfully",
		zap.Float32("power", power))
	return nil
}
