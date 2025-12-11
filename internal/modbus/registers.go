package modbus

import (
	"powerkonnekt/ems/internal/database"

	"github.com/simonvetter/modbus"
)

// RegisterMap holds the register mapping information
type RegisterMap struct {
	// This could be extended for more complex register mappings
}

// NewRegisterMap creates a new register map
func NewRegisterMap() *RegisterMap {
	return &RegisterMap{}
}

// set32BitReg sets 32-bit value as two consecutive registers
func set32BitReg(setReg func(uint16, uint16), idx uint16, val uint32) {
	setReg(idx, uint16(val>>16))
	setReg(idx+1, uint16(val))
}

// convertBMSDataToRegisters converts BMS data to registers
func (h *RequestHandler) convertBMSDataToRegisters(
	data database.BMSData,
	startAddr uint16,
	quantity uint16,
) ([]uint16, error) {
	baseAddr := BMSBaseAddr + uint16(data.ID-1)*BMSDataOffset
	offset := startAddr - (baseAddr + BMSDataStartOffset)
	endOffset := offset + quantity

	if offset >= BMSDataLength || endOffset > BMSDataLength {
		return nil, modbus.ErrIllegalDataAddress
	}

	result := make([]uint16, quantity)

	setReg := func(idx uint16, val uint16) {
		if idx >= offset && idx < endOffset {
			result[idx-offset] = val
		}
	}

	setReg(0, uint16(data.Voltage*10))
	setReg(1, uint16(data.Current))
	setReg(2, uint16(data.SOC*10))
	setReg(3, uint16(data.SOH*10))
	setReg(4, uint16(data.MaxCellVoltage*1000))
	setReg(5, uint16(data.MinCellVoltage*1000))
	setReg(6, uint16(data.AvgCellVoltage*1000))
	setReg(7, uint16(data.MaxCellTemperature))
	setReg(8, uint16(data.MinCellTemperature))
	setReg(9, uint16(data.AvgCellTemperature))
	setReg(10, uint16(data.MaxChargeCurrent))
	setReg(11, uint16(data.MaxDischargeCurrent))
	setReg(12, uint16(data.MaxChargePower))
	setReg(13, uint16(data.MaxDischargePower))
	setReg(14, uint16(data.Power))
	setReg(15, uint16(data.ChargeCapacity))
	setReg(16, uint16(data.DischargeCapacity))
	setReg(17, uint16(data.MaxChargeVoltage*10))
	setReg(18, uint16(data.MaxDischargeVoltage*10))
	setReg(19, uint16(data.InsulationResistancePos))
	setReg(20, uint16(data.InsulationResistanceNeg))

	return result, nil
}

