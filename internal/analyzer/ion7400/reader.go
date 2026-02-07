package ion7400

import "fmt"

// readBaseData reads the base data from the energy analyzer
func (s *Service) readBaseData() error {
	data1, err := s.client.ReadHoldingRegisters(s.ctx, BaseDataStartAddr, BaseDataLength)
	if err != nil {
		return fmt.Errorf("failed to read base registers: %w", err)
	}

	data2, err := s.client.ReadHoldingRegisters(s.ctx, PowerFactorDataStartAddr, PowerFactorDataLength)
	if err != nil {
		return fmt.Errorf("failed to read power factor registers: %w", err)
	}

	data3, err := s.client.ReadHoldingRegisters(s.ctx, EnergyDataStartAddr, EnergyDataLength)
	if err != nil {
		return fmt.Errorf("failed to read energy registers: %w", err)
	}

	baseData := parseBaseData(data1)
	powerFactorData := parsePowerFactorData(data2)
	energyData := parseEnergyData(data3)

	s.mutex.Lock()
	s.lastData = baseData
	s.lastData.PowerFactorL1 = powerFactorData.PowerFactorL1
	s.lastData.PowerFactorL2 = powerFactorData.PowerFactorL2
	s.lastData.PowerFactorL3 = powerFactorData.PowerFactorL3
	s.lastData.PowerFactorAvg = powerFactorData.PowerFactorAvg
	s.lastData.ActiveEnergyExport = energyData.ActiveEnergyExport
	s.lastData.ActiveEnergyImport = energyData.ActiveEnergyImport
	s.lastData.ReactiveEnergyExport = energyData.ReactiveEnergyExport
	s.lastData.ReactiveEnergyImport = energyData.ReactiveEnergyImport
	s.lastData.ApparentEnergyExport = energyData.ApparentEnergyExport
	s.lastData.ApparentEnergyImport = energyData.ApparentEnergyImport
	s.mutex.Unlock()

	return nil
}
