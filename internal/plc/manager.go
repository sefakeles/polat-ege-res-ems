package plc

import (
	"fmt"
	"maps"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Manager manages multiple PLC services
type Manager struct {
	log *zap.Logger

	mutex    sync.RWMutex
	services map[int]*Service
}

// NewManager creates a new PLC manager
func NewManager(configs []config.PLCConfig, influxDB *database.InfluxDB, alarmManager *alarm.Manager, logger *zap.Logger) *Manager {
	managerLogger := logger.With(zap.String("component", "plc_manager"))

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

// Start starts all PLC services
func (m *Manager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		if err := service.Start(); err != nil {
			m.log.Error("Failed to start PLC service", zap.Int("id", id), zap.Error(err))
			return fmt.Errorf("failed to start PLC service %d: %w", id, err)
		}
	}

	return nil
}

// Stop stops all PLC services
func (m *Manager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, service := range m.services {
		service.Stop()
	}
}

// GetService returns a specific PLC service
func (m *Manager) GetService(id int) (*Service, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	service, exists := m.services[id]
	if !exists {
		return nil, fmt.Errorf("PLC service %d not found", id)
	}

	return service, nil
}

// GetAllServices returns all PLC services
func (m *Manager) GetAllServices() map[int]*Service {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	services := make(map[int]*Service)
	maps.Copy(services, m.services)

	return services
}

// GetAggregatedData returns aggregated data from all PLC services
func (m *Manager) GetAggregatedData() map[int]database.PLCData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	data := make(map[int]database.PLCData)
	for id, service := range m.services {
		data[id] = service.GetLatestPLCData()
	}

	return data
}

// HasAnyProtectionRelayFaults checks if any PLC has protection relay faults
func (m *Manager) HasAnyProtectionRelayFaults() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, service := range m.services {
		if service.HasProtectionRelayFaults() {
			return true
		}
	}

	return false
}

// GetAllFaultedRelays returns all faulted relays from all PLCs
func (m *Manager) GetAllFaultedRelays() map[int][]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	faults := make(map[int][]string)
	for id, service := range m.services {
		faulted := service.GetFaultedRelays()
		if len(faulted) > 0 {
			faults[id] = faulted
		}
	}

	return faults
}
