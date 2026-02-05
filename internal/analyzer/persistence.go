package analyzer

import (
	"time"

	"go.uber.org/zap"
)

// persistenceLoop handles data persistence
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

// persistData saves current data to InfluxDB
func (s *Service) persistData() {
	s.mutex.RLock()
	dataToSave := s.lastData
	s.mutex.RUnlock()

	if dataToSave.Timestamp.IsZero() {
		return
	}

	if err := s.influxDB.WriteAnalyzerData(dataToSave); err != nil {
		s.log.Error("Failed to save energy analyzer data to InfluxDB", zap.Error(err))
	}
}
