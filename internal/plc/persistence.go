package plc

import (
	"time"

	"go.uber.org/zap"
)

// persistenceLoop handles data persistence to InfluxDB
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
			s.mutex.RLock()
			plcData := s.lastPLCData
			s.mutex.RUnlock()

			if !plcData.Timestamp.IsZero() {
				if err := s.influxDB.WritePLCData(plcData); err != nil {
					s.log.Error("Failed to write PLC data to InfluxDB", zap.Error(err))
				}
			}

			// Calculate next aligned time and reset timer
			nextTick = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextTick))
		}
	}
}
