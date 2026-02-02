package modbus

import (
	"context"
	"fmt"
	"sync"

	"github.com/simonvetter/modbus"
	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/pcs"
)

// Server represents the Modbus TCP server
type Server struct {
	server  *modbus.ModbusServer
	handler *RequestHandler
	config  config.ModbusServerConfig
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	log     *zap.Logger

	mutex     sync.RWMutex
	isRunning bool
}

// NewServer creates a new Modbus TCP server
func NewServer(
	cfg config.ModbusServerConfig,
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	alarmManager *alarm.Manager,
	controlLogic *control.Logic,
	logger *zap.Logger,
) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create server-specific logger
	serverLogger := logger.With(
		zap.String("component", "modbus_server"),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Uint("max_clients", cfg.MaxClients),
	)

	// Create request handler
	handler := NewRequestHandler(bmsManager, pcsManager, alarmManager, controlLogic, logger)

	// Create server configuration
	serverConfig := &modbus.ServerConfiguration{
		URL:        fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port),
		Timeout:    cfg.Timeout,
		MaxClients: cfg.MaxClients,
	}

	serverLogger.Info("Creating Modbus TCP server",
		zap.String("url", serverConfig.URL),
		zap.Duration("timeout", cfg.Timeout))

	// Create Modbus server
	modbusServer, err := modbus.NewServer(serverConfig, handler)
	if err != nil {
		cancel()
		serverLogger.Error("Failed to create Modbus server", zap.Error(err))
		return nil, fmt.Errorf("failed to create Modbus server: %w", err)
	}

	return &Server{
		server:  modbusServer,
		handler: handler,
		config:  cfg,
		ctx:     ctx,
		cancel:  cancel,
		log:     serverLogger,
	}, nil
}

// Start starts the Modbus server
func (s *Server) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		s.log.Warn("Modbus server start requested but already running")
		return fmt.Errorf("Modbus server is already running")
	}

	s.log.Info("Starting Modbus TCP server")

	if err := s.server.Start(); err != nil {
		s.log.Error("Failed to start Modbus server", zap.Error(err))
		return fmt.Errorf("failed to start Modbus server: %w", err)
	}

	s.isRunning = true
	s.log.Info("Modbus TCP server started successfully")

	return nil
}

// Stop stops the Modbus server
func (s *Server) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return
	}

	s.log.Info("Stopping Modbus TCP server")

	s.cancel()
	s.wg.Wait()

	if s.server != nil {
		s.server.Stop()
	}

	s.isRunning = false
	s.log.Info("Modbus TCP server stopped successfully")
}

// IsRunning returns whether the server is currently running
func (s *Server) IsRunning() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.isRunning
}

// GetStats returns server statistics
func (s *Server) GetStats() map[string]any {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	stats := map[string]any{
		"running":     s.isRunning,
		"host":        s.config.Host,
		"port":        s.config.Port,
		"max_clients": s.config.MaxClients,
		"timeout":     s.config.Timeout.String(),
	}

	return stats
}
