package health

import (
	"fmt"

	"go.uber.org/fx"

	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/internal/plc"
	"powerkonnekt/ems/internal/windfarm"
)

// Module provides health check functionality to the Fx application
var Module = fx.Module("health",
	fx.Provide(ProvideHealthService),
)

// ProvideHealthService creates and configures a health service instance
func ProvideHealthService(
	bmsMgr *bms.Manager,
	pcsMgr *pcs.Manager,
	plcMgr *plc.Manager,
	windFarmMgr *windfarm.Manager,
) *HealthService {
	healthService := NewHealthService()

	// Register health checkers for all BMS instances
	bmsServices := bmsMgr.GetAllServices()
	for bmsID, bmsService := range bmsServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("bms_%d", bmsID), bmsService))
	}

	// Register health checkers for all PCS instances
	pcsServices := pcsMgr.GetAllServices()
	for pcsID, pcsService := range pcsServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("pcs_%d", pcsID), pcsService))
	}

	// Register health checkers for all PLC instances
	plcServices := plcMgr.GetAllServices()
	for plcID, plcService := range plcServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("plc_%d", plcID), plcService))
	}

	// Register health checkers for all WindFarm instances
	windFarmServices := windFarmMgr.GetAllServices()
	for wfID, wfService := range windFarmServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("windfarm_%d", wfID), wfService))
	}

	return healthService
}
