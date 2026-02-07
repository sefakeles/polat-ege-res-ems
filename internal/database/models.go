package database

import "time"

// BMSStatusData represents BMS status data
type BMSStatusData struct {
	Timestamp        time.Time `json:"timestamp"`
	ID               int       `json:"id"`
	Heartbeat        uint16    `json:"heartbeat"`
	HVStatus         uint16    `json:"hv_status"`
	SystemStatus     uint16    `json:"system_status"`
	ConnectedRacks   uint16    `json:"connected_racks"`
	TotalRacks       uint16    `json:"total_racks"`
	StepChargeStatus uint16    `json:"step_charge_status"`
}

// BMSData represents BMS data
type BMSData struct {
	Timestamp                 time.Time `json:"timestamp"`
	ID                        int       `json:"id"`
	Voltage                   float32   `json:"voltage"`
	Current                   int16     `json:"current"`
	SOC                       float32   `json:"soc"`
	SOH                       float32   `json:"soh"`
	MaxCellVoltage            float32   `json:"max_cell_voltage"`
	MinCellVoltage            float32   `json:"min_cell_voltage"`
	AvgCellVoltage            float32   `json:"avg_cell_voltage"`
	MaxCellTemperature        int16     `json:"max_cell_temperature"`
	MinCellTemperature        int16     `json:"min_cell_temperature"`
	AvgCellTemperature        int16     `json:"avg_cell_temperature"`
	MaxChargeCurrent          int16     `json:"max_charge_current"`
	MaxDischargeCurrent       int16     `json:"max_discharge_current"`
	MaxChargePower            int16     `json:"max_charge_power"`
	MaxDischargePower         int16     `json:"max_discharge_power"`
	Power                     int16     `json:"power"`
	ChargeCapacity            uint16    `json:"charge_capacity"`
	DischargeCapacity         uint16    `json:"discharge_capacity"`
	MaxChargeVoltage          float32   `json:"max_charge_voltage"`
	MinDischargeVoltage       float32   `json:"min_discharge_voltage"`
	InsulationDetectionStatus uint16    `json:"insulation_detection_status"`
	InsulationResistancePos   uint16    `json:"insulation_resistance_pos"`
	InsulationResistanceNeg   uint16    `json:"insulation_resistance_neg"`
}

// BMSRackStatusData represents BMS rack status information
type BMSRackStatusData struct {
	Timestamp            time.Time `json:"timestamp"`
	ID                   int       `json:"id"`
	Number               uint8     `json:"number"`
	PreChargeRelayStatus uint16    `json:"pre_charge_relay_status"`
	PositiveRelayStatus  uint16    `json:"positive_relay_status"`
	NegativeRelayStatus  uint16    `json:"negative_relay_status"`
	HVStatus             uint16    `json:"hv_status"`
	SOCMaintenanceStatus uint16    `json:"soc_maintenance_status"`
	StepChargeStatus     uint16    `json:"step_charge_status"`
}

// BMSRackData represents BMS rack-level data
type BMSRackData struct {
	Timestamp            time.Time `json:"timestamp"`
	ID                   int       `json:"id"`
	Number               uint8     `json:"number"`
	Voltage              float32   `json:"voltage"`
	VoltageInside        float32   `json:"voltage_inside"`
	Current              float32   `json:"current"`
	SOC                  float32   `json:"soc"`
	SOH                  float32   `json:"soh"`
	MaxCellVoltage       float32   `json:"max_cell_voltage"`
	MinCellVoltage       float32   `json:"min_cell_voltage"`
	AvgCellVoltage       float32   `json:"avg_cell_voltage"`
	MaxCellTemperature   int16     `json:"max_cell_temperature"`
	MinCellTemperature   int16     `json:"min_cell_temperature"`
	AvgCellTemperature   int16     `json:"avg_cell_temperature"`
	MaxChargeCurrent     float32   `json:"max_charge_current"`
	MaxDischargeCurrent  float32   `json:"max_discharge_current"`
	MaxChargePower       float32   `json:"max_charge_power"`
	MaxDischargePower    float32   `json:"max_discharge_power"`
	Power                float32   `json:"power"`
	MaxVoltageCellNo     uint8     `json:"max_voltage_cell_no"`
	MaxVoltageModuleNo   uint8     `json:"max_voltage_module_no"`
	MinVoltageCellNo     uint8     `json:"min_voltage_cell_no"`
	MinVoltageModuleNo   uint8     `json:"min_voltage_module_no"`
	MaxTempModuleNo      uint16    `json:"max_temp_module_no"`
	MinTempModuleNo      uint16    `json:"min_temp_module_no"`
	ChargeCapacity       float32   `json:"charge_capacity"`
	DischargeCapacity    float32   `json:"discharge_capacity"`
	MaxSelfDischargeRate float32   `json:"max_self_discharge_rate"`
	MinSelfDischargeRate float32   `json:"min_self_discharge_rate"`
	AvgSelfDischargeRate float32   `json:"avg_self_discharge_rate"`
	TotalChargeEnergy    float32   `json:"total_charge_energy"`
	TotalDischargeEnergy float32   `json:"total_discharge_energy"`
	CycleCount           uint16    `json:"cycle_count"`
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
	AutoproducerCB   bool `json:"autoproducer_cb"`
}

