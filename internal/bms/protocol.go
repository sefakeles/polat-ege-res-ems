package bms

// MODBUS Register addresses for Jinko
const (
	// BMS Data
	BMSDataStartAddr = 32
	BMSDataLength    = 24

	// BMS Status Data
	BMSStatusDataStartAddr = 768
	BMSStatusDataLength    = 7

	// Rack Data
	BMSRackDataStartAddr = 100
	BMSRackDataOffset    = 3000
	BMSRackDataLength    = 51

	// Alarms
	BMSAlarmStartAddr     = 0
	BMSAlarmLength        = 32
	BMSRackAlarmStartAddr = 1024
	BMSRackAlarmOffset    = 1024
	BMSRackAlarmLength    = 16

	// Cell Data
	CellVoltageBaseAddr = 191
	CellTempBaseAddr    = 891

	// Control
	HeartbeatRegister      = 896
	BreakerControlRegister = 897
	SystemResetRegister    = 908

	// Cell organization constants
	CellsPerModule       = 48
	TempSensorsPerModule = 12
)

// Run Commands
const (
	CommandDisable = 0
	CommandEnable  = 1
)

// Control Commands
const (
	ControlNoOperation = 0
	ControlReset       = 1
	ControlClose       = 1
	ControlOpen        = 2
	ControlOn          = 1
	ControlOff         = 2
)

// State Codes
const (
	StateInitial     = 0
	StateNormal      = 1
	StateCharging    = 2
	StateDischarging = 3
	StateWarning     = 4
	StateFault       = 5
)

// AlarmDefinition defines the properties of an alarm
type AlarmDefinition struct {
	Message  string
	Severity string
}

