package bms

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
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
	alarmManager *alarm.Manager,
	logger *zap.Logger,
) *Manager {
	return NewManager(cfg.BMS, influxDB, alarmManager, logger)
}

// RegisterLifecycle registers lifecycle hooks for the BMS manager
func RegisterLifecycle(lc fx.Lifecycle, manager *Manager, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting BMS Manager")
			if err := manager.Start(); err != nil {
				logger.Error("Failed to start BMS Manager", zap.Error(err))
				return err
			}
			logger.Info("BMS Manager started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping BMS Manager")
			manager.Stop()
			logger.Info("BMS Manager stopped successfully")
			return nil
		},
	})
}
