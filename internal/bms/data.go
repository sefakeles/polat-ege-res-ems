package bms

import (
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/utils"
)

// ParseBMSData converts raw MODBUS data to BMSData structure
func ParseBMSData(data []byte, id int) database.BMSData {
	if len(data) < BMSDataLength*2 {
		return database.BMSData{Timestamp: time.Now()}
	}

	return database.BMSData{
		Timestamp:              time.Now(),
		ID:                     id,
		CircuitBreakerStatus:   utils.FromBytes[uint16](data[0:2]),                                     // 1 - Circuit breaker status
		Voltage:                utils.Scale(utils.FromBytes[uint16](data[2:4]), float32(0.1)),          // 2 - Voltage (0.1V)
		Current:                utils.Scale(utils.FromBytes[uint16](data[4:6]), float32(0.1)) - 1600.0, // 3 - Current (0.1A, offset -1600)
		SOC:                    utils.FromBytes[uint16](data[6:8]),                                     // 4 - SOC (%)
		SOH:                    utils.FromBytes[uint16](data[8:10]),                                    // 5 - SOH (%)
		MaxCellVoltage:         utils.Scale(utils.FromBytes[uint16](data[10:12]), float32(0.001)),      // 6 - Max cell voltage (0.001V)
		MaxVoltageRackNo:       utils.FromBytes[uint16](data[12:14]),                                   // 7 - Max voltage rack number
		MaxVoltageCellNo:       utils.FromBytes[uint16](data[14:16]),                                   // 8 - Max voltage cell number
		MinCellVoltage:         utils.Scale(utils.FromBytes[uint16](data[16:18]), float32(0.001)),      // 9 - Min cell voltage (0.001V)
		MinVoltageRackNo:       utils.FromBytes[uint16](data[18:20]),                                   // 10 - Min voltage rack number
		MinVoltageCellNo:       utils.FromBytes[uint16](data[20:22]),                                   // 11 - Min voltage cell number
		MaxCellTemperature:     utils.FromBytes[int16](data[22:24]) - 40,                               // 12 - Max cell temperature (°C, offset -40)
		MaxTempRackNo:          utils.FromBytes[uint16](data[24:26]),                                   // 13 - Max temperature rack number
		MaxTempCellNo:          utils.FromBytes[uint16](data[26:28]),                                   // 14 - Max temperature cell number
		MinCellTemperature:     utils.FromBytes[int16](data[28:30]) - 40,                               // 15 - Min cell temperature (°C, offset -40)
		MinTempRackNo:          utils.FromBytes[uint16](data[30:32]),                                   // 16 - Min temperature rack number
		MinTempCellNo:          utils.FromBytes[uint16](data[32:34]),                                   // 17 - Min temperature cell number
		TotalChargeEnergy:      utils.Scale(utils.FromBytes[uint32](data[34:38]), float32(0.1)),        // 18-19 - Total charge energy (0.1kWh)
		TotalDischargeEnergy:   utils.Scale(utils.FromBytes[uint32](data[38:42]), float32(0.1)),        // 20-21 - Total discharge energy (0.1kWh)
		ChargeCapacity:         utils.Scale(utils.FromBytes[uint32](data[50:54]), float32(0.1)),        // 26-27 - Charge capacity (0.1kWh)
		DischargeCapacity:      utils.Scale(utils.FromBytes[uint32](data[54:58]), float32(0.1)),        // 28-29 - Discharge capacity (0.1kWh)
		AvailableDischargeTime: utils.FromBytes[uint16](data[58:60]),                                   // 30 - Available discharge time (min)
		AvailableChargeTime:    utils.FromBytes[uint16](data[60:62]),                                   // 31 - Available charge time (min)
		MaxDischargePower:      utils.Scale(utils.FromBytes[uint16](data[62:64]), float32(0.1)),        // 32 - Max discharge power (0.1kW)
		MaxChargePower:         utils.Scale(utils.FromBytes[uint16](data[64:66]), float32(0.1)),        // 33 - Max charge power (0.1kW)
		MaxDischargeCurrent:    utils.Scale(utils.FromBytes[uint16](data[66:68]), float32(0.1)),        // 34 - Max discharge current (0.1A)
		MaxChargeCurrent:       utils.Scale(utils.FromBytes[uint16](data[68:70]), float32(0.1)),        // 35 - Max charge current (0.1A)
		DischargeTimesToday:    utils.FromBytes[uint16](data[70:72]),                                   // 36 - Discharge times today
		ChargeTimesToday:       utils.FromBytes[uint16](data[72:74]),                                   // 37 - Charge times today
		DischargeEnergyToday:   utils.Scale(utils.FromBytes[uint32](data[74:78]), float32(0.1)),        // 38-39 - Discharge energy today (0.1kWh)
		ChargeEnergyToday:      utils.Scale(utils.FromBytes[uint32](data[78:82]), float32(0.1)),        // 40-41 - Charge energy today (0.1kWh)
		Temperature:            utils.FromBytes[int16](data[82:84]) - 40,                               // 42 - Temperature (°C, offset -40)
		State:                  utils.FromBytes[uint16](data[84:86]),                                   // 43 - State
		ChargeDischargeState:   utils.FromBytes[uint16](data[86:88]),                                   // 44 - Charge/discharge state
		InsulationResistance:   utils.FromBytes[uint16](data[88:90]),                                   // 45 - Insulation resistance (kΩ)
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
		MaxDischargeVoltage:  utils.Scale(utils.FromBytes[uint16](data[8:10]), float32(0.1)),           // 104 - Max discharge voltage (0.1V)
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
