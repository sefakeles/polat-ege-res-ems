package windfarm

import (
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/utils"
)

// ParseMeasuringData converts raw MODBUS data to WindFarmMeasuringData structure
// Expects data starting from register 700 (MeasuringDataStartAddr)
func ParseMeasuringData(data []byte, id int) database.WindFarmMeasuringData {
	if len(data) < MeasuringDataLength*2 {
		return database.WindFarmMeasuringData{Timestamp: time.Now(), ID: id}
	}

	return database.WindFarmMeasuringData{
		Timestamp:                 time.Now(),
		ID:                        id,
		ActivePowerNCP:            utils.Scale(utils.FromBytes[int16](data[0:2]), float32(0.01)),    // 700 - Active power (MW), scale 0.01
		ReactivePowerNCP:          utils.Scale(utils.FromBytes[int16](data[2:4]), float32(0.01)),    // 701 - Reactive power (Mvar), scale 0.01
		VoltageNCP:                utils.Scale(utils.FromBytes[int16](data[4:6]), float32(0.01)),    // 702 - Voltage (kV), scale 0.01
		CurrentNCP:                utils.Scale(utils.FromBytes[int16](data[6:8]), float32(0.1)),     // 703 - Current (A), scale 0.1
		PowerFactorNCP:            utils.Scale(utils.FromBytes[int16](data[8:10]), float32(0.001)),  // 704 - Power factor, scale 0.001
		WECAvailability:           utils.FromBytes[uint16](data[10:12]),                             // 705 - WEC availability (%)
		FrequencyNCP:              utils.Scale(utils.FromBytes[uint16](data[12:14]), float32(0.01)), // 706 - Frequency (Hz), scale 0.01
		WindSpeed:                 utils.Scale(utils.FromBytes[uint16](data[40:42]), float32(0.01)), // 720 - Wind speed (m/s), scale 0.01
		WindDirection:             utils.FromBytes[uint16](data[42:44]),                             // 721 - Wind direction (°)
		PossibleWECPower:          utils.Scale(utils.FromBytes[int16](data[44:46]), float32(0.01)),  // 722 - Possible WEC power (MW), scale 0.01
		WECCommunication:          utils.FromBytes[uint16](data[48:50]),                             // 724 - WEC communication (%)
		RelativePowerAvailability: utils.Scale(utils.FromBytes[int16](data[80:82]), float32(0.01)),  // 740 - Relative power avail (%)
		AbsolutePowerAvailability: utils.Scale(utils.FromBytes[int16](data[82:84]), float32(0.01)),  // 741 - Absolute power avail (MW)
		RelativeMinReactivePower:  utils.Scale(utils.FromBytes[int16](data[84:86]), float32(0.01)),  // 742 - Relative min Q (%)
		AbsoluteMinReactivePower:  utils.Scale(utils.FromBytes[int16](data[86:88]), float32(0.01)),  // 743 - Absolute min Q (MVar)
		RelativeMaxReactivePower:  utils.Scale(utils.FromBytes[int16](data[88:90]), float32(0.01)),  // 744 - Relative max Q (%)
		AbsoluteMaxReactivePower:  utils.Scale(utils.FromBytes[int16](data[90:92]), float32(0.01)),  // 745 - Absolute max Q (MVar)
	}
}

// ParseStatusData converts raw MODBUS data to WindFarmStatusData structure
// Expects data starting from register 649 (ReturnValuesStartAddr)
func ParseStatusData(data []byte, id int) database.WindFarmStatusData {
	if len(data) < ReturnValuesLength*2 {
		return database.WindFarmStatusData{Timestamp: time.Now(), ID: id}
	}

	fcuOnlineStatus := utils.FromBytes[uint16](data[2:4])     // 650 - FCU online status
	activePowerMode := utils.FromBytes[uint16](data[44:46])   // 671 - Active power mode currently used
	reactivePowerMode := utils.FromBytes[uint16](data[46:48]) // 672 - Reactive power mode currently used
	windFarmStatus := utils.FromBytes[uint16](data[60:62])    // 679 - Wind farm start/stop mirror
	rapidDownward := utils.FromBytes[uint16](data[80:82])     // 689 - Rapid downward signal mirror

	return database.WindFarmStatusData{
		Timestamp:                 time.Now(),
		ID:                        id,
		FCUOnline:                 fcuOnlineStatus == FCUOnline,
		FCUHeartbeatCounter:       utils.FromBytes[uint16](data[0:2]), // 649 - heartbeat counter
		ActivePowerControlMode:    activePowerMode,
		ReactivePowerControlMode:  reactivePowerMode,
		WindFarmRunning:           windFarmStatus == WindFarmStart,
		RapidDownwardSignalActive: rapidDownward == RapidDownwardOn,
	}
}

