package bms

import (
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/utils"
)

// ParseBMSData converts raw MODBUS data to BMSData structure
func ParseBMSData(data []byte, id int) database.BMSData {
	if len(data) < BMSDataLength*2 {
		return database.BMSData{Timestamp: time.Now(), ID: id}
	}

	return database.BMSData{
		Timestamp:               time.Now(),
		ID:                      id,
		Voltage:                 utils.Scale(utils.FromBytes[uint16](data[0:2]), float32(0.1)),     // 32 - Voltage (0.1V)
		Current:                 utils.FromBytes[int16](data[2:4]) - 20000,                         // 33 - Current (A, offset -20000)
		SOC:                     utils.Scale(utils.FromBytes[uint16](data[4:6]), float32(0.1)),     // 34 - SOC (0.1%)
		SOH:                     utils.Scale(utils.FromBytes[uint16](data[6:8]), float32(0.1)),     // 35 - SOH (0.1%)
		MaxCellVoltage:          utils.Scale(utils.FromBytes[uint16](data[8:10]), float32(0.001)),  // 36 - Max cell voltage (0.001V)
		MinCellVoltage:          utils.Scale(utils.FromBytes[uint16](data[10:12]), float32(0.001)), // 37 - Min cell voltage (0.001V)
		AvgCellVoltage:          utils.Scale(utils.FromBytes[uint16](data[12:14]), float32(0.001)), // 38 - Average cell voltage (0.001V)
		MaxCellTemperature:      utils.FromBytes[int16](data[14:16]) - 50,                          // 39 - Max cell temperature (°C, offset -50)
		MinCellTemperature:      utils.FromBytes[int16](data[16:18]) - 50,                          // 40 - Min cell temperature (°C, offset -50)
		AvgCellTemperature:      utils.FromBytes[int16](data[18:20]) - 50,                          // 41 - Average cell temperature (°C, offset -50)
		MaxChargeCurrent:        utils.FromBytes[int16](data[20:22]) - 20000,                       // 42 - Max charge current (A, offset -20000)
		MaxDischargeCurrent:     utils.FromBytes[int16](data[22:24]) - 20000,                       // 43 - Max discharge current (A, offset -20000)
		MaxChargePower:          utils.FromBytes[int16](data[24:26]) - 20000,                       // 44 - Max charge power (kW, offset -20000)
		MaxDischargePower:       utils.FromBytes[int16](data[26:28]) - 20000,                       // 45 - Max discharge power (kW, offset -20000)
		Power:                   utils.FromBytes[int16](data[28:30]) - 20000,                       // 46 - Power (kW, offset -20000)
		ChargeCapacity:          utils.FromBytes[uint16](data[34:36]),                              // 49 - Charge capacity (kWh)
		DischargeCapacity:       utils.FromBytes[uint16](data[36:38]),                              // 50 - Discharge capacity (kWh)
		MaxChargeVoltage:        utils.Scale(utils.FromBytes[uint16](data[38:40]), float32(0.1)),   // 51 - Max charge voltage (0.1V)
		MinDischargeVoltage:     utils.Scale(utils.FromBytes[uint16](data[40:42]), float32(0.1)),   // 52 - Min discharge voltage (0.1V)
		InsulationResistancePos: utils.FromBytes[uint16](data[44:46]),                              // 54 - Insulation resistance positive (kΩ)
		InsulationResistanceNeg: utils.FromBytes[uint16](data[46:48]),                              // 55 - Insulation resistance negative (kΩ)
	}
}

// ParseBMSStatusData converts raw MODBUS data to BMSStatusData structure
func ParseBMSStatusData(data []byte, id int) database.BMSStatusData {
	if len(data) < BMSStatusDataLength*2 {
		return database.BMSStatusData{Timestamp: time.Now(), ID: id}
	}

	return database.BMSStatusData{
		Timestamp:      time.Now(),
		ID:             id,
		Heartbeat:      utils.FromBytes[uint16](data[0:2]),   // 768 - Heartbeat
		HVStatus:       utils.FromBytes[uint16](data[2:4]),   // 769 - HV Status
		SystemStatus:   utils.FromBytes[uint16](data[4:6]),   // 770 - System Status
		ConnectedRacks: utils.FromBytes[uint16](data[8:10]),  // 772 - Connected Racks
		TotalRacks:     utils.FromBytes[uint16](data[10:12]), // 773 - Total Racks
	}
}

