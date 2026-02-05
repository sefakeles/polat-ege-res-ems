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

// GetLatestData returns the latest aggregated wind farm data
func (s *Service) GetLatestData() database.WindFarmData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return database.WindFarmData{
		MeasuringData: s.lastMeasuringData,
		StatusData:    s.lastStatusData,
		SetpointData:  s.lastSetpointData,
		WeatherData:   s.lastWeatherData,
	}
}

// GetLatestMeasuringData returns the latest measuring data
func (s *Service) GetLatestMeasuringData() database.WindFarmMeasuringData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastMeasuringData
}

// GetLatestStatusData returns the latest status data
func (s *Service) GetLatestStatusData() database.WindFarmStatusData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastStatusData
}

// GetLatestSetpointData returns the latest setpoint data
func (s *Service) GetLatestSetpointData() database.WindFarmSetpointData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastSetpointData
}

// GetLatestWeatherData returns the latest weather data
func (s *Service) GetLatestWeatherData() database.WindFarmWeatherData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastWeatherData
}

// GetCommandState returns the current command state
func (s *Service) GetCommandState() database.WindFarmCommandState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.commandState
}

// GetDataUpdateChan returns the data update channel
func (s *Service) GetDataUpdateChan() <-chan struct{} {
	return s.dataUpdateChan
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.client.IsConnected()
}

// IsFCUOnline returns whether the FCU is online
func (s *Service) IsFCUOnline() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastStatusData.FCUOnline
}

// GetConfig returns the service configuration
func (s *Service) GetConfig() config.WindFarmConfig {
	return s.config
}
