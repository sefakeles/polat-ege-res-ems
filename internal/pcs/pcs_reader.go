package pcs

import (
	"fmt"
	"sync"

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

// readPCSData reads all PCS data registers concurrently
func (s *Service) readPCSData() error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr error

	// Read all register blocks concurrently
	readFuncs := []struct {
		name string
		fn   func() error
	}{
		{"status", s.readStatusData},
		{"equipment", s.readEquipmentData},
		{"environment", s.readEnvironmentData},
		{"dc_source", s.readDCSourceData},
		{"grid", s.readGridData},
		{"counter", s.readCounterData},
		{"faults", s.readFaults},
		{"warnings", s.readWarnings},
	}

	wg.Add(len(readFuncs))

	for _, rf := range readFuncs {
		go func(name string, fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				s.log.Error("Failed to read "+name+" data", logger.Err(err))
				mu.Lock()
				lastErr = err
				mu.Unlock()
			}
		}(rf.name, rf.fn)
	}

	wg.Wait()

	return lastErr
}
