package control

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/pcs"
)

// Module provides control logic functionality to the Fx application
var Module = fx.Module("control",
	fx.Provide(ProvideLogic),
)

// ProvideLogic creates and provides a control logic instance
func ProvideLogic(
	cfg *config.Config,
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	logger *zap.Logger,
) *Logic {
	return NewLogic(cfg.EMS, bmsManager, pcsManager, logger)
}
