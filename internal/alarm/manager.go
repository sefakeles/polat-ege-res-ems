package alarm

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
)

// Manager handles alarm processing and management
type Manager struct {
	postgresDB   *database.PostgresDB
	activeAlarms map[string]database.BMSAlarmData
	mutex        sync.RWMutex
	log          *zap.Logger
}

// NewManager creates a new alarm manager
func NewManager(postgresDB *database.PostgresDB, logger *zap.Logger) *Manager {
	// Create manager-specific logger
	managerLogger := logger.With(
		zap.String("component", "alarm_manager"),
	)

	manager := &Manager{
		postgresDB:   postgresDB,
		activeAlarms: make(map[string]database.BMSAlarmData),
		log:          managerLogger,
	}

	// Load existing active alarms from PostgreSQL
	manager.loadActiveAlarms()

	return manager
}

// loadActiveAlarms loads active alarms from PostgreSQL into memory
func (m *Manager) loadActiveAlarms() {
	records, err := m.postgresDB.GetActiveAlarms()
	if err != nil {
		m.log.Error("Failed to load active alarms", zap.Error(err))
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, record := range records {
		alarm := database.BMSAlarmData{
			Timestamp: record.Timestamp,
			AlarmType: record.AlarmType,
			Severity:  record.Severity,
			AlarmCode: record.AlarmCode,
			Message:   record.Message,
			Active:    record.Active,
		}
		alarmKey := m.getAlarmKey(alarm)
		m.activeAlarms[alarmKey] = alarm
	}

	m.log.Info("Active alarms loaded from PostgreSQL",
		zap.Int("count", len(m.activeAlarms)))
}

// ProcessAlarm processes a new alarm
func (m *Manager) ProcessAlarm(alarm database.BMSAlarmData) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	alarmKey := m.getAlarmKey(alarm)

	logFields := []zap.Field{
		zap.String("alarm_type", alarm.AlarmType),
		zap.Uint16("alarm_code", alarm.AlarmCode),
		zap.String("severity", alarm.Severity),
		zap.String("message", alarm.Message),
		zap.Bool("active", alarm.Active),
	}

	if alarm.Active {
		if _, exists := m.activeAlarms[alarmKey]; !exists {
			// New alarm
			m.activeAlarms[alarmKey] = alarm

			// Save to PostgreSQL
			if err := m.postgresDB.SaveAlarm(alarm); err != nil {
				m.log.Error("Failed to save alarm to PostgreSQL",
					append(logFields, zap.Error(err))...)
			}

			if alarm.Severity == "HIGH" {
				m.log.Error("NEW CRITICAL ALARM", logFields...)
			} else {
				m.log.Warn("NEW ALARM", logFields...)
			}
		}
	} else {
		if existingAlarm, exists := m.activeAlarms[alarmKey]; exists {
			// Alarm cleared
			delete(m.activeAlarms, alarmKey)

			// Update the existing alarm in PostgreSQL to set active = false
			// First, get the active alarm record from PostgreSQL
			records, err := m.postgresDB.GetAlarmsByType(existingAlarm.AlarmType, true)
			if err != nil {
				m.log.Error("Failed to get active alarms for update",
					append(logFields, zap.Error(err))...)
			} else {
				// Find the matching record and update it
				for _, record := range records {
					if record.AlarmCode == existingAlarm.AlarmCode {
						if err := m.postgresDB.UpdateAlarmStatus(record.ID, false); err != nil {
							m.log.Error("Failed to update alarm status to inactive",
								append(logFields, zap.Error(err))...)
						}
						break
					}
				}
			}

			if existingAlarm.Severity == "HIGH" {
				m.log.Info("CRITICAL ALARM CLEARED", logFields...)
			} else {
				m.log.Info("ALARM CLEARED", logFields...)
			}
		}
	}
}

