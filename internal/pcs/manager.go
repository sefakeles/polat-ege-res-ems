package pcs

import (
	"fmt"
	"maps"
	"sync"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
)

// Manager manages multiple PCS services
type Manager struct {
	services map[int]*Service
	mutex    sync.RWMutex
	log      logger.Logger
}

// NewManager creates a new PCS manager
func NewManager(configs []config.PCSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager) *Manager {
	managerLogger := logger.With(logger.String("component", "pcs_manager"))

	manager := &Manager{
		services: make(map[int]*Service),
		log:      managerLogger,
	}

	for _, cfg := range configs {
		service := NewService(cfg, influxDB, alarmManager)
		manager.services[cfg.ID] = service
	}

	return manager
}

// Start starts all PCS services
func (m *Manager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		if err := service.Start(); err != nil {
			m.log.Error("Failed to start PCS service", logger.Int("id", id), logger.Err(err))
			return fmt.Errorf("failed to start PCS service %d: %w", id, err)
		}
		m.log.Info("PCS service started", logger.Int("id", id))
	}

	return nil
}

// Stop stops all PCS services
func (m *Manager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		service.Stop()
		m.log.Info("PCS service stopped", logger.Int("id", id))
	}
}

// GetService returns a specific PCS service
func (m *Manager) GetService(id int) (*Service, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	service, exists := m.services[id]
	if !exists {
		return nil, fmt.Errorf("PCS service %d not found", id)
	}

	return service, nil
}

// GetAllServices returns all PCS services
func (m *Manager) GetAllServices() map[int]*Service {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	services := make(map[int]*Service)
	maps.Copy(services, m.services)

	return services
}

// StartStopCommandAll sends start/stop command to all PCS concurrently
func (m *Manager) StartStopCommandAll(start bool) error {
	m.mutex.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}
	m.mutex.RUnlock()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr error
	errCount := 0

	wg.Add(len(services))

	for _, service := range services {
		go func(svc *Service) {
			defer wg.Done()
			if err := svc.StartStopCommand(start); err != nil {
				mu.Lock()
				lastErr = err
				errCount++
				mu.Unlock()
			}
		}(service)
	}

	wg.Wait()

	if errCount > 0 {
		m.log.Error("Failed to send start/stop command to some PCS units",
			logger.Int("failed_count", errCount),
			logger.Int("total_count", len(services)),
			logger.Err(lastErr))
		return fmt.Errorf("failed to send command to %d/%d PCS units: %w", errCount, len(services), lastErr)
	}

	return nil
}

// SetActivePowerCommandAll sends active power command to all PCS concurrently
func (m *Manager) SetActivePowerCommandAll(power float32) error {
	m.mutex.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}
	m.mutex.RUnlock()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr error
	errCount := 0

	wg.Add(len(services))

	for _, service := range services {
		go func(svc *Service) {
			defer wg.Done()
			if err := svc.SetActivePowerCommand(power); err != nil {
				mu.Lock()
				lastErr = err
				errCount++
				mu.Unlock()
			}
		}(service)
	}

	wg.Wait()

	if errCount > 0 {
		m.log.Error("Failed to send active power command to some PCS units",
			logger.Int("failed_count", errCount),
			logger.Int("total_count", len(services)),
			logger.Float32("power", power),
			logger.Err(lastErr))
		return fmt.Errorf("failed to send command to %d/%d PCS units: %w", errCount, len(services), lastErr)
	}

	return nil
}

// SetReactivePowerCommandAll sends reactive power command to all PCS concurrently
func (m *Manager) SetReactivePowerCommandAll(power float32) error {
	m.mutex.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}
	m.mutex.RUnlock()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr error
	errCount := 0

	wg.Add(len(services))

	for _, service := range services {
		go func(svc *Service) {
			defer wg.Done()
			if err := svc.SetReactivePowerCommand(power); err != nil {
				mu.Lock()
				lastErr = err
				errCount++
				mu.Unlock()
			}
		}(service)
	}

	wg.Wait()

	if errCount > 0 {
		m.log.Error("Failed to send reactive power command to some PCS units",
			logger.Int("failed_count", errCount),
			logger.Int("total_count", len(services)),
			logger.Float32("power", power),
			logger.Err(lastErr))
		return fmt.Errorf("failed to send command to %d/%d PCS units: %w", errCount, len(services), lastErr)
	}

	return nil
}

// ResetSystemAll sends reset command to all PCS concurrently
func (m *Manager) ResetSystemAll() error {
	m.mutex.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, service := range m.services {
		services = append(services, service)
	}
	m.mutex.RUnlock()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr error
	errCount := 0

	wg.Add(len(services))

	for _, service := range services {
		go func(svc *Service) {
			defer wg.Done()
			if err := svc.ResetSystem(); err != nil {
				mu.Lock()
				lastErr = err
				errCount++
				mu.Unlock()
			}
		}(service)
	}

	wg.Wait()

	if errCount > 0 {
		m.log.Error("Failed to send reset command to some PCS units",
			logger.Int("failed_count", errCount),
			logger.Int("total_count", len(services)),
			logger.Err(lastErr))
		return fmt.Errorf("failed to send command to %d/%d PCS units: %w", errCount, len(services), lastErr)
	}

	return nil
}
