package pcs

import "powerkonnekt/ems/internal/database"

// IsConnected returns the connection status
func (s *Service) IsConnected() bool {
	return s.client.IsConnected()
}

// GetBaseDataUpdateChannel returns the channel that signals when new base data is available
func (s *Service) GetBaseDataUpdateChannel() <-chan struct{} {
	return s.dataUpdateChan
}

// GetLatestPCSData returns the latest PCS data
func (s *Service) GetLatestPCSStatusData() database.PCSStatusData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastStatusData
}

// GetLatestPCSEquipmentData returns the latest PCS equipment data
func (s *Service) GetLatestPCSEquipmentData() database.PCSEquipmentData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastEquipmentData
}

// GetLatestPCSEnvironmentData returns the latest PCS environment data
func (s *Service) GetLatestPCSEnvironmentData() database.PCSEnvironmentData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastEnvironmentData
}

// GetLatestPCSDCSourceData returns the latest PCS DC source data
func (s *Service) GetLatestPCSDCSourceData() database.PCSDCSourceData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastDCSourceData
}

// GetLatestPCSGridData returns the latest PCS grid data
func (s *Service) GetLatestPCSGridData() database.PCSGridData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastGridData
}

// GetLatestPCSCounterData returns the latest PCS counter data
func (s *Service) GetLatestPCSCounterData() database.PCSCounterData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastCounterData
}

func (s *Service) GetLatestPCSData() database.PCSData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return database.PCSData{
		StatusData:      s.lastStatusData,
		EquipmentData:   s.lastEquipmentData,
		EnvironmentData: s.lastEnvironmentData,
		DCSourceData:    s.lastDCSourceData,
		GridData:        s.lastGridData,
		CounterData:     s.lastCounterData,
	}
}

// GetCommandState returns the current command state
func (s *Service) GetCommandState() database.PCSCommandState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.commandState
}
