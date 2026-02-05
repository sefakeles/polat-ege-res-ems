package bms

import (
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/utils"
)

// ParseBMSStatusData converts raw MODBUS data to BMSStatusData structure
func ParseBMSStatusData(data []byte, id int) database.BMSStatusData {
	if len(data) < BMSStatusDataLength*2 {
		return database.BMSStatusData{Timestamp: time.Now(), ID: id}
	}

	return database.BMSStatusData{
		Timestamp:        time.Now(),
		ID:               id,
		Heartbeat:        utils.FromBytes[uint16](data[0:2]),   // 768 - Heartbeat
		HVStatus:         utils.FromBytes[uint16](data[2:4]),   // 769 - High voltage status
		SystemStatus:     utils.FromBytes[uint16](data[4:6]),   // 770 - System status
		ConnectedRacks:   utils.FromBytes[uint16](data[8:10]),  // 772 - Connected racks
		TotalRacks:       utils.FromBytes[uint16](data[10:12]), // 773 - Total racks
		StepChargeStatus: utils.FromBytes[uint16](data[12:14]), // 774 - Step charge status
	}
}

// ParseBMSData converts raw MODBUS data to BMSData structure
func ParseBMSData(data []byte, id int) database.BMSData {
	if len(data) < BMSDataLength*2 {
		return database.BMSData{Timestamp: time.Now(), ID: id}
	}

	return database.BMSData{
		Timestamp:                 time.Now(),
		ID:                        id,
		Voltage:                   utils.Scale(utils.FromBytes[uint16](data[0:2]), float32(0.1)),     // 32 - Voltage (0.1V)
		Current:                   utils.FromBytes[int16](data[2:4]) - 20000,                         // 33 - Current (A, offset -20000)
		SOC:                       utils.Scale(utils.FromBytes[uint16](data[4:6]), float32(0.1)),     // 34 - SOC (0.1%)
		SOH:                       utils.Scale(utils.FromBytes[uint16](data[6:8]), float32(0.1)),     // 35 - SOH (0.1%)
		MaxCellVoltage:            utils.Scale(utils.FromBytes[uint16](data[8:10]), float32(0.001)),  // 36 - Max cell voltage (0.001V)
		MinCellVoltage:            utils.Scale(utils.FromBytes[uint16](data[10:12]), float32(0.001)), // 37 - Min cell voltage (0.001V)
		AvgCellVoltage:            utils.Scale(utils.FromBytes[uint16](data[12:14]), float32(0.001)), // 38 - Average cell voltage (0.001V)
		MaxCellTemperature:        utils.FromBytes[int16](data[14:16]) - 50,                          // 39 - Max cell temperature (°C, offset -50)
		MinCellTemperature:        utils.FromBytes[int16](data[16:18]) - 50,                          // 40 - Min cell temperature (°C, offset -50)
		AvgCellTemperature:        utils.FromBytes[int16](data[18:20]) - 50,                          // 41 - Average cell temperature (°C, offset -50)
		MaxChargeCurrent:          utils.FromBytes[int16](data[20:22]) - 20000,                       // 42 - Max charge current (A, offset -20000)
		MaxDischargeCurrent:       utils.FromBytes[int16](data[22:24]) - 20000,                       // 43 - Max discharge current (A, offset -20000)
		MaxChargePower:            utils.FromBytes[int16](data[24:26]) - 20000,                       // 44 - Max charge power (kW, offset -20000)
		MaxDischargePower:         utils.FromBytes[int16](data[26:28]) - 20000,                       // 45 - Max discharge power (kW, offset -20000)
		Power:                     utils.FromBytes[int16](data[28:30]) - 20000,                       // 46 - Power (kW, offset -20000)
		ChargeCapacity:            utils.FromBytes[uint16](data[34:36]),                              // 49 - Charge capacity (kWh)
		DischargeCapacity:         utils.FromBytes[uint16](data[36:38]),                              // 50 - Discharge capacity (kWh)
		MaxChargeVoltage:          utils.Scale(utils.FromBytes[uint16](data[38:40]), float32(0.1)),   // 51 - Max charge voltage (0.1V)
		MinDischargeVoltage:       utils.Scale(utils.FromBytes[uint16](data[40:42]), float32(0.1)),   // 52 - Min discharge voltage (0.1V)
		InsulationDetectionStatus: utils.FromBytes[uint16](data[42:44]),                              // 53 - Insulation detection status
		InsulationResistancePos:   utils.FromBytes[uint16](data[44:46]),                              // 54 - Insulation resistance positive (kΩ)
		InsulationResistanceNeg:   utils.FromBytes[uint16](data[46:48]),                              // 55 - Insulation resistance negative (kΩ)
	}
}

