package fcr

import (
	"context"
	"fmt"
	"sync"
	"time"

	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/pkg/logger"
)

// Service represents the FCR-N service that integrates with EMS
type Service struct {
	controller *FCRNController
	pcsManager *pcs.Manager
	bmsManager *bms.Manager
	config     config.FCRNConfig
	log        logger.Logger

	// Frequency measurement
	lastFrequency   float64
	frequencySource FrequencySource

	// State
	mutex     sync.RWMutex
	isRunning bool
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// FrequencySource defines the interface for frequency measurement
type FrequencySource interface {
	GetFrequency() (float64, error)
	Subscribe(callback func(float64)) error
}

// NewService creates a new FCR-N service
func NewService(cfg config.FCRNConfig, pcsManager *pcs.Manager, bmsManager *bms.Manager, freqSource FrequencySource) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid FCR-N configuration: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	serviceLog := logger.With(
		logger.String("service", "fcrn"),
		logger.Float64("capacity", cfg.Capacity),
	)

	service := &Service{
		pcsManager:      pcsManager,
		bmsManager:      bmsManager,
		frequencySource: freqSource,
		config:          cfg,
		log:             serviceLog,
		ctx:             ctx,
		cancel:          cancel,
	}

	// Create controller with power command callback
	controller := NewFCRNController(cfg, service.sendPowerCommand)
	service.controller = controller

	return service, nil
}

// Start starts the FCR-N service
func (s *Service) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return fmt.Errorf("FCR-N service already running")
	}

	if !s.config.Enabled {
		s.log.Info("FCR-N service is disabled in configuration")
		return nil
	}

	s.log.Info("Starting FCR-N service")

	// Start controller
	if err := s.controller.Start(); err != nil {
		return fmt.Errorf("failed to start FCR-N controller: %w", err)
	}

	// Start frequency monitoring
	s.wg.Go(s.frequencyMonitorLoop)

	// Start SOC monitoring
	s.wg.Go(s.socMonitorLoop)

	// Start telemetry (if enabled)
	if s.config.EnableTelemetry {
		s.wg.Go(s.telemetryLoop)
	}

	s.isRunning = true
	s.log.Info("FCR-N service started successfully")

	return nil
}

// Stop stops the FCR-N service
func (s *Service) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return fmt.Errorf("FCR-N service not running")
	}

	s.log.Info("Stopping FCR-N service")

	// Stop controller
	if err := s.controller.Stop(); err != nil {
		s.log.Warn("Error stopping FCR-N controller", logger.Err(err))
	}

	// Cancel context and wait for goroutines
	s.cancel()
	s.wg.Wait()

	s.isRunning = false
	s.log.Info("FCR-N service stopped")

	return nil
}

// ActivateFCRN activates FCR-N provision
func (s *Service) ActivateFCRN() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.isRunning {
		return fmt.Errorf("FCR-N service not running")
	}

	return s.controller.Activate()
}

// DeactivateFCRN deactivates FCR-N provision
func (s *Service) DeactivateFCRN() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.isRunning {
		return fmt.Errorf("FCR-N service not running")
	}

	return s.controller.Deactivate()
}

// GetState returns the current FCR-N state
func (s *Service) GetState() FCRNState {
	return s.controller.GetState()
}

// SetCapacity updates the FCR-N capacity
func (s *Service) SetCapacity(capacity float64) error {
	return s.controller.SetCapacity(capacity)
}

// SetDroop updates the droop setting
func (s *Service) SetDroop(droop float64) error {
	return s.controller.SetDroop(droop)
}

// GetMaintainedCapacity returns the maintained (available) capacity
func (s *Service) GetMaintainedCapacity() float64 {
	return s.controller.GetMaintainedCapacity()
}

// SetTestFrequency sets the frequency for test frequency source
func (s *Service) SetTestFrequency(frequency float64) error {
	// Check if using test frequency source
	testSource, ok := s.frequencySource.(*TestFrequencySource)
	if !ok {
		return fmt.Errorf("not using test frequency source (current: %T)", s.frequencySource)
	}

	// Set the frequency - this will trigger the subscriber callback
	testSource.SetFrequency(frequency)

	// Immediately update the controller with new frequency
	s.updateFrequency()

	return nil
}

