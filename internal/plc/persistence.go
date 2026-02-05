package plc

import (
	"time"

	"go.uber.org/zap"
)

// persistenceLoop periodically writes data to InfluxDB
func (s *Service) persistenceLoop() {
	interval := s.config.PersistInterval

	// Calculate first aligned time and create timer
	nextTick := time.Now().Truncate(interval).Add(interval)
	timer := time.NewTimer(time.Until(nextTick))
	defer timer.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-timer.C:
			startTime := time.Now()
			s.persistData()
			if duration := time.Since(startTime); duration > interval {
				s.log.Warn("Data persistence exceeded persist interval",
					zap.Duration("duration", duration),
					zap.Duration("interval", interval))
			}

			// Calculate next aligned time and reset timer
			nextTick = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextTick))
		}
	}
}

// persistData writes all data to InfluxDB
func (s *Service) persistData() {
	s.mutex.RLock()
	plcData := s.lastPLCData
	s.mutex.RUnlock()

	if !plcData.Timestamp.IsZero() {
		if err := s.influxDB.WritePLCData(plcData); err != nil {
			s.log.Error("Failed to write PLC data to InfluxDB", zap.Error(err))
		}
	}
}
