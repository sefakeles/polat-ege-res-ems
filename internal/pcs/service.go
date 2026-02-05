package pcs

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/modbus"
)

// Service represents the PCS service
type Service struct {
	config       config.PCSConfig
	influxDB     *database.InfluxDB
	alarmManager *alarm.Manager
	client       *modbus.Client
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	log          *zap.Logger

	dataUpdateChan chan struct{}

	mutex               sync.RWMutex
	lastStatusData      database.PCSStatusData
	lastEquipmentData   database.PCSEquipmentData
	lastEnvironmentData database.PCSEnvironmentData
	lastDCSourceData    database.PCSDCSourceData
	lastGridData        database.PCSGridData
	lastCounterData     database.PCSCounterData
	commandState        database.PCSCommandState
	previousAlarmStates map[string]bool
	heartbeatCount      uint16
}

// NewService creates a new PCS service
func NewService(cfg config.PCSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager, logger *zap.Logger) *Service {
	client := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		zap.String("service", "pcs"),
		zap.Int("id", cfg.ID),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return &Service{
		config:              cfg,
		influxDB:            influxDB,
		alarmManager:        alarmManager,
		client:              client,
		ctx:                 ctx,
		cancel:              cancel,
		log:                 serviceLogger,
		dataUpdateChan:      make(chan struct{}, 1),
		previousAlarmStates: make(map[string]bool),
	}
}

// Start starts the PCS service
func (s *Service) Start() error {
	s.wg.Go(s.pollLoop)
	s.wg.Go(s.heartbeatLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("PCS service started")

	return nil
}

// Stop stops the PCS service
func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	s.client.Disconnect()
	s.log.Info("PCS service stopped")
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.client.IsConnected()
}

// GetDataUpdateChannel returns the channel that signals when new data is available
func (s *Service) GetDataUpdateChannel() <-chan struct{} {
	return s.dataUpdateChan
}

// GetLatestPCSData returns the latest PCS data
func (s *Service) GetLatestPCSStatusData() database.PCSStatusData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastStatusData
}

// GetLatestPCSEquipmentData returns the latest PCS equipment data
func (s *Service) GetLatestPCSEquipmentData() database.PCSEquipmentData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastEquipmentData
}

// GetLatestPCSEnvironmentData returns the latest PCS environment data
func (s *Service) GetLatestPCSEnvironmentData() database.PCSEnvironmentData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastEnvironmentData
}

// GetLatestPCSDCSourceData returns the latest PCS DC source data
func (s *Service) GetLatestPCSDCSourceData() database.PCSDCSourceData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastDCSourceData
}

// GetLatestPCSGridData returns the latest PCS grid data
func (s *Service) GetLatestPCSGridData() database.PCSGridData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastGridData
}

// GetLatestPCSCounterData returns the latest PCS counter data
func (s *Service) GetLatestPCSCounterData() database.PCSCounterData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastCounterData
}

func (s *Service) GetLatestPCSData() database.PCSData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return database.PCSData{
		StatusData:      s.lastStatusData,
		EquipmentData:   s.lastEquipmentData,
		EnvironmentData: s.lastEnvironmentData,
		DCSourceData:    s.lastDCSourceData,
		GridData:        s.lastGridData,
		CounterData:     s.lastCounterData,
	}
}

// GetCommandState returns the current command state
func (s *Service) GetCommandState() database.PCSCommandState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.commandState
}
