package database

import (
	"fmt"
	"time"

	"powerkonnekt/ems/internal/config"
	"powerkonnekt/ems/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// PostgresDB represents the PostgreSQL connection for alarms
type PostgresDB struct {
	db  *gorm.DB
	log logger.Logger
}

// AlarmRecord represents the alarm table structure
type AlarmRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Timestamp time.Time `gorm:"index" json:"timestamp"`
	AlarmType string    `gorm:"index;size:50" json:"alarm_type"`
	Severity  string    `gorm:"index;size:20" json:"severity"`
	AlarmCode uint16    `json:"alarm_code"`
	Message   string    `gorm:"size:500" json:"message"`
	Active    bool      `gorm:"index" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for AlarmRecord
func (AlarmRecord) TableName() string {
	return "alarms"
}

// InitializePostgreSQL initializes the PostgreSQL connection for alarms
func InitializePostgreSQL(cfg config.PostgreSQLConfig) (*PostgresDB, error) {
	// Create database-specific logger
	dbLogger := logger.With(
		logger.String("database", "postgresql"),
		logger.String("host", cfg.Host),
		logger.Int("port", cfg.Port),
		logger.String("database", cfg.Database),
	)

	dbLogger.Info("Initializing PostgreSQL connection")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		dbLogger.Error("Failed to connect to PostgreSQL", logger.Err(err))
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		dbLogger.Error("Failed to get underlying sql.DB", logger.Err(err))
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		dbLogger.Error("Failed to ping PostgreSQL", logger.Err(err))
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	postgresDB := &PostgresDB{
		db:  db,
		log: dbLogger,
	}

	// Auto-migrate the schema
	if err := postgresDB.migrate(); err != nil {
		dbLogger.Error("Failed to migrate database", logger.Err(err))
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	dbLogger.Info("PostgreSQL connection established successfully",
		logger.Int("max_idle", cfg.MaxIdle),
		logger.Int("max_open", cfg.MaxOpen))
	return postgresDB, nil
}

// migrate creates or updates the database schema
func (p *PostgresDB) migrate() error {
	p.log.Info("Running database migration")

	err := p.db.AutoMigrate(&AlarmRecord{})
	if err != nil {
		p.log.Error("Database migration failed", logger.Err(err))
		return err
	}

	p.log.Info("Database migration completed successfully")
	return nil
}

// Close closes the PostgreSQL connection
func (p *PostgresDB) Close() error {
	p.log.Info("Closing PostgreSQL connection")

	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		p.log.Error("Failed to close PostgreSQL connection", logger.Err(err))
	} else {
		p.log.Info("PostgreSQL connection closed successfully")
	}

	return err
}

// SaveAlarm saves an alarm to PostgreSQL
func (p *PostgresDB) SaveAlarm(alarm BMSAlarmData) error {
	record := AlarmRecord{
		Timestamp: alarm.Timestamp,
		AlarmType: alarm.AlarmType,
		Severity:  alarm.Severity,
		AlarmCode: alarm.AlarmCode,
		Message:   alarm.Message,
		Active:    alarm.Active,
	}

	err := p.db.Create(&record).Error
	if err != nil {
		p.log.Error("Failed to save alarm",
			logger.Err(err),
			logger.String("alarm_type", alarm.AlarmType),
			logger.Uint16("alarm_code", alarm.AlarmCode))
		return err
	}

	return nil
}

// GetActiveAlarms retrieves all active alarms
func (p *PostgresDB) GetActiveAlarms() ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	err := p.db.Where("active = ?", true).
		Order("timestamp desc").
		Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get active alarms", logger.Err(err))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmHistory retrieves alarm history with pagination
func (p *PostgresDB) GetAlarmHistory(limit int, offset int) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	err := p.db.Order("timestamp desc").
		Limit(limit).
		Offset(offset).
		Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarm history",
			logger.Err(err),
			logger.Int("limit", limit),
			logger.Int("offset", offset))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmsByType retrieves alarms by type
func (p *PostgresDB) GetAlarmsByType(alarmType string, active bool) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	query := p.db.Where("alarm_type = ?", alarmType)
	if active {
		query = query.Where("active = ?", true)
	}
	err := query.Order("timestamp desc").Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarms by type",
			logger.Err(err),
			logger.String("alarm_type", alarmType),
			logger.Bool("active", active))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmsBySeverity retrieves alarms by severity
func (p *PostgresDB) GetAlarmsBySeverity(severity string, active bool) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	query := p.db.Where("severity = ?", severity)
	if active {
		query = query.Where("active = ?", true)
	}
	err := query.Order("timestamp desc").Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarms by severity",
			logger.Err(err),
			logger.String("severity", severity),
			logger.Bool("active", active))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmsInTimeRange retrieves alarms within a time range
func (p *PostgresDB) GetAlarmsInTimeRange(start, end time.Time) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	err := p.db.Where("timestamp BETWEEN ? AND ?", start, end).
		Order("timestamp desc").
		Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarms in time range",
			logger.Err(err),
			logger.Time("start", start),
			logger.Time("end", end))
		return nil, err
	}

	return alarms, nil
}

// UpdateAlarmStatus updates the active status of an alarm
func (p *PostgresDB) UpdateAlarmStatus(id uint, active bool) error {
	err := p.db.Model(&AlarmRecord{}).
		Where("id = ?", id).
		Update("active", active).Error
	if err != nil {
		p.log.Error("Failed to update alarm status",
			logger.Err(err),
			logger.Uint("id", id),
			logger.Bool("active", active))
		return err
	}

	return nil
}

// DeleteOldAlarms deletes alarms older than the specified duration
func (p *PostgresDB) DeleteOldAlarms(olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)

	result := p.db.Where("timestamp < ? AND active = ?", cutoffTime, false).
		Delete(&AlarmRecord{})

	if result.Error != nil {
		p.log.Error("Failed to delete old alarms",
			logger.Err(result.Error),
			logger.Duration("older_than", olderThan))
		return result.Error
	}

	p.log.Info("Old alarms deleted",
		logger.Int64("deleted_count", result.RowsAffected),
		logger.Duration("older_than", olderThan))

	return nil
}

// GetAlarmCount returns the count of alarms based on criteria
func (p *PostgresDB) GetAlarmCount(active *bool, severity string) (int64, error) {
	query := p.db.Model(&AlarmRecord{})

	if active != nil {
		query = query.Where("active = ?", *active)
	}

	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		logFields := []logger.Field{logger.Err(err)}
		if active != nil {
			logFields = append(logFields, logger.Bool("active", *active))
		}
		if severity != "" {
			logFields = append(logFields, logger.String("severity", severity))
		}
		p.log.Error("Failed to get alarm count", logFields...)
		return 0, err
	}

	logFields := []logger.Field{logger.Int64("count", count)}
	if active != nil {
		logFields = append(logFields, logger.Bool("active", *active))
	}
	if severity != "" {
		logFields = append(logFields, logger.String("severity", severity))
	}

	return count, nil
}

// HealthCheck checks if PostgreSQL is accessible
func (p *PostgresDB) HealthCheck() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		p.log.Error("Failed to get database connection for health check", logger.Err(err))
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		p.log.Error("PostgreSQL health check failed", logger.Err(err))
	}

	return err
}
