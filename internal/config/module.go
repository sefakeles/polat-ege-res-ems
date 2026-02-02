package config

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

// Module provides configuration to the Fx application
var Module = fx.Module("config",
	fx.Provide(
		ProvideValidator,
		ProvideConfig,
	),
)

// ProvideValidator creates and provides a new validator instance
func ProvideValidator() *validator.Validate {
	return NewValidator()
}

// ProvideConfig creates and provides a new configuration instance
func ProvideConfig(validate *validator.Validate) (*Config, error) {
	return NewConfig(validate)
}
