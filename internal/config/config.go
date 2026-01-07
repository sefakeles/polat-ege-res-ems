package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config represents the complete application configuration
type Config struct {
	PCS          []PCSConfig        `mapstructure:"pcs" validate:"required,min=1,dive"`
	BMS          []BMSConfig        `mapstructure:"bms" validate:"required,min=1,dive"`
	PLC          []PLCConfig        `mapstructure:"plc" validate:"required,min=1,dive"`
	WindFarm     []WindFarmConfig   `mapstructure:"windfarm" validate:"required,min=1,dive"`
	EMS          EMSConfig          `mapstructure:"ems" validate:"required"`
	InfluxDB     InfluxDBConfig     `mapstructure:"influxdb" validate:"required"`
	PostgreSQL   PostgreSQLConfig   `mapstructure:"postgresql" validate:"required"`
	ModbusServer ModbusServerConfig `mapstructure:"modbus_server" validate:"required"`
	Logger       LoggerConfig       `mapstructure:"logger" validate:"required"`
}

// PCSConfig contains PCS-specific configuration
type PCSConfig struct {
	ID                int           `mapstructure:"id" validate:"required,min=1"`
	Host              string        `mapstructure:"host" validate:"required,hostname_rfc1123|ip"`
	Port              int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	SlaveID           byte          `mapstructure:"slave_id" validate:"required,min=1,max=255"`
	Timeout           time.Duration `mapstructure:"timeout" validate:"required"`
	ReconnectDelay    time.Duration `mapstructure:"reconnect_delay" validate:"required"`
	PollInterval      time.Duration `mapstructure:"poll_interval" validate:"required"`
	HeartbeatInterval time.Duration `mapstructure:"heartbeat_interval" validate:"required"`
	PersistInterval   time.Duration `mapstructure:"persist_interval" validate:"required"`
}

// BMSConfig contains BMS-specific configuration
type BMSConfig struct {
	ID                int           `mapstructure:"id" validate:"required,min=1"`
	Host              string        `mapstructure:"host" validate:"required,hostname_rfc1123|ip"`
	Port              int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	SlaveID           byte          `mapstructure:"slave_id" validate:"required,min=1,max=255"`
	Timeout           time.Duration `mapstructure:"timeout" validate:"required"`
	ReconnectDelay    time.Duration `mapstructure:"reconnect_delay" validate:"required"`
	PollInterval      time.Duration `mapstructure:"poll_interval" validate:"required"`
	CellDataInterval  time.Duration `mapstructure:"cell_data_interval" validate:"required"`
	HeartbeatInterval time.Duration `mapstructure:"heartbeat_interval" validate:"required"`
	PersistInterval   time.Duration `mapstructure:"persist_interval" validate:"required"`
	RackCount         int           `mapstructure:"rack_count" validate:"required,min=1,max=20"`
	ModulesPerRack    int           `mapstructure:"modules_per_rack" validate:"required,min=1,max=8"`
	EnableCellData    bool          `mapstructure:"enable_cell_data"`
}

// PLCConfig contains PLC-specific configuration
type PLCConfig struct {
	ID              int           `mapstructure:"id" validate:"required,min=1"`
	Host            string        `mapstructure:"host" validate:"required,hostname_rfc1123|ip"`
	Port            int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	SlaveID         byte          `mapstructure:"slave_id" validate:"required,min=1,max=255"`
	Timeout         time.Duration `mapstructure:"timeout" validate:"required"`
	ReconnectDelay  time.Duration `mapstructure:"reconnect_delay" validate:"required"`
	PollInterval    time.Duration `mapstructure:"poll_interval" validate:"required"`
	PersistInterval time.Duration `mapstructure:"persist_interval" validate:"required"`
}

// WindFarmConfig contains Wind Farm (ENERCON FCU) specific configuration
type WindFarmConfig struct {
	ID                int           `mapstructure:"id" validate:"required,min=1"`
	Host              string        `mapstructure:"host" validate:"required,hostname_rfc1123|ip"`
	Port              int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	SlaveID           byte          `mapstructure:"slave_id" validate:"required,min=1,max=255"`
	Timeout           time.Duration `mapstructure:"timeout" validate:"required"`
	ReconnectDelay    time.Duration `mapstructure:"reconnect_delay" validate:"required"`
	PollInterval      time.Duration `mapstructure:"poll_interval" validate:"required"`
	HeartbeatInterval time.Duration `mapstructure:"heartbeat_interval" validate:"required"`
	PersistInterval   time.Duration `mapstructure:"persist_interval" validate:"required"`
}

// EMSConfig contains EMS-specific configuration
type EMSConfig struct {
	HTTPPort          int     `mapstructure:"http_port" validate:"required,min=1,max=65535"`
	MaxSOC            float32 `mapstructure:"max_soc" validate:"required,min=0,max=100,gtfield=MinSOC"`
	MinSOC            float32 `mapstructure:"min_soc" validate:"required,min=0,max=100"`
	MaxChargePower    float32 `mapstructure:"max_charge_power" validate:"required,min=0"`
	MaxDischargePower float32 `mapstructure:"max_discharge_power" validate:"required,min=0"`
}

