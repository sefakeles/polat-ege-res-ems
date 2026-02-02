package bms

import "fmt"

// readBMSData reads BMS data
func (s *Service) readBMSData() error {
	data, err := s.baseClient.ReadHoldingRegisters(s.ctx, BMSDataStartAddr, BMSDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	bmsData := ParseBMSData(data, s.config.ID)

	s.mutex.Lock()
	s.lastBMSData = bmsData
	s.mutex.Unlock()

	return nil
}

// readBMSStatusData reads BMS status data
func (s *Service) readBMSStatusData() error {
	data, err := s.baseClient.ReadHoldingRegisters(s.ctx, BMSStatusDataStartAddr, BMSStatusDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	bmsStatusData := ParseBMSStatusData(data, s.config.ID)

	s.mutex.Lock()
	s.lastBMSStatusData = bmsStatusData
	s.mutex.Unlock()

	return nil
}

// readBMSRackData reads BMS rack data
func (s *Service) readBMSRackData(rackNo uint8) error {
	startAddr := GetRackDataStartAddr(rackNo)

	data, err := s.baseClient.ReadHoldingRegisters(s.ctx, startAddr, BMSRackDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	bmsRackData := ParseBMSRackData(data, s.config.ID, rackNo)

	s.mutex.Lock()
	s.lastBMSRackData[rackNo-1] = bmsRackData
	s.mutex.Unlock()

	return nil
}

// readAlarms reads alarms
func (s *Service) readAlarms() error {
	data, err := s.baseClient.ReadHoldingRegisters(s.ctx, BMSAlarmStartAddr, BMSAlarmLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	s.processAlarms(data)

	// !Read alarms for each rack
	/*for rackNo := uint8(1); rackNo <= uint8(s.config.RackCount); rackNo++ {
		startAddr := GetRackAlarmStartAddr(rackNo)

		rackAlarmData, err := s.baseClient.ReadHoldingRegisters(s.ctx, startAddr, BMSRackAlarmLength)
		if err != nil {
			s.log.Info("Failed to read alarms",
				logger.Err(err),
				logger.Uint8("rack_no", rackNo))
			continue
		}

		s.processRackAlarms(rackAlarmData, rackNo)
	}*/

	return nil
}
