package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"powerkonnekt/ems/internal/config"
)

// PostgreSQL represents the PostgreSQL connection for alarms
type PostgreSQL struct {
	db  *gorm.DB
	log *zap.Logger
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
func InitializePostgreSQL(cfg config.PostgreSQLConfig, logger *zap.Logger) (*PostgreSQL, error) {
	// Create database-specific logger
	dbLogger := logger.With(
		zap.String("database", "postgresql"),
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Database),
	)

	dbLogger.Info("Initializing PostgreSQL connection")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		dbLogger.Error("Failed to connect to PostgreSQL", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		dbLogger.Error("Failed to get underlying sql.DB", zap.Error(err))
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		dbLogger.Error("Failed to ping PostgreSQL", zap.Error(err))
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	postgreSQL := &PostgreSQL{
		db:  db,
		log: dbLogger,
	}

	// Auto-migrate the schema
	if err := postgreSQL.migrate(); err != nil {
		dbLogger.Error("Failed to migrate database", zap.Error(err))
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	dbLogger.Info("PostgreSQL connection established successfully",
		zap.Int("max_idle", cfg.MaxIdle),
		zap.Int("max_open", cfg.MaxOpen))
	return postgreSQL, nil
}

// migrate creates or updates the database schema
func (p *PostgreSQL) migrate() error {
	p.log.Info("Running database migration")

	err := p.db.AutoMigrate(&AlarmRecord{})
	if err != nil {
		p.log.Error("Database migration failed", zap.Error(err))
		return err
	}

	p.log.Info("Database migration completed successfully")
	return nil
}

// Close closes the PostgreSQL connection
func (p *PostgreSQL) Close() error {
	p.log.Info("Closing PostgreSQL connection")

	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		p.log.Error("Failed to close PostgreSQL connection", zap.Error(err))
	} else {
		p.log.Info("PostgreSQL connection closed successfully")
	}

	return err
}

// SaveAlarm saves an alarm to PostgreSQL
func (p *PostgreSQL) SaveAlarm(alarm BMSAlarmData) error {
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
			zap.Error(err),
			zap.String("alarm_type", alarm.AlarmType),
			zap.Uint16("alarm_code", alarm.AlarmCode))
		return err
	}

	return nil
}

// GetActiveAlarms retrieves all active alarms
func (p *PostgreSQL) GetActiveAlarms() ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	err := p.db.Where("active = ?", true).
		Order("timestamp desc").
		Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get active alarms", zap.Error(err))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmHistory retrieves alarm history with pagination
func (p *PostgreSQL) GetAlarmHistory(limit int, offset int) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	err := p.db.Order("timestamp desc").
		Limit(limit).
		Offset(offset).
		Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarm history",
			zap.Error(err),
			zap.Int("limit", limit),
			zap.Int("offset", offset))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmsByType retrieves alarms by type
func (p *PostgreSQL) GetAlarmsByType(alarmType string, active bool) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	query := p.db.Where("alarm_type = ?", alarmType)
	if active {
		query = query.Where("active = ?", true)
	}
	err := query.Order("timestamp desc").Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarms by type",
			zap.Error(err),
			zap.String("alarm_type", alarmType),
			zap.Bool("active", active))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmsBySeverity retrieves alarms by severity
func (p *PostgreSQL) GetAlarmsBySeverity(severity string, active bool) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	query := p.db.Where("severity = ?", severity)
	if active {
		query = query.Where("active = ?", true)
	}
	err := query.Order("timestamp desc").Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarms by severity",
			zap.Error(err),
			zap.String("severity", severity),
			zap.Bool("active", active))
		return nil, err
	}

	return alarms, nil
}

// GetAlarmsInTimeRange retrieves alarms within a time range
func (p *PostgreSQL) GetAlarmsInTimeRange(start, end time.Time) ([]AlarmRecord, error) {
	var alarms []AlarmRecord
	err := p.db.Where("timestamp BETWEEN ? AND ?", start, end).
		Order("timestamp desc").
		Find(&alarms).Error
	if err != nil {
		p.log.Error("Failed to get alarms in time range",
			zap.Error(err),
			zap.Time("start", start),
			zap.Time("end", end))
		return nil, err
	}

	return alarms, nil
}

// UpdateAlarmStatus updates the active status of an alarm
func (p *PostgreSQL) UpdateAlarmStatus(id uint, active bool) error {
	err := p.db.Model(&AlarmRecord{}).
		Where("id = ?", id).
		Update("active", active).Error
	if err != nil {
		p.log.Error("Failed to update alarm status",
			zap.Error(err),
			zap.Uint("id", id),
			zap.Bool("active", active))
		return err
	}

	return nil
}

// DeactivateAllAlarms deactivates all active alarms in a single query
func (p *PostgreSQL) DeactivateAllAlarms() (int64, error) {
	result := p.db.Model(&AlarmRecord{}).
		Where("active = ?", true).
		Update("active", false)

	if result.Error != nil {
		p.log.Error("Failed to deactivate all alarms", zap.Error(result.Error))
		return 0, result.Error
	}

	p.log.Info("Deactivated all active alarms", zap.Int64("count", result.RowsAffected))
	return result.RowsAffected, nil
}

// DeleteOldAlarms deletes alarms older than the specified duration
func (p *PostgreSQL) DeleteOldAlarms(olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)

	result := p.db.Where("timestamp < ? AND active = ?", cutoffTime, false).
		Delete(&AlarmRecord{})

	if result.Error != nil {
		p.log.Error("Failed to delete old alarms",
			zap.Error(result.Error),
			zap.Duration("older_than", olderThan))
		return result.Error
	}

	p.log.Info("Old alarms deleted",
		zap.Int64("deleted_count", result.RowsAffected),
		zap.Duration("older_than", olderThan))

	return nil
}

// GetAlarmCount returns the count of alarms based on criteria
func (p *PostgreSQL) GetAlarmCount(active *bool, severity string) (int64, error) {
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
		logFields := []zap.Field{zap.Error(err)}
		if active != nil {
			logFields = append(logFields, zap.Bool("active", *active))
		}
		if severity != "" {
			logFields = append(logFields, zap.String("severity", severity))
		}
		p.log.Error("Failed to get alarm count", logFields...)
		return 0, err
	}

	logFields := []zap.Field{zap.Int64("count", count)}
	if active != nil {
		logFields = append(logFields, zap.Bool("active", *active))
	}
	if severity != "" {
		logFields = append(logFields, zap.String("severity", severity))
	}

	return count, nil
}

// HealthCheck checks if PostgreSQL is accessible
func (p *PostgreSQL) HealthCheck() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		p.log.Error("Failed to get database connection for health check", zap.Error(err))
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		p.log.Error("PostgreSQL health check failed", zap.Error(err))
	}

	return err
}
