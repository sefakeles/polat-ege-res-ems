package ems

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
)

// Module provides EMS lifecycle management to the Fx application
var Module = fx.Module("ems",
	fx.Provide(ProvideEMS),
	fx.Invoke(RegisterLifecycle),
)

// ProvideEMS creates and provides an EMS instance
func ProvideEMS(cfg *config.Config, controlLogic *control.Logic, logger *zap.Logger) *EMS {
	return New(cfg.EMS, controlLogic, logger)
}

// RegisterLifecycle registers lifecycle hooks for EMS
func RegisterLifecycle(lc fx.Lifecycle, emsInstance *EMS) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return emsInstance.Start()
		},
		OnStop: func(ctx context.Context) error {
			emsInstance.Stop(ctx)
			return nil
		},
	})
}
