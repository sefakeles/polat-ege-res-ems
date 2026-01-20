package bms

import (
	"context"
	"sync"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
	"powerkonnekt/ems/pkg/modbus"
)

// Service represents the BMS service
type Service struct {
	config       config.BMSConfig
	influxDB     *database.InfluxDB
	alarmManager *alarm.Manager
	baseClient   *modbus.Client
	cellClient   *modbus.Client
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mutex        sync.RWMutex
	log          logger.Logger

	// Channels to signal new data availability
	baseDataUpdateChan chan struct{}
	cellDataUpdateChan chan struct{}

	// Data storage
	lastBMSData          database.BMSData
	lastBMSStatusData    database.BMSStatusData
	lastBMSRackData      []database.BMSRackData
	lastCellVoltages     [][]database.BMSCellVoltageData
	lastCellTemperatures [][]database.BMSCellTemperatureData
	commandState         database.BMSCommandState

	// Heartbeat counter
	heartbeatCount uint16
}

// NewService creates a new BMS service
func NewService(cfg config.BMSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager) *Service {
	baseClient := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	cellClient := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)

	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		logger.String("service", "bms"),
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
	)

	return &Service{
		config:               cfg,
		influxDB:             influxDB,
		alarmManager:         alarmManager,
		baseClient:           baseClient,
		cellClient:           cellClient,
		ctx:                  ctx,
		cancel:               cancel,
		log:                  serviceLogger,
		baseDataUpdateChan:   make(chan struct{}, 1),
		cellDataUpdateChan:   make(chan struct{}, 1),
		lastBMSRackData:      make([]database.BMSRackData, cfg.RackCount),
		lastCellVoltages:     make([][]database.BMSCellVoltageData, cfg.RackCount),
		lastCellTemperatures: make([][]database.BMSCellTemperatureData, cfg.RackCount),
	}
}

// Start starts the BMS service
func (s *Service) Start() error {
	s.wg.Go(s.baseDataPollLoop)
	s.wg.Go(s.cellDataPollLoop)
	s.wg.Go(s.heartbeatLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("BMS service started",
		logger.Int("rack_count", s.config.RackCount),
		logger.Bool("enable_cell_data", s.config.EnableCellData))

	return nil
}

// Stop stops the BMS service
func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	s.baseClient.Disconnect()
	s.cellClient.Disconnect()
	s.log.Info("BMS service stopped")
}
