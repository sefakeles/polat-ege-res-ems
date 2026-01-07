package fcr

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/pkg/logger"
)

const (
	// FCR-N frequency limits (Hz)
	FreqNominal      = 50.0 // Nominal frequency
	FreqFCRNMin      = 49.9 // Lower limit for FCR-N activation
	FreqFCRNMax      = 50.1 // Upper limit for FCR-N activation
	FreqFCRNDeadband = 0.1  // Full activation at ±0.1 Hz

	// NEM/AEM power factors (as fraction of capacity)
	NEMPowerFactor = 0.10 // NEM power is 0.10 * C_FCR-N

	// Rolling average window size (5 minutes at 1 second resolution)
	RollingAverageWindowSize = 300
)

// FCRNState represents the current state of FCR-N controller
type FCRNState struct {
	Enabled            bool      `json:"enabled"`
	Active             bool      `json:"active"`              // FCR-N is actively providing reserves
	NEMActive          bool      `json:"nem_active"`          // NEM is currently active
	AEMActive          bool      `json:"aem_active"`          // AEM is currently active
	FrequencyMeasured  float64   `json:"frequency_measured"`  // Hz
	FrequencyReference float64   `json:"frequency_reference"` // Hz (reference for FCR calculation, affected by AEM)
	ActivatedPower     float64   `json:"activated_power"`     // kW
	Baseline           float64   `json:"baseline"`            // kW (power without FCR activation)
	NEMPower           float64   `json:"nem_power"`           // kW (NEM power adjustment)
	TotalPower         float64   `json:"total_power"`         // kW (baseline + activated + NEM)
	Capacity           float64   `json:"capacity"`            // kW (sold FCR-N capacity)
	Droop              float64   `json:"droop"`               // % (droop setting)
	SOC                float64   `json:"soc"`                 // % (State of Charge)
	EnduranceUpward    float64   `json:"endurance_upward"`    // minutes
	EnduranceDownward  float64   `json:"endurance_downward"`  // minutes
	NEMCurrent         float64   `json:"nem_current"`         // Current NEM activation level (-1 to +1)
	LastUpdate         time.Time `json:"last_update"`
}

// FCRNController implements FCR-N control logic
type FCRNController struct {
	config config.FCRNConfig
	state  FCRNState
	mutex  sync.RWMutex
	log    logger.Logger

	// Context for lifecycle management
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Callbacks for power commands
	powerCommandCallback func(power float64) error

	// Historical data for rolling averages
	frequencyHistory []float64
	maxHistorySize   int

	// NEM/AEM state tracking
	nemAllowedHistory  []float64 // History for NEM_Allowed rolling average
	aemFreqHistory     []float64 // History for AEM frequency reference rolling average
	lastNEMUpdateTime  time.Time // Last time NEM/AEM samples were added
	lastAEMUpdateTime  time.Time // Last time AEM sample was added
	nemPreviousState   int       // Previous NEM_Allowed state for hysteresis (-1, 0, 1)
	nemActiveDirection int       // Direction of active NEM: -1 = charging, 0 = none, 1 = discharging

	// SOC thresholds (calculated from capacity and reservoir)
	socEnableNEMLower  float64
	socDisableNEMLower float64
	socEnableAEMLower  float64
	socDisableAEMLower float64
	socEnableNEMUpper  float64
	socDisableNEMUpper float64
	socEnableAEMUpper  float64
	socDisableAEMUpper float64
}