// frequencyMonitorLoop monitors grid frequency
func (s *Service) frequencyMonitorLoop() {
	ticker := time.NewTicker(s.config.FrequencyUpdateRate)
	defer ticker.Stop()

	s.log.Info("Frequency monitoring loop started")

	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("Frequency monitoring loop stopped")
			return

		case <-ticker.C:
			s.updateFrequency()
		}
	}
}

// updateFrequency reads and updates frequency from source
func (s *Service) updateFrequency() {
	// If using PCS frequency source, trigger update
	if pcsSource, ok := s.frequencySource.(*PCSFrequencySource); ok {
		if err := pcsSource.UpdateFromPCS(); err != nil {
			s.log.Error("Failed to update frequency from PCS", logger.Err(err))
			return
		}
	}
	// Test frequency source doesn't need manual update - it's set via API

	frequency, err := s.frequencySource.GetFrequency()
	if err != nil {
		s.log.Error("Failed to get frequency", logger.Err(err))
		return
	}

	s.mutex.Lock()
	s.lastFrequency = frequency
	s.mutex.Unlock()

	// Update controller with new frequency
	s.controller.UpdateFrequency(frequency)
}

// socMonitorLoop monitors battery SOC
func (s *Service) socMonitorLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	s.log.Info("SOC monitoring loop started")

	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("SOC monitoring loop stopped")
			return

		case <-ticker.C:
			s.updateSOC()
		}
	}
}

// updateSOC reads and updates SOC from BESS
func (s *Service) updateSOC() error {
	if s.pcsManager == nil {
		s.log.Debug("BESS service not available for SOC monitoring")
		return nil
	}

	// Get data from BESS service
	bmsService, err := s.bmsManager.GetService(int(s.config.PCSNumber))
	if err != nil {
		return fmt.Errorf("failed to get BMS service %d: %w", s.config.PCSNumber, err)
	}
	data := bmsService.GetLatestBMSData()

	soc := data.SOC
	if soc >= 60.0 || soc <= 40.0 {
		soc = 50.0
	}

	// Update controller with SOC
	// ! Converting to float64 for controller
	s.controller.UpdateSOC(float64(soc))
	return nil
}

// sendPowerCommand sends power command to BESS
func (s *Service) sendPowerCommand(power float64) error {
	if s.pcsManager == nil {
		s.log.Warn("BESS service not available, cannot send power command",
			logger.Float64("power", power))
		return fmt.Errorf("BESS service not available")
	}

	powerPercent := (power / s.config.Capacity) * 100

	// Send power command to all PCS concurrently
	if err := s.pcsManager.SetActivePowerCommandAll(float32(powerPercent)); err != nil {
		return fmt.Errorf("failed to send power command: %w", err)
	}

	s.log.Debug("Power command sent",
		logger.Float64("power_mw", power))

	return nil
}

// telemetryLoop sends telemetry data to TSO
func (s *Service) telemetryLoop() {
	ticker := time.NewTicker(s.config.TelemetryInterval)
	defer ticker.Stop()

	s.log.Info("Telemetry loop started")

	for {
		select {
		case <-s.ctx.Done():
			s.log.Info("Telemetry loop stopped")
			return

		case <-ticker.C:
			s.sendTelemetry()
		}
	}
}

// sendTelemetry sends telemetry data
func (s *Service) sendTelemetry() {
	state := s.controller.GetState()

	// Log telemetry data
	s.log.Debug("FCR-N telemetry",
		logger.Float64("frequency", state.FrequencyMeasured),
		logger.Float64("activated_power", state.ActivatedPower),
		logger.Float64("total_power", state.TotalPower),
		logger.Float64("maintained_capacity", s.GetMaintainedCapacity()),
		logger.Float64("soc", state.SOC),
		logger.Float64("endurance_up", state.EnduranceUpward),
		logger.Float64("endurance_down", state.EnduranceDownward))

	// TODO: Send to external telemetry system (MQTT, API, etc.)
}

// HealthCheck checks if the service is healthy
func (s *Service) HealthCheck() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.isRunning {
		return fmt.Errorf("FCR-N service not running")
	}

	// Check if frequency data is recent
	state := s.controller.GetState()
	if time.Since(state.LastUpdate) > 5*time.Second {
		return fmt.Errorf("no recent FCR-N updates")
	}

	return nil
}
