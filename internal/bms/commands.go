package bms

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// heartbeatLoop periodically updates heartbeat register in the BMS
func (s *Service) heartbeatLoop() {
	ticker := time.NewTicker(s.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if !s.baseClient.IsConnected() {
				continue
			}

			if err := s.updateHeartbeat(); err != nil {
				s.log.Error("Error updating heartbeat", zap.Error(err))
				continue
			}
		}
	}
}

// updateHeartbeat updates the heartbeat register in the BMS
func (s *Service) updateHeartbeat() error {
	s.mutex.Lock()
	if s.heartbeatCount == 16 {
		s.heartbeatCount = 0
	}
	heartbeatValue := s.heartbeatCount
	s.heartbeatCount++
	s.mutex.Unlock()

	err := s.baseClient.WriteSingleRegister(s.ctx, HeartbeatRegister, heartbeatValue)
	if err != nil {
		return fmt.Errorf("failed to write register: %w", err)
	}

	return nil
}

// ResetSystem sends a command to reset the BMS
func (s *Service) ResetSystem() error {
	return s.baseClient.WriteSingleRegister(s.ctx, SystemResetRegister, ControlReset)
}

// ControlMainBreaker sends a command to open or close the main breaker
func (s *Service) ControlMainBreaker(action uint16) error {
	var start bool
	var logAction string
	if action == ControlOn {
		logAction = "close"
		start = true
	} else {
		logAction = "open"
		start = false
	}

	s.mutex.Lock()
	s.commandState.StartStopCommand = start
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	err := s.baseClient.WriteSingleRegister(s.ctx, BreakerControlRegister, action)
	if err != nil {
		return fmt.Errorf("failed to %s circuit breaker: %w", logAction, err)
	}

	return nil
}
