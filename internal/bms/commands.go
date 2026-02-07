package bms

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"powerkonnekt/ems/pkg/utils"
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
			if s.systemClient.IsConnected() {
				if err := s.updateHeartbeat(); err != nil {
					s.log.Error("Error updating heartbeat", zap.Error(err))
				}
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

	err := s.systemClient.WriteSingleRegister(s.ctx, HeartbeatRegister, heartbeatValue)
	if err != nil {
		return fmt.Errorf("failed to write register: %w", err)
	}

	return nil
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

	err := s.systemClient.WriteSingleRegister(s.ctx, BreakerControlRegister, action)
	if err != nil {
		return fmt.Errorf("failed to %s circuit breaker: %w", logAction, err)
	}

	return nil
}

// ResetSystem sends a fault clear command to the BMS
func (s *Service) ResetSystem() error {
	return s.systemClient.WriteSingleRegister(s.ctx, FaultClearRegister, ControlReset)
}

// ControlInsulationDetection sends a command to turn on or off BMS insulation detection
func (s *Service) ControlInsulationDetection(action uint16) error {
	if action != InsulationControlOn && action != InsulationControlOff {
		return fmt.Errorf("invalid insulation control action: %d", action)
	}

	err := s.systemClient.WriteSingleRegister(s.ctx, InsulationControlRegister, action)
	if err != nil {
		return fmt.Errorf("failed to control insulation detection: %w", err)
	}

	s.log.Info("Insulation detection control executed",
		zap.Uint16("action", action))

	return nil
}

// ControlRackDisable sends a command to enable or disable a specific rack (1-48)
func (s *Service) ControlRackDisable(rackNo uint8, disable bool) error {
	if rackNo < 1 || rackNo > 48 {
		return fmt.Errorf("invalid rack number: %d (must be 1-48)", rackNo)
	}

	// Determine which register and bit position
	var register uint16
	var bitPos uint16
	switch {
	case rackNo <= 16:
		register = RackDisableRegister1
		bitPos = uint16(rackNo - 1)
	case rackNo <= 32:
		register = RackDisableRegister2
		bitPos = uint16(rackNo - 17)
	default:
		register = RackDisableRegister3
		bitPos = uint16(rackNo - 33)
	}

	// Read current register value
	data, err := s.systemClient.ReadHoldingRegisters(s.ctx, register, 1)
	if err != nil {
		return fmt.Errorf("failed to read rack disable register: %w", err)
	}

	currentValue := utils.FromBytes[uint16](data)

	// Set or clear the bit
	if disable {
		currentValue |= 1 << bitPos
	} else {
		currentValue &^= 1 << bitPos
	}

	err = s.systemClient.WriteSingleRegister(s.ctx, register, currentValue)
	if err != nil {
		return fmt.Errorf("failed to control rack %d: %w", rackNo, err)
	}

	s.log.Info("Rack disable control executed",
		zap.Uint8("rack_no", rackNo),
		zap.Bool("disable", disable))

	return nil
}

// ControlStepCharge sends a command to control step-charge mode
func (s *Service) ControlStepCharge(action uint16) error {
	if action > StepChargeControlEnable {
		return fmt.Errorf("invalid step-charge action: %d", action)
	}

	err := s.systemClient.WriteSingleRegister(s.ctx, StepChargeControlRegister, action)
	if err != nil {
		return fmt.Errorf("failed to control step-charge: %w", err)
	}

	s.log.Info("Step-charge control executed",
		zap.Uint16("action", action))

	return nil
}