// NewFCRNController creates a new FCR-N controller
func NewFCRNController(cfg config.FCRNConfig, powerCallback func(power float64) error) *FCRNController {
	ctx, cancel := context.WithCancel(context.Background())

	controllerLog := logger.With(logger.String("component", "fcrn_controller"))

	controller := &FCRNController{
		config: cfg,
		state: FCRNState{
			Enabled:            false,
			Active:             false,
			NEMActive:          false,
			AEMActive:          false,
			FrequencyMeasured:  FreqNominal,
			FrequencyReference: FreqNominal,
			Capacity:           cfg.Capacity,
			Droop:              cfg.Droop,
			Baseline:           0.0,
		},
		log:                  controllerLog,
		ctx:                  ctx,
		cancel:               cancel,
		powerCommandCallback: powerCallback,
		maxHistorySize:       10,
		frequencyHistory:     make([]float64, 0, 10),
		nemAllowedHistory:    make([]float64, 0, RollingAverageWindowSize),
		aemFreqHistory:       make([]float64, 0, RollingAverageWindowSize),
		nemPreviousState:     0,
		nemActiveDirection:   0,
	}

	// Calculate SOC thresholds based on capacity and reservoir
	controller.calculateSOCThresholds()

	return controller
}

// calculateSOCThresholds calculates SOC thresholds according to ENTSO-E Table 11
func (c *FCRNController) calculateSOCThresholds() {
	if !c.config.EnableEnergyManagement || c.config.ReservoirCapacity <= 0 {
		return
	}

	capacity := c.config.Capacity           // C_FCR-N in MW
	reservoir := c.config.ReservoirCapacity // E in MWh

	// Lower thresholds (for low SOC - need to charge)
	c.socEnableAEMLower = (capacity * 5.0 / 60.0) / reservoir
	c.socDisableAEMLower = (capacity * 10.0 / 60.0) / reservoir
	c.socEnableNEMLower = (capacity * 30.0 / 60.0) / reservoir
	c.socDisableNEMLower = (capacity * 57.5 / 60.0) / reservoir

	// Upper thresholds (for high SOC - need to discharge)
	c.socEnableAEMUpper = 1.0 - (capacity*5.0/60.0)/reservoir
	c.socDisableAEMUpper = 1.0 - (capacity*10.0/60.0)/reservoir
	c.socEnableNEMUpper = 1.0 - (capacity*30.0/60.0)/reservoir
	c.socDisableNEMUpper = 1.0 - (capacity*57.5/60.0)/reservoir

	c.log.Info("Calculated SOC thresholds",
		logger.Float64("nem_enable_lower", c.socEnableNEMLower),
		logger.Float64("nem_disable_lower", c.socDisableNEMLower),
		logger.Float64("aem_enable_lower", c.socEnableAEMLower),
		logger.Float64("aem_disable_lower", c.socDisableAEMLower),
		logger.Float64("nem_enable_upper", c.socEnableNEMUpper),
		logger.Float64("nem_disable_upper", c.socDisableNEMUpper),
		logger.Float64("aem_enable_upper", c.socEnableAEMUpper),
		logger.Float64("aem_disable_upper", c.socDisableAEMUpper))
}

// Start starts the FCR-N controller
func (c *FCRNController) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.state.Enabled {
		return fmt.Errorf("FCR-N controller already started")
	}

	c.log.Info("Starting FCR-N controller",
		logger.Float64("capacity_mw", c.config.Capacity),
		logger.Float64("droop_percent", c.config.Droop))

	c.state.Enabled = true
	c.state.Active = false
	c.state.NEMActive = false
	c.state.AEMActive = false
	c.state.LastUpdate = time.Now()

	// Start control loop
	c.wg.Add(1)
	go c.controlLoop()

	return nil
}

// Stop stops the FCR-N controller
func (c *FCRNController) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.state.Enabled {
		return fmt.Errorf("FCR-N controller not running")
	}

	c.log.Info("Stopping FCR-N controller")

	// Reset all power values before stopping
	c.state.ActivatedPower = 0.0
	c.state.NEMPower = 0.0
	c.state.NEMCurrent = 0.0
	c.state.TotalPower = 0.0

	// Send baseline power command to BESS
	if c.powerCommandCallback != nil {
		if err := c.powerCommandCallback(c.state.Baseline); err != nil {
			c.log.Error("Failed to send stop power command",
				logger.Err(err),
				logger.Float64("power", c.state.Baseline))
			// Don't return error - still proceed with stop
		}
	}

	c.cancel()
	c.wg.Wait()

	c.state.Enabled = false
	c.state.Active = false
	c.state.NEMActive = false
	c.state.AEMActive = false

	c.log.Info("FCR-N controller stopped, returned to baseline",
		logger.Float64("baseline", c.state.Baseline))

	return nil
}

