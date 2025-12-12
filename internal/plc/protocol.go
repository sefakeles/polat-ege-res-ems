package plc

// MODBUS Register addresses for PLC
const (
	// Status Data (Read from PLC)
	CircuitBreakerPositionsAddr = 7 // Circuit breaker positions
	MVCircuitBreakerAddr        = 8 // MV-side circuit breaker positions
	ProtectionRelayStatusAddr   = 9 // Protection relay status

	// Control Registers (Write to PLC)
	AuxCBControlAddr          = 10 // Auxiliary CB control
	MVAuxTransformerCBAddr    = 11 // MV Auxiliary Transformer CB control
	Transformer1CBControlAddr = 12 // Transformer 1 CB control
	Transformer2CBControlAddr = 13 // Transformer 2 CB control
	Transformer3CBControlAddr = 14 // Transformer 3 CB control
	Transformer4CBControlAddr = 15 // Transformer 4 CB control
	AutoproducerCBControlAddr = 16 // Autoproducer CB control

	// Data length for reading
	StatusDataLength = 3 // Addresses 7, 8, 9
)

// Control Commands
const (
	ControlNoOperation = 0
	ControlClose       = 1
	ControlOpen        = 2
)

// Circuit Breaker Bit Positions (Address 7)
const (
	BitAuxiliaryCB = 0
	BitPCS1CB      = 1
	BitPCS2CB      = 2
	BitPCS3CB      = 3
	BitPCS4CB      = 4
	BitBMS1CB      = 5
	BitBMS2CB      = 6
	BitBMS3CB      = 7
	BitBMS4CB      = 8
)

// MV Circuit Breaker Bit Positions (Address 8)
const (
	BitMVAuxTransformerCB = 0
	BitTransformer1CB     = 1
	BitTransformer2CB     = 2
	BitTransformer3CB     = 3
	BitTransformer4CB     = 4
	BitAutoproducerCB     = 5
)

// Protection Relay Bit Positions (Address 9)
const (
	BitMVAuxTransformerRelay = 0
	BitTransformer1Relay     = 1
	BitTransformer2Relay     = 2
	BitTransformer3Relay     = 3
	BitTransformer4Relay     = 4
)

// CircuitBreakerStatus represents the status of a circuit breaker
type CircuitBreakerStatus struct {
	AuxiliaryCB bool `json:"auxiliary_cb"`
	PCS1CB      bool `json:"pcs1_cb"`
	PCS2CB      bool `json:"pcs2_cb"`
	PCS3CB      bool `json:"pcs3_cb"`
	PCS4CB      bool `json:"pcs4_cb"`
	BMS1CB      bool `json:"bms1_cb"`
	BMS2CB      bool `json:"bms2_cb"`
	BMS3CB      bool `json:"bms3_cb"`
	BMS4CB      bool `json:"bms4_cb"`
}

// MVCircuitBreakerStatus represents MV-side circuit breaker status
type MVCircuitBreakerStatus struct {
	AuxTransformerCB bool `json:"aux_transformer_cb"`
	Transformer1CB   bool `json:"transformer1_cb"`
	Transformer2CB   bool `json:"transformer2_cb"`
	Transformer3CB   bool `json:"transformer3_cb"`
	Transformer4CB   bool `json:"transformer4_cb"`
	AutoproducerCB   bool `json:"autoproducer_cb"`
}

// ProtectionRelayStatus represents protection relay status
type ProtectionRelayStatus struct {
	AuxTransformerFault bool `json:"aux_transformer_fault"`
	Transformer1Fault   bool `json:"transformer1_fault"`
	Transformer2Fault   bool `json:"transformer2_fault"`
	Transformer3Fault   bool `json:"transformer3_fault"`
	Transformer4Fault   bool `json:"transformer4_fault"`
}

// GetCircuitBreakerName returns human-readable name for circuit breaker
func GetCircuitBreakerName(bit uint8) string {
	names := map[uint8]string{
		BitAuxiliaryCB: "Auxiliary CB",
		BitPCS1CB:      "PCS 1 CB",
		BitPCS2CB:      "PCS 2 CB",
		BitPCS3CB:      "PCS 3 CB",
		BitPCS4CB:      "PCS 4 CB",
		BitBMS1CB:      "BMS 1 CB",
		BitBMS2CB:      "BMS 2 CB",
		BitBMS3CB:      "BMS 3 CB",
		BitBMS4CB:      "BMS 4 CB",
	}
	if name, exists := names[bit]; exists {
		return name
	}
	return "Unknown CB"
}

// GetMVCircuitBreakerName returns human-readable name for MV circuit breaker
func GetMVCircuitBreakerName(bit uint8) string {
	names := map[uint8]string{
		BitMVAuxTransformerCB: "MV Aux Transformer CB",
		BitTransformer1CB:     "Transformer 1 CB",
		BitTransformer2CB:     "Transformer 2 CB",
		BitTransformer3CB:     "Transformer 3 CB",
		BitTransformer4CB:     "Transformer 4 CB",
		BitAutoproducerCB:     "Autoproducer CB",
	}
	if name, exists := names[bit]; exists {
		return name
	}
	return "Unknown MV CB"
}

// GetProtectionRelayName returns human-readable name for protection relay
func GetProtectionRelayName(bit uint8) string {
	names := map[uint8]string{
		BitMVAuxTransformerRelay: "MV Aux Transformer Relay",
		BitTransformer1Relay:     "Transformer 1 Relay",
		BitTransformer2Relay:     "Transformer 2 Relay",
		BitTransformer3Relay:     "Transformer 3 Relay",
		BitTransformer4Relay:     "Transformer 4 Relay",
	}
	if name, exists := names[bit]; exists {
		return name
	}
	return "Unknown Relay"
}
