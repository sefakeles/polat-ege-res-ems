package bms

import (
	"fmt"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/modbus"
)

// readBMSStatusData reads BMS status data
func (s *Service) readBMSStatusData() error {
	data, err := s.systemClient.ReadHoldingRegisters(s.ctx, BMSStatusDataStartAddr, BMSStatusDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	bmsStatusData := parseBMSStatusData(data, s.config.ID)

	s.mutex.Lock()
	s.lastBMSStatusData = bmsStatusData
	s.mutex.Unlock()

	return nil
}

// readBMSData reads BMS data
func (s *Service) readBMSData() error {
	data, err := s.systemClient.ReadHoldingRegisters(s.ctx, BMSDataStartAddr, BMSDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	bmsData := parseBMSData(data, s.config.ID)

	s.mutex.Lock()
	s.lastBMSData = bmsData
	s.mutex.Unlock()

	return nil
}

// readBMSRackStatusData reads BMS rack status data
func (s *Service) readBMSRackStatusData(rackNo uint8) error {
	startAddr := GetRackStatusDataStartAddr(rackNo)

	data, err := s.systemClient.ReadHoldingRegisters(s.ctx, startAddr, BMSRackStatusDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	bmsRackStatusData := parseBMSRackStatusData(data, s.config.ID, rackNo)

	s.mutex.Lock()
	s.lastBMSRackStatusData[rackNo-1] = bmsRackStatusData
	s.mutex.Unlock()

	return nil
}

// readBMSRackData reads BMS rack data
func (s *Service) readBMSRackData(rackNo uint8) error {
	startAddr := GetRackDataStartAddr(rackNo)

	data, err := s.systemClient.ReadHoldingRegisters(s.ctx, startAddr, BMSRackDataLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	bmsRackData := parseBMSRackData(data, s.config.ID, rackNo)

	s.mutex.Lock()
	s.lastBMSRackData[rackNo-1] = bmsRackData
	s.mutex.Unlock()

	return nil
}

// readAlarms reads alarms
func (s *Service) readAlarms() error {
	data, err := s.systemClient.ReadHoldingRegisters(s.ctx, BMSAlarmStartAddr, BMSAlarmLength)
	if err != nil {
		return fmt.Errorf("failed to read registers: %w", err)
	}

	s.processAlarms(data)

	// Read alarms for each rack
	for rackNo := uint8(1); rackNo <= uint8(s.config.RackCount); rackNo++ {
		startAddr := GetRackAlarmStartAddr(rackNo)

		rackAlarmData, err := s.systemClient.ReadHoldingRegisters(s.ctx, startAddr, BMSRackAlarmLength)
		if err != nil {
			s.log.Error("Failed to read rack alarms",
				zap.Error(err),
				zap.Uint8("rack_no", rackNo))
			continue
		}

		s.processRackAlarms(rackAlarmData, rackNo)
	}

	return nil
}

// readCellData reads cell voltages and temperatures for a specific rack
func (s *Service) readCellData(rackNo uint8) error {
	// Read cell voltages
	if err := s.readCellVoltages(rackNo); err != nil {
		s.log.Error("Failed to read cell voltages",
			zap.Error(err),
			zap.Uint8("rack_no", rackNo))
	}

	// Read cell temperatures
	if err := s.readCellTemperatures(rackNo); err != nil {
		s.log.Error("Failed to read cell temperatures",
			zap.Error(err),
			zap.Uint8("rack_no", rackNo))
	}

	return nil
}

// readCellVoltages reads all cell voltages for a rack using chunked requests
func (s *Service) readCellVoltages(rackNo uint8) error {
	// Get the starting MODBUS address for this rack's cell voltages
	startAddr := GetCellVoltageStartAddr(rackNo)

	// Calculate total cells based on config
	totalCells := s.GetTotalCellsPerRack()

	// Pre-allocate slice for all cells in this rack
	allCells := make([]database.BMSCellVoltageData, 0, totalCells)

	// Calculate how many chunks we need to read all cells
	chunks := CalculateReadChunks(totalCells, modbus.MaxRegistersPerRead)

	// Read cells in chunks to avoid MODBUS limitations
	for chunk := range chunks {
		select {
		case <-s.ctx.Done():
			return nil
		default:
		}

		// Calculate which cells to read in this chunk
		startCell := uint16(chunk * modbus.MaxRegistersPerRead)
		cellsInChunk := modbus.MaxRegistersPerRead

		// Last chunk might have fewer cells
		if chunk == chunks-1 {
			cellsInChunk = totalCells - (chunk * modbus.MaxRegistersPerRead)
		}

		// Calculate MODBUS address for this chunk
		chunkAddr := startAddr + startCell

		// Use ReadHoldingRegisters for cell voltage data
		data, err := s.cellClient.ReadHoldingRegisters(s.ctx, chunkAddr, uint16(cellsInChunk))
		if err != nil {
			return fmt.Errorf("failed to read cell voltage chunk %d: %w", chunk, err)
		}

		// Parse raw bytes into structured cell data with rack and module info
		cells := parseCellVoltages(data, s.config.ID, startCell+1, rackNo)

		// Add this chunk's cells to our collection
		allCells = append(allCells, cells...)
	}

	s.mutex.Lock()
	s.lastCellVoltages[rackNo-1] = allCells
	s.mutex.Unlock()

	return nil
}

// readCellTemperatures reads all cell temperatures for a rack using chunked requests
func (s *Service) readCellTemperatures(rackNo uint8) error {
	// Get the starting MODBUS address for this rack's cell temperatures
	startAddr := GetCellTempStartAddr(rackNo)

	// Calculate total sensors based on config
	totalSensors := s.GetTotalTempSensorsPerRack()

	// Pre-allocate slice for all temperature sensors in this rack
	allSensors := make([]database.BMSCellTemperatureData, 0, totalSensors)

	// Calculate how many chunks we need to read all sensors
	chunks := CalculateReadChunks(totalSensors, modbus.MaxRegistersPerRead)

	// Read sensors in chunks to avoid MODBUS limitations
	for chunk := range chunks {
		select {
		case <-s.ctx.Done():
			return nil
		default:
		}

		// Calculate which sensors to read in this chunk
		startSensor := uint16(chunk * modbus.MaxRegistersPerRead)
		sensorsInChunk := modbus.MaxRegistersPerRead

		// Last chunk might have fewer sensors
		if chunk == chunks-1 {
			sensorsInChunk = totalSensors - (chunk * modbus.MaxRegistersPerRead)
		}

		// Calculate MODBUS address for this chunk
		chunkAddr := startAddr + startSensor

		// Use ReadHoldingRegisters for cell temperature data
		data, err := s.cellClient.ReadHoldingRegisters(s.ctx, chunkAddr, uint16(sensorsInChunk))
		if err != nil {
			return fmt.Errorf("failed to read cell temperature chunk %d: %w", chunk, err)
		}

		// Parse raw bytes into structured sensor data with rack and module info
		sensors := parseCellTemperatures(data, s.config.ID, startSensor+1, rackNo)

		// Add this chunk's sensors to our collection
		allSensors = append(allSensors, sensors...)
	}

	s.mutex.Lock()
	s.lastCellTemperatures[rackNo-1] = allSensors
	s.mutex.Unlock()

	return nil
}
