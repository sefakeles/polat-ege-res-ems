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

	ticker := time.NewTicker(s.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if !s.baseClient.IsConnected() {
				s.handleBaseClientConnectionError()
			}

			if err := s.readBaseData(); err != nil {
				s.log.Error("Error reading base data", zap.Error(err))
				continue
			}

			// Signal that new base data is available
			select {
			case s.baseDataUpdateChan <- struct{}{}:
			default:
				// Channel full, skip signal
			}
		}
	}
}

// cellDataPollLoop periodically reads cell data from the BMS
func (s *Service) cellDataPollLoop() {
	if !s.config.EnableCellData {
		return
	}

	if err := s.cellClient.Connect(s.ctx); err != nil {
		s.log.Warn("Initial cell Modbus connection failed", zap.Error(err))
	}

	ticker := time.NewTicker(s.config.CellDataInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			if !s.cellClient.IsConnected() {
				s.handleCellClientConnectionError()
			}

			if err := s.readCellDataForAllRacks(); err != nil {
				s.log.Error("Error reading cell data", zap.Error(err))
				continue
			}

			// Signal that new cell data is available
			select {
			case s.cellDataUpdateChan <- struct{}{}:
			default:
				// Channel full, skip signal
			}
		}
	}
}

// handleBaseClientConnectionError attempts to reconnect to the BMS
func (s *Service) handleBaseClientConnectionError() {
	s.log.Warn("BMS base client connection lost, initiating reconnection procedure")
	s.baseClient.Disconnect()

	reconnectAttempts := 0
	for !s.baseClient.IsConnected() {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(s.config.ReconnectDelay):
			reconnectAttempts++
			if err := s.baseClient.Connect(s.ctx); err != nil {
				s.log.Error("Failed to reconnect to BMS base client",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts),
					zap.Duration("retry_delay", s.config.ReconnectDelay))
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
	for !s.cellClient.IsConnected() {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(s.config.ReconnectDelay):
			reconnectAttempts++
			if err := s.cellClient.Connect(s.ctx); err != nil {
				s.log.Error("Failed to reconnect to BMS cell client",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts),
					zap.Duration("retry_delay", s.config.ReconnectDelay))
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
