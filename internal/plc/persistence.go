package plc

import (
	"time"

	"powerkonnekt/ems/pkg/logger"
)

// persistenceLoop handles data persistence to InfluxDB
func (s *Service) persistenceLoop() {
	ticker := time.NewTicker(s.config.PersistInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mutex.RLock()
			plcData := s.lastPLCData
			s.mutex.RUnlock()

			if !plcData.Timestamp.IsZero() {
				if err := s.influxDB.WritePLCData(plcData); err != nil {
					s.log.Error("Failed to write PLC data to InfluxDB", logger.Err(err))
				}
			}
		}
	}
}
