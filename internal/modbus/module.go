package modbus

import (
	"context"

	"go.uber.org/fx"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/pkg/logger"
)

// Module provides Modbus server functionality to the Fx application
var Module = fx.Module("modbus",
	fx.Provide(ProvideServer),
	fx.Invoke(RegisterLifecycle),
)

// ProvideServer creates and provides a Modbus server instance
func ProvideServer(
	cfg *config.Config,
	bmsMgr *bms.Manager,
	pcsMgr *pcs.Manager,
	alarmMgr *alarm.Manager,
	controlLogic *control.Logic,
) (*Server, error) {
	return NewServer(cfg.ModbusServer, bmsMgr, pcsMgr, alarmMgr, controlLogic)
}

// RegisterLifecycle registers lifecycle hooks for the Modbus server
func RegisterLifecycle(lc fx.Lifecycle, server *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Modbus Server")
			if err := server.Start(); err != nil {
				logger.Error("Failed to start Modbus Server", logger.Err(err))
				return err
			}
			logger.Info("Modbus Server started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Modbus Server")
			server.Stop()
			logger.Info("Modbus Server stopped successfully")
			return nil
		},
	})
}
