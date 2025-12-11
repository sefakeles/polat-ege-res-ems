package pcs

// MODBUS Register addresses for Power Electronics
const (
	// Status Data
	StatusDataStartAddr = 1003
	StatusDataLength    = 1

	// Equipment Data
	EquipmentDataStartAddr = 631
	EquipmentDataLength    = 10

	// Environment Data
	EnvironmentDataStartAddr = 1104
	EnvironmentDataLength    = 1

	// DC Source Data
	DCSourceDataStartAddr = 1372
	DCSourceDataLength    = 16

	// Grid Data
	GridDataStartAddr = 4300
	GridDataLength    = 33

	// Counter Data
	CounterDataStartAddr = 539
	CounterDataLength    = 25

	// Fault Data
	FaultDataStartAddr = 1450
	FaultDataLength    = 16

	// Warning Data
	WarningDataStartAddr = 1512
	WarningDataLength    = 20

	// PCS Status Data
	PCSStatusDataStartAddr = 22001
	PCSStatusDataLength    = 72

	// Control
	CmdStartStopRegister     = 38
	CmdActivePowerRegister   = 862
	CmdReactivePowerRegister = 867
	HeartbeatRegister        = 8027
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

// PCS Status Codes
const (
	StatusPowerUp      = 0
	StatusInit         = 1
	StatusOFF          = 2
	StatusStandbyDC    = 3
	StatusPrechargeDC  = 4
	StatusSoftchargeDC = 5
	StatusReadyDC      = 6
	StatusStandbyAC    = 7
	StatusBlackstart   = 8
	StatusReady        = 9
	StatusWait         = 10
	StatusPreON        = 11
	StatusON           = 12
	StatusStopping     = 13
	StatusReadyAC      = 14
	StatusDiagnosticAC = 15
	StatusDischarge    = 16
	StatusFault        = 17
)

// AlarmDefinition defines the properties of an alarm
type AlarmDefinition struct {
	Message  string
	Severity string
}

// alarmDefinitions contains all alarm definitions
var alarmDefinitions = map[uint16]AlarmDefinition{
	1:   {"Watchdog", "HIGH"},
	2:   {"HW Vdc", "HIGH"},
	3:   {"SoftCharge", "HIGH"},
	4:   {"Discharge", "HIGH"},
	5:   {"High Vac", "HIGH"},
	6:   {"Low Vac", "HIGH"},
	7:   {"High frequency", "HIGH"},
	8:   {"Low frequency", "HIGH"},
	10:  {"Comms SCI FPGA-DSP", "HIGH"},
	11:  {"Active Anti-island", "HIGH"},
	13:  {"No modules", "HIGH"},
	14:  {"Drive-Select DSP", "HIGH"},
	15:  {"Synchronization", "HIGH"},
	23:  {"Unbalanced Vac", "HIGH"},
	25:  {"Low Vdc", "HIGH"},
	27:  {"Modules start fault", "HIGH"},
	28:  {"Passive Anti-island", "HIGH"},
	31:  {"Selfdiagnosis fault", "HIGH"},
	32:  {"Error Mod selfdiagnosis", "HIGH"},
	33:  {"Unable reconnect", "HIGH"},
	35:  {"MT Premag", "HIGH"},
	40:  {"Internal overtemperature", "HIGH"},
	41:  {"GFDI", "HIGH"},
	43:  {"Emergency stop", "HIGH"},
	44:  {"Drive-Select MCU", "HIGH"},
	45:  {"General Isolation", "HIGH"},
	46:  {"Data Fault", "HIGH"},
	47:  {"Watchdog uP", "HIGH"},
	48:  {"Internal comms", "HIGH"},
	49:  {"IMD Selfdiagnosis Error", "HIGH"},
	51:  {"PPC Comms", "HIGH"},
	54:  {"DU Overcurrent", "HIGH"},
	55:  {"External HW/CGB", "HIGH"},
	56:  {"Remote emergency stop", "HIGH"},
	58:  {"Control SW Mismatch", "HIGH"},
	59:  {"Module SW Mismatch", "HIGH"},
	62:  {"DU Comms", "HIGH"},
	63:  {"IMD Earth Connection", "HIGH"},
	64:  {"Invalid HMC", "HIGH"},
	65:  {"Overcurrent AC", "HIGH"},
	66:  {"Module VDC Unbalance", "HIGH"},
	69:  {"Overvoltage DC+", "HIGH"},
	70:  {"Overvoltage DC-", "HIGH"},
	71:  {"Desaturation R1(H)", "HIGH"},
	72:  {"Desaturation R2(H)", "HIGH"},
	73:  {"Desaturation R2(L)", "HIGH"},
	74:  {"Desaturation R3(L)", "HIGH"},
	75:  {"Desaturation S1(H)", "HIGH"},
	76:  {"Desaturation S2(H)", "HIGH"},
	77:  {"Desaturation S2(L)", "HIGH"},
	78:  {"Desaturation S3(L)", "HIGH"},
	79:  {"Desaturation T1(H)", "HIGH"},
	80:  {"Desaturation T2(H)", "HIGH"},
	81:  {"Desaturation T2(L)", "HIGH"},
	82:  {"Desaturation T3(L)", "HIGH"},
	83:  {"Multiple desaturation", "HIGH"},
	84:  {"Comunications", "HIGH"},
	85:  {"Timeout soft charge", "HIGH"},
	86:  {"Feedback breaker", "HIGH"},
	90:  {"T° IGBT's HF", "HIGH"},
	95:  {"Source PCB", "HIGH"},
	98:  {"Idc derivation", "HIGH"},
	99:  {"Med. Iac derivation", "HIGH"},
	102: {"Unbalanced current", "HIGH"},
	103: {"IGBT's temperature", "HIGH"},
	104: {"PCB temperature", "HIGH"},
	106: {"Mod. Vdc unbalanced", "HIGH"},
	107: {"Iac not reached", "MEDIUM"},
	108: {"Mod. High Vdc", "HIGH"},
	109: {"Mod. Low Vdc", "HIGH"},
	112: {"OCAC emergency", "HIGH"},
	113: {"CRC", "HIGH"},
	114: {"Measure Ir", "HIGH"},
	115: {"Measure Is", "HIGH"},
	116: {"Measure It", "HIGH"},
	118: {"Module Disabled", "HIGH"},
	119: {"Desaturation R(L)", "HIGH"},
	120: {"Desaturation R3(H)", "HIGH"},
	121: {"Desaturation S1(L)", "HIGH"},
	122: {"Desaturation S3(H)", "HIGH"},
	123: {"Desaturation T1(L)", "HIGH"},
	124: {"Desaturation T3(H)", "HIGH"},
	126: {"Diff Iac mod run", "HIGH"},
	128: {"Spring loaded", "HIGH"},
	129: {"Module Isolation", "HIGH"},
	130: {"Softcharge Feedback", "HIGH"},
	131: {"Faults limit exceeded", "HIGH"},
	132: {"VDC HW overvoltage", "HIGH"},
	134: {"Module sensor comms", "HIGH"},
	139: {"Mod. low Vdc recover", "MEDIUM"},
	150: {"MV cell wrong state", "HIGH"},
	152: {"DGPT2 Temp failure", "HIGH"},
	153: {"DGPT2 Pressure failure", "HIGH"},
	155: {"DGPT2 oil or gas level", "HIGH"},
	157: {"Feedback AC breaker", "HIGH"},
	158: {"Low pressure SF6", "HIGH"},
	159: {"Overtemperature DU", "HIGH"},
	161: {"LOTO AC", "HIGH"},
	162: {"LOTO DC", "HIGH"},
	163: {"MV Fault", "HIGH"},
	164: {"Excessive Maneuvers", "MEDIUM"},
	165: {"Improper Modulation", "HIGH"},
	166: {"E-STOP", "HIGH"},
	167: {"Open DU", "HIGH"},
	169: {"LCL Feedback", "HIGH"},
	170: {"Tcontactor Feedback", "HIGH"},
	171: {"Filter LC Current", "HIGH"},
	172: {"C LC Filter", "HIGH"},
	173: {"CTs measure", "HIGH"},
	174: {"AC circuit breaker fb", "HIGH"},
	175: {"Overcurrent C", "HIGH"},
	176: {"Multi. Overcurr Events", "HIGH"},
	177: {"Wrong Filt C Config.", "HIGH"},
	190: {"Combined F-V", "HIGH"},
	197: {"HW Overvoltage", "HIGH"},
	198: {"HW Undervoltage", "HIGH"},
	200: {"Overvoltage Lim DC", "HIGH"},
	202: {"Taux overtemperature", "HIGH"},
	203: {"MV overtemperature", "HIGH"},
	204: {"Fan impulsion MT", "MEDIUM"},
	205: {"Fan DU zone", "MEDIUM"},
	206: {"Fan zone filter LC", "MEDIUM"},
	207: {"Fan zone modules", "MEDIUM"},
	208: {"MV crit overtemperature", "HIGH"},
	210: {"uC-FPGA Comms", "HIGH"},
	211: {"Intern Isolation", "HIGH"},
	212: {"Extern Isolation", "HIGH"},
	213: {"Negative DU current", "HIGH"},
	214: {"IMD Comms", "HIGH"},
	215: {"Timeout GFDI", "HIGH"},
	216: {"Timeout IMI Meas", "HIGH"},
	217: {"Critical Temp LC filter", "HIGH"},
	218: {"Overvoltage Lim AC", "HIGH"},
	220: {"PT100 induc not connected", "MEDIUM"},
	221: {"PT100 trafo not connected", "MEDIUM"},
	222: {"NTC C intake not connect", "MEDIUM"},
	223: {"NTC C filter not connec", "MEDIUM"},
	224: {"Ind overtemperature", "HIGH"},
	225: {"Ind crit overtemperature", "HIGH"},
	226: {"IMI Hardware", "HIGH"},
	227: {"Fault Unknown", "HIGH"},
	228: {"Fan setpoint", "MEDIUM"},
	229: {"Timeout IMI Channel", "HIGH"},
	230: {"Fan type not correct DU", "MEDIUM"},
	231: {"Fan type not correct MT", "MEDIUM"},
	232: {"Fan type not correct MOD", "MEDIUM"},
	234: {"Thermal Bus Plus", "HIGH"},
	235: {"T/H sensors comms", "HIGH"},
	236: {"DU sensor comms lost", "HIGH"},
	237: {"Negative union fuse", "HIGH"},
	238: {"AEMO comms lost", "HIGH"},
	240: {"Du BP precurrent", "HIGH"},
	241: {"Delta_T max DU", "HIGH"},
	242: {"Delta_T LC", "HIGH"},
	243: {"Fan saturated setpoint", "MEDIUM"},
	244: {"T LC Slent mode", "MEDIUM"},
	252: {"RFI OverCur Persistent", "HIGH"},
	253: {"Imbalance V", "HIGH"},
	255: {"BAT Open detection", "HIGH"},
}

// warningDefinitions contains all warning definitions
var warningDefinitions = map[uint16]AlarmDefinition{
	5:   {"No start conditions", "LOW"},
	12:  {"Derating thermal IGBT", "LOW"},
	13:  {"Derating thermal adm", "LOW"},
	28:  {"Derating freq temp IGBT", "LOW"},
	29:  {"DGPT2 Temperature alarm", "LOW"},
	30:  {"Overtemperature DU", "LOW"},
	46:  {"No modules active heating", "LOW"},
	53:  {"NO SYNC FAC", "LOW"},
	56:  {"Silent mode derating", "LOW"},
	58:  {"Limit Q per VAC", "LOW"},
	151: {"Lockable module", "LOW"},
	156: {"High resistor isolation", "LOW"},
	158: {"General Isolation", "LOW"},
	159: {"Intern Isolation", "LOW"},
	160: {"Extern Isolation", "LOW"},
	161: {"Fan impulsion MT", "LOW"},
	162: {"Fan DU zone", "LOW"},
	163: {"Fan zone filter LC", "LOW"},
	164: {"Fan zone modules", "LOW"},
	165: {"Fan motor status", "LOW"},
	166: {"Temp sensor read", "LOW"},
	181: {"Aux transf overvoltaje", "LOW"},
	182: {"Therma Bus Plus", "LOW"},
	184: {"MV transformer Temp", "LOW"},
	185: {"IMI not calibrated", "LOW"},
	186: {"Supply low", "LOW"},
	187: {"DU Delta Temp Threshold", "LOW"},
	188: {"DU sensor lost comms", "LOW"},
	190: {"AEMO comms lost", "LOW"},
	191: {"Antipid wrong V", "LOW"},
	192: {"Antipid V not reached", "LOW"},
	193: {"Antipid retries exceed", "LOW"},
	196: {"Fan Incorect Frames", "LOW"},
	252: {"RFI OverCur Persistent", "LOW"},
	253: {"Imbalance V", "LOW"},
	255: {"BAT open detection", "LOW"},
	303: {"RFI OverCurrent", "LOW"},
	310: {"Module Disable", "LOW"},
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

// GetWarningMessage returns warning message based on code
func GetWarningMessage(code uint16) string {
	if def, exists := warningDefinitions[code]; exists {
		return def.Message
	}
	return "Unknown warning"
}

// GetWarningSeverity returns warning severity based on code
func GetWarningSeverity(code uint16) string {
	if def, exists := warningDefinitions[code]; exists {
		return def.Severity
	}
	return "LOW"
}

// GetStatusString returns the human-readable status name
func GetStatusString(status uint16) string {
	switch status {
	case StatusPowerUp:
		return "Power Up"
	case StatusInit:
		return "Init"
	case StatusOFF:
		return "OFF"
	case StatusStandbyDC:
		return "Standby DC"
	case StatusPrechargeDC:
		return "Precharge DC"
	case StatusSoftchargeDC:
		return "Softcharge DC"
	case StatusReadyDC:
		return "Ready DC"
	case StatusStandbyAC:
		return "Standby AC"
	case StatusBlackstart:
		return "Blackstart"
	case StatusReady:
		return "Ready"
	case StatusWait:
		return "Wait"
	case StatusPreON:
		return "Pre ON"
	case StatusON:
		return "ON"
	case StatusStopping:
		return "Stopping"
	case StatusReadyAC:
		return "Ready AC"
	case StatusDiagnosticAC:
		return "Diagnostic AC"
	case StatusDischarge:
		return "Discharge"
	case StatusFault:
		return "Fault"
	default:
		return "Unknown"
	}
}

// IsReadyState checks if the PCS is in the Ready state
func IsReadyState(state uint16) bool {
	return state == StatusReady
}

// IsFaultState checks if the PCS is in the Fault state
func IsFaultState(state uint16) bool {
	return state == StatusFault
}
