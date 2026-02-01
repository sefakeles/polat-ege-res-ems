package control

import (
	"go.uber.org/fx"

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
	bmsMgr *bms.Manager,
	pcsMgr *pcs.Manager,
	cfg *config.Config,
) *Logic {
	return NewLogic(bmsMgr, pcsMgr, cfg.EMS)
}