// Activate activates FCR-N provision with smooth transition
func (c *FCRNController) Activate() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.state.Enabled {
		return fmt.Errorf("FCR-N controller not started")
	}

	if c.state.Active {
		return fmt.Errorf("FCR-N already active")
	}

	c.log.Info("Activating FCR-N provision",
		logger.Float64("current_frequency", c.state.FrequencyMeasured))

	c.state.Active = true

	return nil
}

// Deactivate deactivates FCR-N provision with smooth transition
func (c *FCRNController) Deactivate() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.state.Enabled {
		return fmt.Errorf("FCR-N controller not started")
	}

	if !c.state.Active {
		return fmt.Errorf("FCR-N already inactive")
	}

	c.log.Info("Deactivating FCR-N provision",
		logger.Float64("current_frequency", c.state.FrequencyMeasured))

	// Reset reference frequency
	c.state.FrequencyReference = FreqNominal
	c.state.Active = false
	c.state.NEMActive = false
	c.state.AEMActive = false

	// Clear NEM/AEM history to start fresh on next activation
	c.nemAllowedHistory = c.nemAllowedHistory[:0]
	c.aemFreqHistory = c.aemFreqHistory[:0]
	c.lastNEMUpdateTime = time.Time{}
	c.lastAEMUpdateTime = time.Time{}
	c.nemPreviousState = 0
	c.nemActiveDirection = 0

	// Reset all power values
	c.state.ActivatedPower = 0.0
	c.state.NEMPower = 0.0
	c.state.NEMCurrent = 0.0
	c.state.TotalPower = 0.0

	// Send baseline power command to BESS
	if c.powerCommandCallback != nil {
		if err := c.powerCommandCallback(c.state.Baseline); err != nil {
			c.log.Error("Failed to send deactivation power command",
				logger.Err(err),
				logger.Float64("power", c.state.Baseline))
			// Don't return error - still mark as deactivated
		}
	}

	c.log.Info("FCR-N deactivated, returned to baseline",
		logger.Float64("baseline_kw", c.state.Baseline))

	return nil
}

// UpdateFrequency updates the measured grid frequency
func (c *FCRNController) UpdateFrequency(frequency float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.state.FrequencyMeasured = frequency

	// Add to history for moving average
	c.frequencyHistory = append(c.frequencyHistory, frequency)
	if len(c.frequencyHistory) > c.maxHistorySize {
		c.frequencyHistory = c.frequencyHistory[1:]
	}

	// Update frequency reference if not in AEM mode
	// In AEM mode, reference is calculated separately with rolling average
	if !c.state.AEMActive {
		c.state.FrequencyReference = frequency
	}
}

// UpdateSOC updates the State of Charge
func (c *FCRNController) UpdateSOC(soc float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.state.SOC = soc
}

// UpdateBaseline updates the baseline power
func (c *FCRNController) UpdateBaseline(baseline float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.state.Baseline = baseline
}

// calculateActivatedPower calculates FCR-N activated power based on frequency deviation
func (c *FCRNController) calculateActivatedPower() float64 {
	// Frequency deviation from reference
	// In normal modes: FrequencyReference = FrequencyMeasured (normal FCR response)
	// In AEM mode: FrequencyReference is modified rolling average (gradual reduction)
	deltaF := FreqNominal - c.state.FrequencyReference

	// Linear activation: full capacity at ±0.1 Hz deviation (FreqFCRNDeadband)
	// P = Capacity × (Δf / FreqFCRNDeadband)
	activatedPower := c.state.Capacity * deltaF / FreqFCRNDeadband

	// Limit to capacity (saturation)
	if activatedPower > c.state.Capacity {
		activatedPower = c.state.Capacity
	} else if activatedPower < -c.state.Capacity {
		activatedPower = -c.state.Capacity
	}

	return activatedPower
}