// alarmDefinitions contains all alarm definitions
var alarmDefinitions = map[uint16]AlarmDefinition{
	1:  {"Single cell over-voltage warning level 1", "LOW"},
	2:  {"Single cell over-voltage warning level 2", "LOW"},
	3:  {"Single cell over-voltage warning level 3", "MEDIUM"},
	4:  {"Single cell under-voltage warning level 1", "LOW"},
	5:  {"Single cell under-voltage warning level 2", "LOW"},
	6:  {"Single cell under-voltage warning level 3", "MEDIUM"},
	7:  {"Single cell extreme over-voltage warning", "MEDIUM"},
	8:  {"Single cell extreme under-voltage warning", "MEDIUM"},
	9:  {"Single rack over-voltage warning level 1", "LOW"},
	10: {"Single rack over-voltage warning level 2", "MEDIUM"},
	11: {"Single rack under-voltage warning level 1", "LOW"},
	12: {"Single rack under-voltage warning level 2", "MEDIUM"},
	13: {"Big voltage difference between cells warning", "LOW"},
	14: {"Invalid cell voltage warning", "MEDIUM"},
	15: {"Replacement required due to single cell over discharge", "MEDIUM"},
	16: {"Replacement required due to single cell over charge", "MEDIUM"},
	17: {"Low SOC over discharge warning", "LOW"},
	18: {"Low SOC over discharge warning", "MEDIUM"},
	19: {"Discharge over-current warning level 1", "LOW"},
	20: {"Discharge over-current warning level 2", "LOW"},
	21: {"Discharge over-current warning level 3", "MEDIUM"},
	22: {"Discharge over-current warning level 4", "MEDIUM"},
	23: {"Charge over-current warning level 1", "LOW"},
	24: {"Charge over-current warning level 2", "LOW"},
	25: {"Charge over-current warning level 3", "MEDIUM"},
	26: {"Charge over-current warning level 4", "MEDIUM"},
	27: {"Current sensor failed warning", "MEDIUM"},
	28: {"Single cell over-temperature warning level 1", "LOW"},
	29: {"Single cell over-temperature warning level 2", "LOW"},
	30: {"Single cell over-temperature warning level 3", "MEDIUM"},
	31: {"Single cell over-temperature warning level 4", "MEDIUM"},
	32: {"Single cell under-temperature warning level 1", "LOW"},
	33: {"Single cell under-temperature warning level 2", "LOW"},
	34: {"Single cell under-temperature warning level 3", "MEDIUM"},
	35: {"Big temperature difference between cells warning level 1", "LOW"},
	36: {"Big temperature difference between cells warning level 2", "LOW"},
	37: {"Big temperature difference between cells warning level 3", "MEDIUM"},
	38: {"Single temperature sampling abnormal warning", "MEDIUM"},
	39: {"Multiple temperature sampling abnormal warning", "MEDIUM"},
	40: {"CSC 24V power supply abnormal warning", "MEDIUM"},
	41: {"SBMU 24V power supply abnormal warning", "MEDIUM"},
	42: {"MSD warning", "MEDIUM"},
	43: {"Rack fuse warning", "MEDIUM"},
	44: {"Rack isolate switch warning", "MEDIUM"},
	45: {"Main positive relay sticking warning", "MEDIUM"},
	46: {"Main negative relay sticking warning", "MEDIUM"},
	47: {"Both main positive relay and main negative relay sticking fault", "HIGH"},
	48: {"Main positive relay open circuit warning", "MEDIUM"},
	49: {"Main negative open circuit warning", "MEDIUM"},
	50: {"Battery rack door (travel switch) fault", "HIGH"},
	51: {"Main control box fan warning", "MEDIUM"},
	52: {"Inner communication warning (CCAN)", "MEDIUM"},
	53: {"Inner communication warning (SCAN)", "MEDIUM"},
	54: {"Inner communication warning (MCAN)", "MEDIUM"},
	55: {"Balancing circuit warning", "MEDIUM"},
	56: {"SOC low warning level 1", "LOW"},
	57: {"SOC low warning level 2", "LOW"},
	58: {"HV circuit open circuit warning", "MEDIUM"},
	59: {"Pre-charging failed twice warning", "MEDIUM"},
	60: {"HV+&HV- reversed connection fault", "HIGH"},
	61: {"Thermal runaway caused fire fault", "HIGH"},
	62: {"HV circuit (Fuse) open circuit warning", "MEDIUM"},
	63: {"Big temperature difference between racks", "LOW"},
	64: {"Insulation warning", "LOW"},
	65: {"Insulation fault", "HIGH"},
	66: {"Invalid insulation fault", "HIGH"},
	67: {"The quantity of HV racks less than setting value fault", "HIGH"},
	68: {"SBMU communication lost warning", "MEDIUM"},
	69: {"SBMU communication lost fault", "HIGH"},
	70: {"EMS communication lost fault", "HIGH"},
	71: {"IMM communication lost fault", "HIGH"},
	72: {"Air conditioner communication lost warning", "MEDIUM"},
	73: {"Centralized TMS communication fault", "HIGH"},
	74: {"Centralized TMS communication warning", "LOW"},
	75: {"Centralized TMS fault level 2", "HIGH"},
	76: {"Centralized TMS mode conflict fault", "HIGH"},
	77: {"Centralized TMS fault level 1", "HIGH"},
	78: {"SPD failure warning", "LOW"},
	79: {"AUX Power DCDC failure warning", "LOW"},
	80: {"AUX Power ACDC failure warning", "LOW"},
	81: {"AUX Power failure fault", "HIGH"},
	82: {"Fire system fault level 1", "HIGH"},
	83: {"Fire system fault level 2", "HIGH"},
	84: {"Fire system failure warning", "MEDIUM"},
	85: {"E-STOP fault", "HIGH"},
	86: {"Client E-STOP fault", "HIGH"},
	87: {"Electrical rack door (travel switch) fault", "HIGH"},
	88: {"Electrical rack fan warning", "LOW"},
	89: {"Smoke exhaust ventilation body warning", "MEDIUM"},
	90: {"Smoke exhaust ventilation state fault", "HIGH"},
	91: {"Humidifier1 failure warning", "LOW"},
	92: {"Humidifier2 failure warning", "LOW"},
	93: {"MBMU 24V power supply abnormal fault", "HIGH"},
	94: {"Slave MBMU communication lost fault", "HIGH"},
	95: {"Centralized TMS warning level 3", "MEDIUM"},
	96: {"Step-Charge Mode inconsistency warning", "LOW"},
}

