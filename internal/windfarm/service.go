package windfarm

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
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
	log      *zap.Logger

	dataUpdateChan chan struct{}

	mutex             sync.RWMutex
	lastMeasuringData database.WindFarmMeasuringData
	lastStatusData    database.WindFarmStatusData
	lastSetpointData  database.WindFarmSetpointData
	lastWeatherData   database.WindFarmWeatherData
	commandState      database.WindFarmCommandState
	heartbeatCounter  uint16
}

// NewService creates a new Wind Farm service
func NewService(cfg config.WindFarmConfig, influxDB *database.InfluxDB, logger *zap.Logger) *Service {
	client := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	ctx, cancel := context.WithCancel(context.Background())

	// Create service-specific logger
	serviceLogger := logger.With(
		zap.String("service", "windfarm"),
		zap.Int("id", cfg.ID),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
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
	s.wg.Go(s.dataPollLoop)
	s.wg.Go(s.heartbeatLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("Wind Farm service started")

	return nil
}

// Stop stops the Wind Farm service
func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	s.client.Disconnect()
	s.log.Info("Wind Farm service stopped")
}
