package pcs

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Module provides PCS management functionality to the Fx application
var Module = fx.Module("pcs",
	fx.Provide(ProvideManager),
	fx.Invoke(RegisterLifecycle),
)

// ProvideManager creates and provides a PCS manager instance
func ProvideManager(
	cfg *config.Config,
	influxDB *database.InfluxDB,
	alarmManager *alarm.Manager,
	logger *zap.Logger,
) *Manager {
	return NewManager(cfg.PCS, influxDB, alarmManager, logger)
}

// RegisterLifecycle registers lifecycle hooks for the PCS manager
func RegisterLifecycle(lc fx.Lifecycle, manager *Manager, logger *zap.Logger) {
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
