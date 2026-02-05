package pcs

import (
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

// heartbeatLoop periodically updates heartbeat register in the PCS
func (s *Service) heartbeatLoop() {
	ticker := time.NewTicker(s.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if s.client.IsConnected() {
				if err := s.updateHeartbeat(); err != nil {
					s.log.Error("Error updating heartbeat", zap.Error(err))
				}
			}
		}
	}
}

// updateHeartbeat updates the heartbeat register in the PCS
func (s *Service) updateHeartbeat() error {
	s.mutex.Lock()
	if s.heartbeatCount == math.MaxUint16 {
		s.heartbeatCount = 0
	}
	s.heartbeatCount++
	heartbeatValue := s.heartbeatCount
	s.mutex.Unlock()

	err := s.client.WriteSingleRegister(s.ctx, HeartbeatRegister, heartbeatValue)
	if err != nil {
		return fmt.Errorf("failed to write register: %w", err)
	}

	return nil
}

// ResetSystem sends a command to reset the PCS
func (s *Service) ResetSystem() error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PCS not connected")
	}

	err := s.client.WriteSingleRegister(s.ctx, SystemResetRegister, ControlReset)
	if err != nil {
		return fmt.Errorf("failed to write system reset command: %w", err)
	}

	s.log.Info("PCS reset command sent successfully")

	return nil
}

// StartStopCommand sends a command to start or stop the PCS
func (s *Service) StartStopCommand(start bool) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PCS not connected")
	}

	var value uint16
	var action string
	if start {
		value = 1
		action = "start"
	} else {
		value = 0
		action = "stop"
	}

	err := s.client.WriteSingleRegister(s.ctx, CmdStartStopRegister, value)
	if err != nil {
		return fmt.Errorf("failed to %s PCS: %w", action, err)
	}

	s.mutex.Lock()
	s.commandState.StartStopCommand = start
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("PCS command sent successfully",
		zap.String("action", action),
		zap.Bool("start", start))

	return nil
}

// SetActivePowerCommand sets the active power (kW)
func (s *Service) SetActivePowerCommand(power float32) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PCS not connected")
	}

	// Validate power command range
	const maxPower = 1500.0 // kW
	if power > maxPower || power < -maxPower {
		return fmt.Errorf("active power command out of range: %.1f kW (max: ±%.1f kW)", power, maxPower)
	}

	// Use the standard kW command register
	powerValue := int16(power * 100) // Power in kW
	if err := s.client.WriteSingleRegister(s.ctx, CmdActivePowerRegister, uint16(powerValue)); err != nil {
		return fmt.Errorf("failed to write active power command: %w", err)
	}

	s.mutex.Lock()
	s.commandState.ActivePowerCommand = power
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("PCS active power command set", zap.Float32("power", power))
	return nil
}

// SetReactivePowerCommand sets the reactive power (kVAr)
func (s *Service) SetReactivePowerCommand(power float32) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PCS not connected")
	}

	// Validate power command range
	const maxPower = 1500.0 // kW
	if power > maxPower || power < -maxPower {
		return fmt.Errorf("reactive power command out of range: %.1f kVAr (max: ±%.1f kVAr)", power, maxPower)
	}

	// Use the standard kW command register
	powerValue := int16(power * 100) // Power in kW
	if err := s.client.WriteSingleRegister(s.ctx, CmdReactivePowerRegister, uint16(powerValue)); err != nil {
		return fmt.Errorf("failed to write reactive power command: %w", err)
	}

	s.mutex.Lock()
	s.commandState.ReactivePowerCommand = power
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("PCS reactive power command set", zap.Float32("power", power))
	return nil
}