// calculateEndurance calculates remaining endurance based on equations (27) and (28)
func (c *FCRNController) calculateEndurance() (upward, downward float64) {
	if !c.config.EnableEnergyManagement {
		return 9999.0, 9999.0 // No LER restrictions
	}

	// Get current reservoir state
	currentSOC := c.state.SOC / 100.0                        // Convert to 0-1 range
	currentEnergy := currentSOC * c.config.ReservoirCapacity // E_current [kWh]

	// Reservoir limits
	maxEnergy := c.config.ReservoirCapacity // E_max [kWh]
	minEnergy := 0.0                        // E_min [kWh]

	// Current power setpoint [kW] - includes baseline + FCR activation + NEM
	// Positive = discharging, Negative = charging
	powerSetpoint := c.state.TotalPower

	// FCR-N capacity [kW]
	capacity := c.state.Capacity

	// Reservoir inflow [kW] - zero for BESS
	reservoirInflow := 0.0

	// Equation (27): Upward Endurance (time until empty)
	// L_upwards = |[E_current - E_min] / [P_setpoint + C_upwards - P_inflow]| × 60
	if c.state.AEMActive && currentSOC <= c.socEnableAEMLower {
		upward = 0.0 // In AEM with low SOC, can't discharge
	} else {
		energyAvailable := currentEnergy - minEnergy
		netPowerDrain := powerSetpoint + capacity - reservoirInflow

		if netPowerDrain > 0.001 {
			upward = (energyAvailable / netPowerDrain) * 60.0
		} else {
			upward = 9999.0 // Not draining
		}
	}

	// Equation (28): Downward Endurance (time until full)
	// L_downwards = |[E_max - E_current] / [P_inflow - P_setpoint + C_downwards]| × 60
	if c.state.AEMActive && currentSOC >= c.socEnableAEMUpper {
		downward = 0.0 // In AEM with high SOC, can't charge
	} else {
		energyRemaining := maxEnergy - currentEnergy
		netPowerCharge := reservoirInflow - powerSetpoint + capacity

		if netPowerCharge > 0.001 {
			downward = (energyRemaining / netPowerCharge) * 60.0
		} else {
			downward = 9999.0 // Not charging
		}
	}

	// Ensure non-negative
	if upward < 0 {
		upward = 0.0
	}
	if downward < 0 {
		downward = 0.0
	}

	return upward, downward
}

// isFrequencyInStandardRange checks if frequency is within standard range (49.9 - 50.1 Hz)
func (c *FCRNController) isFrequencyInStandardRange() bool {
	return c.state.FrequencyMeasured >= FreqFCRNMin && c.state.FrequencyMeasured <= FreqFCRNMax
}

// calculateNEMAllowed calculates NEM_Allowed according to equation (14) with proper hysteresis
// CORRECTED: Now properly implements hysteresis using enable/disable thresholds
func (c *FCRNController) calculateNEMAllowed() float64 {
	currentSOC := c.state.SOC / 100.0
	inStandardRange := c.isFrequencyInStandardRange()

	// NEM only operates when frequency is in standard range
	if !inStandardRange {
		// When frequency leaves standard range, NEM_Allowed should become 0
		// This will cause NEM_Current to gradually ramp to 0
		c.nemActiveDirection = 0 // Clear active direction
		return 0.0
	}

	// Implement proper hysteresis for lower thresholds (low SOC - need to charge)
	// Enable charging when SOC drops below enable threshold
	// Keep charging until SOC rises above disable threshold
	switch c.nemActiveDirection {
	case 0:
		// Not currently in NEM - check enable thresholds
		if currentSOC < c.socEnableNEMLower {
			c.nemActiveDirection = -1 // Enable charging
			return -1.0
		} else if currentSOC > c.socEnableNEMUpper {
			c.nemActiveDirection = 1 // Enable discharging
			return 1.0
		}
		return 0.0
	case -1:
		// Currently charging - check disable threshold
		if currentSOC >= c.socDisableNEMLower {
			c.nemActiveDirection = 0 // Disable charging
			return 0.0
		}
		return -1.0 // Continue charging
	case 1:
		// Currently discharging - check disable threshold
		if currentSOC <= c.socDisableNEMUpper {
			c.nemActiveDirection = 0 // Disable discharging
			return 0.0
		}
		return 1.0 // Continue discharging
	}

	return 0.0
}