// ParseSetpointData converts raw MODBUS data to WindFarmSetpointData structure
// Expects data starting from register 649 (ReturnValuesStartAddr)
func ParseSetpointData(data []byte, id int) database.WindFarmSetpointData {
	if len(data) < ReturnValuesLength*2 {
		return database.WindFarmSetpointData{Timestamp: time.Now(), ID: id}
	}

	return database.WindFarmSetpointData{
		Timestamp: time.Now(),
		ID:        id,
		// Setpoint mirrors (commanded values)
		PSetpointMirror:          utils.Scale(utils.FromBytes[int16](data[20:22]), float32(0.01)),   // 659 - P setpoint mirror
		QSetpointMirror:          utils.Scale(utils.FromBytes[int16](data[22:24]), float32(0.01)),   // 660 - Q setpoint mirror
		PowerFactorMirror:        utils.Scale(utils.FromBytes[int16](data[24:26]), float32(0.001)),  // 661 - Power factor mirror
		USetpointMirror:          utils.Scale(utils.FromBytes[int16](data[26:28]), float32(0.01)),   // 662 - U setpoint mirror
		QdUSetpointMirror:        utils.Scale(utils.FromBytes[int16](data[28:30]), float32(0.01)),   // 663 - Q(dU) setpoint mirror
		DPDtMinMirror:            utils.Scale(utils.FromBytes[uint16](data[30:32]), float32(0.001)), // 664 - dP/dt min mirror
		DPDtMaxMirror:            utils.Scale(utils.FromBytes[uint16](data[32:34]), float32(0.001)), // 665 - dP/dt max mirror
		FrequencyReserveCapacity: utils.FromBytes[uint16](data[38:40]),                              // 668 - Frequency reserve capacity
		PfDeadbandMirror:         utils.Scale(utils.FromBytes[uint16](data[40:42]), float32(0.001)), // 669 - P(f) deadband mirror
		PfSlopeMirror:            utils.Scale(utils.FromBytes[uint16](data[42:44]), float32(0.001)), // 670 - P(f) slope mirror
		// Currently used setpoints
		PSetpointCurrent:   utils.Scale(utils.FromBytes[int16](data[48:50]), float32(0.01)),  // 673 - P setpoint current
		QSetpointCurrent:   utils.Scale(utils.FromBytes[int16](data[50:52]), float32(0.01)),  // 674 - Q setpoint current
		PowerFactorCurrent: utils.Scale(utils.FromBytes[int16](data[52:54]), float32(0.001)), // 675 - Power factor current
		USetpointCurrent:   utils.Scale(utils.FromBytes[int16](data[54:56]), float32(0.01)),  // 676 - U setpoint current
		QdUSetpointCurrent: utils.Scale(utils.FromBytes[int16](data[56:58]), float32(0.01)),  // 677 - Q(dU) setpoint current
	}
}

// ParseWeatherData converts raw MODBUS data to WindFarmWeatherData structure
// Expects data starting from register 699 (MeasuringDataStartAddr)
func ParseWeatherData(data []byte, id int) database.WindFarmWeatherData {
	if len(data) < MeasuringDataLength*2 {
		return database.WindFarmWeatherData{Timestamp: time.Now(), ID: id}
	}

	return database.WindFarmWeatherData{
		Timestamp:                time.Now(),
		ID:                       id,
		WindSpeedMeteo:           utils.Scale(utils.FromBytes[uint16](data[60:62]), float32(0.1)),  // 729 - Wind speed meteo (m/s)
		WindDirectionMeteo:       utils.Scale(utils.FromBytes[int16](data[62:64]), float32(0.1)),   // 730 - Wind direction meteo (°)
		OutsideTemperature:       utils.Scale(utils.FromBytes[int16](data[64:66]), float32(0.1)),   // 731 - Outside temp (°C)
		AtmosphericPressure:      utils.FromBytes[uint16](data[66:68]),                             // 732 - Atmospheric pressure (mbar)
		AirHumidity:              utils.Scale(utils.FromBytes[uint16](data[68:70]), float32(0.1)),  // 733 - Air humidity (%)
		RainfallVolume:           utils.Scale(utils.FromBytes[uint16](data[70:72]), float32(0.01)), // 734 - Rainfall (l/m²h)
		SolarRadiation:           utils.Scale(utils.FromBytes[uint16](data[72:74]), float32(0.1)),  // 735 - Solar radiation (W/m²)
		WindFarmCommunication:    utils.FromBytes[uint16](data[74:76]),                             // 736 - Wind farm comm (%)
		WeatherMeasurementsCount: utils.FromBytes[uint16](data[76:78]),                             // 737 - Weather measurements count
	}
}

// ParseFCUMode extracts FCU mode from measuring data
// Expects data starting from register 699 (MeasuringDataStartAddr)
func ParseFCUMode(data []byte) uint16 {
	if len(data) < MeasuringDataLength*2 {
		return 0
	}
	// FCU mode is at register 758, which is offset 118 bytes from 699 (59 registers * 2 bytes)
	return utils.FromBytes[uint16](data[118:120])
}
