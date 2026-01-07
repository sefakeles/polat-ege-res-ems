package database

import (
	"context"
	"fmt"
	"time"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/pkg/logger"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// InfluxDB represents the InfluxDB connection
type InfluxDB struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	queryAPI api.QueryAPI
	config   config.InfluxDBConfig
	log      logger.Logger
}

// InitializeInfluxDB initializes the InfluxDB connection
func InitializeInfluxDB(cfg config.InfluxDBConfig) (*InfluxDB, error) {
	// Create database-specific logger
	dbLogger := logger.With(
		logger.String("database", "influxdb"),
		logger.String("url", cfg.URL),
		logger.String("organization", cfg.Organization),
		logger.String("bucket", cfg.Bucket),
	)

	dbLogger.Info("Initializing InfluxDB connection")

	// Create client with options
	options := influxdb2.DefaultOptions()
	options.SetBatchSize(cfg.BatchSize)
	options.SetFlushInterval(uint(cfg.FlushInterval.Milliseconds()))

	client := influxdb2.NewClientWithOptions(cfg.URL, cfg.Token, options)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	health, err := client.Health(ctx)
	if err != nil {
		dbLogger.Error("Failed to connect to InfluxDB", logger.Err(err))
		return nil, fmt.Errorf("failed to connect to InfluxDB: %w", err)
	}

	if health.Status != "pass" {
		dbLogger.Error("InfluxDB health check failed", logger.String("status", string(health.Status)))
		return nil, fmt.Errorf("InfluxDB health check failed: %s", health.Status)
	}

	writeAPI := client.WriteAPI(cfg.Organization, cfg.Bucket)
	queryAPI := client.QueryAPI(cfg.Organization)

	db := &InfluxDB{
		client:   client,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
		config:   cfg,
		log:      dbLogger,
	}

	dbLogger.Info("InfluxDB connection established successfully",
		logger.Uint("batch_size", cfg.BatchSize),
		logger.Duration("flush_interval", cfg.FlushInterval))
	return db, nil
}

// Close closes the InfluxDB connection
func (db *InfluxDB) Close() error {
	db.log.Info("Closing InfluxDB connection")

	if db.writeAPI != nil {
		db.writeAPI.Flush()
	}
	if db.client != nil {
		db.client.Close()
	}

	db.log.Info("InfluxDB connection closed")
	return nil
}

// HealthCheck checks if InfluxDB is accessible
func (db *InfluxDB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	health, err := db.client.Health(ctx)
	if err != nil {
		db.log.Error("InfluxDB health check failed", logger.Err(err))
		return err
	}

	if health.Status != "pass" {
		db.log.Error("InfluxDB health check status failed", logger.String("status", string(health.Status)))
		return fmt.Errorf("InfluxDB health check failed: %s", health.Status)
	}

	return nil
}

