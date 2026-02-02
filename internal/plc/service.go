package plc

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/modbus"
)

// Service represents the PLC service
type Service struct {
	config       config.PLCConfig
	influxDB     *database.InfluxDB
	alarmManager *alarm.Manager
	client       *modbus.Client
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	log          *zap.Logger

	dataUpdateChan chan struct{}

	mutex               sync.RWMutex
	lastPLCData         database.PLCData
	previousRelayStates map[string]bool
}

// NewService creates a new PLC service
func NewService(cfg config.PLCConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager, logger *zap.Logger) *Service {
	client := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		zap.String("service", "plc"),
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
		previousRelayStates: make(map[string]bool),
	}
}

// Start starts the PLC service
func (s *Service) Start() error {
	s.wg.Go(s.pollLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("PLC service started")

	return nil
}

// Stop stops the PLC service
func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	s.client.Disconnect()
	s.log.Info("PLC service stopped")
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.client.IsConnected()
}

// GetDataUpdateChannel returns the channel that signals when new data is available
func (s *Service) GetDataUpdateChannel() <-chan struct{} {
	return s.dataUpdateChan
}

// GetLatestPLCData returns the latest PLC data
func (s *Service) GetLatestPLCData() database.PLCData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastPLCData
}
