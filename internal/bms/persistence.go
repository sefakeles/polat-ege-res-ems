package bms

import (
	"time"

	"go.uber.org/zap"

	"powerkonnekt/ems/internal/database"
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
			s.persistData()

			// Calculate next aligned time and reset timer
			nextTick = time.Now().Truncate(interval).Add(interval)
			timer.Reset(time.Until(nextTick))
		}
	}
}

// persistData writes all data to InfluxDB
func (s *Service) persistData() {
	s.mutex.RLock()
	bmsData := s.lastBMSData
	bmsStatusData := s.lastBMSStatusData
	bmsRackData := make([]database.BMSRackData, len(s.lastBMSRackData))
	copy(bmsRackData, s.lastBMSRackData)
	bmsRackStatusData := make([]database.BMSRackStatusData, len(s.lastBMSRackStatusData))
	copy(bmsRackStatusData, s.lastBMSRackStatusData)

	// Copy cell data
	cellVoltages := make([][]database.BMSCellVoltageData, len(s.lastCellVoltages))
	cellTemperatures := make([][]database.BMSCellTemperatureData, len(s.lastCellTemperatures))
	for i, cells := range s.lastCellVoltages {
		cellVoltages[i] = make([]database.BMSCellVoltageData, len(cells))
		copy(cellVoltages[i], cells)
	}
	for i, cells := range s.lastCellTemperatures {
		cellTemperatures[i] = make([]database.BMSCellTemperatureData, len(cells))
		copy(cellTemperatures[i], cells)
	}
	s.mutex.RUnlock()

	// Save BMS data to InfluxDB
	if !bmsData.Timestamp.IsZero() {
		if err := s.influxDB.WriteBMSData(bmsData); err != nil {
			s.log.Error("Failed to save BMS data to InfluxDB", zap.Error(err))
		}
	}

	// Save BMS status data to InfluxDB
	if !bmsStatusData.Timestamp.IsZero() {
		if err := s.influxDB.WriteBMSStatusData(bmsStatusData); err != nil {
			s.log.Error("Failed to save BMS status data to InfluxDB", zap.Error(err))
		}
	}

	// Save rack data to InfluxDB
	for _, rack := range bmsRackData {
		if !rack.Timestamp.IsZero() {
			if err := s.influxDB.WriteBMSRackData(rack); err != nil {
				s.log.Error("Failed to save BMS rack data to InfluxDB",
					zap.Error(err),
					zap.Uint8("rack_no", rack.Number))
			}
		}
	}

	// Save rack status data to InfluxDB
	for _, rackStatus := range bmsRackStatusData {
		if !rackStatus.Timestamp.IsZero() {
			if err := s.influxDB.WriteBMSRackStatusData(rackStatus); err != nil {
				s.log.Error("Failed to save BMS rack status data to InfluxDB",
					zap.Error(err),
					zap.Uint8("rack_no", rackStatus.Number))
			}
		}
	}

	// Save cell voltage data to InfluxDB
	for rackNo, cells := range cellVoltages {
		if len(cells) > 0 {
			if err := s.influxDB.WriteBMSCellVoltageData(cells); err != nil {
				s.log.Error("Failed to save cell voltage data to InfluxDB",
					zap.Error(err),
					zap.Int("rack_no", rackNo))
			}
		}
	}

	// Save cell temperature data to InfluxDB
	for rackNo, cells := range cellTemperatures {
		if len(cells) > 0 {
			if err := s.influxDB.WriteBMSCellTemperatureData(cells); err != nil {
				s.log.Error("Failed to save cell temperature data to InfluxDB",
					zap.Error(err),
					zap.Int("rack_no", rackNo))
			}
		}
	}
}
