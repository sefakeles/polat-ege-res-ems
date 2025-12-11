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