// convertPCSDataToRegisters converts PCS data to registers
func (h *RequestHandler) convertPCSDataToRegisters(
	pcsData database.PCSData,
	startAddr uint16,
	quantity uint16,
) ([]uint16, error) {
	baseAddr := PCSBaseAddr + uint16(pcsData.StatusData.ID-1)*PCSDataOffset
	offset := startAddr - (baseAddr + PCSDataStartOffset)
	endOffset := offset + quantity

	if offset >= PCSDataLength || endOffset > PCSDataLength {
		return nil, modbus.ErrIllegalDataAddress
	}

	result := make([]uint16, quantity)

	setReg := func(idx uint16, val uint16) {
		if idx >= offset && idx < endOffset {
			result[idx-offset] = val
		}
	}

	// Status Data (register 0)
	setReg(0, pcsData.StatusData.Status)

	// Equipment Data (registers 1-8)
	setReg(1, pcsData.EquipmentData.LVSwitchStatus)
	setReg(2, pcsData.EquipmentData.MVSwitchStatus)
	setReg(3, pcsData.EquipmentData.MVDisconnectorStatus)
	setReg(4, pcsData.EquipmentData.MVEarthingSwitchStatus)
	setReg(5, pcsData.EquipmentData.DC1SwitchStatus)
	setReg(6, pcsData.EquipmentData.DC2SwitchStatus)
	setReg(7, pcsData.EquipmentData.DC3SwitchStatus)
	setReg(8, pcsData.EquipmentData.DC4SwitchStatus)

	// Environment Data (register 9)
	setReg(9, uint16(pcsData.EnvironmentData.AirInletTemperature))

	// DC Source Data (registers 10-25)
	setReg(10, uint16(pcsData.DCSourceData.DC1Power))
	setReg(11, uint16(pcsData.DCSourceData.DC2Power))
	setReg(12, uint16(pcsData.DCSourceData.DC3Power))
	setReg(13, uint16(pcsData.DCSourceData.DC4Power))
	setReg(14, pcsData.DCSourceData.DC1Current)
	setReg(15, pcsData.DCSourceData.DC2Current)
	setReg(16, pcsData.DCSourceData.DC3Current)
	setReg(17, pcsData.DCSourceData.DC4Current)
	setReg(18, uint16(pcsData.DCSourceData.DC1VoltageExternal*10))
	setReg(19, uint16(pcsData.DCSourceData.DC2VoltageExternal*10))
	setReg(20, uint16(pcsData.DCSourceData.DC3VoltageExternal*10))
	setReg(21, uint16(pcsData.DCSourceData.DC4VoltageExternal*10))

	// Grid Data (registers 22-43)
	// MV Grid Voltages
	setReg(22, uint16(pcsData.GridData.MVGridVoltageAB*10))
	setReg(23, uint16(pcsData.GridData.MVGridVoltageBC*10))
	setReg(24, uint16(pcsData.GridData.MVGridVoltageCA*10))

	// MV Grid Currents
	setReg(25, uint16(pcsData.GridData.MVGridCurrentA*10))
	setReg(26, uint16(pcsData.GridData.MVGridCurrentB*10))
	setReg(27, uint16(pcsData.GridData.MVGridCurrentC*10))

	// MV Grid Power
	setReg(28, uint16(pcsData.GridData.MVGridActivePower))
	setReg(29, uint16(pcsData.GridData.MVGridReactivePower))
	setReg(30, pcsData.GridData.MVGridApparentPower)
	setReg(31, uint16(pcsData.GridData.MVGridCosPhi*1000))

	// LV Grid Voltages
	setReg(32, uint16(pcsData.GridData.LVGridVoltageAB*10))
	setReg(33, uint16(pcsData.GridData.LVGridVoltageBC*10))
	setReg(34, uint16(pcsData.GridData.LVGridVoltageCA*10))

	// LV Grid Currents
	setReg(35, uint16(pcsData.GridData.LVGridCurrentA*10))
	setReg(36, uint16(pcsData.GridData.LVGridCurrentB*10))
	setReg(37, uint16(pcsData.GridData.LVGridCurrentC*10))

	// LV Grid Power
	setReg(38, uint16(pcsData.GridData.LVGridActivePower))
	setReg(39, uint16(pcsData.GridData.LVGridReactivePower))
	setReg(40, pcsData.GridData.LVGridApparentPower)
	setReg(41, uint16(pcsData.GridData.LVGridCosPhi*1000))

	// Grid Frequency (32-bit)
	set32BitReg(setReg, 42, uint32(pcsData.GridData.GridFrequency*10000))

	// Counter Data (registers 44-67)
	// Active Energy
	set32BitReg(setReg, 44, pcsData.CounterData.ActiveEnergyToday)
	set32BitReg(setReg, 46, pcsData.CounterData.ActiveEnergyYesterday)
	set32BitReg(setReg, 48, pcsData.CounterData.ActiveEnergyThisMonth)
	set32BitReg(setReg, 50, pcsData.CounterData.ActiveEnergyLastMonth)
	set32BitReg(setReg, 52, pcsData.CounterData.ActiveEnergyTotal)

	// Consumed Energy
	set32BitReg(setReg, 54, pcsData.CounterData.ConsumedEnergyToday)
	set32BitReg(setReg, 56, pcsData.CounterData.ConsumedEnergyTotal)

	// Reactive Energy
	set32BitReg(setReg, 58, uint32(pcsData.CounterData.ReactiveEnergyToday))
	set32BitReg(setReg, 60, uint32(pcsData.CounterData.ReactiveEnergyYesterday))
	set32BitReg(setReg, 62, uint32(pcsData.CounterData.ReactiveEnergyThisMonth))
	set32BitReg(setReg, 64, uint32(pcsData.CounterData.ReactiveEnergyLastMonth))
	set32BitReg(setReg, 66, uint32(pcsData.CounterData.ReactiveEnergyTotal))

	return result, nil
}
