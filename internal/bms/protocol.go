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
	ControlOn          = 2
	ControlOff         = 3
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
	1:   {"Single cell over-voltage warning level 1", "LOW"},
	2:   {"Single cell over-voltage warning level 2", "LOW"},
	3:   {"Single cell over-voltage warning level 3", "MEDIUM"},
	4:   {"Single cell under-voltage warning level 1", "LOW"},
	5:   {"Single cell under-voltage warning level 2", "LOW"},
	6:   {"Single cell under-voltage warning level 3", "MEDIUM"},
	7:   {"Single cell extreme over-voltage warning", "MEDIUM"},
	8:   {"Single cell extreme under-voltage warning", "MEDIUM"},
	9:   {"Single rack over-voltage warning level 1", "LOW"},
	10:  {"Single rack over-voltage warning level 2", "MEDIUM"},
	11:  {"Single rack under-voltage warning level 1", "LOW"},
	12:  {"Single rack under-voltage warning level 2", "MEDIUM"},
	13:  {"Big voltage difference between cells warning", "LOW"},
	14:  {"Invalid cell voltage warning", "MEDIUM"},
	15:  {"Replacement required due to single cell over discharge", "MEDIUM"},
	16:  {"Replacement required due to single cell over charge", "MEDIUM"},
	17:  {"Low SOC over discharge warning", "LOW"},
	18:  {"Low SOC over discharge warning", "MEDIUM"},
	19:  {"Discharge over-current warning level 1", "LOW"},
	20:  {"Discharge over-current warning level 2", "LOW"},
	21:  {"Discharge over-current warning level 3", "MEDIUM"},
	22:  {"Discharge over-current warning level 4", "MEDIUM"},
	23:  {"Charge over-current warning level 1", "LOW"},
	24:  {"Charge over-current warning level 2", "LOW"},
	25:  {"Charge over-current warning level 3", "MEDIUM"},
	26:  {"Charge over-current warning level 4", "MEDIUM"},
	27:  {"Current sensor failed warning", "MEDIUM"},
	31:  {"Single cell over-temperature warning level 1", "LOW"},
	32:  {"Single cell over-temperature warning level 2", "LOW"},
	33:  {"Single cell over-temperature warning level 3", "MEDIUM"},
	34:  {"Single cell over-temperature warning level 4", "MEDIUM"},
	35:  {"Single cell under-temperature warning level 1", "LOW"},
	36:  {"Single cell under-temperature warning level 2", "LOW"},
	37:  {"Single cell under-temperature warning level 3", "MEDIUM"},
	38:  {"Big temperature difference between cells warning level 1", "LOW"},
	39:  {"Big temperature difference between cells warning level 2", "LOW"},
	40:  {"Big temperature difference between cells warning level 3", "MEDIUM"},
	48:  {"Single temperature sampling abnormal warning", "MEDIUM"},
	49:  {"Multiple temperature sampling abnormal warning", "MEDIUM"},
	74:  {"CSC 24V power supply abnormal warning", "MEDIUM"},
	75:  {"SBMU 24V power supply abnormal warning", "MEDIUM"},
	79:  {"MSD warning", "MEDIUM"},
	80:  {"Rack fuse warning", "MEDIUM"},
	81:  {"Rack isolate switch warning", "MEDIUM"},
	82:  {"Main positive relay sticking warning", "MEDIUM"},
	83:  {"Main negative relay sticking warning", "MEDIUM"},
	84:  {"Both main positive relay and main negative relay sticking fault", "HIGH"},
	85:  {"Main positive relay open circuit warning", "MEDIUM"},
	86:  {"Main negative open circuit warning", "MEDIUM"},
	87:  {"Battery rack door (travel switch) fault", "HIGH"},
	88:  {"Main control box fan warning", "MEDIUM"},
	102: {"Inner communication warning (CCAN)", "MEDIUM"},
	103: {"Inner communication warning (SCAN)", "MEDIUM"},
	104: {"Inner communication warning (MCAN)", "MEDIUM"},
	108: {"Balancing circuit warning", "MEDIUM"},
	110: {"SOC low warning level 1", "LOW"},
	111: {"SOC low warning level 2", "LOW"},
	112: {"HV circuit open circuit warning", "MEDIUM"},
	116: {"Pre-charging failed twice warning", "MEDIUM"},
	122: {"HV+&HV- reversed connection fault", "HIGH"},
	129: {"Thermal runaway caused fire fault", "HIGH"},
	134: {"HV circuit (Fuse) open circuit warning", "MEDIUM"},
	153: {"Big temperature difference between racks", "LOW"},
	257: {"Insulation warning", "LOW"},
	258: {"Insulation fault", "HIGH"},
	259: {"Invalid insulation fault", "HIGH"},
	260: {"The quantity of HV racks less than setting value fault", "HIGH"},
	263: {"SBMU communication lost warning", "MEDIUM"},
	264: {"SBMU communication lost fault", "HIGH"},
	265: {"EMS communication lost fault", "HIGH"},
	266: {"IMM communication lost fault", "HIGH"},
	267: {"Air conditioner communication lost warning", "MEDIUM"},
	268: {"Centralized TMS communication fault", "HIGH"},
	269: {"Centralized TMS communication warning", "LOW"},
	270: {"Centralized TMS fault level 2", "HIGH"},
	271: {"Centralized TMS mode conflict fault", "HIGH"},
	272: {"Centralized TMS fault level 1", "HIGH"},
	273: {"SPD failure warning", "LOW"},
	274: {"AUX Power DCDC failure warning", "LOW"},
	275: {"AUX Power ACDC failure warning", "LOW"},
	276: {"AUX Power failure fault", "HIGH"},
	277: {"Fire system fault level 1", "HIGH"},
	278: {"Fire system fault level 2", "HIGH"},
	279: {"Fire system failure warning", "MEDIUM"},
	280: {"E-STOP fault", "HIGH"},
	281: {"Client E-STOP fault", "HIGH"},
	282: {"Electrical rack door (travel switch) fault", "HIGH"},
	283: {"Electrical rack fan warning", "LOW"},
	284: {"Smoke exhaust ventilation body warning", "MEDIUM"},
	285: {"Smoke exhaust ventilation state fault", "HIGH"},
	286: {"Humidifier1 failure warning", "LOW"},
	289: {"Humidifier2 failure warning", "LOW"},
	308: {"MBMU 24V power supply abnormal fault", "HIGH"},
	309: {"Slave MBMU communication lost fault", "HIGH"},
	312: {"Centralized TMS warning level 3", "MEDIUM"},
	320: {"Step-Charge Mode inconsistency warning", "LOW"},
}

