package plc

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Module provides PLC management functionality to the Fx application
var Module = fx.Module("plc",
	fx.Provide(ProvideManager),
	fx.Invoke(RegisterLifecycle),
)

// ProvideManager creates and provides a PLC manager instance
func ProvideManager(
	cfg *config.Config,
	influxDB *database.InfluxDB,
	alarmManager *alarm.Manager,
	logger *zap.Logger,
) *Manager {
	return NewManager(cfg.PLC, influxDB, alarmManager, logger)
}

// RegisterLifecycle registers lifecycle hooks for the PLC manager
func RegisterLifecycle(lc fx.Lifecycle, manager *Manager, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting PLC Manager")
			if err := manager.Start(); err != nil {
				logger.Error("Failed to start PLC Manager", zap.Error(err))
				return err
			}
			logger.Info("PLC Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping PLC Manager")
			manager.Stop()
			logger.Info("PLC Manager stopped successfully")
			return nil
		},
	})
}
