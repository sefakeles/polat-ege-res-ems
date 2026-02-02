package windfarm

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// dataPollLoop periodically reads data from the Wind Farm FCU
func (s *Service) dataPollLoop() {
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
					s.log.Error("Error reading wind farm data", zap.Error(err))
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

// handleConnectionError attempts to reconnect to the FCU
func (s *Service) handleConnectionError() {
	s.log.Warn("Wind Farm FCU connection lost, initiating reconnection procedure")
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
				s.log.Error("Failed to reconnect to Wind Farm FCU",
					zap.Error(err),
					zap.Int("attempt", reconnectAttempts),
					zap.Duration("retry_delay", s.config.ReconnectDelay))
				timer.Reset(s.config.ReconnectDelay)
			} else {
				s.log.Info("Successfully reconnected to Wind Farm FCU",
					zap.Int("total_attempts", reconnectAttempts),
					zap.Duration("total_downtime", time.Duration(reconnectAttempts)*s.config.ReconnectDelay))
				return
			}
		}
	}
}

// readAllData reads all data from the FCU
func (s *Service) readAllData() error {
	// Read measuring data (registers 700-759)
	if err := s.readMeasuringData(); err != nil {
		return fmt.Errorf("failed to read measuring data: %w", err)
	}

	// Read return values / status data (registers 649-689)
	if err := s.readReturnValues(); err != nil {
		return fmt.Errorf("failed to read return values: %w", err)
	}

	return nil
}

// readMeasuringData reads measuring data from registers 700-759
func (s *Service) readMeasuringData() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, MeasuringDataStartAddr, MeasuringDataLength)
	if err != nil {
		return fmt.Errorf("failed to read measuring data registers: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.lastMeasuringData = ParseMeasuringData(data, s.config.ID)
	s.lastWeatherData = ParseWeatherData(data, s.config.ID)
	s.lastStatusData.FCUMode = ParseFCUMode(data)

	return nil
}

// readReturnValues reads return values / status data from registers 649-689
func (s *Service) readReturnValues() error {
	data, err := s.client.ReadHoldingRegisters(s.ctx, ReturnValuesStartAddr, ReturnValuesLength)
	if err != nil {
		return fmt.Errorf("failed to read return values registers: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Parse status data (preserving FCUMode from measuring data)
	fcuMode := s.lastStatusData.FCUMode
	s.lastStatusData = ParseStatusData(data, s.config.ID)
	s.lastStatusData.FCUMode = fcuMode

	// Parse setpoint data
	s.lastSetpointData = ParseSetpointData(data, s.config.ID)

	return nil
}
