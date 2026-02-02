package bms

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// baseDataPollLoop periodically reads base data from the BMS
func (s *Service) baseDataPollLoop() {
	if err := s.baseClient.Connect(s.ctx); err != nil {
		s.log.Warn("Initial base Modbus connection failed", zap.Error(err))
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
			if !s.baseClient.IsConnected() {
				s.handleBaseClientConnectionError()
			} else {
				startTime := time.Now()
				if err := s.readBaseData(); err != nil {
					s.log.Error("Error reading base data", zap.Error(err))
				} else {
					// Signal that new base data is available
					select {
					case s.baseDataUpdateChan <- struct{}{}:
					default:
						// Channel full, skip signal
					}
				}

				if duration := time.Since(startTime); duration > interval {
					s.log.Warn("Data read exceeded poll interval (system client)",
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

// cellDataPollLoop periodically reads cell data from the BMS
func (s *Service) cellDataPollLoop() {
	if err := s.cellClient.Connect(s.ctx); err != nil {
		s.log.Warn("Initial cell Modbus connection failed", zap.Error(err))
	}

	interval := s.config.CellDataInterval

	// Calculate first aligned time and create timer
	nextTick := time.Now().Truncate(interval).Add(interval)
	timer := time.NewTimer(time.Until(nextTick))
	defer timer.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			if !s.cellClient.IsConnected() {
				s.handleCellClientConnectionError()
			} else {
				startTime := time.Now()
				if err := s.readCellDataForAllRacks(); err != nil {
					s.log.Error("Error reading cell data", zap.Error(err))
				} else {
					// Signal that new cell data is available
					select {
					case s.cellDataUpdateChan <- struct{}{}:
					default:
						// Channel full, skip signal
					}
				}

				if duration := time.Since(startTime); duration > interval {
					s.log.Warn("Data read exceeded poll interval (cell client)",
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

// handleBaseClientConnectionError attempts to reconnect to the BMS
func (s *Service) handleBaseClientConnectionError() {
	s.log.Warn("BMS base client connection lost, initiating reconnection procedure")
	s.baseClient.Disconnect()

	reconnectAttempts := 0
	timer := time.NewTimer(s.config.ReconnectDelay)
	defer timer.Stop()

	for !s.baseClient.IsConnected() {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			reconnectAttempts++
			if err := s.baseClient.Connect(s.ctx); err != nil {
				s.log.Error("Failed to reconnect to BMS base client",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts))
				timer.Reset(s.config.ReconnectDelay)
			} else {
				s.log.Info("Successfully reconnected to BMS base client",
					zap.Int("total_attempts", reconnectAttempts),
					zap.Duration("total_downtime", time.Duration(reconnectAttempts)*s.config.ReconnectDelay))
				return
			}
		}
	}
}

// handleCellClientConnectionError attempts to reconnect to the BMS
func (s *Service) handleCellClientConnectionError() {
	s.log.Warn("BMS cell client connection lost, initiating reconnection procedure")
	s.cellClient.Disconnect()

	reconnectAttempts := 0
	timer := time.NewTimer(s.config.ReconnectDelay)
	defer timer.Stop()

	for !s.cellClient.IsConnected() {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			reconnectAttempts++
			if err := s.cellClient.Connect(s.ctx); err != nil {
				s.log.Error("Failed to reconnect to BMS cell client",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts))
				timer.Reset(s.config.ReconnectDelay)
			} else {
				s.log.Info("Successfully reconnected to BMS cell client",
					zap.Int("total_attempts", reconnectAttempts),
					zap.Duration("total_downtime", time.Duration(reconnectAttempts)*s.config.ReconnectDelay))
				return
			}
		}
	}
}

// readBaseData reads base data
func (s *Service) readBaseData() error {
	// Read BMS data
	if err := s.readBMSData(); err != nil {
		return fmt.Errorf("failed to read BMS data: %w", err)
	}

	// Read BMS status data
	if err := s.readBMSStatusData(); err != nil {
		return fmt.Errorf("failed to read BMS status data: %w", err)
	}

	// Read alarms
	if err := s.readAlarms(); err != nil {
		return fmt.Errorf("failed to read alarms: %w", err)
	}

	for rackNo := uint8(1); rackNo <= uint8(s.config.RackCount); rackNo++ {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
		}

		// Read BMS rack data
		if err := s.readBMSRackData(rackNo); err != nil {
			s.log.Error("Failed to read BMS rack data",
				zap.Error(err),
				zap.Uint8("rack_no", rackNo))
		}
	}

	return nil
}

// readCellDataForAllRacks reads cell data for all racks
func (s *Service) readCellDataForAllRacks() error {
	for rackNo := uint8(1); rackNo <= uint8(s.config.RackCount); rackNo++ {
		// Check if context is cancelled before reading each rack
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
		}

		// Read cell data
		if err := s.readCellData(rackNo); err != nil {
			s.log.Error("Failed to read cell data",
				zap.Error(err),
				zap.Uint8("rack_no", rackNo))
		}
	}

	return nil
}
