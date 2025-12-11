package plc

import (
	"time"

	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/pkg/utils"
)

// ParsePLCData parses raw Modbus data to PLCData structure
func ParsePLCData(data []byte, id int) database.PLCData {
	if len(data) < StatusDataLength*2 {
		return database.PLCData{
			Timestamp: time.Now(),
			ID:        id,
		}
	}

	// Parse circuit breaker positions (address 7)
	cbPositions := utils.FromBytes[uint16](data[0:2])

	// Parse MV circuit breaker positions (address 8)
	mvCBPositions := utils.FromBytes[uint16](data[2:4])

	// Parse protection relay status (address 9)
	relayStatus := utils.FromBytes[uint16](data[4:6])

	return database.PLCData{
		Timestamp:         time.Now(),
		ID:                id,
		CircuitBreakers:   parseCircuitBreakers(cbPositions),
		MVCircuitBreakers: parseMVCircuitBreakers(mvCBPositions),
		ProtectionRelays:  parseProtectionRelays(relayStatus),
	}
}

// parseCircuitBreakers extracts individual circuit breaker states from register value
func parseCircuitBreakers(value uint16) database.CircuitBreakerStatus {
	return database.CircuitBreakerStatus{
		AuxiliaryCB: (value & (1 << BitAuxiliaryCB)) != 0,
		PCS1CB:      (value & (1 << BitPCS1CB)) != 0,
		PCS2CB:      (value & (1 << BitPCS2CB)) != 0,
		PCS3CB:      (value & (1 << BitPCS3CB)) != 0,
		PCS4CB:      (value & (1 << BitPCS4CB)) != 0,
		BMS1CB:      (value & (1 << BitBMS1CB)) != 0,
		BMS2CB:      (value & (1 << BitBMS2CB)) != 0,
		BMS3CB:      (value & (1 << BitBMS3CB)) != 0,
		BMS4CB:      (value & (1 << BitBMS4CB)) != 0,
	}
}

// parseMVCircuitBreakers extracts MV circuit breaker states from register value
func parseMVCircuitBreakers(value uint16) database.MVCircuitBreakerStatus {
	return database.MVCircuitBreakerStatus{
		AuxTransformerCB: (value & (1 << BitMVAuxTransformerCB)) != 0,
		Transformer1CB:   (value & (1 << BitTransformer1CB)) != 0,
		Transformer2CB:   (value & (1 << BitTransformer2CB)) != 0,
		Transformer3CB:   (value & (1 << BitTransformer3CB)) != 0,
		Transformer4CB:   (value & (1 << BitTransformer4CB)) != 0,
	}
}

// parseProtectionRelays extracts protection relay states from register value
func parseProtectionRelays(value uint16) database.ProtectionRelayStatus {
	return database.ProtectionRelayStatus{
		AuxTransformerFault: (value & (1 << BitMVAuxTransformerRelay)) != 0,
		Transformer1Fault:   (value & (1 << BitTransformer1Relay)) != 0,
		Transformer2Fault:   (value & (1 << BitTransformer2Relay)) != 0,
		Transformer3Fault:   (value & (1 << BitTransformer3Relay)) != 0,
		Transformer4Fault:   (value & (1 << BitTransformer4Relay)) != 0,
	}
}
