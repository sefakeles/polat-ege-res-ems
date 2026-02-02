package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes configures all API routes
func SetupRoutes(handlers *Handlers, logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(LoggerMiddleware(logger))
	router.Use(CORSMiddleware())
	router.Use(ErrorHandlerMiddleware(logger))
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// System status
		api.GET("/status", handlers.GetStatus)

		// Telemetry endpoint
		api.GET("/telemetry", handlers.GetTelemetry)

		// Data endpoints
		api.GET("/alarms", handlers.GetAlarms)

		// Schedule endpoint
		api.POST("/schedule", handlers.ReceiveSchedule)

		// Control endpoints
		api.POST("/control/mode", handlers.SetControlMode)
		api.POST("/control/active-power", handlers.SetPowerCommand)
		api.POST("/control/reactive-power", handlers.SetReactivePowerCommand)

		// BMS endpoints
		bmsGroup := api.Group("/bms")
		{
			// Data endpoints
			bmsGroup.GET("/data/:id", handlers.GetBMSData)
			bmsGroup.GET("/racks/:id", handlers.GetBMSRacks)
			bmsGroup.GET("/racks/:id/:rack_no", handlers.GetBMSRackData)
			bmsGroup.GET("/command-state/:id", handlers.GetBMSCommandState)

			// Control endpoints
			bmsGroup.POST("/reset", handlers.BMSReset)
			bmsGroup.POST("/breaker", handlers.BMSBreakerControl)
		}

		// PCS endpoints
		pcsGroup := api.Group("/pcs")
		{
			pcsGroup.GET("/data/:id", handlers.GetPCSData)
			pcsGroup.GET("/command-state/:id", handlers.GetPCSCommandState)
			pcsGroup.POST("/start", handlers.SetPCSStartStop)
			pcsGroup.POST("/reset", handlers.PCSReset)
		}

		// PLC endpoints
		plcGroup := api.Group("/plc")
		{
			// Data endpoints
			plcGroup.GET("/data/:id", handlers.GetPLCData)

			// Control endpoints
			plcGroup.POST("/auxiliary-cb", handlers.ControlAuxiliaryCB)
			plcGroup.POST("/mv-aux-transformer-cb", handlers.ControlMVAuxTransformerCB)
			plcGroup.POST("/transformer-cb", handlers.ControlTransformerCB)
			plcGroup.POST("/autoproducer-cb", handlers.ControlAutoproducerCB)
			plcGroup.POST("/reset-all", handlers.ResetAllCircuitBreakers)
		}

		// Wind Farm endpoints
		windFarmGroup := api.Group("/windfarm")
		{
			// Data endpoints
			windFarmGroup.GET("/data/:id", handlers.GetWindFarmData)
			windFarmGroup.GET("/summary", handlers.GetWindFarmSummary)
			windFarmGroup.GET("/command-state/:id", handlers.GetWindFarmCommandState)

			// Control endpoints
			windFarmGroup.POST("/start", handlers.StartWindFarm)
			windFarmGroup.POST("/stop", handlers.StopWindFarm)
			windFarmGroup.POST("/power-setpoint", handlers.SetWindFarmPowerSetpoint)
			windFarmGroup.POST("/reactive-power-setpoint", handlers.SetWindFarmReactivePowerSetpoint)
			windFarmGroup.POST("/power-factor-setpoint", handlers.SetWindFarmPowerFactorSetpoint)
			windFarmGroup.POST("/rapid-downward", handlers.SetWindFarmRapidDownward)
		}
	}

	return router
}