// ProtectionRelayStatus represents protection relay status
type ProtectionRelayStatus struct {
	AuxTransformerFault bool `json:"aux_transformer_fault"`
	Transformer1Fault   bool `json:"transformer1_fault"`
	Transformer2Fault   bool `json:"transformer2_fault"`
	Transformer3Fault   bool `json:"transformer3_fault"`
	Transformer4Fault   bool `json:"transformer4_fault"`
}

// =============================================================================
// Wind Farm (ENERCON FCU) Data Models
// =============================================================================

// WindFarmData represents aggregated wind farm data
type WindFarmData struct {
	MeasuringData WindFarmMeasuringData `json:"measuring_data"`
	StatusData    WindFarmStatusData    `json:"status_data"`
	SetpointData  WindFarmSetpointData  `json:"setpoint_data"`
	WeatherData   WindFarmWeatherData   `json:"weather_data"`
}

// WindFarmMeasuringData represents real-time measuring data at NCP (Network Connection Point)
type WindFarmMeasuringData struct {
	Timestamp                 time.Time `json:"timestamp"`
	ID                        int       `json:"id"`
	ActivePowerNCP            float32   `json:"active_power_ncp"`            // MW, scale 0.01
	ReactivePowerNCP          float32   `json:"reactive_power_ncp"`          // Mvar, scale 0.01
	VoltageNCP                float32   `json:"voltage_ncp"`                 // kV, scale 0.01
	CurrentNCP                float32   `json:"current_ncp"`                 // A, scale 0.1
	PowerFactorNCP            float32   `json:"power_factor_ncp"`            // scale 0.001
	FrequencyNCP              float32   `json:"frequency_ncp"`               // Hz, scale 0.01
	WECAvailability           uint16    `json:"wec_availability"`            // %
	WindSpeed                 float32   `json:"wind_speed"`                  // m/s, scale 0.01
	WindDirection             uint16    `json:"wind_direction"`              // degrees
	PossibleWECPower          float32   `json:"possible_wec_power"`          // MW, scale 0.01
	WECCommunication          uint16    `json:"wec_communication"`           // %
	RelativePowerAvailability float32   `json:"relative_power_availability"` // %, scale 0.01
	AbsolutePowerAvailability float32   `json:"absolute_power_availability"` // MW, scale 0.01/0.1
	RelativeMinReactivePower  float32   `json:"relative_min_reactive_power"` // %, scale 0.01
	AbsoluteMinReactivePower  float32   `json:"absolute_min_reactive_power"` // MVar, scale 0.01/0.1
	RelativeMaxReactivePower  float32   `json:"relative_max_reactive_power"` // %, scale 0.01
	AbsoluteMaxReactivePower  float32   `json:"absolute_max_reactive_power"` // MVar, scale 0.01/0.1
}

// WindFarmStatusData represents FCU status data
type WindFarmStatusData struct {
	Timestamp                 time.Time `json:"timestamp"`
	ID                        int       `json:"id"`
	FCUOnline                 bool      `json:"fcu_online"`
	FCUMode                   uint16    `json:"fcu_mode"`                    // 0=Standard, 2=Master
	FCUHeartbeatCounter       uint16    `json:"fcu_heartbeat_counter"`       // Increments once per second
	ActivePowerControlMode    uint16    `json:"active_power_control_mode"`   // 0=Open, 1=Closed
	ReactivePowerControlMode  uint16    `json:"reactive_power_control_mode"` // 0=Q, 1=U, 2=CosPhi, 3=Q(dU)
	WindFarmRunning           bool      `json:"wind_farm_running"`           // Start/Stop status
	RapidDownwardSignalActive bool      `json:"rapid_downward_signal_active"`
}

