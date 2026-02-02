package pcs

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// pollLoop periodically reads data from the PCS
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
				if err := s.readAllData(); err != nil {
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
						zap.Duration("read_duration", duration),
						zap.Duration("interval", interval))
				}
			}

			// Calculate next aligned time and reset timer
			nextTick = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextTick))
		}
	}
}

// handleConnectionError attempts to reconnect to the PCS
func (s *Service) handleConnectionError() {
	s.log.Warn("PCS connection lost, initiating reconnection procedure")
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
				s.log.Error("Failed to reconnect to PCS",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts))
				timer.Reset(s.config.ReconnectDelay)
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
