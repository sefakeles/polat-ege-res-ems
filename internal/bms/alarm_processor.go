package bms

import (
	"fmt"
	"time"

	"powerkonnekt/ems/internal/database"
)

// processAlarms processes alarm bits from the given data
func (s *Service) processAlarms(data []byte) {
	timestamp := time.Now()

	for byteIdx, b := range data {
		for bitIdx := range 8 {
			relativeCode := uint16(byteIdx*8 + bitIdx)
			isActive := (b & (1 << bitIdx)) != 0

			alarmCode := relativeCode + 1

			message := GetAlarmMessage(alarmCode)
			severity := GetAlarmSeverity(alarmCode)

			alarm := database.BMSAlarmData{
				Timestamp: timestamp,
				AlarmType: fmt.Sprintf("BMS_%d", s.config.ID),
				AlarmCode: alarmCode,
				Message:   message,
				Severity:  severity,
				Active:    isActive,
			}

			s.alarmManager.ProcessAlarm(alarm)
		}
	}
}

// processRackAlarms processes rack alarm bits from the given data
func (s *Service) processRackAlarms(data []byte, rackNo uint8) {
	baseCode := BMSRackAlarmStartAddr + uint16(rackNo-1)*BMSRackAlarmOffset
	timestamp := time.Now()

	for byteIdx, b := range data {
		for bitIdx := range 8 {
			relativeCode := uint16(byteIdx*8 + bitIdx)
			isActive := (b & (1 << bitIdx)) != 0

			alarmCode := baseCode + relativeCode
			message := GetRackAlarmMessage(relativeCode)
			severity := GetRackAlarmSeverity(relativeCode)

			alarm := database.BMSAlarmData{
				Timestamp: timestamp,
				AlarmType: fmt.Sprintf("BMS_%d_RACK", s.config.ID),
				AlarmCode: alarmCode,
				Message:   message,
				Severity:  severity,
				Active:    isActive,
			}

			s.alarmManager.ProcessAlarm(alarm)
		}
	}
}
