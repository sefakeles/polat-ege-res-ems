package bms

import (
	"context"

	"go.uber.org/fx"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
)

// Module provides BMS management functionality to the Fx application
var Module = fx.Module("bms",
	fx.Provide(ProvideManager),
	fx.Invoke(RegisterLifecycle),
)

// ProvideManager creates and provides a BMS manager instance
func ProvideManager(
	cfg *config.Config,
	influxDB *database.InfluxDB,
	alarmMgr *alarm.Manager,
) *Manager {
	return NewManager(cfg.BMS, influxDB, alarmMgr)
}

// RegisterLifecycle registers lifecycle hooks for the BMS manager
func RegisterLifecycle(lc fx.Lifecycle, mgr *Manager) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting BMS Manager")
			if err := mgr.Start(); err != nil {
				logger.Error("Failed to start BMS Manager", logger.Err(err))
				return err
			}
			logger.Info("BMS Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping BMS Manager")
			mgr.Stop()
			logger.Info("BMS Manager stopped successfully")
			return nil
		},
	})
}
