package plc

import (
	"context"

	"go.uber.org/fx"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
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
	alarmMgr *alarm.Manager,
) *Manager {
	return NewManager(cfg.PLC, influxDB, alarmMgr)
}

// RegisterLifecycle registers lifecycle hooks for the PLC manager
func RegisterLifecycle(lc fx.Lifecycle, mgr *Manager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting PLC Manager")
			if err := mgr.Start(); err != nil {
				logger.Error("Failed to start PLC Manager", logger.Err(err))
				return err
			}
			logger.Info("PLC Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping PLC Manager")
			mgr.Stop()
			logger.Info("PLC Manager stopped successfully")
			return nil
		},
	})
}
