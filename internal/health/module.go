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
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	plcManager *plc.Manager,
	windFarmManager *windfarm.Manager,
) *HealthService {
	healthService := NewService()

	// Register health checkers for all BMS instances
	bmsServices := bmsManager.GetAllServices()
	for bmsID, bmsService := range bmsServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("bms_%d", bmsID), bmsService))
	}

	// Register health checkers for all PCS instances
	pcsServices := pcsManager.GetAllServices()
	for pcsID, pcsService := range pcsServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("pcs_%d", pcsID), pcsService))
	}

	// Register health checkers for all PLC instances
	plcServices := plcManager.GetAllServices()
	for plcID, plcService := range plcServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("plc_%d", plcID), plcService))
	}

	// Register health checkers for all WindFarm instances
	windFarmServices := windFarmManager.GetAllServices()
	for wfID, wfService := range windFarmServices {
		healthService.RegisterChecker(NewServiceChecker(fmt.Sprintf("windfarm_%d", wfID), wfService))
	}

	return healthService
}