// rackAlarmDefinitions contains all rack alarm definitions
var rackAlarmDefinitions = map[uint16]AlarmDefinition{
	0:  {"BCU communication lost", "HIGH"},
	1:  {"Total voltage low", "LOW"},
	2:  {"Total voltage low", "MEDIUM"},
	3:  {"Total voltage low", "HIGH"},
	4:  {"Total voltage high", "LOW"},
	5:  {"Total voltage high", "MEDIUM"},
	6:  {"Total voltage high", "HIGH"},
	7:  {"Current high", "LOW"},
	8:  {"Current high", "MEDIUM"},
	9:  {"Current high", "HIGH"},
	10: {"Cell voltage low", "LOW"},
	11: {"Cell voltage low", "MEDIUM"},
	12: {"Cell voltage low", "HIGH"},
	13: {"Cell voltage high", "LOW"},
	14: {"Cell voltage high", "MEDIUM"},
	15: {"Cell voltage high", "HIGH"},
	16: {"Cell temperature low", "LOW"},
	17: {"Cell temperature low", "MEDIUM"},
	18: {"Cell temperature low", "HIGH"},
	19: {"Cell temperature high", "LOW"},
	20: {"Cell temperature high", "MEDIUM"},
	21: {"Cell temperature high", "HIGH"},
	22: {"Cell SOC low", "LOW"},
	23: {"Cell SOC low", "MEDIUM"},
	24: {"Cell SOC low", "HIGH"},
	25: {"Cell SOC high", "LOW"},
	26: {"Cell SOC high", "MEDIUM"},
	27: {"Cell SOC high", "HIGH"},
	28: {"Cell SOH low", "LOW"},
	29: {"Cell SOH low", "MEDIUM"},
	30: {"Cell SOH low", "HIGH"},
	31: {"Cell voltage imbalance", "LOW"},
	32: {"Cell voltage imbalance", "MEDIUM"},
	33: {"Cell voltage imbalance", "HIGH"},
	34: {"Cell temperature imbalance", "LOW"},
	35: {"Cell temperature imbalance", "MEDIUM"},
	36: {"Cell temperature imbalance", "HIGH"},
	37: {"BMU 1 communication lost", "HIGH"},
	38: {"BMU 2 communication lost", "HIGH"},
	39: {"BMU 3 communication lost", "HIGH"},
	40: {"BMU 4 communication lost", "HIGH"},
	41: {"BMU 5 communication lost", "HIGH"},
	42: {"BMU 6 communication lost", "HIGH"},
	43: {"BMU 7 communication lost", "HIGH"},
	44: {"BMU 8 communication lost", "HIGH"},
	45: {"BMU 9 communication lost", "HIGH"},
	46: {"BMU 10 communication lost", "HIGH"},
	47: {"BMU 11 communication lost", "HIGH"},
	48: {"BMU 12 communication lost", "HIGH"},
	49: {"BMU 13 communication lost", "HIGH"},
	50: {"BMU 14 communication lost", "HIGH"},
	51: {"BMU 15 communication lost", "HIGH"},
	52: {"BMU 16 communication lost", "HIGH"},
	53: {"BMU 17 communication lost", "HIGH"},
	54: {"BMU 18 communication lost", "HIGH"},
	55: {"BMU 19 communication lost", "HIGH"},
	56: {"BMU 20 communication lost", "HIGH"},
	57: {"BMU 21 communication lost", "HIGH"},
	58: {"BMU 22 communication lost", "HIGH"},
	59: {"BMU 23 communication lost", "HIGH"},
	60: {"BMU 24 communication lost", "HIGH"},
	61: {"BMU 25 communication lost", "HIGH"},
	62: {"BMU 26 communication lost", "HIGH"},
	63: {"BMU 27 communication lost", "HIGH"},
	64: {"BMU 28 communication lost", "HIGH"},
	65: {"BMU 29 communication lost", "HIGH"},
	66: {"BMU 30 communication lost", "HIGH"},
	67: {"BMU 31 communication lost", "HIGH"},
	68: {"BMU 32 communication lost", "HIGH"},
	69: {"BMU 33 communication lost", "HIGH"},
	70: {"BMU 34 communication lost", "HIGH"},
	71: {"BMU 35 communication lost", "HIGH"},
	72: {"BMU 36 communication lost", "HIGH"},
	73: {"BMU 37 communication lost", "HIGH"},
	74: {"BMU 38 communication lost", "HIGH"},
	75: {"BMU 39 communication lost", "HIGH"},
	76: {"BMU 40 communication lost", "HIGH"},
	77: {"Terminal temperature high", "LOW"},
	78: {"Terminal temperature high", "MEDIUM"},
	79: {"Terminal temperature high", "HIGH"},
	80: {"Module voltage high", "LOW"},
	81: {"Module voltage high", "MEDIUM"},
	82: {"Module voltage high", "HIGH"},
	83: {"Module voltage low", "LOW"},
	84: {"Module voltage low", "MEDIUM"},
	85: {"Module voltage low", "HIGH"},
	86: {"Cell voltage sensor fault", "HIGH"},
	87: {"Cell temperature sensor fault", "HIGH"},
}

