package bms

import (
	"fmt"
	"time"

	"powerkonnekt/ems/internal/database"
)

// processAlarms processes alarm bits from the given data
func (s *Service) processAlarms(data []byte) {
	timestamp := time.Now()

	// Reverse byte order for every word (2 bytes)
	swappedData := make([]byte, len(data))
	for i := 0; i < len(data); i += 2 {
		if i+1 < len(data) {
			swappedData[i] = data[i+1]
			swappedData[i+1] = data[i]
		} else {
			swappedData[i] = data[i]
		}
	}
	data = swappedData

	for byteIdx, b := range data {
		for bitIdx := range 8 {
			relativeCode := uint16(byteIdx*8 + bitIdx)
			isActive := (b & (1 << bitIdx)) != 0

			alarmType := fmt.Sprintf("BMS_%d", s.config.ID)
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

// processRackAlarms processes rack alarm bits from the given data
func (s *Service) processRackAlarms(data []byte, rackNo uint8) {
	baseCode := BMSRackAlarmStartAddr + uint16(rackNo-1)*BMSRackAlarmOffset
	timestamp := time.Now()

	for byteIdx, b := range data {
		for bitIdx := range 8 {
			relativeCode := uint16(byteIdx*8 + bitIdx)
			isActive := (b & (1 << bitIdx)) != 0

			alarmType := fmt.Sprintf("BMS_%d_RACK", s.config.ID)
			alarmCode := baseCode + relativeCode
			message := GetRackAlarmMessage(relativeCode)
			severity := GetRackAlarmSeverity(relativeCode)

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
