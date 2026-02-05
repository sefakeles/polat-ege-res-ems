package pcs

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
			s.mutex.RLock()
			statusData := s.lastStatusData
			equipmentData := s.lastEquipmentData
			environmentData := s.lastEnvironmentData
			dcSourceData := s.lastDCSourceData
			gridData := s.lastGridData
			counterData := s.lastCounterData
			s.mutex.RUnlock()

			if !statusData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSStatusData(statusData); err != nil {
					s.log.Error("Failed to write status data", zap.Error(err))
				}
			}

			if !equipmentData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSEquipmentData(equipmentData); err != nil {
					s.log.Error("Failed to write equipment data", zap.Error(err))
				}
			}

			if !environmentData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSEnvironmentData(environmentData); err != nil {
					s.log.Error("Failed to write environment data", zap.Error(err))
				}
			}

			if !dcSourceData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSDCSourceData(dcSourceData); err != nil {
					s.log.Error("Failed to write DC source data", zap.Error(err))
				}
			}

			if !gridData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSGridData(gridData); err != nil {
					s.log.Error("Failed to write grid data", zap.Error(err))
				}
			}

			if !counterData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSCounterData(counterData); err != nil {
					s.log.Error("Failed to write counter data", zap.Error(err))
				}
			}

			// Calculate next aligned time and reset timer
			nextTick = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextTick))
		}
	}
}
