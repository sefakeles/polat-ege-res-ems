package pcs

import (
	"fmt"
	"maps"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Manager manages multiple PCS services
type Manager struct {
	log *zap.Logger

	mutex    sync.RWMutex
	services map[int]*Service
}

// NewManager creates a new PCS manager
func NewManager(configs []config.PCSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager, logger *zap.Logger) *Manager {
	managerLogger := logger.With(zap.String("component", "pcs_manager"))

	manager := &Manager{
		services: make(map[int]*Service),
		log:      managerLogger,
	}

	for _, cfg := range configs {
		service := NewService(cfg, influxDB, alarmManager, logger)
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
			m.log.Error("Failed to start PCS service", zap.Int("id", id), zap.Error(err))
			return fmt.Errorf("failed to start PCS service %d: %w", id, err)
		}
	}

	return nil
}

// Stop stops all PCS services
func (m *Manager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		service.Stop()
		m.log.Info("PCS service stopped", zap.Int("id", id))
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
			zap.Int("failed_count", errCount),
			zap.Int("total_count", len(services)),
			zap.Error(lastErr))
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
			zap.Int("failed_count", errCount),
			zap.Int("total_count", len(services)),
			zap.Float32("power", power),
			zap.Error(lastErr))
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
			zap.Int("failed_count", errCount),
			zap.Int("total_count", len(services)),
			zap.Float32("power", power),
			zap.Error(lastErr))
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
			zap.Int("failed_count", errCount),
			zap.Int("total_count", len(services)),
			zap.Error(lastErr))
		return fmt.Errorf("failed to send command to %d/%d PCS units: %w", errCount, len(services), lastErr)
	}

	return nil
}