// InfluxDBConfig contains InfluxDB-specific configuration
type InfluxDBConfig struct {
	URL           string        `mapstructure:"url" validate:"required,url"`
	Token         string        `mapstructure:"token" validate:"required"`
	Organization  string        `mapstructure:"organization" validate:"required"`
	Bucket        string        `mapstructure:"bucket" validate:"required"`
	BatchSize     uint          `mapstructure:"batch_size" validate:"required,min=1"`
	FlushInterval time.Duration `mapstructure:"flush_interval" validate:"required"`
}

// PostgreSQLConfig contains PostgreSQL-specific configuration
type PostgreSQLConfig struct {
	Host     string `mapstructure:"host" validate:"required,hostname_rfc1123|ip"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Database string `mapstructure:"database" validate:"required"`
	SSLMode  string `mapstructure:"ssl_mode" validate:"required,oneof=disable allow prefer require verify-ca verify-full"`
	MaxIdle  int    `mapstructure:"max_idle_connections" validate:"required,min=1"`
	MaxOpen  int    `mapstructure:"max_open_connections" validate:"required,min=1"`
}

// ModbusServerConfig contains Modbus server configuration
type ModbusServerConfig struct {
	Host       string        `mapstructure:"host" validate:"required,hostname_rfc1123|ip"`
	Port       int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	Timeout    time.Duration `mapstructure:"timeout" validate:"required"`
	MaxClients uint          `mapstructure:"max_clients" validate:"required,min=1,max=100"`
}

// LoggerConfig contains logger-specific configuration
type LoggerConfig struct {
	Level  string `mapstructure:"level" validate:"required,oneof=DEBUG INFO WARN ERROR FATAL"`
	Format string `mapstructure:"format" validate:"required,oneof=json console"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Load loads configuration from the specified file path
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set configuration file path and name
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("json")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// Set default values
	setDefaults(v)

	// Enable environment variable support
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("EMS")

	// Explicitly bind all config keys for env variable support
	bindEnvVariables(v)

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal configuration
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// bindEnvVariables explicitly binds all configuration keys to environment variables
func bindEnvVariables(v *viper.Viper) {
	// EMS
	v.BindEnv("ems.http_port")
	v.BindEnv("ems.max_soc")
	v.BindEnv("ems.min_soc")
	v.BindEnv("ems.max_charge_power")
	v.BindEnv("ems.max_discharge_power")

	// InfluxDB
	v.BindEnv("influxdb.url")
	v.BindEnv("influxdb.token")
	v.BindEnv("influxdb.organization")
	v.BindEnv("influxdb.bucket")
	v.BindEnv("influxdb.batch_size")
	v.BindEnv("influxdb.flush_interval")

	// PostgreSQL
	v.BindEnv("postgresql.host")
	v.BindEnv("postgresql.port")
	v.BindEnv("postgresql.username")
	v.BindEnv("postgresql.password")
	v.BindEnv("postgresql.database")
	v.BindEnv("postgresql.ssl_mode")
	v.BindEnv("postgresql.max_idle_connections")
	v.BindEnv("postgresql.max_open_connections")

	// Modbus Server
	v.BindEnv("modbus_server.host")
	v.BindEnv("modbus_server.port")
	v.BindEnv("modbus_server.timeout")
	v.BindEnv("modbus_server.max_clients")

	// Logger
	v.BindEnv("logger.level")
	v.BindEnv("logger.format")
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// EMS defaults
	v.SetDefault("ems.http_port", 8080)
	v.SetDefault("ems.max_soc", 90.0)
	v.SetDefault("ems.min_soc", 10.0)
	v.SetDefault("ems.max_charge_power", 100.0)
	v.SetDefault("ems.max_discharge_power", 100.0)

	// InfluxDB defaults
	v.SetDefault("influxdb.batch_size", 100)
	v.SetDefault("influxdb.flush_interval", 5*time.Second)

	// PostgreSQL defaults
	v.SetDefault("postgresql.port", 5432)
	v.SetDefault("postgresql.ssl_mode", "disable")
	v.SetDefault("postgresql.max_idle_connections", 5)
	v.SetDefault("postgresql.max_open_connections", 10)

	// Modbus server defaults
	v.SetDefault("modbus_server.host", "0.0.0.0")
	v.SetDefault("modbus_server.port", 502)
	v.SetDefault("modbus_server.timeout", 30*time.Second)
	v.SetDefault("modbus_server.max_clients", 10)

	// Logger defaults
	v.SetDefault("logger.level", "INFO")
	v.SetDefault("logger.format", "json")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	return validate.Struct(c)
}
