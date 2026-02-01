package pcs

import (
	"context"

	"go.uber.org/fx"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
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
	alarmMgr *alarm.Manager,
) *Manager {
	return NewManager(cfg.PCS, influxDB, alarmMgr)
}

// RegisterLifecycle registers lifecycle hooks for the PCS manager
func RegisterLifecycle(lc fx.Lifecycle, mgr *Manager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting PCS Manager")
			if err := mgr.Start(); err != nil {
				logger.Error("Failed to start PCS Manager", logger.Err(err))
				return err
			}
			logger.Info("PCS Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping PCS Manager")
			mgr.Stop()
			logger.Info("PCS Manager stopped successfully")
			return nil
		},
	})
}
