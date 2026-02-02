package pcs

import (
	"encoding/binary"
	"fmt"
	"time"

	"powerkonnekt/ems/internal/database"
)

// processFaults processes fault bits from the given data
func (s *Service) processFaults(data []byte) {
	timestamp := time.Now()

	// Process data as 16-bit unsigned integers
	for i := 0; i < len(data); i += 2 {
		if i+1 >= len(data) {
			break // Skip incomplete uint16
		}

		value := binary.BigEndian.Uint16(data[i : i+2])
		wordIdx := i / 2

		for bitIdx := range 16 {
			relativeCode := uint16(wordIdx*16 + bitIdx)
			isActive := (value & (1 << bitIdx)) != 0

			alarmType := fmt.Sprintf("PCS_%d", s.config.ID)
			alarmCode := relativeCode + 1
			message := GetAlarmMessage(alarmCode)
			severity := GetAlarmSeverity(alarmCode)

			if message == "Unknown alarm" {
				continue
			}

			// Create unique alarm key
			alarmKey := fmt.Sprintf("%s_%d", alarmType, alarmCode)

			// Check if alarm state has changed
			s.mutex.Lock()
			previousState, exists := s.previousAlarmStates[alarmKey]
			stateChanged := (!exists && isActive) || (exists && previousState != isActive)
			s.previousAlarmStates[alarmKey] = isActive
			s.mutex.Unlock()

			if stateChanged {
				alarm := database.BMSAlarmData{
					Timestamp: timestamp,
					AlarmType: alarmType,
					AlarmCode: alarmCode,
					Message:   message,
					Severity:  severity,
					Active:    isActive,
				}

				s.alarmManager.SubmitAlarm(alarm)
			}
		}
	}
}

// processWarnings processes warning bits from the given data
func (s *Service) processWarnings(data []byte) {
	timestamp := time.Now()

	// Process data as 16-bit unsigned integers
	for i := 0; i < len(data); i += 2 {
		if i+1 >= len(data) {
			break // Skip incomplete uint16
		}

		value := binary.BigEndian.Uint16(data[i : i+2])
		wordIdx := i / 2

		for bitIdx := range 16 {
			relativeCode := uint16(wordIdx*16 + bitIdx)
			isActive := (value & (1 << bitIdx)) != 0

			alarmType := fmt.Sprintf("PCS_%d_WARNING", s.config.ID)
			alarmCode := relativeCode + 1
			message := GetWarningMessage(alarmCode)
			severity := GetWarningSeverity(alarmCode)

			if message == "Unknown warning" {
				continue
			}

			// Create unique alarm key
			alarmKey := fmt.Sprintf("%s_%d", alarmType, alarmCode)

			// Check if alarm state has changed
			s.mutex.Lock()
			previousState, exists := s.previousAlarmStates[alarmKey]
			stateChanged := (!exists && isActive) || (exists && previousState != isActive)
			s.previousAlarmStates[alarmKey] = isActive
			s.mutex.Unlock()

			if stateChanged {
				warning := database.BMSAlarmData{
					Timestamp: timestamp,
					AlarmType: alarmType,
					AlarmCode: alarmCode,
					Message:   message,
					Severity:  severity,
					Active:    isActive,
				}

				s.alarmManager.SubmitAlarm(warning)
			}
		}
	}
}
