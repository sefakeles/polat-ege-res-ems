package modbus

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/pcs"
)

// Module provides Modbus server functionality to the Fx application
var Module = fx.Module("modbus",
	fx.Provide(ProvideServer),
	fx.Invoke(RegisterLifecycle),
)

// ProvideServer creates and provides a Modbus server instance
func ProvideServer(
	cfg *config.Config,
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	alarmManager *alarm.Manager,
	controlLogic *control.Logic,
	logger *zap.Logger,
) (*Server, error) {
	return NewServer(cfg.ModbusServer, bmsManager, pcsManager, alarmManager, controlLogic, logger)
}

// RegisterLifecycle registers the Modbus server lifecycle hooks with Fx
func RegisterLifecycle(lc fx.Lifecycle, server *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start()
		},
		OnStop: func(ctx context.Context) error {
			server.Stop()
			return nil
		},
	})
}
