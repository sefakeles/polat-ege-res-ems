package windfarm

import (
	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/internal/database"
)

// GetLatestData returns the latest aggregated wind farm data
func (s *Service) GetLatestData() database.WindFarmData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return database.WindFarmData{
		MeasuringData: s.lastMeasuringData,
		StatusData:    s.lastStatusData,
		SetpointData:  s.lastSetpointData,
		WeatherData:   s.lastWeatherData,
	}
}

// GetLatestMeasuringData returns the latest measuring data
func (s *Service) GetLatestMeasuringData() database.WindFarmMeasuringData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastMeasuringData
}

// GetLatestStatusData returns the latest status data
func (s *Service) GetLatestStatusData() database.WindFarmStatusData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastStatusData
}

// GetLatestSetpointData returns the latest setpoint data
func (s *Service) GetLatestSetpointData() database.WindFarmSetpointData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastSetpointData
}

// GetLatestWeatherData returns the latest weather data
func (s *Service) GetLatestWeatherData() database.WindFarmWeatherData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastWeatherData
}

// GetCommandState returns the current command state
func (s *Service) GetCommandState() database.WindFarmCommandState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.commandState
}

// GetDataUpdateChan returns the data update channel
func (s *Service) GetDataUpdateChan() <-chan struct{} {
	return s.dataUpdateChan
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.client.IsConnected()
}

// IsFCUOnline returns whether the FCU is online
func (s *Service) IsFCUOnline() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastStatusData.FCUOnline
}

// GetConfig returns the service configuration
func (s *Service) GetConfig() config.WindFarmConfig {
	return s.config
}
