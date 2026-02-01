package config

import "go.uber.org/fx"

// Module provides configuration to the Fx application
var Module = fx.Module("config",
	fx.Provide(ProvideConfig),
)

// ProvideConfig loads and provides the application configuration
func ProvideConfig() (*Config, error) {
	return Load("configs/config.json")
}
