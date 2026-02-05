package windfarm

import (
	"time"

	"go.uber.org/zap"
)

// persistenceLoop periodically writes data to InfluxDB
func (s *Service) persistenceLoop() {
	interval := s.config.PersistInterval

	// Calculate first aligned time and create timer
	nextTick := time.Now().Truncate(interval).Add(interval)
	timer := time.NewTimer(time.Until(nextTick))
	defer timer.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			s.persistData()

			// Calculate next aligned time and reset timer
			nextTick = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextTick))
		}
	}
}

// persistData writes all data to InfluxDB
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