// WriteBMSData writes BMS data to InfluxDB
func (db *InfluxDB) WriteBMSData(data BMSData) error {
	point := influxdb2.NewPointWithMeasurement("bms").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("voltage", data.Voltage).
		AddField("current", data.Current).
		AddField("soc", data.SOC).
		AddField("soh", data.SOH).
		AddField("max_cell_voltage", data.MaxCellVoltage).
		AddField("min_cell_voltage", data.MinCellVoltage).
		AddField("avg_cell_voltage", data.AvgCellVoltage).
		AddField("max_cell_temperature", data.MaxCellTemperature).
		AddField("min_cell_temperature", data.MinCellTemperature).
		AddField("avg_cell_temperature", data.AvgCellTemperature).
		AddField("max_charge_current", data.MaxChargeCurrent).
		AddField("max_discharge_current", data.MaxDischargeCurrent).
		AddField("max_charge_power", data.MaxChargePower).
		AddField("max_discharge_power", data.MaxDischargePower).
		AddField("power", data.Power).
		AddField("charge_capacity", data.ChargeCapacity).
		AddField("discharge_capacity", data.DischargeCapacity).
		AddField("max_charge_voltage", data.MaxChargeVoltage).
		AddField("max_discharge_voltage", data.MaxDischargeVoltage).
		AddField("insulation_resistance_pos", data.InsulationResistancePos).
		AddField("insulation_resistance_neg", data.InsulationResistanceNeg).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteBMSStatusData writes BMS status data to InfluxDB
func (db *InfluxDB) WriteBMSStatusData(data BMSStatusData) error {
	point := influxdb2.NewPointWithMeasurement("bms").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("hv_status", data.HVStatus).
		AddField("system_status", data.SystemStatus).
		AddField("connected_racks", data.ConnectedRacks).
		AddField("total_racks", data.TotalRacks).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteBMSRackData writes BMS rack data to InfluxDB
func (db *InfluxDB) WriteBMSRackData(data BMSRackData) error {
	point := influxdb2.NewPointWithMeasurement("bms_rack").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddTag("number", fmt.Sprintf("%d", data.Number)).
		AddField("state", data.State).
		AddField("max_charge_power", data.MaxChargePower).
		AddField("max_discharge_power", data.MaxDischargePower).
		AddField("max_charge_voltage", data.MaxChargeVoltage).
		AddField("max_discharge_voltage", data.MaxDischargeVoltage).
		AddField("max_charge_current", data.MaxChargeCurrent).
		AddField("max_discharge_current", data.MaxDischargeCurrent).
		AddField("voltage", data.Voltage).
		AddField("current", data.Current).
		AddField("temperature", data.Temperature).
		AddField("soc", data.SOC).
		AddField("soh", data.SOH).
		AddField("insulation_resistance", data.InsulationResistance).
		AddField("avg_cell_voltage", data.AvgCellVoltage).
		AddField("avg_cell_temperature", data.AvgCellTemperature).
		AddField("max_cell_voltage", data.MaxCellVoltage).
		AddField("max_voltage_cell_no", data.MaxVoltageCellNo).
		AddField("min_cell_voltage", data.MinCellVoltage).
		AddField("min_voltage_cell_no", data.MinVoltageCellNo).
		AddField("max_cell_temperature", data.MaxCellTemperature).
		AddField("max_temp_cell_no", data.MaxTempCellNo).
		AddField("min_cell_temperature", data.MinCellTemperature).
		AddField("min_temp_cell_no", data.MinTempCellNo).
		AddField("total_charge_energy", data.TotalChargeEnergy).
		AddField("total_discharge_energy", data.TotalDischargeEnergy).
		AddField("charge_capacity", data.ChargeCapacity).
		AddField("discharge_capacity", data.DischargeCapacity).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteBMSCellVoltageData writes BMS cell voltage data to InfluxDB
func (db *InfluxDB) WriteBMSCellVoltageData(cells []BMSCellVoltageData) error {
	if len(cells) == 0 {
		return nil
	}

	for _, cell := range cells {
		point := influxdb2.NewPointWithMeasurement("bms_cell").
			AddTag("id", fmt.Sprintf("%d", cell.ID)).
			AddTag("rack_number", fmt.Sprintf("%d", cell.RackNo)).
			AddTag("module_number", fmt.Sprintf("%d", cell.ModuleNo)).
			AddTag("cell_number", fmt.Sprintf("%d", cell.CellNo)).
			AddField("voltage", cell.Voltage).
			SetTime(cell.Timestamp)
		db.writeAPI.WritePoint(point)
	}

	return nil
}

// WriteBMSCellTemperatureData writes BMS cell temperature data to InfluxDB
func (db *InfluxDB) WriteBMSCellTemperatureData(cells []BMSCellTemperatureData) error {
	if len(cells) == 0 {
		return nil
	}

	for _, cell := range cells {
		point := influxdb2.NewPointWithMeasurement("bms_cell").
			AddTag("id", fmt.Sprintf("%d", cell.ID)).
			AddTag("rack_number", fmt.Sprintf("%d", cell.RackNo)).
			AddTag("module_number", fmt.Sprintf("%d", cell.ModuleNo)).
			AddTag("sensor_number", fmt.Sprintf("%d", cell.SensorNo)).
			AddField("temperature", cell.Temperature).
			SetTime(cell.Timestamp)
		db.writeAPI.WritePoint(point)
	}

	return nil
}

// WritePCSStatusData writes PCS status data to InfluxDB
func (db *InfluxDB) WritePCSStatusData(data PCSStatusData) error {
	point := influxdb2.NewPointWithMeasurement("pcs").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("status", data.Status).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WritePCSEquipmentData writes PCS equipment data to InfluxDB
func (db *InfluxDB) WritePCSEquipmentData(data PCSEquipmentData) error {
	point := influxdb2.NewPointWithMeasurement("pcs").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("lv_switch_status", data.LVSwitchStatus).
		AddField("mv_switch_status", data.MVSwitchStatus).
		AddField("mv_disconnector_status", data.MVDisconnectorStatus).
		AddField("mv_earthing_switch_status", data.MVEarthingSwitchStatus).
		AddField("dc1_switch_status", data.DC1SwitchStatus).
		AddField("dc2_switch_status", data.DC2SwitchStatus).
		AddField("dc3_switch_status", data.DC3SwitchStatus).
		AddField("dc4_switch_status", data.DC4SwitchStatus).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WritePCSEnvironmentData writes PCS environment data to InfluxDB
func (db *InfluxDB) WritePCSEnvironmentData(data PCSEnvironmentData) error {
	point := influxdb2.NewPointWithMeasurement("pcs").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("air_inlet_temperature", data.AirInletTemperature).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WritePCSDCSourceData writes PCS DC source data to InfluxDB
func (db *InfluxDB) WritePCSDCSourceData(data PCSDCSourceData) error {
	point := influxdb2.NewPointWithMeasurement("pcs").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("dc1_power", data.DC1Power).
		AddField("dc2_power", data.DC2Power).
		AddField("dc3_power", data.DC3Power).
		AddField("dc4_power", data.DC4Power).
		AddField("dc1_current", data.DC1Current).
		AddField("dc2_current", data.DC2Current).
		AddField("dc3_current", data.DC3Current).
		AddField("dc4_current", data.DC4Current).
		AddField("dc1_voltage_external", data.DC1VoltageExternal).
		AddField("dc2_voltage_external", data.DC2VoltageExternal).
		AddField("dc3_voltage_external", data.DC3VoltageExternal).
		AddField("dc4_voltage_external", data.DC4VoltageExternal).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WritePCSGridData writes PCS grid data to InfluxDB
func (db *InfluxDB) WritePCSGridData(data PCSGridData) error {
	point := influxdb2.NewPointWithMeasurement("pcs").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("mv_grid_voltage_ab", data.MVGridVoltageAB).
		AddField("mv_grid_voltage_bc", data.MVGridVoltageBC).
		AddField("mv_grid_voltage_ca", data.MVGridVoltageCA).
		AddField("mv_grid_current_a", data.MVGridCurrentA).
		AddField("mv_grid_current_b", data.MVGridCurrentB).
		AddField("mv_grid_current_c", data.MVGridCurrentC).
		AddField("mv_grid_active_power", data.MVGridActivePower).
		AddField("mv_grid_reactive_power", data.MVGridReactivePower).
		AddField("mv_grid_apparent_power", data.MVGridApparentPower).
		AddField("mv_grid_cos_phi", data.MVGridCosPhi).
		AddField("lv_grid_voltage_ab", data.LVGridVoltageAB).
		AddField("lv_grid_voltage_bc", data.LVGridVoltageBC).
		AddField("lv_grid_voltage_ca", data.LVGridVoltageCA).
		AddField("lv_grid_current_a", data.LVGridCurrentA).
		AddField("lv_grid_current_b", data.LVGridCurrentB).
		AddField("lv_grid_current_c", data.LVGridCurrentC).
		AddField("lv_grid_active_power", data.LVGridActivePower).
		AddField("lv_grid_reactive_power", data.LVGridReactivePower).
		AddField("lv_grid_apparent_power", data.LVGridApparentPower).
		AddField("lv_grid_cos_phi", data.LVGridCosPhi).
		AddField("grid_frequency", data.GridFrequency).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WritePCSCounterData writes PCS counter data to InfluxDB
func (db *InfluxDB) WritePCSCounterData(data PCSCounterData) error {
	point := influxdb2.NewPointWithMeasurement("pcs").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("active_energy_today", data.ActiveEnergyToday).
		AddField("active_energy_yesterday", data.ActiveEnergyYesterday).
		AddField("active_energy_this_month", data.ActiveEnergyThisMonth).
		AddField("active_energy_last_month", data.ActiveEnergyLastMonth).
		AddField("active_energy_total", data.ActiveEnergyTotal).
		AddField("consumed_energy_today", data.ConsumedEnergyToday).
		AddField("consumed_energy_total", data.ConsumedEnergyTotal).
		AddField("reactive_energy_today", data.ReactiveEnergyToday).
		AddField("reactive_energy_yesterday", data.ReactiveEnergyYesterday).
		AddField("reactive_energy_this_month", data.ReactiveEnergyThisMonth).
		AddField("reactive_energy_last_month", data.ReactiveEnergyLastMonth).
		AddField("reactive_energy_total", data.ReactiveEnergyTotal).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WritePLCData writes PLC data to InfluxDB
func (db *InfluxDB) WritePLCData(data PLCData) error {
	point := influxdb2.NewPointWithMeasurement("plc").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		// Circuit Breakers
		AddField("auxiliary_cb", boolToInt(data.CircuitBreakers.AuxiliaryCB)).
		AddField("pcs1_cb", boolToInt(data.CircuitBreakers.PCS1CB)).
		AddField("pcs2_cb", boolToInt(data.CircuitBreakers.PCS2CB)).
		AddField("pcs3_cb", boolToInt(data.CircuitBreakers.PCS3CB)).
		AddField("pcs4_cb", boolToInt(data.CircuitBreakers.PCS4CB)).
		AddField("bms1_cb", boolToInt(data.CircuitBreakers.BMS1CB)).
		AddField("bms2_cb", boolToInt(data.CircuitBreakers.BMS2CB)).
		AddField("bms3_cb", boolToInt(data.CircuitBreakers.BMS3CB)).
		AddField("bms4_cb", boolToInt(data.CircuitBreakers.BMS4CB)).
		// MV Circuit Breakers
		AddField("mv_aux_transformer_cb", boolToInt(data.MVCircuitBreakers.AuxTransformerCB)).
		AddField("mv_transformer1_cb", boolToInt(data.MVCircuitBreakers.Transformer1CB)).
		AddField("mv_transformer2_cb", boolToInt(data.MVCircuitBreakers.Transformer2CB)).
		AddField("mv_transformer3_cb", boolToInt(data.MVCircuitBreakers.Transformer3CB)).
		AddField("mv_transformer4_cb", boolToInt(data.MVCircuitBreakers.Transformer4CB)).
		AddField("mv_autoproducer_cb", boolToInt(data.MVCircuitBreakers.AutoproducerCB)).
		// Protection Relays
		AddField("relay_aux_transformer_fault", boolToInt(data.ProtectionRelays.AuxTransformerFault)).
		AddField("relay_transformer1_fault", boolToInt(data.ProtectionRelays.Transformer1Fault)).
		AddField("relay_transformer2_fault", boolToInt(data.ProtectionRelays.Transformer2Fault)).
		AddField("relay_transformer3_fault", boolToInt(data.ProtectionRelays.Transformer3Fault)).
		AddField("relay_transformer4_fault", boolToInt(data.ProtectionRelays.Transformer4Fault)).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// boolToInt converts boolean to integer (1 for true, 0 for false)
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// WriteSystemMetrics writes system metrics to InfluxDB
func (db *InfluxDB) WriteSystemMetrics(data SystemMetrics) error {
	point := influxdb2.NewPointWithMeasurement("system_metrics").
		AddField("cpu_usage", data.CPUUsage).
		AddField("memory_usage", data.MemoryUsage).
		AddField("disk_usage", data.DiskUsage).
		AddField("network_rx", data.NetworkRx).
		AddField("network_tx", data.NetworkTx).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteRuntimeMetrics writes runtime metrics to InfluxDB
func (db *InfluxDB) WriteRuntimeMetrics(data RuntimeMetrics) error {
	point := influxdb2.NewPointWithMeasurement("runtime_metrics").
		AddField("uptime_seconds", data.UptimeSeconds).
		AddField("goroutines", data.Goroutines).
		AddField("heap_alloc_mb", data.HeapAllocMB).
		AddField("heap_sys_mb", data.HeapSysMB).
		AddField("heap_idle_mb", data.HeapIdleMB).
		AddField("heap_in_use_mb", data.HeapInUseMB).
		AddField("heap_released_mb", data.HeapReleasedMB).
		AddField("stack_in_use_mb", data.StackInUseMB).
		AddField("stack_sys_mb", data.StackSysMB).
		AddField("gc_runs", data.GCRuns).
		AddField("gc_pause_total_ns", data.GCPauseTotalNs).
		AddField("gc_cpu_fraction", data.GCCPUFraction).
		AddField("next_gc_mb", data.NextGCMB).
		AddField("last_gc_time", data.LastGCTime).
		AddField("mallocs_total", data.MallocsTotal).
		AddField("frees_total", data.FreesTotal).
		AddField("total_alloc_mb", data.TotalAllocMB).
		AddField("lookups_total", data.LookupsTotal).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteWindFarmMeasuringData writes wind farm measuring data to InfluxDB
func (db *InfluxDB) WriteWindFarmMeasuringData(data WindFarmMeasuringData) error {
	point := influxdb2.NewPointWithMeasurement("windfarm_measuring").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("active_power_ncp", data.ActivePowerNCP).
		AddField("reactive_power_ncp", data.ReactivePowerNCP).
		AddField("voltage_ncp", data.VoltageNCP).
		AddField("current_ncp", data.CurrentNCP).
		AddField("power_factor_ncp", data.PowerFactorNCP).
		AddField("frequency_ncp", data.FrequencyNCP).
		AddField("wec_availability", data.WECAvailability).
		AddField("wind_speed", data.WindSpeed).
		AddField("wind_direction", data.WindDirection).
		AddField("possible_wec_power", data.PossibleWECPower).
		AddField("wec_communication", data.WECCommunication).
		AddField("relative_power_availability", data.RelativePowerAvailability).
		AddField("absolute_power_availability", data.AbsolutePowerAvailability).
		AddField("relative_min_reactive_power", data.RelativeMinReactivePower).
		AddField("absolute_min_reactive_power", data.AbsoluteMinReactivePower).
		AddField("relative_max_reactive_power", data.RelativeMaxReactivePower).
		AddField("absolute_max_reactive_power", data.AbsoluteMaxReactivePower).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteWindFarmStatusData writes wind farm status data to InfluxDB
func (db *InfluxDB) WriteWindFarmStatusData(data WindFarmStatusData) error {
	point := influxdb2.NewPointWithMeasurement("windfarm_status").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("fcu_online", data.FCUOnline).
		AddField("fcu_mode", data.FCUMode).
		AddField("fcu_heartbeat_counter", data.FCUHeartbeatCounter).
		AddField("active_power_control_mode", data.ActivePowerControlMode).
		AddField("reactive_power_control_mode", data.ReactivePowerControlMode).
		AddField("wind_farm_running", data.WindFarmRunning).
		AddField("rapid_downward_signal_active", data.RapidDownwardSignalActive).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteWindFarmSetpointData writes wind farm setpoint data to InfluxDB
func (db *InfluxDB) WriteWindFarmSetpointData(data WindFarmSetpointData) error {
	point := influxdb2.NewPointWithMeasurement("windfarm_setpoint").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("p_setpoint_mirror", data.PSetpointMirror).
		AddField("q_setpoint_mirror", data.QSetpointMirror).
		AddField("power_factor_mirror", data.PowerFactorMirror).
		AddField("u_setpoint_mirror", data.USetpointMirror).
		AddField("qdu_setpoint_mirror", data.QdUSetpointMirror).
		AddField("dpdt_min_mirror", data.DPDtMinMirror).
		AddField("dpdt_max_mirror", data.DPDtMaxMirror).
		AddField("frequency_reserve_capacity", data.FrequencyReserveCapacity).
		AddField("pf_deadband_mirror", data.PfDeadbandMirror).
		AddField("pf_slope_mirror", data.PfSlopeMirror).
		AddField("p_setpoint_current", data.PSetpointCurrent).
		AddField("q_setpoint_current", data.QSetpointCurrent).
		AddField("power_factor_current", data.PowerFactorCurrent).
		AddField("u_setpoint_current", data.USetpointCurrent).
		AddField("qdu_setpoint_current", data.QdUSetpointCurrent).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// WriteWindFarmWeatherData writes wind farm weather data to InfluxDB
func (db *InfluxDB) WriteWindFarmWeatherData(data WindFarmWeatherData) error {
	point := influxdb2.NewPointWithMeasurement("windfarm_weather").
		AddTag("id", fmt.Sprintf("%d", data.ID)).
		AddField("wind_speed_meteo", data.WindSpeedMeteo).
		AddField("wind_direction_meteo", data.WindDirectionMeteo).
		AddField("outside_temperature", data.OutsideTemperature).
		AddField("atmospheric_pressure", data.AtmosphericPressure).
		AddField("air_humidity", data.AirHumidity).
		AddField("rainfall_volume", data.RainfallVolume).
		AddField("solar_radiation", data.SolarRadiation).
		AddField("wind_farm_communication", data.WindFarmCommunication).
		AddField("weather_measurements_count", data.WeatherMeasurementsCount).
		SetTime(data.Timestamp)

	db.writeAPI.WritePoint(point)

	return nil
}

// Flush forces writing of any buffered data
func (db *InfluxDB) Flush() {
	db.writeAPI.Flush()
}
