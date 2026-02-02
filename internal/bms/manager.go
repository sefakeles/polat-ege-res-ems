package bms

import (
	"fmt"
	"maps"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Manager manages multiple BMS services
type Manager struct {
	services map[int]*Service
	mutex    sync.RWMutex
	log      *zap.Logger
}

// NewManager creates a new BMS manager
func NewManager(configs []config.BMSConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager, logger *zap.Logger) *Manager {
	managerLogger := logger.With(zap.String("component", "bms_manager"))

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

// Start starts all BMS services
func (m *Manager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		if err := service.Start(); err != nil {
			m.log.Error("Failed to start BMS service", zap.Int("id", id), zap.Error(err))
			return fmt.Errorf("failed to start BMS service %d: %w", id, err)
		}
		m.log.Info("BMS service started", zap.Int("id", id))
	}

	return nil
}

// Stop stops all BMS services
func (m *Manager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		service.Stop()
		m.log.Info("BMS service stopped", zap.Int("id", id))
	}
}

// GetService returns a specific BMS service
func (m *Manager) GetService(id int) (*Service, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	service, exists := m.services[id]
	if !exists {
		return nil, fmt.Errorf("BMS service %d not found", id)
	}

	return service, nil
}

// GetAllServices returns all BMS services
func (m *Manager) GetAllServices() map[int]*Service {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	services := make(map[int]*Service)
	maps.Copy(services, m.services)

	return services
}

// GetAggregatedData returns aggregated data from all BMS services
func (m *Manager) GetAggregatedData() map[int]database.BMSData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	data := make(map[int]database.BMSData)
	for id, service := range m.services {
		data[id] = service.GetLatestBMSData()
	}

	return data
}
