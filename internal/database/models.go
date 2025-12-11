package database

import (
	"time"
)

// BMSData represents BMS data
type BMSData struct {
	Timestamp               time.Time `json:"timestamp"`
	ID                      int       `json:"id"`
	Voltage                 float32   `json:"voltage"`
	Current                 int16     `json:"current"`
	SOC                     float32   `json:"soc"`
	SOH                     float32   `json:"soh"`
	MaxCellVoltage          float32   `json:"max_cell_voltage"`
	MinCellVoltage          float32   `json:"min_cell_voltage"`
	AvgCellVoltage          float32   `json:"avg_cell_voltage"`
	MaxCellTemperature      int16     `json:"max_cell_temperature"`
	MinCellTemperature      int16     `json:"min_cell_temperature"`
	AvgCellTemperature      int16     `json:"avg_cell_temperature"`
	MaxChargeCurrent        int16     `json:"max_charge_current"`
	MaxDischargeCurrent     int16     `json:"max_discharge_current"`
	MaxChargePower          int16     `json:"max_charge_power"`
	MaxDischargePower       int16     `json:"max_discharge_power"`
	Power                   int16     `json:"power"`
	ChargeCapacity          uint16    `json:"charge_capacity"`
	DischargeCapacity       uint16    `json:"discharge_capacity"`
	MaxChargeVoltage        float32   `json:"max_charge_voltage"`
	MaxDischargeVoltage     float32   `json:"max_discharge_voltage"`
	InsulationResistancePos uint16    `json:"insulation_resistance_pos"`
	InsulationResistanceNeg uint16    `json:"insulation_resistance_neg"`
}

// BMSStatusData represents BMS status data
type BMSStatusData struct {
	Timestamp      time.Time `json:"timestamp"`
	ID             int       `json:"id"`
	Heartbeat      uint16    `json:"heartbeat"`
	HVStatus       uint16    `json:"hv_status"`
	SystemStatus   uint16    `json:"system_status"`
	ConnectedRacks uint16    `json:"connected_racks"`
	TotalRacks     uint16    `json:"total_racks"`
}

// BMSRackData represents BMS rack-level data
type BMSRackData struct {
	Timestamp            time.Time `json:"timestamp"`
	ID                   int       `json:"id"`
	Number               uint8     `json:"number"`
	State                uint16    `json:"state"`
	MaxChargePower       float32   `json:"max_charge_power"`
	MaxDischargePower    float32   `json:"max_discharge_power"`
	MaxChargeVoltage     float32   `json:"max_charge_voltage"`
	MaxDischargeVoltage  float32   `json:"max_discharge_voltage"`
	MaxChargeCurrent     float32   `json:"max_charge_current"`
	MaxDischargeCurrent  float32   `json:"max_discharge_current"`
	Voltage              float32   `json:"voltage"`
	Current              float32   `json:"current"`
	Temperature          int16     `json:"temperature"`
	SOC                  uint16    `json:"soc"`
	SOH                  uint16    `json:"soh"`
	InsulationResistance uint16    `json:"insulation_resistance"`
	AvgCellVoltage       float32   `json:"avg_cell_voltage"`
	AvgCellTemperature   int16     `json:"avg_cell_temperature"`
	MaxCellVoltage       float32   `json:"max_cell_voltage"`
	MaxVoltageCellNo     uint16    `json:"max_voltage_cell_no"`
	MinCellVoltage       float32   `json:"min_cell_voltage"`
	MinVoltageCellNo     uint16    `json:"min_voltage_cell_no"`
	MaxCellTemperature   int16     `json:"max_cell_temperature"`
	MaxTempCellNo        uint16    `json:"max_temp_cell_no"`
	MinCellTemperature   int16     `json:"min_cell_temperature"`
	MinTempCellNo        uint16    `json:"min_temp_cell_no"`
	TotalChargeEnergy    float32   `json:"total_charge_energy"`
	TotalDischargeEnergy float32   `json:"total_discharge_energy"`
	ChargeCapacity       float32   `json:"charge_capacity"`
	DischargeCapacity    float32   `json:"discharge_capacity"`
}

// BMSCellVoltageData represents individual cell voltage data
type BMSCellVoltageData struct {
	Timestamp time.Time `json:"timestamp"`
	ID        int       `json:"id"`
	RackNo    uint8     `json:"rack_no"`
	ModuleNo  uint8     `json:"module_no"`
	CellNo    uint16    `json:"cell_no"`
	Voltage   float32   `json:"voltage"`
}

// BMSCellTemperatureData represents individual cell temperature data
type BMSCellTemperatureData struct {
	Timestamp   time.Time `json:"timestamp"`
	ID          int       `json:"id"`
	RackNo      uint8     `json:"rack_no"`
	ModuleNo    uint8     `json:"module_no"`
	SensorNo    uint16    `json:"sensor_no"`
	Temperature int16     `json:"temperature"`
}

