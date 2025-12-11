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

			alarmCode := relativeCode + 1

			message := GetAlarmMessage(alarmCode)
			severity := GetAlarmSeverity(alarmCode)

			alarm := database.BMSAlarmData{
				Timestamp: timestamp,
				AlarmType: fmt.Sprintf("PCS_%d", s.config.ID),
				AlarmCode: alarmCode,
				Message:   message,
				Severity:  severity,
				Active:    isActive,
			}

			s.alarmManager.ProcessAlarm(alarm)
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

			warningCode := relativeCode + 1
			message := GetWarningMessage(warningCode)
			severity := GetWarningSeverity(warningCode)

			warning := database.BMSAlarmData{
				Timestamp: timestamp,
				AlarmType: fmt.Sprintf("PCS_%d_WARNING", s.config.ID),
				AlarmCode: warningCode,
				Message:   message,
				Severity:  severity,
				Active:    isActive,
			}

			s.alarmManager.ProcessAlarm(warning)
		}
	}
}
