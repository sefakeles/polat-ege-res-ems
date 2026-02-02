package windfarm

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Module provides wind farm management functionality to the Fx application
var Module = fx.Module("windfarm",
	fx.Provide(ProvideManager),
	fx.Invoke(RegisterLifecycle),
)

// ProvideManager creates and provides a wind farm manager instance
func ProvideManager(
	cfg *config.Config,
	influxDB *database.InfluxDB,
	logger *zap.Logger,
) *Manager {
	return NewManager(cfg.WindFarm, influxDB, logger)
}

// RegisterLifecycle registers lifecycle hooks for the WindFarm manager
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
