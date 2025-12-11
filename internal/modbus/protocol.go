package modbus

// Modbus Server Register Address Map
const (
	// BMS System Data
	BMSBaseAddr        = 1000
	BMSDataOffset      = 100
	BMSDataStartOffset = 0
	BMSDataLength      = 41

	// PCS Data
	PCSBaseAddr        = 4000
	PCSDataOffset      = 300
	PCSDataStartOffset = 0
	PCSDataLength      = 68

	// Control Command Registers
	CmdBaseAddr             = 1000
	CmdOffset               = 100
	RegStartStopCommand     = 0
	RegActivePowerCommand   = 1
	RegReactivePowerCommand = 2
)
