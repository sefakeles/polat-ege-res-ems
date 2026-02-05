package plc

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
)

// pollLoop periodically reads data from the PLC
func (s *Service) pollLoop() {
	if err := s.client.Connect(s.ctx); err != nil {
		s.log.Warn("Initial Modbus connection failed", zap.Error(err))
	}

	interval := s.config.PollInterval

	// Calculate first aligned time and create timer
	nextTick := time.Now().Truncate(interval).Add(interval)
	timer := time.NewTimer(time.Until(nextTick))
	defer timer.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			if !s.client.IsConnected() {
				s.handleConnectionError()
			} else {
				startTime := time.Now()
				if err := s.readPLCData(); err != nil {
					s.log.Error("Error reading data", zap.Error(err))
				} else {
					// Signal that new data is available
					select {
					case s.dataUpdateChan <- struct{}{}:
					default:
						// Channel full, skip signal
					}
				}

				if duration := time.Since(startTime); duration > interval {
					s.log.Warn("Data read exceeded poll interval",
						zap.Duration("duration", duration),
						zap.Duration("interval", interval))
				}
			}

			// Calculate next aligned time and reset timer
			nextTick = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextTick))
		}
	}
}

// handleConnectionError attempts to reconnect to the PLC
func (s *Service) handleConnectionError() {
	s.log.Warn("PLC connection lost, initiating reconnection procedure")
	s.client.Disconnect()

	reconnectAttempts := 0
	timer := time.NewTimer(s.config.ReconnectDelay)
	defer timer.Stop()

	for !s.client.IsConnected() {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			reconnectAttempts++
			if err := s.client.Connect(s.ctx); err != nil {
				s.log.Error("Failed to reconnect to PLC",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts))
				timer.Reset(s.config.ReconnectDelay)
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
