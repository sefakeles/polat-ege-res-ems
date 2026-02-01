package metrics

import (
	"context"

	"go.uber.org/fx"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
)

// Module provides metrics management functionality to the Fx application
var Module = fx.Module("metrics",
	fx.Provide(ProvideManager),
	fx.Invoke(RegisterLifecycle),
)

// ProvideManager creates and provides a metrics manager instance
func ProvideManager(influxDB *database.InfluxDB) *Manager {
	return NewManager(influxDB)
}

// RegisterLifecycle registers lifecycle hooks for the Metrics manager
func RegisterLifecycle(lc fx.Lifecycle, mgr *Manager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Metrics Manager")
			if err := mgr.Start(); err != nil {
				logger.Error("Failed to start Metrics Manager", logger.Err(err))
				return err
			}
			logger.Info("Metrics Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Metrics Manager")
			mgr.Stop()
			logger.Info("Metrics Manager stopped successfully")
			return nil
		},
	})
}
