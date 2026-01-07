package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"powerkonnekt/ems/internal/alarm"
	"powerkonnekt/ems/internal/bms"
	"powerkonnekt/ems/internal/control"
	"powerkonnekt/ems/internal/database"
	"powerkonnekt/ems/internal/health"
	"powerkonnekt/ems/internal/pcs"
	"powerkonnekt/ems/internal/plc"
	"powerkonnekt/ems/internal/windfarm"
	"powerkonnekt/ems/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Handlers contains all API handlers
type Handlers struct {
	bmsManager      *bms.Manager
	pcsManager      *pcs.Manager
	plcManager      *plc.Manager
	windFarmManager *windfarm.Manager
	alarmManager    *alarm.Manager
	controlLogic    *control.Logic
	healthService   *health.HealthService
	log             logger.Logger
}

// NewHandlers creates a new handlers instance
func NewHandlers(
	bmsManager *bms.Manager,
	pcsManager *pcs.Manager,
	plcManager *plc.Manager,
	windFarmManager *windfarm.Manager,
	alarmManager *alarm.Manager,
	controlLogic *control.Logic,
	healthService *health.HealthService,
) *Handlers {
	// Create handlers-specific logger
	handlersLogger := logger.With(
		logger.String("component", "api_handlers"),
	)

	return &Handlers{
		bmsManager:      bmsManager,
		pcsManager:      pcsManager,
		plcManager:      plcManager,
		windFarmManager: windFarmManager,
		alarmManager:    alarmManager,
		controlLogic:    controlLogic,
		healthService:   healthService,
		log:             handlersLogger,
	}
}

// HealthCheck returns detailed health status
func (h *Handlers) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	results := h.healthService.CheckAll(ctx)
	overallStatus := h.healthService.GetOverallStatus(results)

	response := gin.H{
		"checks": results,
		"status": overallStatus,
	}

	statusCode := http.StatusOK
	switch overallStatus {
	case health.StatusUnhealthy:
		statusCode = http.StatusServiceUnavailable
		h.log.Warn("Health check failed - system unhealthy",
			logger.String("status", string(overallStatus)))
	case health.StatusDegraded:
		statusCode = http.StatusPartialContent
		h.log.Warn("Health check shows degraded status",
			logger.String("status", string(overallStatus)))
	}

	c.JSON(statusCode, response)
}

