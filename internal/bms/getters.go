package bms

import "powerkonnekt/ems/internal/database"

// GetTotalCellsPerRack returns the total number of cells per rack
func (s *Service) GetTotalCellsPerRack() int {
	return s.config.ModulesPerRack * CellsPerModule
}

// GetTotalTempSensorsPerRack returns the total number of temperature sensors per rack
func (s *Service) GetTotalTempSensorsPerRack() int {
	return s.config.ModulesPerRack * TempSensorsPerModule
}

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.baseClient.IsConnected()
}

// GetBaseDataUpdateChannel returns the channel that signals when new base data is available
func (s *Service) GetBaseDataUpdateChannel() <-chan struct{} {
	return s.baseDataUpdateChan
}

// GetCellDataUpdateChannel returns the channel that signals when new cell data is available
func (s *Service) GetCellDataUpdateChannel() <-chan struct{} {
	return s.cellDataUpdateChan
}

// GetLatestBMSData returns the latest BMS data
func (s *Service) GetLatestBMSData() database.BMSData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastBMSData
}

// GetLatestBMSRackData returns the latest BMS rack data
func (s *Service) GetLatestBMSRackData() []database.BMSRackData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return append([]database.BMSRackData(nil), s.lastBMSRackData...)
}

// GetLatestCellVoltageData returns the latest cell voltage data for a specific rack
func (s *Service) GetLatestCellVoltageData(rackNo uint8) []database.BMSCellVoltageData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var voltages []database.BMSCellVoltageData

	if rackNo > 0 && int(rackNo) <= len(s.lastCellVoltages) {
		voltages = make([]database.BMSCellVoltageData, len(s.lastCellVoltages[rackNo-1]))
		copy(voltages, s.lastCellVoltages[rackNo-1])
	}

	return voltages
}

// GetLatestCellTemperatureData returns the latest cell temperature data for a specific rack
func (s *Service) GetLatestCellTemperatureData(rackNo uint8) []database.BMSCellTemperatureData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var temperatures []database.BMSCellTemperatureData

	if rackNo > 0 && int(rackNo) <= len(s.lastCellTemperatures) {
		temperatures = make([]database.BMSCellTemperatureData, len(s.lastCellTemperatures[rackNo-1]))
		copy(temperatures, s.lastCellTemperatures[rackNo-1])
	}

	return temperatures
}

// GetCommandState returns the current command state
func (s *Service) GetCommandState() database.BMSCommandState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.commandState
}
