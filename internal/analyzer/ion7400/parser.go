package ion7400

import (
	"encoding/binary"
	"math"
	"time"

	"powerkonnekt/ems/internal/database"
)

// parseBaseData converts raw MODBUS data to AnalyzerData structure
func parseBaseData(data []byte) database.AnalyzerData {
	if len(data) < BaseDataLength*2 {
		return database.AnalyzerData{Timestamp: time.Now()}
	}

	return database.AnalyzerData{
		Timestamp:        time.Now(),
		CurrentL1:        float32FromBytes(data, 0),   // 2999 - Current A (A)
		CurrentL2:        float32FromBytes(data, 4),   // 3001 - Current B (A)
		CurrentL3:        float32FromBytes(data, 8),   // 3003 - Current C (A)
		CurrentN:         float32FromBytes(data, 12),  // 3005 - Current N (A)
		VoltageL1L2:      float32FromBytes(data, 40),  // 3019 - Voltage A-B (V)
		VoltageL2L3:      float32FromBytes(data, 44),  // 3021 - Voltage B-C (V)
		VoltageL3L1:      float32FromBytes(data, 48),  // 3023 - Voltage C-A (V)
		VoltageLLAvg:     float32FromBytes(data, 52),  // 3025 - Voltage L-L Avg (V)
		VoltageL1:        float32FromBytes(data, 56),  // 3027 - Voltage A-N (V)
		VoltageL2:        float32FromBytes(data, 60),  // 3029 - Voltage B-N (V)
		VoltageL3:        float32FromBytes(data, 64),  // 3031 - Voltage C-N (V)
		VoltageLNAvg:     float32FromBytes(data, 72),  // 3035 - Voltage L-N Avg (V)
		ActivePowerL1:    float32FromBytes(data, 108), // 3053 - Active Power A (W)
		ActivePowerL2:    float32FromBytes(data, 112), // 3055 - Active Power B (W)
		ActivePowerL3:    float32FromBytes(data, 116), // 3057 - Active Power C (W)
		ActivePowerSum:   float32FromBytes(data, 120), // 3059 - Active Power Total (W)
		ReactivePowerL1:  float32FromBytes(data, 124), // 3061 - Reactive Power A (VAr)
		ReactivePowerL2:  float32FromBytes(data, 128), // 3063 - Reactive Power B (VAr)
		ReactivePowerL3:  float32FromBytes(data, 132), // 3065 - Reactive Power C (VAr)
		ReactivePowerSum: float32FromBytes(data, 136), // 3067 - Reactive Power Total (VAr)
		ApparentPowerL1:  float32FromBytes(data, 140), // 3069 - Apparent Power A (VA)
		ApparentPowerL2:  float32FromBytes(data, 144), // 3071 - Apparent Power B (VA)
		ApparentPowerL3:  float32FromBytes(data, 148), // 3073 - Apparent Power C (VA)
		ApparentPowerSum: float32FromBytes(data, 152), // 3075 - Apparent Power Total (VA)
		Frequency:        float32FromBytes(data, 220), // 3109 - Frequency (Hz)
	}
}

// parsePowerFactorData converts raw MODBUS data to AnalyzerData structure
func parsePowerFactorData(data []byte) database.AnalyzerData {
	if len(data) < PowerFactorDataLength*2 {
		return database.AnalyzerData{Timestamp: time.Now()}
	}

	return database.AnalyzerData{
		Timestamp:      time.Now(),
		PowerFactorL1:  float32FromBytes(data, 0),  // 3143 - Power Factor A
		PowerFactorL2:  float32FromBytes(data, 4),  // 3145 - Power Factor B
		PowerFactorL3:  float32FromBytes(data, 8),  // 3147 - Power Factor C
		PowerFactorAvg: float32FromBytes(data, 12), // 3149 - Power Factor Avg
	}
}

// parseEnergyData converts raw MODBUS data to AnalyzerData structure for energy registers
func parseEnergyData(data []byte) database.AnalyzerData {
	if len(data) < EnergyDataLength*2 {
		return database.AnalyzerData{Timestamp: time.Now()}
	}

	return database.AnalyzerData{
		Timestamp:            time.Now(),
		ActiveEnergyExport:   int64FromBytes(data, 0),  // 3203 - Active Energy Export (Wh)
		ActiveEnergyImport:   int64FromBytes(data, 8),  // 3207 - Active Energy Import (Wh)
		ReactiveEnergyExport: int64FromBytes(data, 32), // 3219 - Reactive Energy Export (VARh)
		ReactiveEnergyImport: int64FromBytes(data, 40), // 3223 - Reactive Energy Import (VARh)
		ApparentEnergyExport: int64FromBytes(data, 64), // 3235 - Apparent Energy Export (VAh)
		ApparentEnergyImport: int64FromBytes(data, 72), // 3239 - Apparent Energy Import (VAh)
	}
}

// float32FromBytes converts bytes to float32
func float32FromBytes(data []byte, offset int) float32 {
	if len(data) < offset+4 {
		return 0.0
	}

	// Convert 4 bytes to uint32
	bits := binary.BigEndian.Uint32(data[offset : offset+4])

	// Convert uint32 bits to float32
	return math.Float32frombits(bits)
}

// int64FromBytes converts bytes to int64
func int64FromBytes(data []byte, offset int) int64 {
	if len(data) < offset+8 {
		return 0
	}

	// Convert 8 bytes to int64 using big-endian byte order
	return int64(binary.BigEndian.Uint64(data[offset : offset+8]))
}
