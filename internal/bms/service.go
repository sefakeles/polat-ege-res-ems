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
	systemClient *modbus.Client
	cellClient   *modbus.Client
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	log          *zap.Logger

	systemDataUpdateChan chan struct{}
	cellDataUpdateChan   chan struct{}

	mutex                 sync.RWMutex
	lastBMSData           database.BMSData
	lastBMSStatusData     database.BMSStatusData
	lastBMSRackData       []database.BMSRackData
	lastBMSRackStatusData []database.BMSRackStatusData
	lastCellVoltages      [][]database.BMSCellVoltageData
	lastCellTemperatures  [][]database.BMSCellTemperatureData
	commandState          database.BMSCommandState
	previousAlarmStates   map[string]bool
	heartbeatCount        uint16
}

// NewService creates a new BMS service
func NewService(cfg config.BMSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager, logger *zap.Logger) *Service {
	systemClient := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	cellClient := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)

	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		zap.String("service", "bms"),
		zap.Int("id", cfg.ID),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return &Service{
		config:                cfg,
		influxDB:              influxDB,
		alarmManager:          alarmManager,
		systemClient:          systemClient,
		cellClient:            cellClient,
		ctx:                   ctx,
		cancel:                cancel,
		log:                   serviceLogger,
		systemDataUpdateChan:  make(chan struct{}, 1),
		cellDataUpdateChan:    make(chan struct{}, 1),
		lastBMSRackData:       make([]database.BMSRackData, cfg.RackCount),
		lastBMSRackStatusData: make([]database.BMSRackStatusData, cfg.RackCount),
		lastCellVoltages:      make([][]database.BMSCellVoltageData, cfg.RackCount),
		lastCellTemperatures:  make([][]database.BMSCellTemperatureData, cfg.RackCount),
		previousAlarmStates:   make(map[string]bool),
	}
}

// Start starts the BMS service
func (s *Service) Start() error {
	s.wg.Go(s.systemDataPollLoop)
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
	s.systemClient.Disconnect()
	s.cellClient.Disconnect()
	s.log.Info("BMS service stopped")
}

// GetTotalCellsPerRack returns the total number of cells per rack
func (s *Service) GetTotalCellsPerRack() int {
	return s.config.ModulesPerRack * CellsPerModule
}

// GetTotalTempSensorsPerRack returns the total number of temperature sensors per rack
func (s *Service) GetTotalTempSensorsPerRack() int {
	return s.config.ModulesPerRack * TempSensorsPerModule
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.systemClient.IsConnected()
}

// GetSystemDataUpdateChannel returns the channel that signals when new system data is available
func (s *Service) GetSystemDataUpdateChannel() <-chan struct{} {
	return s.systemDataUpdateChan
}

// GetCellDataUpdateChannel returns the channel that signals when new cell data is available
func (s *Service) GetCellDataUpdateChannel() <-chan struct{} {
	return s.cellDataUpdateChan
}

// GetLatestBMSData returns the latest BMS data
func (s *Service) GetLatestBMSData() database.BMSData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastBMSData
}

// GetLatestBMSStatusData returns the latest BMS status data
func (s *Service) GetLatestBMSStatusData() database.BMSStatusData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastBMSStatusData
}

// GetLatestBMSRackData returns the latest BMS rack data
func (s *Service) GetLatestBMSRackData() []database.BMSRackData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return append([]database.BMSRackData(nil), s.lastBMSRackData...)
}

// GetLatestBMSRackStatusData returns the latest BMS rack status data
func (s *Service) GetLatestBMSRackStatusData() []database.BMSRackStatusData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return append([]database.BMSRackStatusData(nil), s.lastBMSRackStatusData...)
}

// GetLatestCellVoltageData returns the latest cell voltage data for a specific rack
func (s *Service) GetLatestCellVoltageData(rackNo uint8) []database.BMSCellVoltageData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var voltages []database.BMSCellVoltageData

	if rackNo > 0 && int(rackNo) <= len(s.lastCellVoltages) {
		voltages = make([]database.BMSCellVoltageData, len(s.lastCellVoltages[rackNo-1]))
		copy(voltages, s.lastCellVoltages[rackNo-1])
	}

	return voltages
}

// GetLatestCellTemperatureData returns the latest cell temperature data for a specific rack
func (s *Service) GetLatestCellTemperatureData(rackNo uint8) []database.BMSCellTemperatureData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var temperatures []database.BMSCellTemperatureData

	if rackNo > 0 && int(rackNo) <= len(s.lastCellTemperatures) {
		temperatures = make([]database.BMSCellTemperatureData, len(s.lastCellTemperatures[rackNo-1]))
		copy(temperatures, s.lastCellTemperatures[rackNo-1])
	}

	return temperatures
}

// GetCommandState returns the current command state
func (s *Service) GetCommandState() database.BMSCommandState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.commandState
}