// updateNEMCurrent updates NEM_Current as rolling average of NEM_Allowed (equation 15)
// CRITICAL: Only adds samples at 1 second intervals per ENTSO-E requirements
func (c *FCRNController) updateNEMCurrent() {
	now := time.Now()

	// Only add a new sample if at least 1 second has passed since last update
	if now.Sub(c.lastNEMUpdateTime) >= 1*time.Second {
		nemAllowed := c.calculateNEMAllowed()

		// Add to history
		c.nemAllowedHistory = append(c.nemAllowedHistory, nemAllowed)
		if len(c.nemAllowedHistory) > RollingAverageWindowSize {
			c.nemAllowedHistory = c.nemAllowedHistory[1:]
		}

		c.lastNEMUpdateTime = now
	}

	// Always recalculate the rolling average (even if we didn't add a new sample)
	// CRITICAL: Divide by FULL window size (300), not current history length
	// This ensures proper ramp: at 60 seconds, NEM_Current = -60/300 = -0.2 (20%)
	if len(c.nemAllowedHistory) > 0 {
		sum := 0.0
		for _, val := range c.nemAllowedHistory {
			sum += val
		}
		// Divide by RollingAverageWindowSize (300), not len(history)
		c.state.NEMCurrent = sum / float64(RollingAverageWindowSize)
	} else {
		c.state.NEMCurrent = 0.0
	}
}

// updateEnergyManagement updates energy management mode based on SOC and frequency
// CORRECTED: NEM and AEM can now be active simultaneously
// CORRECTED: AEM activation is based on SOC only, not frequency
func (c *FCRNController) updateEnergyManagement() {
	if !c.config.EnableEnergyManagement {
		return
	}

	currentSOC := c.state.SOC / 100.0

	// Store previous states for logging
	oldNEMActive := c.state.NEMActive
	oldAEMActive := c.state.AEMActive

	if !c.state.Active {
		// Don't change state if not active
		c.state.NEMActive = false
		c.state.AEMActive = false
		return
	}

	// Determine if NEM should be active based on NEM_Current
	// NEM is active when NEM_Current is non-zero
	c.state.NEMActive = c.state.NEMCurrent != 0.0

	// Determine if AEM should be active based on SOC
	// CORRECTED: AEM activates based on SOC alone, regardless of frequency
	// This is because AEM is designed for alert state (when frequency is outside standard range)
	inAEMLower := currentSOC <= c.socEnableAEMLower
	inAEMUpper := currentSOC >= c.socEnableAEMUpper

	// Apply hysteresis for AEM
	if c.state.AEMActive {
		// Currently in AEM - check disable thresholds
		if inAEMLower {
			c.state.AEMActive = currentSOC <= c.socDisableAEMLower
		} else if inAEMUpper {
			c.state.AEMActive = currentSOC >= c.socDisableAEMUpper
		} else {
			c.state.AEMActive = false
		}
	} else {
		// Not in AEM - check enable thresholds
		c.state.AEMActive = inAEMLower || inAEMUpper
	}

	// Log state changes
	if c.state.NEMActive != oldNEMActive || c.state.AEMActive != oldAEMActive {
		c.log.Info("Energy management state changed",
			logger.Bool("nem_active", c.state.NEMActive),
			logger.Bool("aem_active", c.state.AEMActive),
			logger.Float64("soc", c.state.SOC),
			logger.Float64("frequency", c.state.FrequencyMeasured),
			logger.Float64("nem_current", c.state.NEMCurrent))
	}
}

