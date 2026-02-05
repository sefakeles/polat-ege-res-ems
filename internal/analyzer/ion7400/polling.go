package ion7400

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// pollLoop handles continuous data polling
func (s *Service) pollLoop() {
	if err := s.client.Connect(s.ctx); err != nil {
		s.log.Error("Initial ION7400 connection failed", zap.Error(err))
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

// handleConnectionError attempts to reconnect to the ION7400
func (s *Service) handleConnectionError() {
	s.log.Warn("ION7400 connection lost, attempting reconnection")
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
				s.log.Error("Failed to reconnect to ION7400",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts))
				timer.Reset(s.config.ReconnectDelay)
			} else {
				s.log.Info("Successfully reconnected to ION7400",
					zap.Int("total_attempts", reconnectAttempts),
					zap.Duration("total_downtime", time.Duration(reconnectAttempts)*s.config.ReconnectDelay))
			}
		}
	}
}

// readAllData reads all necessary data from the ION7400
func (s *Service) readAllData() error {
	if err := s.readBaseData(); err != nil {
		return fmt.Errorf("failed to read base data: %w", err)
	}

	return nil
}
