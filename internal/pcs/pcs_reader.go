package pcs

import (
	"fmt"
	"sync"
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/logger"
)

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

	// Create a single timestamp for all data
	timestamp := time.Now()

	// Temporary storage for concurrent reads
	var (
		statusData      database.PCSStatusData
		equipmentData   database.PCSEquipmentData
		environmentData database.PCSEnvironmentData
		dcSourceData    database.PCSDCSourceData
		gridData        database.PCSGridData
		counterData     database.PCSCounterData
	)

	// Read all register blocks concurrently
	readFuncs := []struct {
		name string
		fn   func() error
	}{
		{"status", func() error {
			data, err := s.client.ReadHoldingRegisters(s.ctx, StatusDataStartAddr, StatusDataLength)
			if err != nil {
				return fmt.Errorf("failed to read registers: %w", err)
			}
			statusData = ParseStatusData(data, s.config.ID, timestamp)
			return nil
		}},
		{"equipment", func() error {
			data, err := s.client.ReadHoldingRegisters(s.ctx, EquipmentDataStartAddr, EquipmentDataLength)
			if err != nil {
				return fmt.Errorf("failed to read registers: %w", err)
			}
			equipmentData = ParseEquipmentData(data, s.config.ID, timestamp)
			return nil
		}},
		{"environment", func() error {
			data, err := s.client.ReadHoldingRegisters(s.ctx, EnvironmentDataStartAddr, EnvironmentDataLength)
			if err != nil {
				return fmt.Errorf("failed to read registers: %w", err)
			}
			environmentData = ParseEnvironmentData(data, s.config.ID, timestamp)
			return nil
		}},
		{"dc_source", func() error {
			data, err := s.client.ReadHoldingRegisters(s.ctx, DCSourceDataStartAddr, DCSourceDataLength)
			if err != nil {
				return fmt.Errorf("failed to read registers: %w", err)
			}
			dcSourceData = ParseDCSourceData(data, s.config.ID, timestamp)
			return nil
		}},
		{"grid", func() error {
			data, err := s.client.ReadHoldingRegisters(s.ctx, GridDataStartAddr, GridDataLength)
			if err != nil {
				return fmt.Errorf("failed to read registers: %w", err)
			}
			gridData = ParseGridData(data, s.config.ID, timestamp)
			return nil
		}},
		{"counter", func() error {
			data, err := s.client.ReadHoldingRegisters(s.ctx, CounterDataStartAddr, CounterDataLength)
			if err != nil {
				return fmt.Errorf("failed to read registers: %w", err)
			}
			counterData = ParseCounterData(data, s.config.ID, timestamp)
			return nil
		}},
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

	// Update all data atomically after all reads complete
	s.mutex.Lock()
	s.lastStatusData = statusData
	s.lastEquipmentData = equipmentData
	s.lastEnvironmentData = environmentData
	s.lastDCSourceData = dcSourceData
	s.lastGridData = gridData
	s.lastCounterData = counterData
	s.mutex.Unlock()

	return lastErr
}
