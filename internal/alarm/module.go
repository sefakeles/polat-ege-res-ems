package alarm

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Module provides alarm management functionality to the Fx application
var Module = fx.Module("alarm",
	fx.Provide(ProvideManager),
	fx.Invoke(RegisterLifecycle),
)

// ProvideManager creates and provides an alarm manager instance
func ProvideManager(
	cfg *config.Config,
	postgreSQL *database.PostgresDB,
	logger *zap.Logger,
) *Manager {
	return NewManager(cfg.Alarm, postgreSQL, logger)
}

// RegisterLifecycle registers lifecycle hooks for the alarm manager
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
