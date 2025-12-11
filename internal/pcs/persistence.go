package pcs

import (
	"time"

	"powerkonnekt/ems/pkg/logger"
)

// persistenceLoop handles data persistence
func (s *Service) persistenceLoop() {
	ticker := time.NewTicker(s.config.PersistInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
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
					s.log.Error("Failed to write status data", logger.Err(err))
				}
			}

			if !equipmentData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSEquipmentData(equipmentData); err != nil {
					s.log.Error("Failed to write equipment data", logger.Err(err))
				}
			}

			if !environmentData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSEnvironmentData(environmentData); err != nil {
					s.log.Error("Failed to write environment data", logger.Err(err))
				}
			}

			if !dcSourceData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSDCSourceData(dcSourceData); err != nil {
					s.log.Error("Failed to write DC source data", logger.Err(err))
				}
			}

			if !gridData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSGridData(gridData); err != nil {
					s.log.Error("Failed to write grid data", logger.Err(err))
				}
			}

			if !counterData.Timestamp.IsZero() {
				if err := s.influxDB.WritePCSCounterData(counterData); err != nil {
					s.log.Error("Failed to write counter data", logger.Err(err))
				}
			}
		}
	}
}
