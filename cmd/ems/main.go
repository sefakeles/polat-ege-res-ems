package main

import (
	"go.uber.org/fx"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/api"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/internal/ems"
	"powerkonnekt/ems/internal/health"
	"powerkonnekt/ems/internal/metrics"
	"powerkonnekt/ems/internal/modbus"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/internal/plc"
	"powerkonnekt/ems/internal/windfarm"
	"powerkonnekt/ems/pkg/logger"
)

func main() {
	app := fx.New(
		// Configuration and Logger
		config.Module,
		logger.Module,

		// Database connections
		database.Module,

		// Core services
		alarm.Module,
		metrics.Module,

		// Device managers
		bms.Module,
		pcs.Module,
		plc.Module,
		windfarm.Module,

		// Control logic
		control.Module,

		// Modbus server
		modbus.Module,

		// Health monitoring
		health.Module,

		// API server
		api.Module,

		// EMS lifecycle management
		ems.Module,
	)

	app.Run()
}
