package bms

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
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
	log          *zap.Logger

	baseDataUpdateChan chan struct{}
	cellDataUpdateChan chan struct{}

	mutex                sync.RWMutex
	lastBMSData          database.BMSData
	lastBMSStatusData    database.BMSStatusData
	lastBMSRackData      []database.BMSRackData
	lastCellVoltages     [][]database.BMSCellVoltageData
	lastCellTemperatures [][]database.BMSCellTemperatureData
	commandState         database.BMSCommandState
	previousAlarmStates  map[string]bool
	heartbeatCount       uint16
}

// NewService creates a new BMS service
func NewService(cfg config.BMSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager, logger *zap.Logger) *Service {
	baseClient := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	cellClient := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)

	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		zap.String("service", "bms"),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
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
		previousAlarmStates:  make(map[string]bool),
	}
}

// Start starts the BMS service
func (s *Service) Start() error {
	s.wg.Go(s.baseDataPollLoop)
	if s.config.EnableCellData {
		s.wg.Go(s.cellDataPollLoop)
	}
	s.wg.Go(s.heartbeatLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("BMS service started",
		zap.Int("rack_count", s.config.RackCount),
		zap.Bool("enable_cell_data", s.config.EnableCellData))

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
