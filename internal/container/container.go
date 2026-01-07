package container

import (
	"fmt"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/internal/fcr"
	"powerkonnekt/ems/internal/metrics"
	"powerkonnekt/ems/internal/modbus"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/internal/plc"
	"powerkonnekt/ems/internal/windfarm"
	"powerkonnekt/ems/pkg/logger"
)

type Container struct {
	Config          *config.Config
	InfluxDB        *database.InfluxDB
	PostgresDB      *database.PostgresDB
	BMSManager      *bms.Manager
	PCSManager      *pcs.Manager
	PLCManager      *plc.Manager
	WindFarmManager *windfarm.Manager
	FCRNService     *fcr.Service
	ControlLogic    *control.Logic
	AlarmManager    *alarm.Manager
	MetricsManager  *metrics.Manager
	ModbusServer    *modbus.Server
	log             logger.Logger
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Create container-specific logger
	containerLogger := logger.With(
		logger.String("component", "container"),
	)

	containerLogger.Info("Initializing dependency injection container")

	// Initialize databases
	influxDB, err := database.InitializeInfluxDB(cfg.InfluxDB)
	if err != nil {
		containerLogger.Error("Failed to initialize InfluxDB", logger.Err(err))
		return nil, fmt.Errorf("failed to initialize InfluxDB: %w", err)
	}

	postgresDB, err := database.InitializePostgreSQL(cfg.PostgreSQL)
	if err != nil {
		containerLogger.Error("Failed to initialize PostgreSQL", logger.Err(err))
		return nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	// Initialize managers
	alarmManager := alarm.NewManager(postgresDB)
	metricsManager := metrics.NewManager(influxDB)

	// Initialize BMS and PCS managers
	bmsManager := bms.NewManager(cfg.BMS, influxDB, alarmManager)
	pcsManager := pcs.NewManager(cfg.PCS, influxDB, alarmManager)
	plcManager := plc.NewManager(cfg.PLC, influxDB, alarmManager)
	windFarmManager := windfarm.NewManager(cfg.WindFarm, influxDB)

	// Initialize FCR-N service if enabled
	var fcrnService *fcr.Service
	if cfg.FCRN.Enabled {
		var freqSource fcr.FrequencySource

		// Create frequency source based on configuration
		switch cfg.FCRN.FrequencySource {
		case "test":
			freqSource = fcr.NewTestFrequencySource()
			containerLogger.Info("Using test frequency source for FCR-N")
		case "pcs":
			freqSource = fcr.NewPCSFrequencySource(pcsManager, uint8(cfg.FCRN.PCSNumber))
			containerLogger.Info("Using PCS frequency source for FCR-N",
				logger.Int("pcs_number", cfg.FCRN.PCSNumber))
		default:
			containerLogger.Warn("Unknown frequency source, using analyzer",
				logger.String("source", cfg.FCRN.FrequencySource))
			freqSource = fcr.NewPCSFrequencySource(pcsManager, uint8(cfg.FCRN.PCSNumber))
		}

		fcrnService, err = fcr.NewService(cfg.FCRN, pcsManager, bmsManager, freqSource)
		if err != nil {
			return nil, fmt.Errorf("failed to create FCR-N service: %w", err)
		}
	}

	// Initialize control logic with managers
	controlLogic := control.NewLogic(bmsManager, pcsManager, cfg.EMS)

	modbusServer, err := modbus.NewServer(cfg.ModbusServer, bmsManager, pcsManager, alarmManager, controlLogic)
	if err != nil {
		containerLogger.Error("Failed to initialize Modbus server", logger.Err(err))
		return nil, fmt.Errorf("failed to initialize Modbus server: %w", err)
	}

	container := &Container{
		Config:          cfg,
		InfluxDB:        influxDB,
		PostgresDB:      postgresDB,
		BMSManager:      bmsManager,
		PCSManager:      pcsManager,
		PLCManager:      plcManager,
		WindFarmManager: windFarmManager,
		FCRNService:     fcrnService,
		ControlLogic:    controlLogic,
		AlarmManager:    alarmManager,
		MetricsManager:  metricsManager,
		ModbusServer:    modbusServer,
		log:             containerLogger,
	}

	containerLogger.Info("Dependency injection container initialized successfully")
	return container, nil
}

func (c *Container) Close() error {
	c.log.Info("Closing container and releasing resources")

	var lastErr error

	if c.InfluxDB != nil {
		if err := c.InfluxDB.Close(); err != nil {
			c.log.Error("Failed to close InfluxDB", logger.Err(err))
			lastErr = err
		}
	}

	if c.PostgresDB != nil {
		if err := c.PostgresDB.Close(); err != nil {
			c.log.Error("Failed to close PostgreSQL", logger.Err(err))
			lastErr = err
		}
	}

	if lastErr != nil {
		c.log.Error("Container closed with errors", logger.Err(lastErr))
	} else {
		c.log.Info("Container closed successfully")
	}

	return lastErr
}