// WindFarmSetpointData represents setpoint values (both commanded and current)
type WindFarmSetpointData struct {
	Timestamp time.Time `json:"timestamp"`
	ID        int       `json:"id"`
	// Commanded setpoints (mirrors)
	PSetpointMirror          float32 `json:"p_setpoint_mirror"`          // %, scale 0.01
	QSetpointMirror          float32 `json:"q_setpoint_mirror"`          // %, scale 0.01
	PowerFactorMirror        float32 `json:"power_factor_mirror"`        // scale 0.001
	USetpointMirror          float32 `json:"u_setpoint_mirror"`          // %, scale 0.01
	QdUSetpointMirror        float32 `json:"qdu_setpoint_mirror"`        // %, scale 0.01
	DPDtMinMirror            float32 `json:"dpdt_min_mirror"`            // p.u./min, scale 0.001
	DPDtMaxMirror            float32 `json:"dpdt_max_mirror"`            // p.u./min, scale 0.001
	FrequencyReserveCapacity uint16  `json:"frequency_reserve_capacity"` // %
	PfDeadbandMirror         float32 `json:"pf_deadband_mirror"`         // Hz, scale 0.001
	PfSlopeMirror            float32 `json:"pf_slope_mirror"`            // p.u./Hz, scale 0.001
	// Currently used setpoints
	PSetpointCurrent   float32 `json:"p_setpoint_current"`   // %, scale 0.01
	QSetpointCurrent   float32 `json:"q_setpoint_current"`   // %, scale 0.01
	PowerFactorCurrent float32 `json:"power_factor_current"` // scale 0.001
	USetpointCurrent   float32 `json:"u_setpoint_current"`   // %, scale 0.01
	QdUSetpointCurrent float32 `json:"qdu_setpoint_current"` // %, scale 0.01
}

// WindFarmWeatherData represents weather/meteo data
type WindFarmWeatherData struct {
	Timestamp                time.Time `json:"timestamp"`
	ID                       int       `json:"id"`
	WindSpeedMeteo           float32   `json:"wind_speed_meteo"`        // m/s, scale 0.1
	WindDirectionMeteo       float32   `json:"wind_direction_meteo"`    // degrees, scale 0.1
	OutsideTemperature       float32   `json:"outside_temperature"`     // °C, scale 0.1
	AtmosphericPressure      uint16    `json:"atmospheric_pressure"`    // mbar
	AirHumidity              float32   `json:"air_humidity"`            // %, scale 0.1
	RainfallVolume           float32   `json:"rainfall_volume"`         // l/m²h, scale 0.01
	SolarRadiation           float32   `json:"solar_radiation"`         // W/m², scale 0.1
	WindFarmCommunication    uint16    `json:"wind_farm_communication"` // %
	WeatherMeasurementsCount uint16    `json:"weather_measurements_count"`
}

// WindFarmCommandState represents the current command state for the wind farm
type WindFarmCommandState struct {
	LastUpdated              time.Time `json:"last_updated"`
	HeartbeatCounter         uint16    `json:"heartbeat_counter"`
	ActivePowerControlMode   uint16    `json:"active_power_control_mode"`
	ReactivePowerControlMode uint16    `json:"reactive_power_control_mode"`
	PSetpoint                float32   `json:"p_setpoint"`
	QSetpoint                float32   `json:"q_setpoint"`
	PowerFactorSetpoint      float32   `json:"power_factor_setpoint"`
	USetpoint                float32   `json:"u_setpoint"`
	WindFarmStartStop        uint16    `json:"wind_farm_start_stop"`
	RapidDownwardSignal      uint16    `json:"rapid_downward_signal"`
}

