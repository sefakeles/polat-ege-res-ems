package windfarm

import (
	"context"

	"go.uber.org/fx"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
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
) *Manager {
	return NewManager(cfg.WindFarm, influxDB)
}

// RegisterLifecycle registers lifecycle hooks for the WindFarm manager
func RegisterLifecycle(lc fx.Lifecycle, mgr *Manager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting WindFarm Manager")
			if err := mgr.Start(); err != nil {
				logger.Error("Failed to start WindFarm Manager", logger.Err(err))
				return err
			}
			logger.Info("WindFarm Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping WindFarm Manager")
			mgr.Stop()
			logger.Info("WindFarm Manager stopped successfully")
			return nil
		},
	})
}
