package config

import (
	"fmt"
	"time"
)

// FCRNConfig represents FCR-N configuration
type FCRNConfig struct {
	// Basic settings
	Enabled  bool    `mapstructure:"enabled"`  // Enable/disable FCR-N
	Capacity float64 `mapstructure:"capacity"` // kW - sold FCR-N capacity
	Droop    float64 `mapstructure:"droop"`    // % - droop setting (Hz/%)

	// Power limits
	MaxPower float64 `mapstructure:"max_power"` // kW - maximum power output
	MinPower float64 `mapstructure:"min_power"` // kW - minimum power output

	// Control parameters
	UpdateInterval time.Duration `mapstructure:"update_interval"` // Control loop update interval

	// Energy management (for LER - Limited Energy Reservoir entities)
	EnableEnergyManagement bool    `mapstructure:"enable_energy_management"` // Enable NEM/AEM
	ReservoirCapacity      float64 `mapstructure:"reservoir_capacity"`       // kWh - battery capacity
	MinSOC                 float64 `mapstructure:"min_soc"`                  // % - minimum state of charge
	MaxSOC                 float64 `mapstructure:"max_soc"`                  // % - maximum state of charge

	// Frequency measurement
	FrequencySource       string        `mapstructure:"frequency_source"`        // Source of frequency measurement
	FrequencyUpdateRate   time.Duration `mapstructure:"frequency_update_rate"`   // How often to update frequency
	FrequencySmoothWindow int           `mapstructure:"frequency_smooth_window"` // Moving average window
	PCSNumber             int           `mapstructure:"pcs_number"`              // PCS number for frequency measurement

	// Activation behavior
	SmoothActivation     bool          `mapstructure:"smooth_activation"`      // Enable smooth activation
	ActivationRampTime   time.Duration `mapstructure:"activation_ramp_time"`   // Time to ramp to full activation
	DeactivationRampTime time.Duration `mapstructure:"deactivation_ramp_time"` // Time to ramp down

	// Telemetry and logging
	EnableTelemetry    bool          `mapstructure:"enable_telemetry"`     // Send data to TSO
	TelemetryInterval  time.Duration `mapstructure:"telemetry_interval"`   // How often to send telemetry
	LogLevel           string        `mapstructure:"log_level"`            // Log level for FCR-N controller
	DataLoggingEnabled bool          `mapstructure:"data_logging_enabled"` // Enable data logging
}

// Validate validates FCR-N configuration
func (c *FCRNConfig) Validate() error {
	if c.Capacity <= 0 {
		return fmt.Errorf("capacity must be positive")
	}

	if c.Droop < 0 || c.Droop > 100 {
		return fmt.Errorf("droop must be between 0 and 100")
	}

	if c.MaxPower <= c.MinPower {
		return fmt.Errorf("max_power must be greater than min_power")
	}

	if c.EnableEnergyManagement {
		if c.ReservoirCapacity <= 0 {
			return fmt.Errorf("reservoir_capacity must be positive when energy management is enabled")
		}

		if c.MinSOC < 0 || c.MinSOC > 100 {
			return fmt.Errorf("min_soc must be between 0 and 100")
		}

		if c.MaxSOC < 0 || c.MaxSOC > 100 {
			return fmt.Errorf("max_soc must be between 0 and 100")
		}

		if c.MinSOC >= c.MaxSOC {
			return fmt.Errorf("min_soc must be less than max_soc")
		}

		// Check if we have enough energy for 1 hour endurance (FCR-N requirement)
		requiredEnergy := c.Capacity * 1.0 // 1 hour * capacity
		if c.ReservoirCapacity < requiredEnergy {
			return fmt.Errorf("reservoir_capacity (%.2f kWh) is less than required for 1 hour endurance (%.2f kWh)",
				c.ReservoirCapacity, requiredEnergy)
		}
	}

	return nil
}
