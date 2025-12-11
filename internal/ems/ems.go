package ems

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"powerkonnekt/ems/internal/api"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/container"
	"powerkonnekt/ems/internal/health"
	"powerkonnekt/ems/pkg/logger"
)

// EMS represents the main EMS application
type EMS struct {
	container     *container.Container
	healthService *health.HealthService
	httpServer    *http.Server
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// New creates a new EMS instance
func New(cfg *config.Config) (*EMS, error) {
	// Create dependency injection container
	cont, err := container.NewContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create health service
	healthService := health.NewHealthService()

	// Register health checkers for all BMS instances
	bmsServices := cont.BMSManager.GetAllServices()
	for bmsID, bmsService := range bmsServices {
		healthService.RegisterChecker(health.NewServiceChecker(fmt.Sprintf("bms_%s", bmsID), bmsService))
	}

	// Register health checkers for all PCS instances
	pcsServices := cont.PCSManager.GetAllServices()
	for pcsID, pcsService := range pcsServices {
		healthService.RegisterChecker(health.NewServiceChecker(fmt.Sprintf("pcs_%s", pcsID), pcsService))
	}

	// Register health checkers for databases
	healthService.RegisterChecker(health.NewDatabaseChecker("influxdb", cont.InfluxDB))
	healthService.RegisterChecker(health.NewDatabaseChecker("postgresql", cont.PostgresDB))

	ems := &EMS{
		container:     cont,
		healthService: healthService,
		ctx:           ctx,
		cancel:        cancel,
	}

	// Setup HTTP server
	ems.setupHTTPServer()

	return ems, nil
}

// Start starts the EMS
func (e *EMS) Start() error {
	// Start metrics collection
	if err := e.container.MetricsManager.Start(); err != nil {
		return fmt.Errorf("failed to start metrics manager: %w", err)
	}

	// Start services
	if err := e.container.BMSManager.Start(); err != nil {
		return fmt.Errorf("failed to start BMS services: %w", err)
	}

	if err := e.container.PCSManager.Start(); err != nil {
		return fmt.Errorf("failed to start PCS services: %w", err)
	}

	// Start Modbus server
	if err := e.container.ModbusServer.Start(); err != nil {
		return fmt.Errorf("failed to start Modbus server: %w", err)
	}

	// Start reactive control logic (event-driven)
	e.wg.Go(e.reactiveControlLoop)

	// Start HTTP server
	e.wg.Go(func() {
		logger.Info("Starting HTTP server", logger.Int("port", e.container.Config.EMS.HTTPPort))
		if err := e.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", logger.Err(err))
		}
	})

	return nil
}

// Stop stops the EMS
func (e *EMS) Stop() {
	logger.Info("Stopping EMS services")

	// Stop HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown HTTP server gracefully", logger.Err(err))
	}

	// Stop Modbus server
	e.container.ModbusServer.Stop()

	// Stop services
	e.cancel()
	e.wg.Wait()

	e.container.BMSManager.Stop()
	e.container.PCSManager.Stop()
	e.container.MetricsManager.Stop()

	// Close databases
	if err := e.container.Close(); err != nil {
		logger.Error("Failed to close databases", logger.Err(err))
	}

	logger.Info("EMS services stopped successfully")
}

// reactiveControlLoop runs reactive control logic triggered by data updates
func (e *EMS) reactiveControlLoop() {
	bessUpdateChan := e.container.ControlLogic.GetBESSUpdateChannel()

	// Also run periodic control as a safety fallback
	fallbackTicker := time.NewTicker(5 * time.Second)
	defer fallbackTicker.Stop()

	logger.Info("Starting reactive control loop")

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-bessUpdateChan:
			// BESS data updated, execute control immediately
			// e.container.ControlLogic.ExecuteControl()
		case <-fallbackTicker.C:
			// Safety fallback - ensure control runs at least once per second
			e.container.ControlLogic.ExecuteControl()
		}
	}
}

// setupHTTPServer initializes the HTTP server with routes and handlers
func (e *EMS) setupHTTPServer() {
	handlers := api.NewHandlers(
		e.container.BMSManager,
		e.container.PCSManager,
		e.container.AlarmManager,
		e.container.ControlLogic,
		e.healthService,
	)
	router := api.SetupRoutes(handlers)

	e.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", e.container.Config.EMS.HTTPPort),
		Handler: router,
	}
}