// ParseBMSRackStatusData converts raw MODBUS data to BMSRackStatusData structure
func ParseBMSRackStatusData(data []byte, id int, rackNo uint8) database.BMSRackStatusData {
	if len(data) < BMSRackStatusDataLength*2 {
		return database.BMSRackStatusData{
			Timestamp: time.Now(),
			ID:        id,
			Number:    rackNo,
		}
	}

	return database.BMSRackStatusData{
		Timestamp:            time.Now(),
		ID:                   id,
		Number:               rackNo,
		PreChargeRelayStatus: utils.FromBytes[uint16](data[0:2]),   // 1040 - Pre-charge relay status
		PositiveRelayStatus:  utils.FromBytes[uint16](data[2:4]),   // 1041 - Positive relay status
		NegativeRelayStatus:  utils.FromBytes[uint16](data[4:6]),   // 1042 - Negative relay status
		HVStatus:             utils.FromBytes[uint16](data[6:8]),   // 1043 - High voltage status
		SOCMaintenanceStatus: utils.FromBytes[uint16](data[8:10]),  // 1044 - SOC maintenance status
		StepChargeStatus:     utils.FromBytes[uint16](data[10:12]), // 1045 - Step charge status
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
		VoltageOut:           utils.Scale(utils.FromBytes[uint16](data[0:2]), float32(0.1)),         // 1056 - Voltage outside (0.1V)
		Voltage:              utils.Scale(utils.FromBytes[uint16](data[2:4]), float32(0.1)),         // 1057 - Voltage inside (0.1V)
		Current:              utils.Scale(utils.FromBytes[int16](data[4:6])-20000, float32(0.1)),    // 1058 - Current (0.1A, offset -20000)
		SOC:                  utils.Scale(utils.FromBytes[uint16](data[6:8]), float32(0.1)),         // 1059 - SOC (0.1%)
		SOH:                  utils.Scale(utils.FromBytes[uint16](data[8:10]), float32(0.1)),        // 1060 - SOH (0.1%)
		MaxCellVoltage:       utils.Scale(utils.FromBytes[uint16](data[10:12]), float32(0.001)),     // 1061 - Max cell voltage (0.001V)
		MinCellVoltage:       utils.Scale(utils.FromBytes[uint16](data[12:14]), float32(0.001)),     // 1062 - Min cell voltage (0.001V)
		AvgCellVoltage:       utils.Scale(utils.FromBytes[uint16](data[14:16]), float32(0.001)),     // 1063 - Average cell voltage (0.001V)
		MaxCellTemperature:   utils.FromBytes[int16](data[16:18]) - 50,                              // 1064 - Max cell temperature (°C, offset -50)
		MinCellTemperature:   utils.FromBytes[int16](data[18:20]) - 50,                              // 1065 - Min cell temperature (°C, offset -50)
		AvgCellTemperature:   utils.FromBytes[int16](data[20:22]) - 50,                              // 1066 - Average cell temperature (°C, offset -50)
		MaxChargeCurrent:     utils.Scale(utils.FromBytes[int16](data[22:24])-20000, float32(0.1)),  // 1067 - Max charge current (A, offset -20000)
		MaxDischargeCurrent:  utils.Scale(utils.FromBytes[int16](data[24:26])-20000, float32(0.1)),  // 1068 - Max discharge current (A, offset -20000)
		MaxChargePower:       utils.Scale(utils.FromBytes[int16](data[26:28])-20000, float32(0.1)),  // 1069 - Max charge power (kW, offset -20000)
		MaxDischargePower:    utils.Scale(utils.FromBytes[int16](data[28:30])-20000, float32(0.1)),  // 1070 - Max discharge power (kW, offset -20000)
		Power:                utils.Scale(utils.FromBytes[int16](data[30:32])-20000, float32(0.1)),  // 1071 - Power (kW, offset -20000)
		MaxVoltageModuleNo:   data[33],                                                              // 1072 - Max voltage module number (high byte)
		MaxVoltageCellNo:     data[32],                                                              // 1072 - Max voltage cell number (low byte)
		MinVoltageModuleNo:   data[35],                                                              // 1073 - Min voltage module number (high byte)
		MinVoltageCellNo:     data[34],                                                              // 1073 - Min voltage cell number (low byte)
		MaxTempCellNo:        utils.FromBytes[uint16](data[36:38]),                                  // 1074 - Max temperature cell number
		MinTempCellNo:        utils.FromBytes[uint16](data[38:40]),                                  // 1075 - Min temperature cell number
		ChargeCapacity:       utils.Scale(utils.FromBytes[uint16](data[40:42]), float32(0.1)),       // 1076 - Charge capacity (kWh)
		DischargeCapacity:    utils.Scale(utils.FromBytes[uint16](data[42:44]), float32(0.1)),       // 1077 - Discharge capacity (kWh)
		MaxSelfDischargeRate: utils.Scale(utils.FromBytes[uint16](data[50:52]), float32(0.1)),       // 1081 - Max Cell Self Discharging rate (%)
		MinSelfDischargeRate: utils.Scale(utils.FromBytes[uint16](data[64:66]), float32(0.1)),       // 1088 - Min Cell Self Discharging rate (%)
		AvgSelfDischargeRate: utils.Scale(utils.FromBytes[uint16](data[66:68]), float32(0.1)),       // 1089 - Avg Cell Self Discharging rate (%)
		TotalChargeEnergy:    utils.Scale(utils.FromBytesCDAB[uint32](data[104:108]), float32(0.1)), // 1108-1109 - Total charge energy (kWh)
		TotalDischargeEnergy: utils.Scale(utils.FromBytesCDAB[uint32](data[108:112]), float32(0.1)), // 1110-1111 - Total discharge energy (kWh)
		CycleCount:           utils.FromBytes[uint16](data[120:122]),                                // 1116 - Cycle count
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
