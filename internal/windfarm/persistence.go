package windfarm

import (
	"time"

	"go.uber.org/zap"
)

// persistenceLoop periodically writes data to InfluxDB
func (s *Service) persistenceLoop() {
	ticker := time.NewTicker(s.config.PersistInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if err := s.persistData(); err != nil {
				s.log.Error("Error persisting wind farm data", zap.Error(err))
			}
		}
	}
}

// persistData writes all wind farm data to InfluxDB
func (s *Service) persistData() error {
	s.mutex.RLock()
	measuringData := s.lastMeasuringData
	statusData := s.lastStatusData
	setpointData := s.lastSetpointData
	weatherData := s.lastWeatherData
	s.mutex.RUnlock()

	// Persist measuring data
	if err := s.influxDB.WriteWindFarmMeasuringData(measuringData); err != nil {
		s.log.Error("Failed to write measuring data", zap.Error(err))
	}

	// Persist status data
	if err := s.influxDB.WriteWindFarmStatusData(statusData); err != nil {
		s.log.Error("Failed to write status data", zap.Error(err))
	}

	// Persist setpoint data
	if err := s.influxDB.WriteWindFarmSetpointData(setpointData); err != nil {
		s.log.Error("Failed to write setpoint data", zap.Error(err))
	}

	// Persist weather data
	if err := s.influxDB.WriteWindFarmWeatherData(weatherData); err != nil {
		s.log.Error("Failed to write weather data", zap.Error(err))
	}

	return nil
}
