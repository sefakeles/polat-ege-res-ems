package metrics

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
)

// Module provides metrics management functionality to the Fx application
var Module = fx.Module("metrics",
	fx.Provide(ProvideManager),
	fx.Invoke(RegisterLifecycle),
)

// ProvideManager creates and provides a metrics manager instance
func ProvideManager(influxDB *database.InfluxDB, logger *zap.Logger) *Manager {
	return NewManager(influxDB, logger)
}

// RegisterLifecycle registers lifecycle hooks for the Metrics manager
func RegisterLifecycle(lc fx.Lifecycle, manager *Manager, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Metrics Manager")
			if err := manager.Start(); err != nil {
				logger.Error("Failed to start Metrics Manager", zap.Error(err))
				return err
			}
			logger.Info("Metrics Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Metrics Manager")
			manager.Stop()
			logger.Info("Metrics Manager stopped successfully")
			return nil
		},
	})
}