// rackAlarmDefinitions contains all rack alarm definitions
var rackAlarmDefinitions = map[uint16]AlarmDefinition{
	0:   {"Single cell over-voltage warning level 1", "LOW"},
	1:   {"Single cell over-voltage warning level 2", "LOW"},
	2:   {"Single cell over-voltage warning level 3", "MEDIUM"},
	3:   {"Single cell under-voltage warning level 1", "LOW"},
	4:   {"Single cell under-voltage warning level 2", "LOW"},
	5:   {"Single cell under-voltage warning level 3", "MEDIUM"},
	6:   {"Single cell extreme over-voltage fault", "MEDIUM"},
	7:   {"Single cell extreme under-voltage fault", "MEDIUM"},
	8:   {"Single rack over-voltage warning level 1", "LOW"},
	9:   {"Single rack over-voltage warning level 2", "MEDIUM"},
	10:  {"Single rack under-voltage warning level 1", "LOW"},
	11:  {"Single rack under-voltage warning level 2", "MEDIUM"},
	12:  {"Big voltage difference between cells warning", "LOW"},
	13:  {"Invalid cell voltage fault", "MEDIUM"},
	14:  {"Replacement required due to single cell over discharge", "HIGH"},
	15:  {"Replacement required due to single cell over charge", "HIGH"},
	16:  {"Low SOC over discharge warning", "LOW"},
	17:  {"Low SOC over discharge fault", "MEDIUM"},
	18:  {"Discharge over-current warning level 1", "LOW"},
	19:  {"Discharge over-current warning level 2", "LOW"},
	20:  {"Discharge over-current warning level 3", "MEDIUM"},
	21:  {"Discharge over-current fault level 4", "MEDIUM"},
	22:  {"Charge over-current warning level 1", "LOW"},
	23:  {"Charge over-current warning level 2", "LOW"},
	24:  {"Charge over-current warning level 3", "MEDIUM"},
	25:  {"Charge over-current fault level 4", "MEDIUM"},
	26:  {"Current sensor failed fault", "MEDIUM"},
	30:  {"Single cell over-temperature warning level 1", "LOW"},
	31:  {"Single cell over-temperature warning level 2", "LOW"},
	32:  {"Single cell over-temperature warning level 3", "MEDIUM"},
	33:  {"Single cell over-temperature fault level 4", "MEDIUM"},
	34:  {"Single cell under-temperature warning level 1", "LOW"},
	35:  {"Single cell under-temperature warning level 2", "LOW"},
	36:  {"Single cell under-temperature warning level 3", "MEDIUM"},
	37:  {"Big temperature difference between cells warning level 1", "LOW"},
	38:  {"Big temperature difference between cells warning level 2", "LOW"},
	39:  {"Big temperature difference between cells warning level 3", "MEDIUM"},
	47:  {"Single temperature sampling abnormal fault", "MEDIUM"},
	48:  {"Multiple temperature sampling abnormal fault", "MEDIUM"},
	73:  {"CSC 24V power supply abnormal fault", "MEDIUM"},
	74:  {"SBMU 24V power supply abnormal fault", "MEDIUM"},
	78:  {"MSD fault", "MEDIUM"},
	79:  {"Rack fuse fault", "MEDIUM"},
	80:  {"Rack isolate switch fault", "MEDIUM"},
	81:  {"Main positive relay sticking fault", "MEDIUM"},
	82:  {"Main negative relay sticking fault", "MEDIUM"},
	83:  {"Both main positive relay and main negative relay sticking fault", "HIGH"},
	84:  {"Main positive relay open circuit fault", "MEDIUM"},
	85:  {"Main negative open circuit fault", "MEDIUM"},
	86:  {"Battery rack door (travel switch) fault", "HIGH"},
	87:  {"Main control box fan fault", "MEDIUM"},
	101: {"Inner communication fault (CCAN)", "MEDIUM"},
	102: {"Inner communication fault (SCAN)", "MEDIUM"},
	103: {"Inner communication fault (MCAN)", "MEDIUM"},
	107: {"Balancing circuit fault", "MEDIUM"},
	109: {"SOC low warning level 1", "LOW"},
	110: {"SOC low warning level 2", "LOW"},
	111: {"HV circuit open circuit fault", "MEDIUM"},
	115: {"Pre-charging failed twice fault", "HIGH"},
	121: {"HV+&HV- reversed connection fault", "HIGH"},
	128: {"Thermal runaway caused fire fault", "HIGH"},
	133: {"HV circuit (Fuse) open circuit fault", "MEDIUM"},
	152: {"Big temperature difference between racks", "LOW"},
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
