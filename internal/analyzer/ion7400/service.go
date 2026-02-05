package ion7400

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/modbus"
)

// Service represents the ION7400 service
type Service struct {
	config   config.AnalyzerConfig
	influxDB *database.InfluxDB
	client   *modbus.Client
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	log      *zap.Logger

	dataUpdateChan chan struct{}

	mutex    sync.RWMutex
	lastData database.AnalyzerData
}

// NewService creates a new ION7400 service
func NewService(cfg config.AnalyzerConfig, influxDB *database.InfluxDB, logger *zap.Logger) *Service {
	client := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	ctx, cancel := context.WithCancel(context.Background())

	serviceLogger := logger.With(
		zap.String("service", "ion7400"),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port))

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

// Start starts the ION7400 service
func (s *Service) Start() error {
	s.wg.Go(s.pollLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("ION7400 service started",
		zap.Duration("poll_interval", s.config.PollInterval),
		zap.Duration("persist_interval", s.config.PersistInterval))

	return nil
}

// Stop stops the ION7400 service
func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	s.client.Disconnect()
	s.log.Info("ION7400 service stopped")
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.client.IsConnected()
}

// GetDataUpdateChannel returns the channel that signals when new data is available
func (s *Service) GetDataUpdateChannel() <-chan struct{} {
	return s.dataUpdateChan
}

// GetLatestData returns the latest ION7400 data
func (s *Service) GetLatestData() database.AnalyzerData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastData
}
