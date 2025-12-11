package pcs

import (
	"fmt"

	"powerkonnekt/ems/pkg/logger"
)

// readStatusData reads status data registers
func (s *Service) readStatusData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, StatusDataStartAddr, StatusDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	statusData := ParseStatusData(data, s.config.ID)

	s.mutex.Lock()
	s.lastStatusData = statusData
	s.mutex.Unlock()

	return nil
}

// readEquipmentData reads equipment data registers
func (s *Service) readEquipmentData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, EquipmentDataStartAddr, EquipmentDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	equipmentData := ParseEquipmentData(data, s.config.ID)

	s.mutex.Lock()
	s.lastEquipmentData = equipmentData
	s.mutex.Unlock()

	return nil
}

// readEnvironmentData reads environment data registers
func (s *Service) readEnvironmentData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, EnvironmentDataStartAddr, EnvironmentDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	environmentData := ParseEnvironmentData(data, s.config.ID)

	s.mutex.Lock()
	s.lastEnvironmentData = environmentData
	s.mutex.Unlock()

	return nil
}

// readDCSourceData reads DC source data registers
func (s *Service) readDCSourceData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, DCSourceDataStartAddr, DCSourceDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	dcSourceData := ParseDCSourceData(data, s.config.ID)

	s.mutex.Lock()
	s.lastDCSourceData = dcSourceData
	s.mutex.Unlock()

	return nil
}

// readGridData reads grid data registers
func (s *Service) readGridData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, GridDataStartAddr, GridDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	gridData := ParseGridData(data, s.config.ID)

	s.mutex.Lock()
	s.lastGridData = gridData
	s.mutex.Unlock()

	return nil
}

// readCounterData reads counter data registers
func (s *Service) readCounterData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, CounterDataStartAddr, CounterDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	counterData := ParseCounterData(data, s.config.ID)

	s.mutex.Lock()
	s.lastCounterData = counterData
	s.mutex.Unlock()

	return nil
}

// readFaults reads fault registers
func (s *Service) readFaults() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, FaultDataStartAddr, FaultDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	s.processFaults(data)

	return nil
}

// readWarnings reads warning registers
func (s *Service) readWarnings() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, WarningDataStartAddr, WarningDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	s.processWarnings(data)

	return nil
}

// readPCSData reads all PCS data registers
func (s *Service) readPCSData() error {
	var lastErr error

	if err := s.readStatusData(); err != nil {
		s.log.Error("Failed to read status data", logger.Err(err))
		lastErr = err
	}

	if err := s.readEquipmentData(); err != nil {
		s.log.Error("Failed to read equipment data", logger.Err(err))
		lastErr = err
	}

	if err := s.readEnvironmentData(); err != nil {
		s.log.Error("Failed to read environment data", logger.Err(err))
		lastErr = err
	}

	if err := s.readDCSourceData(); err != nil {
		s.log.Error("Failed to read DC source data", logger.Err(err))
		lastErr = err
	}

	if err := s.readGridData(); err != nil {
		s.log.Error("Failed to read grid data", logger.Err(err))
		lastErr = err
	}

	if err := s.readCounterData(); err != nil {
		s.log.Error("Failed to read counter data", logger.Err(err))
		lastErr = err
	}

	if err := s.readFaults(); err != nil {
		s.log.Error("Failed to read faults", logger.Err(err))
		lastErr = err
	}

	if err := s.readWarnings(); err != nil {
		s.log.Error("Failed to read warnings", logger.Err(err))
		lastErr = err
	}

	return lastErr
}
