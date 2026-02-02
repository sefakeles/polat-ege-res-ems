package pcs

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// pollLoop periodically reads data from the PCS
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

			if err := s.readAllData(); err != nil {
				s.log.Error("Error reading data", zap.Error(err))
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

// handleConnectionError attempts to reconnect to the PCS
func (s *Service) handleConnectionError() {
	s.log.Warn("PCS connection lost, initiating reconnection procedure")
	s.client.Disconnect()

	reconnectAttempts := 0
	for !s.client.IsConnected() {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(s.config.ReconnectDelay):
			reconnectAttempts++
			if err := s.client.Connect(s.ctx); err != nil {
				s.log.Error("Failed to reconnect to PCS",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts),
					zap.Duration("retry_delay", s.config.ReconnectDelay))
			} else {
				s.log.Info("Successfully reconnected to PCS",
					zap.Int("total_attempts", reconnectAttempts),
					zap.Duration("total_downtime", time.Duration(reconnectAttempts)*s.config.ReconnectDelay))
				return
			}
		}
	}
}

// readAllData reads all data
func (s *Service) readAllData() error {
	// Read PCS data
	if err := s.readPCSData(); err != nil {
		return fmt.Errorf("failed to read PCS data: %w", err)
	}

	return nil
}
