package plc

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
)

// pollLoop periodically reads data from the PLC
func (s *Service) pollLoop() {
	ticker := time.NewTicker(s.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if !s.client.IsConnected() {
				s.handleConnectionError()
			}

			if err := s.readPLCData(); err != nil {
				s.log.Error("Error reading PLC data", zap.Error(err))
				continue
			}

			// Signal that new data is available
			select {
			case s.dataUpdateChan <- struct{}{}:
			default:
				// Channel full, skip signal
			}
		}
	}
}

// handleConnectionError attempts to reconnect to the PLC
func (s *Service) handleConnectionError() {
	s.log.Warn("PLC connection lost, initiating reconnection procedure")
	s.client.Disconnect()

	reconnectAttempts := 0
	for !s.client.IsConnected() {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(s.config.ReconnectDelay):
			reconnectAttempts++
			if err := s.client.Connect(s.ctx); err != nil {
				s.log.Error("Failed to reconnect to PLC",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts),
					zap.Duration("retry_delay", s.config.ReconnectDelay))
			} else {
				s.log.Info("Successfully reconnected to PLC",
					zap.Int("total_attempts", reconnectAttempts),
					zap.Duration("total_downtime", time.Duration(reconnectAttempts)*s.config.ReconnectDelay))
				return
			}
		}
	}
}

// readPLCData reads status data from the PLC
func (s *Service) readPLCData() error {
	// Read circuit breaker positions, MV circuit breakers, and protection relays
	// These are consecutive registers starting at address 7
	data, err := s.client.ReadHoldingRegisters(s.ctx, CircuitBreakerPositionsAddr, StatusDataLength)
	if err != nil {
		return fmt.Errorf("failed to read PLC registers: %w", err)
	}

	plcData := ParsePLCData(data, s.config.ID)

	s.mutex.Lock()
	s.lastPLCData = plcData
	s.mutex.Unlock()

	// Check for protection relay faults and create alarms
	s.checkProtectionRelayFaults(plcData)

	return nil
}

// checkProtectionRelayFaults checks for protection relay faults and creates alarms
func (s *Service) checkProtectionRelayFaults(data database.PLCData) {
	timestamp := time.Now()

	relayFaults := map[string]bool{
		"MV Aux Transformer Relay": data.ProtectionRelays.AuxTransformerFault,
		"Transformer 1 Relay":      data.ProtectionRelays.Transformer1Fault,
		"Transformer 2 Relay":      data.ProtectionRelays.Transformer2Fault,
		"Transformer 3 Relay":      data.ProtectionRelays.Transformer3Fault,
		"Transformer 4 Relay":      data.ProtectionRelays.Transformer4Fault,
	}

	alarmCode := uint16(1)
	for relayName, hasFault := range relayFaults {
		// Check if state has changed
		previousState, exists := s.previousRelayStates[relayName]
		stateChanged := !exists || previousState != hasFault

		// Only process if state changed
		if stateChanged {
			alarm := database.BMSAlarmData{
				Timestamp: timestamp,
				AlarmType: fmt.Sprintf("PLC_%d_RELAY", s.config.ID),
				AlarmCode: alarmCode,
				Message:   fmt.Sprintf("%s Fault", relayName),
				Severity:  "HIGH",
				Active:    hasFault,
			}

			if s.alarmManager != nil {
				s.alarmManager.SubmitAlarm(alarm)
			}

			// Update previous state
			s.previousRelayStates[relayName] = hasFault
		}

		alarmCode++
	}
}
