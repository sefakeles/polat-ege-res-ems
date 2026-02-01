package ems

import (
	"go.uber.org/fx"
)

// Module provides EMS core functionality to the Fx application
// Note: Lifecycle management is now handled by individual modules (BMS, PCS, PLC, WindFarm, Modbus)
var Module = fx.Module("ems")