type PCSData struct {
	StatusData      PCSStatusData      `json:"status_data"`
	EquipmentData   PCSEquipmentData   `json:"equipment_data"`
	EnvironmentData PCSEnvironmentData `json:"environment_data"`
	DCSourceData    PCSDCSourceData    `json:"dc_source_data"`
	GridData        PCSGridData        `json:"grid_data"`
	CounterData     PCSCounterData     `json:"counter_data"`
}

type PCSStatusData struct {
	Timestamp time.Time `json:"timestamp"`
	ID        int       `json:"id"`
	Status    uint16    `json:"status"`
}

type PCSEquipmentData struct {
	Timestamp              time.Time `json:"timestamp"`
	ID                     int       `json:"id"`
	LVSwitchStatus         uint16    `json:"lv_switch_status"`
	MVSwitchStatus         uint16    `json:"mv_switch_status"`
	MVDisconnectorStatus   uint16    `json:"mv_disconnector_status"`
	MVEarthingSwitchStatus uint16    `json:"mv_earthing_switch_status"`
	DC1SwitchStatus        uint16    `json:"dc1_switch_status"`
	DC2SwitchStatus        uint16    `json:"dc2_switch_status"`
	DC3SwitchStatus        uint16    `json:"dc3_switch_status"`
	DC4SwitchStatus        uint16    `json:"dc4_switch_status"`
}

type PCSEnvironmentData struct {
	Timestamp           time.Time `json:"timestamp"`
	ID                  int       `json:"id"`
	AirInletTemperature int16     `json:"air_inlet_temperature"`
}

type PCSDCSourceData struct {
	Timestamp          time.Time `json:"timestamp"`
	ID                 int       `json:"id"`
	DC1Power           int16     `json:"dc1_power"`
	DC2Power           int16     `json:"dc2_power"`
	DC3Power           int16     `json:"dc3_power"`
	DC4Power           int16     `json:"dc4_power"`
	DC1Current         uint16    `json:"dc1_current"`
	DC2Current         uint16    `json:"dc2_current"`
	DC3Current         uint16    `json:"dc3_current"`
	DC4Current         uint16    `json:"dc4_current"`
	DC1VoltageExternal float32   `json:"dc1_voltage_external"`
	DC2VoltageExternal float32   `json:"dc2_voltage_external"`
	DC3VoltageExternal float32   `json:"dc3_voltage_external"`
	DC4VoltageExternal float32   `json:"dc4_voltage_external"`
}

type PCSGridData struct {
	Timestamp           time.Time `json:"timestamp"`
	ID                  int       `json:"id"`
	MVGridVoltageAB     float32   `json:"mv_grid_voltage_ab"`
	MVGridVoltageBC     float32   `json:"mv_grid_voltage_bc"`
	MVGridVoltageCA     float32   `json:"mv_grid_voltage_ca"`
	MVGridCurrentA      float32   `json:"mv_grid_current_a"`
	MVGridCurrentB      float32   `json:"mv_grid_current_b"`
	MVGridCurrentC      float32   `json:"mv_grid_current_c"`
	MVGridActivePower   int16     `json:"mv_grid_active_power"`
	MVGridReactivePower int16     `json:"mv_grid_reactive_power"`
	MVGridApparentPower uint16    `json:"mv_grid_apparent_power"`
	MVGridCosPhi        float32   `json:"mv_grid_cos_phi"`
	LVGridVoltageAB     float32   `json:"lv_grid_voltage_ab"`
	LVGridVoltageBC     float32   `json:"lv_grid_voltage_bc"`
	LVGridVoltageCA     float32   `json:"lv_grid_voltage_ca"`
	LVGridCurrentA      float32   `json:"lv_grid_current_a"`
	LVGridCurrentB      float32   `json:"lv_grid_current_b"`
	LVGridCurrentC      float32   `json:"lv_grid_current_c"`
	LVGridActivePower   int16     `json:"lv_grid_active_power"`
	LVGridReactivePower int16     `json:"lv_grid_reactive_power"`
	LVGridApparentPower uint16    `json:"lv_grid_apparent_power"`
	LVGridCosPhi        float32   `json:"lv_grid_cos_phi"`
	GridFrequency       float32   `json:"grid_frequency"`
}

type PCSCounterData struct {
	Timestamp               time.Time `json:"timestamp"`
	ID                      int       `json:"id"`
	ActiveEnergyToday       uint32    `json:"active_energy_today"`
	ActiveEnergyYesterday   uint32    `json:"active_energy_yesterday"`
	ActiveEnergyThisMonth   uint32    `json:"active_energy_this_month"`
	ActiveEnergyLastMonth   uint32    `json:"active_energy_last_month"`
	ActiveEnergyTotal       uint32    `json:"active_energy_total"`
	ConsumedEnergyToday     uint32    `json:"consumed_energy_today"`
	ConsumedEnergyTotal     uint32    `json:"consumed_energy_total"`
	ReactiveEnergyToday     int32     `json:"reactive_energy_today"`
	ReactiveEnergyYesterday int32     `json:"reactive_energy_yesterday"`
	ReactiveEnergyThisMonth int32     `json:"reactive_energy_this_month"`
	ReactiveEnergyLastMonth int32     `json:"reactive_energy_last_month"`
	ReactiveEnergyTotal     int32     `json:"reactive_energy_total"`
}

