package main

import (
	"go.uber.org/fx"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/analyzer/ion7400"
	"powerkonnekt/ems/internal/api"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/internal/ems"
	"powerkonnekt/ems/internal/health"
	"powerkonnekt/ems/internal/logger"
	"powerkonnekt/ems/internal/metrics"
	"powerkonnekt/ems/internal/modbus"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/internal/plc"
	"powerkonnekt/ems/internal/windfarm"
)

func main() {
	app := fx.New(
		// Configuration
		config.Module,

		// Logging
		logger.Module,
		logger.FxLogger,

		// Database
		database.Module,

		// Core services
		alarm.Module,
		metrics.Module,

		// Device managers
		bms.Module,
		pcs.Module,
		plc.Module,
		windfarm.Module,
		ion7400.Module,

		// Control logic
		control.Module,

		// Modbus server
		modbus.Module,

		// Health monitoring
		health.Module,

		// API server
		api.Module,

		// EMS
		ems.Module,
	)

	app.Run()
}