// GetStatus returns system status
func (h *Handlers) GetStatus(c *gin.Context) {
	service, err := h.bmsManager.GetService(1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bmsData := service.GetLatestBMSData()
	bmsStatusData := service.GetLatestBMSStatusData()
	bmsRackData := service.GetLatestBMSRackData()
	activeAlarms := h.alarmManager.GetActiveAlarms()

	status := gin.H{
		"control_mode":         h.controlLogic.GetMode(),
		"active_power_control": h.controlLogic.GetActivePowerControl(),
		"bess_connected":       service.IsConnected(),
		"bms_soc":              bmsData.SOC,
		"bms_soh":              bmsData.SOH,
		"bms_voltage":          bmsData.Voltage,
		"bms_current":          bmsData.Current,
		"bms_state":            bms.GetStateDescription(bmsStatusData.SystemStatus),
		"active_alarms":        len(activeAlarms),
		"rack_count":           len(bmsRackData),
		"critical_alarms":      h.alarmManager.HasCriticalAlarms(),
	}

	c.JSON(http.StatusOK, status)
}

// GetBMSData returns BMS data
func (h *Handlers) GetBMSData(c *gin.Context) {
	bmsID := c.Param("id")
	bmsIDInt, err := strconv.Atoi(bmsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BMS ID"})
		return
	}

	service, err := h.bmsManager.GetService(bmsIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bmsData := service.GetLatestBMSData()
	bmsStatusData := service.GetLatestBMSStatusData()
	bmsRackData := service.GetLatestBMSRackData()

	// Create BMS data response with state description instead of numeric value
	type BMSDataResponse struct {
		database.BMSData
		State string `json:"state"`
	}

	bmsDataResponse := BMSDataResponse{
		BMSData: bmsData,
		State:   bms.GetStateDescription(bmsStatusData.SystemStatus),
	}

	// Create BMS rack data response with state description instead of numeric value
	type BMSRackDataResponse struct {
		database.BMSRackData
		State string `json:"state"`
	}

	bmsRackDataResponse := make([]BMSRackDataResponse, len(bmsRackData))
	for i, rackData := range bmsRackData {
		bmsRackDataResponse[i] = BMSRackDataResponse{
			BMSRackData: rackData,
			State:       bms.GetStateDescription(rackData.State),
		}
	}

	response := gin.H{
		"bms_data":      bmsDataResponse,
		"bms_rack_data": bmsRackDataResponse,
		"bms_connected": service.IsConnected(),
	}

	c.JSON(http.StatusOK, response)
}

// GetBMSRacks returns all BMS rack data
func (h *Handlers) GetBMSRacks(c *gin.Context) {
	bmsID := c.Param("id")
	bmsIDInt, err := strconv.Atoi(bmsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BMS ID"})
		return
	}

	service, err := h.bmsManager.GetService(bmsIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bmsRackData := service.GetLatestBMSRackData()

	response := gin.H{
		"rack_count": len(bmsRackData),
		"racks":      bmsRackData,
	}

	c.JSON(http.StatusOK, response)
}

// GetBMSRackData returns specific BMS rack data
func (h *Handlers) GetBMSRackData(c *gin.Context) {
	bmsID := c.Param("id")
	bmsIDInt, err := strconv.Atoi(bmsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BMS ID"})
		return
	}

	service, err := h.bmsManager.GetService(bmsIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	rackNoParam := c.Param("rack_no")
	rackNo, err := strconv.Atoi(rackNoParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rack number"})
		return
	}

	bmsRackData := service.GetLatestBMSRackData()

	if rackNo < 1 || rackNo > len(bmsRackData) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rack not found"})
		return
	}

	c.JSON(http.StatusOK, bmsRackData[rackNo-1])
}

// GetBMSCommandState returns BMS command state
func (h *Handlers) GetBMSCommandState(c *gin.Context) {
	bmsID := c.Param("id")
	bmsIDInt, err := strconv.Atoi(bmsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BMS ID"})
		return
	}

	service, err := h.bmsManager.GetService(bmsIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	bmsCommandState := service.GetCommandState()

	c.JSON(http.StatusOK, gin.H{
		"command_state": bmsCommandState,
	})
}

// GetPCSData returns PCS data
func (h *Handlers) GetPCSData(c *gin.Context) {
	pcsID := c.Param("id")
	pcsIDInt, err := strconv.Atoi(pcsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PCS ID"})
		return
	}

	service, err := h.pcsManager.GetService(pcsIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	pcsData := service.GetLatestPCSStatusData()

	c.JSON(http.StatusOK, gin.H{
		"data": pcsData,
	})
}

// GetPCSCommandState returns PCS command state
func (h *Handlers) GetPCSCommandState(c *gin.Context) {
	pcsID := c.Param("id")
	pcsIDInt, err := strconv.Atoi(pcsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PCS ID"})
		return
	}

	service, err := h.pcsManager.GetService(pcsIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	pcsCommandState := service.GetCommandState()

	c.JSON(http.StatusOK, gin.H{
		"command_state": pcsCommandState,
	})
}

// GetAlarms returns alarm information
func (h *Handlers) GetAlarms(c *gin.Context) {
	// Get query parameters
	alarmType := c.Query("type")
	severity := c.Query("severity")
	active := c.Query("active")

	var alarms []any

	if active == "false" {
		// Get alarm history
		limit := 100
		offset := 0
		if l := c.Query("limit"); l != "" {
			if parsed, parseErr := strconv.Atoi(l); parseErr == nil {
				limit = parsed
			}
		}
		if o := c.Query("offset"); o != "" {
			if parsed, parseErr := strconv.Atoi(o); parseErr == nil {
				offset = parsed
			}
		}

		history, err := h.alarmManager.GetAlarmHistory(limit, offset)
		if err != nil {
			h.log.Error("Failed to get alarm history",
				logger.Err(err),
				logger.Int("limit", limit),
				logger.Int("offset", offset))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, alarm := range history {
			alarms = append(alarms, alarm)
		}
	} else {
		// Get active alarms
		activeAlarms := h.alarmManager.GetActiveAlarms()

		// Filter by type and severity if specified
		for _, alarm := range activeAlarms {
			if alarmType != "" && alarm.AlarmType != alarmType {
				continue
			}
			if severity != "" && alarm.Severity != severity {
				continue
			}
			alarms = append(alarms, alarm)
		}
	}

	response := gin.H{
		"alarms":      alarms,
		"total_count": len(alarms),
		"timestamp":   time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// SetControlMode sets the control mode
func (h *Handlers) SetControlMode(c *gin.Context) {
	var request struct {
		Mode string `json:"mode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Warn("Invalid control mode request", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Control mode change requested",
		logger.String("requested_mode", request.Mode))

	// Validate mode
	validModes := []string{"AUTO", "MANUAL", "MAINTENANCE"}
	isValid := false
	for _, mode := range validModes {
		if request.Mode == mode {
			isValid = true
			break
		}
	}

	if !isValid {
		h.log.Warn("Invalid control mode requested",
			logger.String("mode", request.Mode))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode. Valid modes: AUTO, MANUAL, MAINTENANCE"})
		return
	}

	h.controlLogic.SetMode(request.Mode)

	h.log.Info("Control mode changed successfully",
		logger.String("mode", request.Mode))

	c.JSON(http.StatusOK, gin.H{
		"message": "Control mode set successfully",
		"mode":    request.Mode,
	})
}

// SetPCSStartStop starts or stops the PCS
func (h *Handlers) SetPCSStartStop(c *gin.Context) {
	var req struct {
		ID    int   `json:"id" binding:"required"`
		Start *bool `json:"start" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.pcsManager.GetService(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.StartStopCommand(*req.Start); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	action := "stopped"
	if *req.Start {
		action = "started"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("PCS %s successfully", action),
		"start":   *req.Start,
	})
}

// SetPowerCommand sets manual power command
func (h *Handlers) SetPowerCommand(c *gin.Context) {
	var request struct {
		ID    int      `json:"id" binding:"required"`
		Power *float32 `json:"power" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Warn("Invalid power command request", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.pcsManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if request.Power == nil {
		h.log.Warn("Power command request missing power field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "power field is required"})
		return
	}

	h.log.Info("Manual power command requested",
		logger.Float32("power", *request.Power))

	// Execute manual power command
	if err := service.SetActivePowerCommand(*request.Power); err != nil {
		h.log.Error("Manual power command failed",
			logger.Err(err),
			logger.Float32("power", *request.Power))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Manual power command executed successfully",
		logger.Float32("power", *request.Power))

	c.JSON(http.StatusOK, gin.H{
		"message": "Power command executed successfully",
		"power":   *request.Power,
	})
}

// SetReactivePowerCommand sets manual reactive power command
func (h *Handlers) SetReactivePowerCommand(c *gin.Context) {
	var request struct {
		ID    int      `json:"id" binding:"required"`
		Power *float32 `json:"power" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.log.Warn("Invalid reactive power command request", logger.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.pcsManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if request.Power == nil {
		h.log.Warn("Reactive power command request missing power field")
		c.JSON(http.StatusBadRequest, gin.H{"error": "power field is required"})
		return
	}

	h.log.Info("Manual reactive power command requested",
		logger.Float32("power", *request.Power))

	// Execute manual power command
	if err := service.SetReactivePowerCommand(*request.Power); err != nil {
		h.log.Error("Manual reactive power command failed",
			logger.Err(err),
			logger.Float32("power", *request.Power))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Manual power command executed successfully",
		logger.Float32("power", *request.Power))

	c.JSON(http.StatusOK, gin.H{
		"message": "Power command executed successfully",
		"power":   *request.Power,
	})
}

// BMS Control Handlers

// BMSReset resets the BMS system
func (h *Handlers) BMSReset(c *gin.Context) {
	var request struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.bmsManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.ResetSystem(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("BMS system reset executed",
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{"message": "BMS system reset executed"})
}

// PCSReset resets the PCS system
func (h *Handlers) PCSReset(c *gin.Context) {
	var request struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.pcsManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.ResetSystem(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("PCS system reset executed",
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{"message": "PCS system reset executed"})
}

// BMSBreakerControl controls the main breaker
func (h *Handlers) BMSBreakerControl(c *gin.Context) {
	var request struct {
		ID     int    `json:"id" binding:"required"`
		Action string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.bmsManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var action uint16
	switch request.Action {
	case "OPEN":
		action = bms.ControlOff
	case "CLOSE":
		action = bms.ControlOn
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Valid actions: OPEN, CLOSE"})
		return
	}

	if err := service.ControlMainBreaker(action); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("BMS breaker control executed",
		logger.String("action", request.Action),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": "Breaker control executed",
		"action":  request.Action,
	})
}

// GetPLCData returns PLC data
func (h *Handlers) GetPLCData(c *gin.Context) {
	plcID := c.Param("id")
	plcIDInt, err := strconv.Atoi(plcID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PLC ID"})
		return
	}

	service, err := h.plcManager.GetService(plcIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	plcData := service.GetLatestPLCData()
	cbStatus := service.GetCircuitBreakerStatus()
	mvCBStatus := service.GetMVCircuitBreakerStatus()
	protectionRelayStatus := service.GetProtectionRelayStatus()

	response := gin.H{
		"data":                    plcData,
		"circuit_breakers":        cbStatus,
		"mv_circuit_breakers":     mvCBStatus,
		"protection_relay_status": protectionRelayStatus,
		"connected":               service.IsConnected(),
		"relay_faults":            service.HasProtectionRelayFaults(),
		"faulted_relays":          service.GetFaultedRelays(),
	}

	c.JSON(http.StatusOK, response)
}

// ControlAuxiliaryCB controls the auxiliary circuit breaker
func (h *Handlers) ControlAuxiliaryCB(c *gin.Context) {
	var request struct {
		ID    int   `json:"id" binding:"required"`
		Close *bool `json:"close" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.plcManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.ControlAuxiliaryCB(*request.Close); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	action := "opened"
	if *request.Close {
		action = "closed"
	}

	h.log.Info("Auxiliary CB control executed",
		logger.String("action", action),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Auxiliary CB %s successfully", action),
		"close":   *request.Close,
	})
}

// ControlMVAuxTransformerCB controls the MV auxiliary transformer circuit breaker
func (h *Handlers) ControlMVAuxTransformerCB(c *gin.Context) {
	var request struct {
		ID    int   `json:"id" binding:"required"`
		Close *bool `json:"close" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.plcManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.ControlMVAuxTransformerCB(*request.Close); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	action := "opened"
	if *request.Close {
		action = "closed"
	}

	h.log.Info("MV Aux Transformer CB control executed",
		logger.String("action", action),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("MV Aux Transformer CB %s successfully", action),
		"close":   *request.Close,
	})
}

// ControlTransformerCB controls a transformer circuit breaker
func (h *Handlers) ControlTransformerCB(c *gin.Context) {
	var request struct {
		ID            int   `json:"id" binding:"required"`
		TransformerNo uint8 `json:"transformer_no" binding:"required,min=1,max=4"`
		Close         *bool `json:"close" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.plcManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.ControlTransformerCB(request.TransformerNo, *request.Close); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	action := "opened"
	if *request.Close {
		action = "closed"
	}

	h.log.Info("Transformer CB control executed",
		logger.Uint8("transformer_no", request.TransformerNo),
		logger.String("action", action),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message":        fmt.Sprintf("Transformer %d CB %s successfully", request.TransformerNo, action),
		"transformer_no": request.TransformerNo,
		"close":          *request.Close,
	})
}

// ControlAutoproducerCB controls the autoproducer circuit breaker
func (h *Handlers) ControlAutoproducerCB(c *gin.Context) {
	var request struct {
		ID    int   `json:"id" binding:"required"`
		Close *bool `json:"close" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.plcManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.ControlAutoproducerCB(*request.Close); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	action := "opened"
	if *request.Close {
		action = "closed"
	}

	h.log.Info("Autoproducer CB control executed",
		logger.String("action", action),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Autoproducer CB %s successfully", action),
		"close":   *request.Close,
	})
}

// ResetAllCircuitBreakers opens all circuit breakers (emergency function)
func (h *Handlers) ResetAllCircuitBreakers(c *gin.Context) {
	var request struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.plcManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.ResetAllCircuitBreakers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Warn("Emergency: All circuit breakers reset",
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": "All circuit breakers opened successfully",
	})
}

// Wind Farm Handlers

// GetWindFarmData returns wind farm data
func (h *Handlers) GetWindFarmData(c *gin.Context) {
	windFarmID := c.Param("id")
	windFarmIDInt, err := strconv.Atoi(windFarmID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wind farm ID"})
		return
	}

	service, err := h.windFarmManager.GetService(windFarmIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	data := service.GetLatestData()

	c.JSON(http.StatusOK, gin.H{
		"data":      data,
		"connected": service.IsConnected(),
		"fcu_online": service.IsFCUOnline(),
	})
}

// GetWindFarmSummary returns aggregated data from all wind farms
func (h *Handlers) GetWindFarmSummary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_active_power":   h.windFarmManager.GetTotalActivePower(),
		"total_reactive_power": h.windFarmManager.GetTotalReactivePower(),
		"total_possible_power": h.windFarmManager.GetTotalPossiblePower(),
		"average_wind_speed":   h.windFarmManager.GetAverageWindSpeed(),
		"service_count":        h.windFarmManager.GetServiceCount(),
		"all_fcus_online":      h.windFarmManager.AreAllFCUsOnline(),
	})
}

// GetWindFarmCommandState returns wind farm command state
func (h *Handlers) GetWindFarmCommandState(c *gin.Context) {
	windFarmID := c.Param("id")
	windFarmIDInt, err := strconv.Atoi(windFarmID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wind farm ID"})
		return
	}

	service, err := h.windFarmManager.GetService(windFarmIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	commandState := service.GetCommandState()

	c.JSON(http.StatusOK, gin.H{
		"command_state": commandState,
	})
}

// StartWindFarm starts a wind farm
func (h *Handlers) StartWindFarm(c *gin.Context) {
	var request struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.windFarmManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.StartWindFarm(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Wind farm start command executed",
		logger.Int("id", request.ID),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": "Wind farm start command sent successfully",
	})
}

// StopWindFarm stops a wind farm
func (h *Handlers) StopWindFarm(c *gin.Context) {
	var request struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.windFarmManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.StopWindFarm(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Wind farm stop command executed",
		logger.Int("id", request.ID),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": "Wind farm stop command sent successfully",
	})
}

// SetWindFarmPowerSetpoint sets the active power setpoint for a wind farm
func (h *Handlers) SetWindFarmPowerSetpoint(c *gin.Context) {
	var request struct {
		ID       int      `json:"id" binding:"required"`
		Setpoint *float32 `json:"setpoint" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.windFarmManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.SetPowerSetpoint(*request.Setpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Wind farm power setpoint set",
		logger.Int("id", request.ID),
		logger.Float32("setpoint", *request.Setpoint),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message":  "Power setpoint set successfully",
		"setpoint": *request.Setpoint,
	})
}

// SetWindFarmReactivePowerSetpoint sets the reactive power setpoint for a wind farm
func (h *Handlers) SetWindFarmReactivePowerSetpoint(c *gin.Context) {
	var request struct {
		ID       int      `json:"id" binding:"required"`
		Setpoint *float32 `json:"setpoint" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.windFarmManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.SetReactivePowerSetpoint(*request.Setpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Wind farm reactive power setpoint set",
		logger.Int("id", request.ID),
		logger.Float32("setpoint", *request.Setpoint),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message":  "Reactive power setpoint set successfully",
		"setpoint": *request.Setpoint,
	})
}

// SetWindFarmPowerFactorSetpoint sets the power factor setpoint for a wind farm
func (h *Handlers) SetWindFarmPowerFactorSetpoint(c *gin.Context) {
	var request struct {
		ID       int      `json:"id" binding:"required"`
		Setpoint *float32 `json:"setpoint" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.windFarmManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.SetPowerFactorSetpoint(*request.Setpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Wind farm power factor setpoint set",
		logger.Int("id", request.ID),
		logger.Float32("setpoint", *request.Setpoint),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message":  "Power factor setpoint set successfully",
		"setpoint": *request.Setpoint,
	})
}

// SetWindFarmRapidDownward sets the rapid downward signal for a wind farm
func (h *Handlers) SetWindFarmRapidDownward(c *gin.Context) {
	var request struct {
		ID int   `json:"id" binding:"required"`
		On *bool `json:"on" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.windFarmManager.GetService(request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := service.SetRapidDownwardSignal(*request.On); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := "deactivated"
	if *request.On {
		status = "activated"
	}

	h.log.Info("Wind farm rapid downward signal set",
		logger.Int("id", request.ID),
		logger.Bool("on", *request.On),
		logger.String("client_ip", c.ClientIP()))

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Rapid downward signal %s successfully", status),
		"on":      *request.On,
	})
}