func (m *Manager) getAlarmKey(alarm database.BMSAlarmData) string {
	return fmt.Sprintf("%s_%d", alarm.AlarmType, alarm.AlarmCode)
}

// GetActiveAlarms returns all active alarms
func (m *Manager) GetActiveAlarms() []database.BMSAlarmData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	alarms := make([]database.BMSAlarmData, 0, len(m.activeAlarms))
	for _, alarm := range m.activeAlarms {
		alarms = append(alarms, alarm)
	}

	return alarms
}

// GetAlarmHistory returns alarm history from PostgreSQL
func (m *Manager) GetAlarmHistory(limit int, offset int) ([]database.AlarmRecord, error) {
	records, err := m.postgresDB.GetAlarmHistory(limit, offset)
	if err != nil {
		m.log.Error("Failed to get alarm history",
			zap.Error(err),
			zap.Int("limit", limit),
			zap.Int("offset", offset))
		return nil, err
	}

	return records, nil
}

// GetActiveAlarmsByType returns active alarms by type
func (m *Manager) GetActiveAlarmsByType(alarmType string) []database.BMSAlarmData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var alarms []database.BMSAlarmData
	for _, alarm := range m.activeAlarms {
		if alarm.AlarmType == alarmType {
			alarms = append(alarms, alarm)
		}
	}

	return alarms
}

// GetActiveAlarmsBySeverity returns active alarms by severity
func (m *Manager) GetActiveAlarmsBySeverity(severity string) []database.BMSAlarmData {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var alarms []database.BMSAlarmData
	for _, alarm := range m.activeAlarms {
		if alarm.Severity == severity {
			alarms = append(alarms, alarm)
		}
	}

	return alarms
}

// HasCriticalAlarms checks if there are critical alarms
func (m *Manager) HasCriticalAlarms() bool {
	criticalAlarms := m.GetActiveAlarmsBySeverity("HIGH")
	hasCritical := len(criticalAlarms) > 0

	if hasCritical {
		m.log.Warn("Critical alarms detected", zap.Int("count", len(criticalAlarms)))
	}

	return hasCritical
}

// GetAlarmsByTimeRange returns alarms in a specific time range
func (m *Manager) GetAlarmsByTimeRange(start, end time.Time) ([]database.AlarmRecord, error) {
	records, err := m.postgresDB.GetAlarmsInTimeRange(start, end)
	if err != nil {
		m.log.Error("Failed to get alarms by time range",
			zap.Error(err),
			zap.Time("start", start),
			zap.Time("end", end))
		return nil, err
	}

	return records, nil
}

// GetAlarmCount returns the count of alarms based on criteria
func (m *Manager) GetAlarmCount(active *bool, severity string) (int64, error) {
	logFields := []zap.Field{}
	if active != nil {
		logFields = append(logFields, zap.Bool("active", *active))
	}
	if severity != "" {
		logFields = append(logFields, zap.String("severity", severity))
	}

	count, err := m.postgresDB.GetAlarmCount(active, severity)
	if err != nil {
		m.log.Error("Failed to get alarm count",
			append(logFields, zap.Error(err))...)
		return 0, err
	}

	return count, nil
}

// CleanupOldAlarms removes old inactive alarms
func (m *Manager) CleanupOldAlarms(olderThan time.Duration) error {
	m.log.Info("Starting alarm cleanup",
		zap.Duration("older_than", olderThan))

	err := m.postgresDB.DeleteOldAlarms(olderThan)
	if err != nil {
		m.log.Error("Failed to cleanup old alarms",
			zap.Error(err),
			zap.Duration("older_than", olderThan))
		return err
	}

	m.log.Info("Alarm cleanup completed successfully",
		zap.Duration("older_than", olderThan))
	return nil
}

// RefreshActiveAlarms reloads active alarms from PostgreSQL
func (m *Manager) RefreshActiveAlarms() {
	m.log.Info("Refreshing active alarms from PostgreSQL")
	m.loadActiveAlarms()
}
