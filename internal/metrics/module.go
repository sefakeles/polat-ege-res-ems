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
func ProvideManager(lc fx.Lifecycle, influxDB *database.InfluxDB, logger *zap.Logger) *Manager {
	return NewManager(influxDB, logger)
}

// RegisterLifecycle registers lifecycle hooks for the metrics manager
func RegisterLifecycle(lc fx.Lifecycle, manager *Manager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return manager.Start()
		},
		OnStop: func(ctx context.Context) error {
			manager.Stop()
			return nil
		},
	})
}
