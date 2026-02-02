package windfarm

import (
	"fmt"
	"maps"
	"sync"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// Manager manages multiple Wind Farm services
type Manager struct {
	services map[int]*Service
	mutex    sync.RWMutex
	log      *zap.Logger
}

// NewManager creates a new Wind Farm manager
func NewManager(configs []config.WindFarmConfig, influxDB *database.InfluxDB, logger *zap.Logger) *Manager {
	managerLogger := logger.With(zap.String("component", "windfarm_manager"))

	manager := &Manager{
		services: make(map[int]*Service),
		log:      managerLogger,
	}

	for _, cfg := range configs {
		service := NewService(cfg, influxDB, logger)
		manager.services[cfg.ID] = service
	}

	return manager
}

// Start starts all Wind Farm services
func (m *Manager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		if err := service.Start(); err != nil {
			m.log.Error("Failed to start Wind Farm service", zap.Int("id", id), zap.Error(err))
			return fmt.Errorf("failed to start Wind Farm service %d: %w", id, err)
		}
	}

	return nil
}

// Stop stops all Wind Farm services
func (m *Manager) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for id, service := range m.services {
		service.Stop()
		m.log.Info("Wind Farm service stopped", zap.Int("id", id))
	}
}

// GetService returns a specific Wind Farm service
func (m *Manager) GetService(id int) (*Service, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	service, exists := m.services[id]
	if !exists {
		return nil, fmt.Errorf("Wind Farm service %d not found", id)
	}

	return service, nil
}

// GetAllServices returns all Wind Farm services
func (m *Manager) GetAllServices() map[int]*Service {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	services := make(map[int]*Service)
	maps.Copy(services, m.services)

	return services
}

// GetAggregatedData returns aggregated data from all Wind Farm services
func (m *Manager) GetAggregatedData() map[int]database.WindFarmData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	data := make(map[int]database.WindFarmData)
	for id, service := range m.services {
		data[id] = service.GetLatestData()
	}

	return data
}

// GetAggregatedMeasuringData returns aggregated measuring data from all services
func (m *Manager) GetAggregatedMeasuringData() map[int]database.WindFarmMeasuringData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	data := make(map[int]database.WindFarmMeasuringData)
	for id, service := range m.services {
		data[id] = service.GetLatestMeasuringData()
	}

	return data
}

// GetTotalActivePower returns the total active power from all wind farms
func (m *Manager) GetTotalActivePower() float32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var total float32
	for _, service := range m.services {
		data := service.GetLatestMeasuringData()
		total += data.ActivePowerNCP
	}

	return total
}

// GetTotalReactivePower returns the total reactive power from all wind farms
func (m *Manager) GetTotalReactivePower() float32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var total float32
	for _, service := range m.services {
		data := service.GetLatestMeasuringData()
		total += data.ReactivePowerNCP
	}

	return total
}

// GetTotalPossiblePower returns the total possible WEC power from all wind farms
func (m *Manager) GetTotalPossiblePower() float32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var total float32
	for _, service := range m.services {
		data := service.GetLatestMeasuringData()
		total += data.PossibleWECPower
	}

	return total
}

// GetAverageWindSpeed returns the average wind speed across all wind farms
func (m *Manager) GetAverageWindSpeed() float32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.services) == 0 {
		return 0
	}

	var total float32
	for _, service := range m.services {
		data := service.GetLatestMeasuringData()
		total += data.WindSpeed
	}

	return total / float32(len(m.services))
}

// GetServiceCount returns the number of configured wind farm services
func (m *Manager) GetServiceCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.services)
}

// AreAllFCUsOnline checks if all FCUs are online
func (m *Manager) AreAllFCUsOnline() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, service := range m.services {
		if !service.IsFCUOnline() {
			return false
		}
	}

	return len(m.services) > 0
}
