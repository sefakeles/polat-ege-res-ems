package ion7400

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Module provides ION7400 management functionality to the Fx application
var Module = fx.Module("ion7400",
	fx.Provide(ProvideService),
	fx.Invoke(RegisterLifecycle),
)

// ProvideService creates and provides an ION7400 service instance
func ProvideService(
	cfg *config.Config,
	influxDB *database.InfluxDB,
	logger *zap.Logger,
) *Service {
	return NewService(cfg.ION7400, influxDB, logger)
}

// RegisterLifecycle registers lifecycle hooks for the ION7400 service
func RegisterLifecycle(lc fx.Lifecycle, service *Service) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return service.Start()
		},
		OnStop: func(ctx context.Context) error {
			service.Stop()
			return nil
		},
	})
}