// GetAlarmMessage returns alarm message based on code
func GetAlarmMessage(code uint16) string {
	if def, exists := alarmDefinitions[code]; exists {
		return def.Message
	}
	return "Unknown alarm"
}

// GetAlarmSeverity returns alarm severity based on code
func GetAlarmSeverity(code uint16) string {
	if def, exists := alarmDefinitions[code]; exists {
		return def.Severity
	}
	return "LOW"
}

// GetRackAlarmMessage returns the alarm message for a rack alarm
func GetRackAlarmMessage(relativeCode uint16) string {
	if def, exists := rackAlarmDefinitions[relativeCode]; exists {
		return def.Message
	}
	return "Unknown alarm"
}

// GetRackAlarmSeverity returns the severity for a rack alarm
func GetRackAlarmSeverity(relativeCode uint16) string {
	if def, exists := rackAlarmDefinitions[relativeCode]; exists {
		return def.Severity
	}
	return "LOW"
}

// GetStateDescription returns human-readable state description
func GetStateDescription(state uint16) string {
	switch state {
	case StateInitial:
		return "Initial"
	case StateNormal:
		return "Normal"
	case StateCharging:
		return "Charging"
	case StateDischarging:
		return "Discharging"
	case StateWarning:
		return "Warning"
	case StateFault:
		return "Fault"
	default:
		return "Unknown"
	}
}

// IsNormalState checks if the state is normal
func IsNormalState(state uint16) bool {
	return state == StateNormal
}

// IsChargingState checks if the state is charging
func IsChargingState(state uint16) bool {
	return state == StateCharging
}

// IsDischargingState checks if the state is discharging
func IsDischargingState(state uint16) bool {
	return state == StateDischarging
}

// IsWarningState checks if the state is warning
func IsWarningState(state uint16) bool {
	return state == StateWarning
}

// IsFaultState checks if the state is fault
func IsFaultState(state uint16) bool {
	return state == StateFault
}

// GetRackDataStartAddr returns the starting address for data of a specific rack
func GetRackDataStartAddr(rackNo uint8) uint16 {
	return BMSRackDataStartAddr + uint16(rackNo-1)*BMSRackDataOffset
}

// GetRackAlarmStartAddr returns the starting address for alarm of a specific rack
func GetRackAlarmStartAddr(rackNo uint8) uint16 {
	return BMSRackAlarmStartAddr + uint16(rackNo-1)*BMSRackAlarmOffset
}

// GetCellVoltageStartAddr returns the starting address for cell voltages of a specific rack
func GetCellVoltageStartAddr(rackNo uint8) uint16 {
	return CellVoltageBaseAddr + uint16(rackNo-1)*BMSRackDataOffset
}

// GetCellTempStartAddr returns the starting address for cell temperatures of a specific rack
func GetCellTempStartAddr(rackNo uint8) uint16 {
	return CellTempBaseAddr + uint16(rackNo-1)*BMSRackDataOffset
}

// CalculateReadChunks calculates the number of chunks needed to read data
func CalculateReadChunks(registerCount, maxRegistersPerRead int) int {
	return (registerCount + maxRegistersPerRead - 1) / maxRegistersPerRead
}