// calculateNEMPower calculates the NEM power adjustment according to equation (16)
func (c *FCRNController) calculateNEMPower() float64 {
	if !c.state.NEMActive || c.state.NEMCurrent == 0.0 {
		return 0.0
	}

	// P_NEM^FCR-N = 0.34 * C_FCR-N * NEM_Current
	// where NEM_Current is between -1 and +1
	nemPower := NEMPowerFactor * c.state.Capacity * c.state.NEMCurrent

	return nemPower
}

// updateAEMReference updates the AEM frequency reference according to equations (18)-(20)
// CORRECTED: Now uses simple arithmetic mean, not weighted average
// CRITICAL: Only adds samples at 1 second intervals per ENTSO-E requirements
func (c *FCRNController) updateAEMReference() {
	if !c.state.AEMActive {
		// If not in AEM, reference follows measured frequency
		c.state.FrequencyReference = c.state.FrequencyMeasured
		// Clear history when leaving AEM
		c.aemFreqHistory = c.aemFreqHistory[:0]
		c.lastAEMUpdateTime = time.Time{} // Reset timer
		return
	}

	now := time.Now()

	// Only add a new sample if at least 1 second has passed since last update
	if now.Sub(c.lastAEMUpdateTime) >= 1*time.Second {
		// Calculate f_AEM according to equation (19) and (20)
		var fAEM float64
		fMeasured := c.state.FrequencyMeasured

		// For FCR-N: saturate measured frequency to standard range (equation 20)
		if fMeasured > FreqFCRNMax {
			fAEM = FreqFCRNMax
		} else if fMeasured < FreqFCRNMin {
			fAEM = FreqFCRNMin
		} else {
			fAEM = fMeasured
		}

		// Add to history
		c.aemFreqHistory = append(c.aemFreqHistory, fAEM)
		if len(c.aemFreqHistory) > RollingAverageWindowSize {
			c.aemFreqHistory = c.aemFreqHistory[1:]
		}

		c.lastAEMUpdateTime = now
	}

	// Calculate reference frequency as rolling mean (equation 18)
	// CORRECTED: Use simple arithmetic mean, not weighted average
	if len(c.aemFreqHistory) > 0 {
		sum := 0.0
		for _, val := range c.aemFreqHistory {
			sum += val
		}
		// Simple average - divide by actual number of samples
		c.state.FrequencyReference = sum / float64(len(c.aemFreqHistory))
	} else {
		// No history yet - use nominal frequency
		c.state.FrequencyReference = FreqNominal
	}
}

// controlLoop is the main control loop
func (c *FCRNController) controlLoop() {
	defer c.wg.Done()

	ticker := time.NewTicker(c.config.UpdateInterval)
	defer ticker.Stop()

	c.log.Info("FCR-N control loop started")

	for {
		select {
		case <-c.ctx.Done():
			c.log.Info("FCR-N control loop stopped")
			return

		case <-ticker.C:
			c.executeControlCycle()
		}
	}
}

