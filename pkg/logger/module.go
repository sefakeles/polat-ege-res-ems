package logger

import (
	"go.uber.org/fx"

	"powerkonnekt/ems/internal/config"
)

// Module provides logger functionality to the Fx application
var Module = fx.Module("logger",
	fx.Invoke(InitLogger),
)

// InitLogger initializes the logger with configuration
func InitLogger(cfg *config.Config) error {
	return InitializeWithConfig(Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	})
}
