package windfarm

// MODBUS Register addresses for ENERCON FCU (Farm Control Unit)
// Based on ENERCON Modbus TCP Protocol Document

const (
	// Heartbeat counter
	HeartbeatAddr = 600

	// Setpoint registers
	PSetpointAddr           = 610
	QSetpointAddr           = 611
	PowerFactorSetpointAddr = 612

	// Wind farm control
	WindFarmStartStopAddr   = 629
	RapidDownwardSignalAddr = 639
)

const (
	// Return values length
	ReturnValuesStartAddr = 649
	ReturnValuesLength    = 41 // 649-689
)

const (
	// Measuring data length
	MeasuringDataStartAddr = 700
	MeasuringDataLength    = 60
)

// Wind Farm Control Commands
const (
	WindFarmStart = 0 // Start wind farm
	WindFarmStop  = 1 // Stop wind farm
)

// Rapid Downward Signal
const (
	RapidDownwardOff = 0 // Rapid downward signal off
	RapidDownwardOn  = 1 // Rapid downward signal on
)

// FCU Status
const (
	FCUOffline = 0 // FCU offline
	FCUOnline  = 1 // FCU online
)
