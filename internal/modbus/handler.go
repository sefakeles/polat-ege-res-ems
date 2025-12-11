package modbus

import (
	"sync"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/pkg/logger"

	"github.com/simonvetter/modbus"
)

// RequestHandler implements the modbus.RequestHandler interface
type RequestHandler struct {
	bmsManager   *bms.Manager
	pcsManager   *pcs.Manager
	alarmManager *alarm.Manager
	controlLogic *control.Logic
	registers    *RegisterMap
	mutex        sync.RWMutex
	log          logger.Logger
}

// NewRequestHandler creates a new Modbus request handler
func NewRequestHandler(
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	alarmManager *alarm.Manager,
	controlLogic *control.Logic,
) *RequestHandler {
	// Create handler-specific logger
	handlerLogger := logger.With(
		logger.String("component", "modbus_handler"),
	)

	return &RequestHandler{
		bmsManager:   bmsManager,
		pcsManager:   pcsManager,
		alarmManager: alarmManager,
		controlLogic: controlLogic,
		registers:    NewRegisterMap(),
		log:          handlerLogger,
	}
}

// HandleCoils handles coil read/write requests
func (h *RequestHandler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	return nil, modbus.ErrIllegalFunction
}

// HandleDiscreteInputs handles discrete input read requests
func (h *RequestHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	return nil, modbus.ErrIllegalFunction
}

// HandleHoldingRegisters handles holding register read/write requests
func (h *RequestHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Handle write requests
	if req.IsWrite {
		return h.handleHoldingRegistersWrite(req)
	}

	// Handle read requests
	return h.handleHoldingRegistersRead(req)
}

// HandleInputRegisters handles input register read requests
func (h *RequestHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	switch {
	case req.Addr >= BMSBaseAddr && req.Addr < PCSBaseAddr:
		return h.handleBMSInputRegisters(req.Addr, req.Quantity)
	case req.Addr >= PCSBaseAddr:
		return h.handlePCSInputRegisters(req.Addr, req.Quantity)
	default:
		h.log.Warn("Address out of range",
			logger.Uint16("address", req.Addr))
		return nil, modbus.ErrIllegalDataAddress
	}
}

// handleBMSInputRegisters handles BMS input register reads
func (h *RequestHandler) handleBMSInputRegisters(addr uint16, quantity uint16) ([]uint16, error) {
	bmsNo := h.getBMSNumberFromServerAddress(addr)
	if bmsNo == 0 {
		h.log.Warn("Invalid BMS number",
			logger.Uint8("bms_no", bmsNo))
		return nil, modbus.ErrIllegalDataAddress
	}

	bmsBaseAddr := BMSBaseAddr + uint16(bmsNo-1)*BMSDataOffset
	offsetInBMS := addr - bmsBaseAddr

	if offsetInBMS >= BMSDataStartOffset && offsetInBMS < BMSDataStartOffset+BMSDataLength {
		service, err := h.bmsManager.GetService(int(bmsNo))
		if err != nil {
			h.log.Warn("BMS service not found",
				logger.Uint8("bms_no", bmsNo),
				logger.Err(err))
			return nil, modbus.ErrIllegalDataAddress
		}

		bmsData := service.GetLatestBMSData()
		return h.convertBMSDataToRegisters(bmsData, addr, quantity)
	}

	h.log.Warn("Illegal data address requested",
		logger.Uint16("address", addr),
		logger.Uint8("bms_no", bmsNo))
	return nil, modbus.ErrIllegalDataAddress
}

// handlePCSInputRegisters handles PCS input register reads
func (h *RequestHandler) handlePCSInputRegisters(addr uint16, quantity uint16) ([]uint16, error) {
	pcsNo := h.getPCSNumberFromServerAddress(addr)
	if pcsNo == 0 {
		h.log.Warn("Invalid PCS number",
			logger.Uint8("pcs_no", pcsNo))
		return nil, modbus.ErrIllegalDataAddress
	}

	pcsBaseAddr := PCSBaseAddr + uint16(pcsNo-1)*PCSDataOffset
	offsetInPCS := addr - pcsBaseAddr

	if offsetInPCS >= PCSDataStartOffset && offsetInPCS < PCSDataStartOffset+PCSDataLength {
		service, err := h.pcsManager.GetService(int(pcsNo))
		if err != nil {
			h.log.Warn("PCS service not found",
				logger.Uint8("pcs_no", pcsNo),
				logger.Err(err))
			return nil, modbus.ErrIllegalDataAddress
		}

		pcsData := service.GetLatestPCSData()
		return h.convertPCSDataToRegisters(pcsData, addr, quantity)
	}

	h.log.Warn("Illegal data address requested",
		logger.Uint16("address", addr),
		logger.Uint8("pcs_no", pcsNo))
	return nil, modbus.ErrIllegalDataAddress
}

