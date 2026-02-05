package analyzer

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
		CurrentL1:        float32FromBytes(data, 0),   // 3000 - Current A (A)
		CurrentL2:        float32FromBytes(data, 4),   // 3002 - Current B (A)
		CurrentL3:        float32FromBytes(data, 8),   // 3004 - Current C (A)
		CurrentN:         float32FromBytes(data, 12),  // 3006 - Current N (A)
		VoltageL1L2:      float32FromBytes(data, 40),  // 3020 - Voltage A-B (V)
		VoltageL2L3:      float32FromBytes(data, 44),  // 3022 - Voltage B-C (V)
		VoltageL3L1:      float32FromBytes(data, 48),  // 3024 - Voltage C-A (V)
		VoltageLLAvg:     float32FromBytes(data, 52),  // 3026 - Voltage L-L Avg (V)
		VoltageL1:        float32FromBytes(data, 56),  // 3028 - Voltage A-N (V)
		VoltageL2:        float32FromBytes(data, 60),  // 3030 - Voltage B-N (V)
		VoltageL3:        float32FromBytes(data, 64),  // 3032 - Voltage C-N (V)
		VoltageLNAvg:     float32FromBytes(data, 72),  // 3036 - Voltage L-N Avg (V)
		ActivePowerL1:    float32FromBytes(data, 108), // 3054 - Active Power A (W)
		ActivePowerL2:    float32FromBytes(data, 112), // 3056 - Active Power B (W)
		ActivePowerL3:    float32FromBytes(data, 116), // 3058 - Active Power C (W)
		ActivePowerSum:   float32FromBytes(data, 120), // 3060 - Active Power Total (W)
		ReactivePowerL1:  float32FromBytes(data, 124), // 3062 - Reactive Power A (VAr)
		ReactivePowerL2:  float32FromBytes(data, 128), // 3064 - Reactive Power B (VAr)
		ReactivePowerL3:  float32FromBytes(data, 132), // 3066 - Reactive Power C (VAr)
		ReactivePowerSum: float32FromBytes(data, 136), // 3068 - Reactive Power Total (VAr)
		ApparentPowerL1:  float32FromBytes(data, 140), // 3070 - Apparent Power A (VA)
		ApparentPowerL2:  float32FromBytes(data, 144), // 3072 - Apparent Power B (VA)
		ApparentPowerL3:  float32FromBytes(data, 148), // 3074 - Apparent Power C (VA)
		ApparentPowerSum: float32FromBytes(data, 152), // 3076 - Apparent Power Total (VA)
		Frequency:        float32FromBytes(data, 220), // 3110 - Frequency (Hz)
	}
}

// float32FromBytes converts bytes to float32
func float32FromBytes(data []byte, offset int) float32 {
	if len(data) < offset+4 {
		return 0.0
	}

	// Read two 16-bit words in big endian byte order
	word1 := binary.BigEndian.Uint16(data[offset : offset+2])
	word2 := binary.BigEndian.Uint16(data[offset+2 : offset+4])

	// Combine with little endian word order (swap the words)
	bits := uint32(word2)<<16 | uint32(word1)

	// Convert uint32 bits to float32
	return math.Float32frombits(bits)
}
