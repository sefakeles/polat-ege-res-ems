package windfarm

import (
	"context"
	"sync"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
	"powerkonnekt/ems/pkg/modbus"
)

// Service represents the Wind Farm (FCU) service
type Service struct {
	config   config.WindFarmConfig
	influxDB *database.InfluxDB
	client   *modbus.Client
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mutex    sync.RWMutex
	log      logger.Logger

	// Channel to signal new data availability
	dataUpdateChan chan struct{}

	// Data storage
	lastMeasuringData database.WindFarmMeasuringData
	lastStatusData    database.WindFarmStatusData
	lastSetpointData  database.WindFarmSetpointData
	lastWeatherData   database.WindFarmWeatherData
	commandState      database.WindFarmCommandState

	// Heartbeat counter for heartbeat
	heartbeatCounter uint16
}

// NewService creates a new Wind Farm service
func NewService(cfg config.WindFarmConfig, influxDB *database.InfluxDB) *Service {
	client := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		logger.String("service", "windfarm"),
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
	)

	return &Service{
		config:         cfg,
		influxDB:       influxDB,
		client:         client,
		ctx:            ctx,
		cancel:         cancel,
		log:            serviceLogger,
		dataUpdateChan: make(chan struct{}, 1),
	}
}

// Start starts the Wind Farm service
func (s *Service) Start() error {
	if err := s.client.Connect(s.ctx); err != nil {
		s.log.Warn("Initial Modbus connection failed", logger.Err(err))
	}

	s.wg.Go(s.dataPollLoop)
	s.wg.Go(s.heartbeatLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("Wind Farm service started",
		logger.Int("id", s.config.ID),
		logger.String("host", s.config.Host))

	return nil
}

// Stop stops the Wind Farm service
func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	s.client.Disconnect()
	s.log.Info("Wind Farm service stopped")
}
