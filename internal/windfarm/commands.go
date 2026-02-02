package windfarm

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// heartbeatLoop sends heartbeat updates
func (s *Service) heartbeatLoop() {
	ticker := time.NewTicker(s.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if !s.client.IsConnected() {
				continue
			}

			if err := s.sendHeartbeat(); err != nil {
				s.log.Error("Error sending heartbeat", zap.Error(err))
			}
		}
	}
}

// sendHeartbeat sends the heartbeat to maintain communication
func (s *Service) sendHeartbeat() error {
	s.mutex.Lock()
	s.heartbeatCounter++
	counter := s.heartbeatCounter
	s.mutex.Unlock()

	if err := s.client.WriteSingleRegister(s.ctx, HeartbeatAddr, counter); err != nil {
		return fmt.Errorf("failed to write heartbeat: %w", err)
	}

	s.log.Debug("Heartbeat sent", zap.Uint16("counter", counter))
	return nil
}

// SetPowerSetpoint sets the active power setpoint (0-100%)
func (s *Service) SetPowerSetpoint(setpoint float32) error {
	if setpoint < 0 || setpoint > 100 {
		return fmt.Errorf("power setpoint must be between 0 and 100, got %f", setpoint)
	}

	// Scale: 0.01, so multiply by 100 for register value
	value := uint16(setpoint * 100)

	if err := s.client.WriteSingleRegister(s.ctx, PSetpointAddr, value); err != nil {
		return fmt.Errorf("failed to write P setpoint: %w", err)
	}

	s.mutex.Lock()
	s.commandState.PSetpoint = setpoint
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("Power setpoint set", zap.Float32("setpoint", setpoint))
	return nil
}

// SetReactivePowerSetpoint sets the reactive power setpoint (-100% to 100%)
func (s *Service) SetReactivePowerSetpoint(setpoint float32) error {
	if setpoint < -100 || setpoint > 100 {
		return fmt.Errorf("reactive power setpoint must be between -100 and 100, got %f", setpoint)
	}

	// Scale: 0.01, so multiply by 100 for register value (signed)
	value := uint16(int16(setpoint * 100))

	if err := s.client.WriteSingleRegister(s.ctx, QSetpointAddr, value); err != nil {
		return fmt.Errorf("failed to write Q setpoint: %w", err)
	}

	s.mutex.Lock()
	s.commandState.QSetpoint = setpoint
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("Reactive power setpoint set", zap.Float32("setpoint", setpoint))
	return nil
}

// SetPowerFactorSetpoint sets the power factor setpoint
func (s *Service) SetPowerFactorSetpoint(setpoint float32) error {
	if setpoint < -1 || setpoint > 1 {
		return fmt.Errorf("power factor setpoint must be between -1 and 1, got %f", setpoint)
	}

	// Scale: 0.001, so multiply by 1000 for register value
	value := uint16(int16(setpoint * 1000))

	if err := s.client.WriteSingleRegister(s.ctx, PowerFactorSetpointAddr, value); err != nil {
		return fmt.Errorf("failed to write power factor setpoint: %w", err)
	}

	s.mutex.Lock()
	s.commandState.PowerFactorSetpoint = setpoint
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("Power factor setpoint set", zap.Float32("setpoint", setpoint))
	return nil
}

// StartWindFarm sends the start command to the wind farm
func (s *Service) StartWindFarm() error {
	if err := s.client.WriteSingleRegister(s.ctx, WindFarmStartStopAddr, WindFarmStart); err != nil {
		return fmt.Errorf("failed to send start command: %w", err)
	}

	s.mutex.Lock()
	s.commandState.WindFarmStartStop = WindFarmStart
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("Wind farm start command sent")
	return nil
}

// StopWindFarm sends the stop command to the wind farm
func (s *Service) StopWindFarm() error {
	if err := s.client.WriteSingleRegister(s.ctx, WindFarmStartStopAddr, WindFarmStop); err != nil {
		return fmt.Errorf("failed to send stop command: %w", err)
	}

	s.mutex.Lock()
	s.commandState.WindFarmStartStop = WindFarmStop
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("Wind farm stop command sent")
	return nil
}

// SetRapidDownwardSignal sets the rapid downward signal on or off
func (s *Service) SetRapidDownwardSignal(on bool) error {
	value := uint16(RapidDownwardOff)
	if on {
		value = RapidDownwardOn
	}

	if err := s.client.WriteSingleRegister(s.ctx, RapidDownwardSignalAddr, value); err != nil {
		return fmt.Errorf("failed to write rapid downward signal: %w", err)
	}

	s.mutex.Lock()
	s.commandState.RapidDownwardSignal = value
	s.commandState.LastUpdated = time.Now()
	s.mutex.Unlock()

	s.log.Info("Rapid downward signal set", zap.Bool("on", on))
	return nil
}