// handleHoldingRegistersRead handles holding register read requests
func (h *RequestHandler) handleHoldingRegistersRead(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	addr := req.Addr
	quantity := req.Quantity

	// Validate quantity
	if quantity == 0 || quantity > 125 {
		return nil, modbus.ErrIllegalDataValue
	}

	// Calculate PCS number from command address
	if addr < CmdBaseAddr {
		h.log.Warn("Read attempt from invalid command address",
			logger.Uint16("address", addr))
		return nil, modbus.ErrIllegalDataAddress
	}

	pcsNo := uint8((addr-CmdBaseAddr)/CmdOffset) + 1
	cmdOffset := (addr - CmdBaseAddr) % CmdOffset

	// Get PCS service
	service, err := h.pcsManager.GetService(int(pcsNo))
	if err != nil {
		h.log.Warn("PCS service not found for command read",
			logger.Uint8("pcs_no", pcsNo),
			logger.Err(err))
		return nil, modbus.ErrIllegalDataAddress
	}

	// Get command state for this PCS
	cmdState := service.GetCommandState()

	result := make([]uint16, quantity)

	for i := range quantity {
		currentOffset := cmdOffset + uint16(i)

		switch currentOffset {
		case RegStartStopCommand:
			// Return start/stop command state
			if cmdState.StartStopCommand {
				result[i] = 1
			} else {
				result[i] = 0
			}

		case RegActivePowerCommand:
			// Return active power command (kW * 10, signed int16)
			powerValue := int16(cmdState.ActivePowerCommand * 10)
			result[i] = uint16(powerValue)

		case RegReactivePowerCommand:
			// Return reactive power command (kVAr * 10, signed int16)
			powerValue := int16(cmdState.ReactivePowerCommand * 10)
			result[i] = uint16(powerValue)

		default:
			h.log.Warn("Read attempt from unsupported holding register",
				logger.Uint16("address", addr+uint16(i)),
				logger.Uint8("pcs_no", pcsNo),
				logger.Uint16("cmd_offset", currentOffset))
			return nil, modbus.ErrIllegalDataAddress
		}
	}

	return result, nil
}

// handleHoldingRegistersWrite handles holding register write requests (commands)
func (h *RequestHandler) handleHoldingRegistersWrite(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	addr := req.Addr
	values := req.Args

	// Calculate PCS number from command address
	if addr < CmdBaseAddr {
		h.log.Warn("Write attempt to invalid command address",
			logger.Uint16("address", addr))
		return nil, modbus.ErrIllegalDataAddress
	}

	pcsNo := uint8((addr-CmdBaseAddr)/CmdOffset) + 1
	cmdOffset := (addr - CmdBaseAddr) % CmdOffset

	// Get PCS service
	service, err := h.pcsManager.GetService(int(pcsNo))
	if err != nil {
		h.log.Warn("PCS service not found for command",
			logger.Uint8("pcs_no", pcsNo),
			logger.Err(err))
		return nil, modbus.ErrIllegalDataAddress
	}

	switch cmdOffset {
	case RegStartStopCommand:
		// Start/Stop command
		if len(values) < 1 {
			return nil, modbus.ErrIllegalDataValue
		}

		start := values[0] != 0

		h.log.Info("Modbus start/stop command received",
			logger.Uint8("pcs_no", pcsNo),
			logger.Bool("start", start))

		if err := service.StartStopCommand(start); err != nil {
			h.log.Error("Failed to execute Modbus start/stop command",
				logger.Uint8("pcs_no", pcsNo),
				logger.Err(err),
				logger.Bool("start", start))
			return nil, modbus.ErrServerDeviceFailure
		}

		h.log.Info("Modbus start/stop command executed successfully",
			logger.Uint8("pcs_no", pcsNo),
			logger.Bool("start", start))
		return values, nil

	case RegActivePowerCommand:
		// Active power command (kW * 10, signed int16)
		if len(values) < 1 {
			return nil, modbus.ErrIllegalDataValue
		}

		powerValue := int16(values[0])
		power := float32(powerValue) / 10.0

		h.log.Info("Modbus active power command received",
			logger.Uint8("pcs_no", pcsNo),
			logger.Float32("power", power))

		if err := service.SetActivePowerCommand(power); err != nil {
			h.log.Error("Failed to execute Modbus active power command",
				logger.Uint8("pcs_no", pcsNo),
				logger.Err(err),
				logger.Float32("power", power))
			return nil, modbus.ErrServerDeviceFailure
		}

		h.log.Info("Modbus active power command executed successfully",
			logger.Uint8("pcs_no", pcsNo),
			logger.Float32("power", power))
		return values, nil

	case RegReactivePowerCommand:
		// Reactive power command (kVAr * 10, signed int16)
		if len(values) < 1 {
			return nil, modbus.ErrIllegalDataValue
		}

		powerValue := int16(values[0])
		power := float32(powerValue) / 10.0

		h.log.Info("Modbus reactive power command received",
			logger.Uint8("pcs_no", pcsNo),
			logger.Float32("power", power))

		if err := service.SetReactivePowerCommand(power); err != nil {
			h.log.Error("Failed to execute Modbus reactive power command",
				logger.Uint8("pcs_no", pcsNo),
				logger.Err(err),
				logger.Float32("power", power))
			return nil, modbus.ErrServerDeviceFailure
		}

		h.log.Info("Modbus reactive power command executed successfully",
			logger.Uint8("pcs_no", pcsNo),
			logger.Float32("power", power))
		return values, nil

	default:
		h.log.Warn("Write attempt to unsupported holding register",
			logger.Uint16("address", addr),
			logger.Uint8("pcs_no", pcsNo),
			logger.Uint16("cmd_offset", uint16(cmdOffset)))
		return nil, modbus.ErrIllegalDataAddress
	}
}

// getBMSNumberFromServerAddress calculates BMS number from server Modbus address
func (h *RequestHandler) getBMSNumberFromServerAddress(addr uint16) uint8 {
	if addr < BMSBaseAddr {
		return 0
	}
	return uint8((addr-BMSBaseAddr)/BMSDataOffset) + 1
}

// getPCSNumberFromServerAddress calculates PCS number from server Modbus address
func (h *RequestHandler) getPCSNumberFromServerAddress(addr uint16) uint8 {
	if addr < PCSBaseAddr {
		return 0
	}
	return uint8((addr-PCSBaseAddr)/PCSDataOffset) + 1
}
