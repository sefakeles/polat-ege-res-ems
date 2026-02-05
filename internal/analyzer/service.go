package analyzer

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/modbus"
)

// Service represents the energy analyzer service
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

// NewService creates a new energy analyzer service
func NewService(cfg config.AnalyzerConfig, influxDB *database.InfluxDB, logger *zap.Logger) *Service {
	client := modbus.NewClient(cfg.Host, cfg.Port, cfg.SlaveID, cfg.Timeout)
	ctx, cancel := context.WithCancel(context.Background())

	serviceLogger := logger.With(
		zap.String("service", "analyzer"),
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

// Start starts the energy analyzer service
func (s *Service) Start() error {
	s.wg.Go(s.pollLoop)
	s.wg.Go(s.persistenceLoop)

	s.log.Info("Energy analyzer service started",
		zap.Duration("poll_interval", s.config.PollInterval),
		zap.Duration("persist_interval", s.config.PersistInterval))

	return nil
}

// Stop stops the energy analyzer service
func (s *Service) Stop() {
	s.cancel()
	s.wg.Wait()
	s.client.Disconnect()
	s.log.Info("Energy analyzer service stopped")
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.client.IsConnected()
}

// GetDataUpdateChannel returns the channel that signals when new data is available
func (s *Service) GetDataUpdateChannel() <-chan struct{} {
	return s.dataUpdateChan
}

// GetLatestData returns the latest energy analyzer data
func (s *Service) GetLatestData() database.AnalyzerData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastData
}

// GetSystemHealth returns system health information
func (s *Service) GetSystemHealth() map[string]any {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	health := map[string]any{
		"connected":   s.client.IsConnected(),
		"frequency":   s.lastData.Frequency,
		"voltage_l1":  s.lastData.VoltageL1,
		"voltage_l2":  s.lastData.VoltageL2,
		"voltage_l3":  s.lastData.VoltageL3,
		"last_update": s.lastData.Timestamp,
	}

	return health
}

// GetPowerBalance returns power balance information
func (s *Service) GetPowerBalance() map[string]float32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	balance := map[string]float32{
		"active_power_sum":   s.lastData.ActivePowerSum,
		"reactive_power_sum": s.lastData.ReactivePowerSum,
		"apparent_power_sum": s.lastData.ApparentPowerSum,
		"power_factor_avg":   s.lastData.PowerFactorAvg,
		"frequency":          s.lastData.Frequency,
	}

	return balance
}
