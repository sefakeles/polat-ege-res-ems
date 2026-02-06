package pcs

import (
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/utils"
)

// parseStatusData parses status data registers
func parseStatusData(data []byte, id int, timestamp time.Time) database.PCSStatusData {
	if len(data) < StatusDataLength*2 {
		return database.PCSStatusData{
			Timestamp: timestamp,
			ID:        id,
		}
	}

	return database.PCSStatusData{
		Timestamp: timestamp,
		ID:        id,
		Status:    utils.FromBytes[uint16](data[0:2]), // 1003 - Status
	}
}

// parseEquipmentData parses equipment data registers
func parseEquipmentData(data []byte, id int, timestamp time.Time) database.PCSEquipmentData {
	if len(data) < EquipmentDataLength*2 {
		return database.PCSEquipmentData{
			Timestamp: timestamp,
			ID:        id,
		}
	}

	return database.PCSEquipmentData{
		Timestamp:              timestamp,
		ID:                     id,
		LVSwitchStatus:         utils.FromBytes[uint16](data[0:2]),   // 631 - LV switch status
		MVSwitchStatus:         utils.FromBytes[uint16](data[4:6]),   // 633 - MV switch status
		MVDisconnectorStatus:   utils.FromBytes[uint16](data[8:10]),  // 635 - MV disconnector status
		MVEarthingSwitchStatus: utils.FromBytes[uint16](data[10:12]), // 636 - MV earthing switch status
		DC1SwitchStatus:        utils.FromBytes[uint16](data[12:14]), // 637 - DC switch status of module 1
		DC2SwitchStatus:        utils.FromBytes[uint16](data[14:16]), // 638 - DC switch status of module 2
		DC3SwitchStatus:        utils.FromBytes[uint16](data[16:18]), // 639 - DC switch status of module 3
		DC4SwitchStatus:        utils.FromBytes[uint16](data[18:20]), // 640 - DC switch status of module 4
	}
}

// parseEnvironmentData parses environment data registers
func parseEnvironmentData(data []byte, id int, timestamp time.Time) database.PCSEnvironmentData {
	if len(data) < EnvironmentDataLength*2 {
		return database.PCSEnvironmentData{
			Timestamp: timestamp,
			ID:        id,
		}
	}

	return database.PCSEnvironmentData{
		Timestamp:           timestamp,
		ID:                  id,
		AirInletTemperature: utils.FromBytes[int16](data[0:2]), // 1104 - Air inlet temperature (°C)
	}
}

// parseDCSourceData parses DC source data registers
func parseDCSourceData(data []byte, id int, timestamp time.Time) database.PCSDCSourceData {
	if len(data) < DCSourceDataLength*2 {
		return database.PCSDCSourceData{
			Timestamp: timestamp,
			ID:        id,
		}
	}

	return database.PCSDCSourceData{
		Timestamp:  timestamp,
		ID:         id,
		DC1Power:   utils.FromBytes[int16](data[0:2]),    // 1372 - DC power of busbar 1 (kW)
		DC2Power:   utils.FromBytes[int16](data[2:4]),    // 1373 - DC power of busbar 2 (kW)
		DC3Power:   utils.FromBytes[int16](data[4:6]),    // 1374 - DC power of busbar 3 (kW)
		DC4Power:   utils.FromBytes[int16](data[6:8]),    // 1375 - DC power of busbar 4 (kW)
		DC1Current: utils.FromBytes[uint16](data[8:10]),  // 1376 - DC current of busbar 1 (A)
		DC2Current: utils.FromBytes[uint16](data[10:12]), // 1377 - DC current of busbar 2 (A)
		DC3Current: utils.FromBytes[uint16](data[12:14]), // 1378 - DC current of busbar 3 (A)
		DC4Current: utils.FromBytes[uint16](data[14:16]), // 1379 - DC current of busbar 4 (A)
	}
}