// ParseBMSRackData converts raw MODBUS data to BMSRackData structure
func ParseBMSRackData(data []byte, id int, rackNo uint8) database.BMSRackData {
	if len(data) < BMSRackDataLength*2 {
		return database.BMSRackData{
			Timestamp: time.Now(),
			ID:        id,
			Number:    rackNo,
		}
	}

	return database.BMSRackData{
		Timestamp:            time.Now(),
		ID:                   id,
		Number:               rackNo,
		State:                utils.FromBytes[uint16](data[0:2]),                                       // 100 - State
		MaxChargePower:       utils.Scale(utils.FromBytes[uint16](data[2:4]), float32(0.1)),            // 101 - Max charge power (0.1kW)
		MaxDischargePower:    utils.Scale(utils.FromBytes[uint16](data[4:6]), float32(0.1)),            // 102 - Max discharge power (0.1kW)
		MaxChargeVoltage:     utils.Scale(utils.FromBytes[uint16](data[6:8]), float32(0.1)),            // 103 - Max charge voltage (0.1V)
		MinDischargeVoltage:  utils.Scale(utils.FromBytes[uint16](data[8:10]), float32(0.1)),           // 104 - Min discharge voltage (0.1V)
		MaxChargeCurrent:     utils.Scale(utils.FromBytes[uint16](data[10:12]), float32(0.1)),          // 105 - Max charge current (0.1A)
		MaxDischargeCurrent:  utils.Scale(utils.FromBytes[uint16](data[12:14]), float32(0.1)),          // 106 - Max discharge current (0.1A)
		Voltage:              utils.Scale(utils.FromBytes[uint16](data[30:32]), float32(0.1)),          // 115 - Voltage (0.1V)
		Current:              utils.Scale(utils.FromBytes[uint16](data[32:34]), float32(0.1)) - 1600.0, // 116 - Current (0.1A, offset -1600)
		Temperature:          utils.FromBytes[int16](data[34:36]) - 40,                                 // 117 - Temperature (°C, offset -40)
		SOC:                  utils.FromBytes[uint16](data[36:38]),                                     // 118 - SOC (%)
		SOH:                  utils.FromBytes[uint16](data[38:40]),                                     // 119 - SOH (%)
		InsulationResistance: utils.FromBytes[uint16](data[40:42]),                                     // 120 - Insulation resistance (kΩ)
		AvgCellVoltage:       utils.Scale(utils.FromBytes[uint16](data[42:44]), float32(0.001)),        // 121 - Average cell voltage (0.001V)
		AvgCellTemperature:   utils.FromBytes[int16](data[44:46]) - 40,                                 // 122 - Average cell temperature (°C, offset -40)
		MaxCellVoltage:       utils.Scale(utils.FromBytes[uint16](data[46:48]), float32(0.001)),        // 123 - Max cell voltage (0.001V)
		MaxVoltageCellNo:     utils.FromBytes[uint16](data[48:50]),                                     // 124 - Max voltage cell number
		MinCellVoltage:       utils.Scale(utils.FromBytes[uint16](data[50:52]), float32(0.001)),        // 125 - Min cell voltage (0.001V)
		MinVoltageCellNo:     utils.FromBytes[uint16](data[52:54]),                                     // 126 - Min voltage cell number
		MaxCellTemperature:   utils.FromBytes[int16](data[54:56]) - 40,                                 // 127 - Max cell temperature (°C, offset -40)
		MaxTempCellNo:        utils.FromBytes[uint16](data[56:58]),                                     // 128 - Max temperature cell number
		MinCellTemperature:   utils.FromBytes[int16](data[58:60]) - 40,                                 // 129 - Min cell temperature (°C, offset -40)
		MinTempCellNo:        utils.FromBytes[uint16](data[60:62]),                                     // 130 - Min temperature cell number
		TotalChargeEnergy:    utils.Scale(utils.FromBytes[uint32](data[78:82]), float32(0.1)),          // 139-140 - Total charge energy (0.1kWh)
		TotalDischargeEnergy: utils.Scale(utils.FromBytes[uint32](data[82:86]), float32(0.1)),          // 141-142 - Total discharge energy (0.1kWh)
		ChargeCapacity:       utils.Scale(utils.FromBytes[uint32](data[94:98]), float32(0.1)),          // 147-148 - Charge capacity (0.1kWh)
		DischargeCapacity:    utils.Scale(utils.FromBytes[uint32](data[98:102]), float32(0.1)),         // 149-150 - Discharge capacity (0.1kWh)
	}
}

// ParseCellVoltages converts raw MODBUS data to cell voltage array
func ParseCellVoltages(data []byte, id int, startCellNo uint16, rackNo uint8) []database.BMSCellVoltageData {
	if len(data) < 2 {
		return nil
	}

	cellCount := len(data) / 2
	cells := make([]database.BMSCellVoltageData, cellCount)

	// Use the same timestamp for all cells in this batch
	timestamp := time.Now()

	for i := range cellCount {
		voltage := utils.Scale(utils.FromBytes[uint16](data[i*2:(i+1)*2]), float32(0.001))

		// Calculate current cell number (1-based)
		currentCellNo := startCellNo + uint16(i)

		// Calculate module number (1-based): cells 1-48 = module 1, 49-96 = module 2, etc.
		moduleNo := uint8((currentCellNo-1)/CellsPerModule) + 1

		cells[i] = database.BMSCellVoltageData{
			Timestamp: timestamp,
			ID:        id,
			RackNo:    rackNo,
			ModuleNo:  moduleNo,
			CellNo:    currentCellNo,
			Voltage:   voltage,
		}
	}

	return cells
}

// ParseCellTemperatures converts raw MODBUS data to cell temperature array
func ParseCellTemperatures(data []byte, id int, startSensorNo uint16, rackNo uint8) []database.BMSCellTemperatureData {
	if len(data) < 2 {
		return nil
	}

	sensorCount := len(data) / 2
	sensors := make([]database.BMSCellTemperatureData, sensorCount)

	// Use the same timestamp for all sensors in this batch
	timestamp := time.Now()

	for i := range sensorCount {
		// Parse temperature with offset: Unit:°C, offset: -40°C
		temperature := utils.FromBytes[int16](data[i*2:(i+1)*2]) - 40

		// Calculate current sensor number (1-based, 1-60 total)
		currentSensorNo := startSensorNo + uint16(i)

		// Calculate module number: sensors 1-12 = module 1, 13-24 = module 2, etc.
		moduleNo := uint8((currentSensorNo-1)/TempSensorsPerModule) + 1

		sensors[i] = database.BMSCellTemperatureData{
			Timestamp:   timestamp,
			ID:          id,
			RackNo:      rackNo,
			ModuleNo:    moduleNo,
			SensorNo:    currentSensorNo,
			Temperature: temperature,
		}
	}

	return sensors
}
