package modbus

import (
	"context"
	"fmt"
	"sync"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/pkg/logger"

	"github.com/simonvetter/modbus"
)

// Server represents the Modbus TCP server
type Server struct {
	server    *modbus.ModbusServer
	handler   *RequestHandler
	config    config.ModbusServerConfig
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	isRunning bool
	mutex     sync.RWMutex
	log       logger.Logger
}

// NewServer creates a new Modbus TCP server
func NewServer(
	cfg config.ModbusServerConfig,
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	alarmManager *alarm.Manager,
	controlLogic *control.Logic,
) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create server-specific logger
	serverLogger := logger.With(
		logger.String("component", "modbus_server"),
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
		logger.Uint("max_clients", cfg.MaxClients),
	)

	// Create request handler
	handler := NewRequestHandler(bmsManager, pcsManager, alarmManager, controlLogic)

	// Create server configuration
	serverConfig := &modbus.ServerConfiguration{
		URL:        fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port),
		Timeout:    cfg.Timeout,
		MaxClients: cfg.MaxClients,
	}

	serverLogger.Info("Creating Modbus TCP server",
		logger.String("url", serverConfig.URL),
		logger.Duration("timeout", cfg.Timeout))

	// Create Modbus server
	modbusServer, err := modbus.NewServer(serverConfig, handler)
	if err != nil {
		cancel()
		serverLogger.Error("Failed to create Modbus server", logger.Err(err))
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
		s.log.Error("Failed to start Modbus server", logger.Err(err))
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