// executeControlCycle executes one control cycle
func (c *FCRNController) executeControlCycle() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.state.Enabled || !c.state.Active {
		return
	}

	// Update NEM_Current rolling average
	c.updateNEMCurrent()

	// Update energy management states (NEM/AEM)
	c.updateEnergyManagement()

	// Update AEM frequency reference
	c.updateAEMReference()

	// Calculate activated power (uses frequency reference, which is affected by AEM)
	c.state.ActivatedPower = c.calculateActivatedPower()

	// Calculate NEM power adjustment
	c.state.NEMPower = c.calculateNEMPower()

	// Calculate total power: baseline + FCR activation + NEM adjustment
	c.state.TotalPower = c.state.Baseline + c.state.ActivatedPower + c.state.NEMPower

	// Apply power limits
	c.state.TotalPower = c.applyPowerLimits(c.state.TotalPower)

	// Calculate endurance
	c.state.EnduranceUpward, c.state.EnduranceDownward = c.calculateEndurance()

	// Update timestamp
	c.state.LastUpdate = time.Now()

	// Send power command
	if c.powerCommandCallback != nil {
		if err := c.powerCommandCallback(c.state.TotalPower); err != nil {
			c.log.Error("Failed to send power command",
				logger.Err(err),
				logger.Float64("power", c.state.TotalPower))
		}
	}

	// Log state periodically
	if time.Now().Second()%10 == 0 {
		c.log.Info("FCR-N control cycle",
			logger.Bool("nem_active", c.state.NEMActive),
			logger.Bool("aem_active", c.state.AEMActive),
			logger.Float64("droop", c.state.Droop),
			logger.Float64("frequency", c.state.FrequencyMeasured),
			logger.Float64("freq_reference", c.state.FrequencyReference),
			logger.Float64("activated_power", c.state.ActivatedPower),
			logger.Float64("nem_power", c.state.NEMPower),
			logger.Float64("nem_current", c.state.NEMCurrent),
			logger.Int("nem_history_size", len(c.nemAllowedHistory)),
			logger.Float64("total_power", c.state.TotalPower),
			logger.Float64("soc", c.state.SOC))
	}

	// Log NEM ramp progression when ramping
	if c.state.NEMActive && math.Abs(c.state.NEMCurrent) > 0.01 && math.Abs(c.state.NEMCurrent) < 0.99 {
		if time.Now().Second()%5 == 0 {
			c.log.Info("NEM ramping",
				logger.Float64("nem_current", c.state.NEMCurrent),
				logger.Float64("nem_power", c.state.NEMPower),
				logger.Int("samples", len(c.nemAllowedHistory)),
				logger.Float64("progress_pct", math.Abs(c.state.NEMCurrent)*100))
		}
	}
}

// applyPowerLimits applies power limits based on capacity and NEM reservation
func (c *FCRNController) applyPowerLimits(power float64) float64 {
	// Maximum power must account for both FCR capacity and NEM reservation
	// Total reserved power = C_FCR-N + NEM reservation (0.34 * C_FCR-N)
	maxPowerReservation := c.state.Capacity * (1.0 + NEMPowerFactor)

	maxPowerUp := maxPowerReservation
	maxPowerDown := -maxPowerReservation

	if power > maxPowerUp {
		return maxPowerUp
	} else if power < maxPowerDown {
		return maxPowerDown
	}

	return power
}

// GetState returns the current state
func (c *FCRNController) GetState() FCRNState {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.state
}

// SetCapacity updates the FCR-N capacity
func (c *FCRNController) SetCapacity(capacity float64) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if capacity < 0 {
		return fmt.Errorf("capacity must be positive")
	}

	c.log.Info("Updating FCR-N capacity",
		logger.Float64("old_capacity", c.state.Capacity),
		logger.Float64("new_capacity", capacity))

	c.state.Capacity = capacity
	c.config.Capacity = capacity

	// Recalculate SOC thresholds
	c.calculateSOCThresholds()

	return nil
}

// SetDroop updates the droop setting
func (c *FCRNController) SetDroop(droop float64) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if droop <= 0 || droop > 100 {
		return fmt.Errorf("droop must be between 0 and 100")
	}

	c.log.Info("Updating FCR-N droop",
		logger.Float64("old_droop", c.state.Droop),
		logger.Float64("new_droop", droop))

	c.state.Droop = droop
	c.config.Droop = droop

	return nil
}

// GetMaintainedCapacity calculates the maintained capacity (available capacity)
func (c *FCRNController) GetMaintainedCapacity() float64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	maxPower := c.config.MaxPower
	minPower := c.config.MinPower
	baseline := c.state.Baseline

	capacityUp := maxPower - baseline
	capacityDown := baseline - minPower

	maintainedCapacity := math.Min(capacityUp, capacityDown)
	maintainedCapacity = math.Min(maintainedCapacity, c.state.Capacity)

	return math.Max(0, maintainedCapacity)
}
