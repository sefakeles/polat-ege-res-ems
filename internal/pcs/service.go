package pcs

import (
	"context"
	"sync"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
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
	mutex        sync.RWMutex
	log          logger.Logger

	// Channel to signal new data availability
	dataUpdateChan chan struct{}

	// Data storage
	lastStatusData      database.PCSStatusData
	lastEquipmentData   database.PCSEquipmentData
	lastEnvironmentData database.PCSEnvironmentData
	lastDCSourceData    database.PCSDCSourceData
	lastGridData        database.PCSGridData
	lastCounterData     database.PCSCounterData
	commandState        database.PCSCommandState

	// Heartbeat counter
	heartbeatCount uint16
}

// NewService creates a new PCS service
func NewService(cfg config.PCSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager) *Service {
	client := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		logger.String("service", "pcs"),
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
	)

	return &Service{
		config:         cfg,
		influxDB:       influxDB,
		alarmManager:   alarmManager,
		client:         client,
		ctx:            ctx,
		cancel:         cancel,
		log:            serviceLogger,
		dataUpdateChan: make(chan struct{}, 1),
	}
}

// Start starts the PCS service
func (s *Service) Start() error {
	if err := s.client.Connect(s.ctx); err != nil {
		s.log.Warn("Initial Modbus connection failed", logger.Err(err))
	}

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
