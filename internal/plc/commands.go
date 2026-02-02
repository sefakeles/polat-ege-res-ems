package plc

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
)

// ControlAuxiliaryCB controls the auxiliary circuit breaker
func (s *Service) ControlAuxiliaryCB(close bool) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PLC not connected")
	}

	var command uint16
	var action string
	if close {
		command = ControlClose
		action = "close"
	} else {
		command = ControlOpen
		action = "open"
	}

	err := s.client.WriteSingleRegister(s.ctx, AuxCBControlAddr, command)
	if err != nil {
		return fmt.Errorf("failed to %s auxiliary CB: %w", action, err)
	}

	s.log.Info("Auxiliary CB command sent successfully",
		zap.String("action", action),
		zap.Bool("close", close))

	return nil
}

// ControlMVAuxTransformerCB controls the MV auxiliary transformer circuit breaker
func (s *Service) ControlMVAuxTransformerCB(close bool) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PLC not connected")
	}

	var command uint16
	var action string
	if close {
		command = ControlClose
		action = "close"
	} else {
		command = ControlOpen
		action = "open"
	}

	err := s.client.WriteSingleRegister(s.ctx, MVAuxTransformerCBAddr, command)
	if err != nil {
		return fmt.Errorf("failed to %s MV aux transformer CB: %w", action, err)
	}

	s.log.Info("MV Aux Transformer CB command sent successfully",
		zap.String("action", action),
		zap.Bool("close", close))

	return nil
}

// ControlTransformerCB controls a transformer circuit breaker (1-4)
func (s *Service) ControlTransformerCB(transformerNo uint8, close bool) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PLC not connected")
	}

	if transformerNo < 1 || transformerNo > 4 {
		return fmt.Errorf("invalid transformer number: %d (must be 1-4)", transformerNo)
	}

	var command uint16
	var action string
	if close {
		command = ControlClose
		action = "close"
	} else {
		command = ControlOpen
		action = "open"
	}

	// Calculate register address based on transformer number
	// Transformer 1 = address 12, Transformer 2 = 13, etc.
	registerAddr := Transformer1CBControlAddr + uint16(transformerNo-1)

	err := s.client.WriteSingleRegister(s.ctx, registerAddr, command)
	if err != nil {
		return fmt.Errorf("failed to %s transformer %d CB: %w", action, transformerNo, err)
	}

	s.log.Info("Transformer CB command sent successfully",
		zap.Uint8("transformer_no", transformerNo),
		zap.String("action", action),
		zap.Bool("close", close))

	return nil
}

func (s *Service) ControlAutoproducerCB(close bool) error {
	if !s.client.IsConnected() {
		return fmt.Errorf("PLC not connected")
	}

	var command uint16
	var action string
	if close {
		command = ControlClose
		action = "close"
	} else {
		command = ControlOpen
		action = "open"
	}

	err := s.client.WriteSingleRegister(s.ctx, AutoproducerCBControlAddr, command)
	if err != nil {
		return fmt.Errorf("failed to %s autoproducer CB: %w", action, err)
	}

	s.log.Info("Autoproducer CB command sent successfully",
		zap.String("action", action),
		zap.Bool("close", close))

	return nil
}

// GetCircuitBreakerStatus returns the current status of all circuit breakers
func (s *Service) GetCircuitBreakerStatus() database.CircuitBreakerStatus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastPLCData.CircuitBreakers
}

// GetMVCircuitBreakerStatus returns the current status of MV circuit breakers
func (s *Service) GetMVCircuitBreakerStatus() database.MVCircuitBreakerStatus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastPLCData.MVCircuitBreakers
}

// GetProtectionRelayStatus returns the current status of protection relays
func (s *Service) GetProtectionRelayStatus() database.ProtectionRelayStatus {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.lastPLCData.ProtectionRelays
}

// HasProtectionRelayFaults checks if any protection relay has faults
func (s *Service) HasProtectionRelayFaults() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	relays := s.lastPLCData.ProtectionRelays
	return relays.AuxTransformerFault ||
		relays.Transformer1Fault ||
		relays.Transformer2Fault ||
		relays.Transformer3Fault ||
		relays.Transformer4Fault
}

// GetFaultedRelays returns a list of faulted relay names
func (s *Service) GetFaultedRelays() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var faulted []string
	relays := s.lastPLCData.ProtectionRelays

	if relays.AuxTransformerFault {
		faulted = append(faulted, "MV Aux Transformer Relay")
	}
	if relays.Transformer1Fault {
		faulted = append(faulted, "Transformer 1 Relay")
	}
	if relays.Transformer2Fault {
		faulted = append(faulted, "Transformer 2 Relay")
	}
	if relays.Transformer3Fault {
		faulted = append(faulted, "Transformer 3 Relay")
	}
	if relays.Transformer4Fault {
		faulted = append(faulted, "Transformer 4 Relay")
	}

	return faulted
}

// ResetAllCircuitBreakers attempts to open all circuit breakers (emergency function)
func (s *Service) ResetAllCircuitBreakers() error {
	s.log.Warn("Emergency: Opening all circuit breakers")

	var lastErr error

	// Open auxiliary CB
	if err := s.ControlAuxiliaryCB(false); err != nil {
		s.log.Error("Failed to open auxiliary CB", zap.Error(err))
		lastErr = err
	}
	time.Sleep(100 * time.Millisecond)

	// Open MV aux transformer CB
	if err := s.ControlMVAuxTransformerCB(false); err != nil {
		s.log.Error("Failed to open MV aux transformer CB", zap.Error(err))
		lastErr = err
	}
	time.Sleep(100 * time.Millisecond)

	// Open all transformer CBs
	for i := uint8(1); i <= 4; i++ {
		if err := s.ControlTransformerCB(i, false); err != nil {
			s.log.Error("Failed to open transformer CB",
				zap.Uint8("transformer_no", i),
				zap.Error(err))
			lastErr = err
		}
		time.Sleep(100 * time.Millisecond)
	}

	if lastErr != nil {
		return fmt.Errorf("failed to open all circuit breakers: %w", lastErr)
	}

	s.log.Info("All circuit breakers opened successfully")
	return nil
}