// AnalyzerData represents energy analyzer data
type AnalyzerData struct {
	Timestamp        time.Time `json:"timestamp"`
	VoltageL1        float32   `json:"voltage_l1"`
	VoltageL2        float32   `json:"voltage_l2"`
	VoltageL3        float32   `json:"voltage_l3"`
	VoltageLNAvg     float32   `json:"voltage_ln_avg"`
	VoltageL1L2      float32   `json:"voltage_l1l2"`
	VoltageL2L3      float32   `json:"voltage_l2l3"`
	VoltageL3L1      float32   `json:"voltage_l3l1"`
	VoltageLLAvg     float32   `json:"voltage_ll_avg"`
	CurrentL1        float32   `json:"current_l1"`
	CurrentL2        float32   `json:"current_l2"`
	CurrentL3        float32   `json:"current_l3"`
	CurrentN         float32   `json:"current_n"`
	ActivePowerL1    float32   `json:"active_power_l1"`
	ActivePowerL2    float32   `json:"active_power_l2"`
	ActivePowerL3    float32   `json:"active_power_l3"`
	ActivePowerSum   float32   `json:"active_power_sum"`
	ApparentPowerL1  float32   `json:"apparent_power_l1"`
	ApparentPowerL2  float32   `json:"apparent_power_l2"`
	ApparentPowerL3  float32   `json:"apparent_power_l3"`
	ApparentPowerSum float32   `json:"apparent_power_sum"`
	ReactivePowerL1  float32   `json:"reactive_power_l1"`
	ReactivePowerL2  float32   `json:"reactive_power_l2"`
	ReactivePowerL3  float32   `json:"reactive_power_l3"`
	ReactivePowerSum float32   `json:"reactive_power_sum"`
	PowerFactorL1    float32   `json:"power_factor_l1"`
	PowerFactorL2    float32   `json:"power_factor_l2"`
	PowerFactorL3    float32   `json:"power_factor_l3"`
	PowerFactorAvg   float32   `json:"power_factor_avg"`
	Frequency        float32   `json:"frequency"`
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

// TelemetryResponse represents the complete telemetry response
type TelemetryResponse struct {
	ParkName         string         `json:"park-name"`
	Timestamp        string         `json:"timestamp"`
	PowerplantStatus int            `json:"powerplant-status"`
	GenerationData   GenerationData `json:"generation-data"`
	BESSData         BESSData       `json:"bess-data"`
	POIData          POIData        `json:"poi-data"`
}

// GenerationData represents generation unit data
type GenerationData struct {
	TotalActivePowerMW                  float64 `json:"total-active-power-mw"`
	TotalReactivePowerMvar              float64 `json:"total-reactive-power-mvar"`
	AmbientTemperatureCelcius           float64 `json:"ambient-temperature-celcius"`
	CurrentMaximumActivePowerSetpointMW float64 `json:"current-maximum-active-power-setpoint-mw"`
}

// BESSData represents battery energy storage system data
type BESSData struct {
	TotalSOCMWh                    float64 `json:"total-soc-mwh"`
	TotalSOCPercentage             float64 `json:"total-soc-percentage"`
	TotalSOHPercentage             float64 `json:"total-soh-percentage"`
	TotalAvailableCapacityMWh      float64 `json:"total-available-capacity-mwh"`
	MaxAvailableChargingPowerMW    float64 `json:"max-available-charging-power-mw"`
	MaxAvailableDischargingPowerMW float64 `json:"max-available-discharging-power-mw"`
	CurrentActivePowerMW           float64 `json:"current-active-power-mw"`
	CurrentActivePowerSetpointMW   float64 `json:"current-active-power-setpoint-mw"`
}

// POIData represents point of injection data
type POIData struct {
	CurrentPOIActivePowerMW   float64 `json:"current-poi-active-power-mw"`
	CurrentPOIReactivePowerMW float64 `json:"current-poi-reactive-power-mw"`
}

// ScheduleDataPoint represents a single data point in the schedule
type ScheduleDataPoint struct {
	Timestamp               string  `json:"timestamp"`
	GenPCurtailmentSchedule float64 `json:"gen-p-curtailment-schedule"`
	GenPTradeSchedule       float64 `json:"gen-p-trade-schedule"`
	BessPTradeSchedule      float64 `json:"bess-p-trade-schedule"`
	PlantModeOfOperation    int     `json:"plant-mode-of-operation"`
}

// ScheduleRequest represents the schedule request payload
type ScheduleRequest struct {
	MsgID          string              `json:"msg-id"`
	ParkName       string              `json:"park-name"`
	MessageVersion string              `json:"message-version"`
	VersionDate    string              `json:"version-date"`
	SPSeconds      int                 `json:"sp-seconds"`
	Data           []ScheduleDataPoint `json:"data"`
}

// ScheduleResponse represents the schedule response
type ScheduleResponse struct {
	MsgID         string  `json:"msg-id"`
	Status        bool    `json:"status"`
	StatusMessage *string `json:"status-message"`
}