// parseGridData parses grid data registers
func parseGridData(data []byte, id int, timestamp time.Time) database.PCSGridData {
	if len(data) < GridDataLength*2 {
		return database.PCSGridData{
			Timestamp: timestamp,
			ID:        id,
		}
	}

	return database.PCSGridData{
		Timestamp:           timestamp,
		ID:                  id,
		MVGridVoltageAB:     utils.Scale(utils.FromBytes[uint16](data[0:2]), float32(0.1)),      // 4300 - MV grid voltage AB (0.1V)
		MVGridVoltageBC:     utils.Scale(utils.FromBytes[uint16](data[2:4]), float32(0.1)),      // 4301 - MV grid voltage BC (0.1V)
		MVGridVoltageCA:     utils.Scale(utils.FromBytes[uint16](data[4:6]), float32(0.1)),      // 4302 - MV grid voltage CA (0.1V)
		MVGridCurrentA:      utils.Scale(utils.FromBytes[uint16](data[6:8]), float32(0.1)),      // 4303 - MV grid current A (0.1A)
		MVGridCurrentB:      utils.Scale(utils.FromBytes[uint16](data[8:10]), float32(0.1)),     // 4304 - MV grid current B (0.1A)
		MVGridCurrentC:      utils.Scale(utils.FromBytes[uint16](data[10:12]), float32(0.1)),    // 4305 - MV grid current C (0.1A)
		MVGridActivePower:   utils.FromBytes[int16](data[12:14]),                                // 4306 - Active power (kW)
		MVGridReactivePower: utils.FromBytes[int16](data[14:16]),                                // 4307 - Reactive power (kVAr)
		MVGridApparentPower: utils.FromBytes[uint16](data[16:18]),                               // 4308 - Apparent power (kVA)
		MVGridCosPhi:        utils.Scale(utils.FromBytes[uint16](data[18:20]), float32(0.001)),  // 4309 - Grid cos phi (0.001)
		LVGridVoltageAB:     utils.Scale(utils.FromBytes[uint16](data[40:42]), float32(0.1)),    // 4320 - LV grid voltage AB (0.1V)
		LVGridVoltageBC:     utils.Scale(utils.FromBytes[uint16](data[42:44]), float32(0.1)),    // 4321 - LV grid voltage BC (0.1V)
		LVGridVoltageCA:     utils.Scale(utils.FromBytes[uint16](data[44:46]), float32(0.1)),    // 4322 - LV grid voltage CA (0.1V)
		LVGridCurrentA:      utils.Scale(utils.FromBytes[uint16](data[46:48]), float32(0.1)),    // 4323 - LV grid current A (0.1A)
		LVGridCurrentB:      utils.Scale(utils.FromBytes[uint16](data[48:50]), float32(0.1)),    // 4324 - LV grid current B (0.1A)
		LVGridCurrentC:      utils.Scale(utils.FromBytes[uint16](data[50:52]), float32(0.1)),    // 4325 - LV grid current C (0.1A)
		LVGridActivePower:   utils.FromBytes[int16](data[52:54]),                                // 4326 - Active power (kW)
		LVGridReactivePower: utils.FromBytes[int16](data[54:56]),                                // 4327 - Reactive power (kVAr)
		LVGridApparentPower: utils.FromBytes[uint16](data[56:58]),                               // 4328 - Apparent power (kVA)
		LVGridCosPhi:        utils.Scale(utils.FromBytes[uint16](data[58:60]), float32(0.001)),  // 4329 - Grid cos phi (0.001)
		GridFrequency:       utils.Scale(utils.FromBytes[uint32](data[62:66]), float32(0.0001)), // 4331-4332 - Grid frequency (0.0001Hz)
	}
}

// parseCounterData parses counter data registers
func parseCounterData(data []byte, id int, timestamp time.Time) database.PCSCounterData {
	if len(data) < CounterDataLength*2 {
		return database.PCSCounterData{
			Timestamp: timestamp,
			ID:        id,
		}
	}

	return database.PCSCounterData{
		Timestamp:               timestamp,
		ID:                      id,
		ActiveEnergyToday:       utils.FromBytes[uint32](data[0:4]),   // 539-540 - Today's active energy (kWh)
		ActiveEnergyYesterday:   utils.FromBytes[uint32](data[4:8]),   // 541-542 - Yesterday's active energy (kWh)
		ActiveEnergyThisMonth:   utils.FromBytes[uint32](data[8:12]),  // 543-544 - This month's active energy (kWh)
		ActiveEnergyLastMonth:   utils.FromBytes[uint32](data[12:16]), // 545-546 - Last month's active energy (kWh)
		ActiveEnergyTotal:       utils.FromBytes[uint32](data[16:20]), // 547-548 - Total active energy (kWh)
		ConsumedEnergyToday:     utils.FromBytes[uint32](data[20:24]), // 549-550 - Today's consumed energy (kWh)
		ConsumedEnergyTotal:     utils.FromBytes[uint32](data[24:28]), // 551-552 - Total consumed energy (kWh)
		ReactiveEnergyToday:     utils.FromBytes[int32](data[30:34]),  // 554-555 - Today's reactive energy (kVArh)
		ReactiveEnergyYesterday: utils.FromBytes[int32](data[34:38]),  // 556-557 - Yesterday's reactive energy (kVArh)
		ReactiveEnergyThisMonth: utils.FromBytes[int32](data[38:42]),  // 558-559 - This month's reactive energy (kVArh)
		ReactiveEnergyLastMonth: utils.FromBytes[int32](data[42:46]),  // 560-561 - Last month's reactive energy (kVArh)
		ReactiveEnergyTotal:     utils.FromBytes[int32](data[46:50]),  // 562-563 - Total reactive energy (kVArh)
	}
}