// PCSCommandState represents the current command state
type PCSCommandState struct {
	LastUpdated          time.Time `json:"last_updated"`
	StartStopCommand     bool      `json:"start_stop_command"`
	ActivePowerCommand   float32   `json:"active_power_command"`
	ReactivePowerCommand float32   `json:"reactive_power_command"`
}

// BMSCommandState represents the current command state
type BMSCommandState struct {
	LastUpdated      time.Time `json:"last_updated"`
	StartStopCommand bool      `json:"start_stop_command"`
}

// BMSAlarmData represents BMS alarm information
type BMSAlarmData struct {
	Timestamp time.Time `json:"timestamp"`
	AlarmType string    `json:"alarm_type"`
	AlarmCode uint16    `json:"alarm_code"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Active    bool      `json:"active"`
}

// PLCData represents PLC data
type PLCData struct {
	Timestamp         time.Time              `json:"timestamp"`
	ID                int                    `json:"id"`
	CircuitBreakers   CircuitBreakerStatus   `json:"circuit_breakers"`
	MVCircuitBreakers MVCircuitBreakerStatus `json:"mv_circuit_breakers"`
	ProtectionRelays  ProtectionRelayStatus  `json:"protection_relays"`
}

// CircuitBreakerStatus represents the status of all circuit breakers
type CircuitBreakerStatus struct {
	AuxiliaryCB bool `json:"auxiliary_cb"`
	PCS1CB      bool `json:"pcs1_cb"`
	PCS2CB      bool `json:"pcs2_cb"`
	PCS3CB      bool `json:"pcs3_cb"`
	PCS4CB      bool `json:"pcs4_cb"`
	BMS1CB      bool `json:"bms1_cb"`
	BMS2CB      bool `json:"bms2_cb"`
	BMS3CB      bool `json:"bms3_cb"`
	BMS4CB      bool `json:"bms4_cb"`
}

// MVCircuitBreakerStatus represents MV-side circuit breaker status
type MVCircuitBreakerStatus struct {
	AuxTransformerCB bool `json:"aux_transformer_cb"`
	Transformer1CB   bool `json:"transformer1_cb"`
	Transformer2CB   bool `json:"transformer2_cb"`
	Transformer3CB   bool `json:"transformer3_cb"`
	Transformer4CB   bool `json:"transformer4_cb"`
}

// ProtectionRelayStatus represents protection relay status
type ProtectionRelayStatus struct {
	AuxTransformerFault bool `json:"aux_transformer_fault"`
	Transformer1Fault   bool `json:"transformer1_fault"`
	Transformer2Fault   bool `json:"transformer2_fault"`
	Transformer3Fault   bool `json:"transformer3_fault"`
	Transformer4Fault   bool `json:"transformer4_fault"`
}

// SystemMetrics represents system performance metrics
type SystemMetrics struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUUsage    float32   `json:"cpu_usage"`
	MemoryUsage float32   `json:"memory_usage"`
	DiskUsage   float32   `json:"disk_usage"`
	NetworkRx   uint64    `json:"network_rx"`
	NetworkTx   uint64    `json:"network_tx"`
}

// RuntimeMetrics represents application runtime performance metrics
type RuntimeMetrics struct {
	Timestamp time.Time `json:"timestamp"`

	// General metrics
	UptimeSeconds float64 `json:"uptime_seconds"`
	Goroutines    int     `json:"goroutines"`

	// Memory metrics (in MB)
	HeapAllocMB    float64 `json:"heap_alloc_mb"`
	HeapSysMB      float64 `json:"heap_sys_mb"`
	HeapIdleMB     float64 `json:"heap_idle_mb"`
	HeapInUseMB    float64 `json:"heap_in_use_mb"`
	HeapReleasedMB float64 `json:"heap_released_mb"`
	StackInUseMB   float64 `json:"stack_in_use_mb"`
	StackSysMB     float64 `json:"stack_sys_mb"`

	// GC metrics
	GCRuns         uint32  `json:"gc_runs"`
	GCPauseTotalNs uint64  `json:"gc_pause_total_ns"`
	GCCPUFraction  float64 `json:"gc_cpu_fraction"`
	NextGCMB       float64 `json:"next_gc_mb"`
	LastGCTime     int64   `json:"last_gc_time"`

	// Allocation metrics
	MallocsTotal uint64  `json:"mallocs_total"`
	FreesTotal   uint64  `json:"frees_total"`
	TotalAllocMB float64 `json:"total_alloc_mb"`
	LookupsTotal uint64  `json:"lookups_total"`
}
