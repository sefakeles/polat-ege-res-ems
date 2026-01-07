package api

import (
	"net/http"
	"powerkonnekt/ems/internal/fcr"
	"powerkonnekt/ems/pkg/logger"

	"github.com/gin-gonic/gin"
)

// FCRNHandlers handles FCR-N related API requests
type FCRNHandlers struct {
	service *fcr.Service
	log     logger.Logger
}

// NewFCRNHandlers creates new FCR-N API handlers
func NewFCRNHandlers(service *fcr.Service) *FCRNHandlers {
	return &FCRNHandlers{
		service: service,
		log: logger.With(
			logger.String("component", "fcrn_api"),
		),
	}
}

// RegisterRoutes registers FCR-N routes
func (h *FCRNHandlers) RegisterRoutes(router *gin.RouterGroup) {
	fcrn := router.Group("/fcrn")
	{
		fcrn.GET("/status", h.GetStatus)
		fcrn.GET("/state", h.GetState)
		fcrn.POST("/activate", h.Activate)
		fcrn.POST("/deactivate", h.Deactivate)
		fcrn.POST("/capacity", h.SetCapacity)
		fcrn.POST("/droop", h.SetDroop)
		fcrn.GET("/maintained-capacity", h.GetMaintainedCapacity)

		// Test endpoints
		test := fcrn.Group("/test")
		{
			test.POST("/frequency", h.SetTestFrequency)
		}
	}
}

// GetStatus returns FCR-N service status
// @Summary Get FCR-N status
// @Description Get current FCR-N service status
// @Tags FCR-N
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/status [get]
func (h *FCRNHandlers) GetStatus(c *gin.Context) {
	if err := h.service.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, ErrorResponse{
			Error:   "FCR-N service unavailable",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "FCR-N service is running",
	})
}

// GetState returns current FCR-N state
// @Summary Get FCR-N state
// @Description Get detailed FCR-N controller state
// @Tags FCR-N
// @Accept json
// @Produce json
// @Success 200 {object} fcr.FCRNState
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/state [get]
func (h *FCRNHandlers) GetState(c *gin.Context) {
	state := h.service.GetState()
	c.JSON(http.StatusOK, state)
}

// ActivateRequest represents FCR-N activation request
type ActivateRequest struct {
	// Optional parameters can be added here
}

// Activate activates FCR-N provision
// @Summary Activate FCR-N
// @Description Activate FCR-N provision with smooth transition
// @Tags FCR-N
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/activate [post]
func (h *FCRNHandlers) Activate(c *gin.Context) {
	h.log.Info("FCR-N activation requested")

	if err := h.service.ActivateFCRN(); err != nil {
		h.log.Error("Failed to activate FCR-N", logger.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to activate FCR-N",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("FCR-N activated successfully")
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "FCR-N activated successfully",
	})
}

// Deactivate deactivates FCR-N provision
// @Summary Deactivate FCR-N
// @Description Deactivate FCR-N provision with smooth transition
// @Tags FCR-N
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/deactivate [post]
func (h *FCRNHandlers) Deactivate(c *gin.Context) {
	h.log.Info("FCR-N deactivation requested")

	if err := h.service.DeactivateFCRN(); err != nil {
		h.log.Error("Failed to deactivate FCR-N", logger.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to deactivate FCR-N",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("FCR-N deactivated successfully")
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "FCR-N deactivated successfully",
	})
}

// SetCapacityRequest represents capacity update request
type SetCapacityRequest struct {
	Capacity float64 `json:"capacity" binding:"required,gt=0"` // kW
}

// SetCapacity updates FCR-N capacity
// @Summary Set FCR-N capacity
// @Description Update the FCR-N sold capacity
// @Tags FCR-N
// @Accept json
// @Produce json
// @Param request body SetCapacityRequest true "Capacity update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/capacity [put]
func (h *FCRNHandlers) SetCapacity(c *gin.Context) {
	var req SetCapacityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("FCR-N capacity update requested",
		logger.Float64("new_capacity", req.Capacity))

	if err := h.service.SetCapacity(req.Capacity); err != nil {
		h.log.Error("Failed to set FCR-N capacity", logger.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to set capacity",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("FCR-N capacity updated successfully",
		logger.Float64("capacity", req.Capacity))
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Capacity updated successfully",
		Data: map[string]interface{}{
			"capacity": req.Capacity,
		},
	})
}

// SetDroopRequest represents droop update request
type SetDroopRequest struct {
	Droop float64 `json:"droop" binding:"required,gt=0,lte=100"` // %
}

// SetDroop updates FCR-N droop setting
// @Summary Set FCR-N droop
// @Description Update the FCR-N droop setting
// @Tags FCR-N
// @Accept json
// @Produce json
// @Param request body SetDroopRequest true "Droop update request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/droop [put]
func (h *FCRNHandlers) SetDroop(c *gin.Context) {
	var req SetDroopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("FCR-N droop update requested",
		logger.Float64("new_droop", req.Droop))

	if err := h.service.SetDroop(req.Droop); err != nil {
		h.log.Error("Failed to set FCR-N droop", logger.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to set droop",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("FCR-N droop updated successfully",
		logger.Float64("droop", req.Droop))
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Droop updated successfully",
		Data: map[string]interface{}{
			"droop": req.Droop,
		},
	})
}

// GetMaintainedCapacity returns the maintained (available) capacity
// @Summary Get maintained capacity
// @Description Get current FCR-N maintained (available) capacity
// @Tags FCR-N
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/maintained-capacity [get]
func (h *FCRNHandlers) GetMaintainedCapacity(c *gin.Context) {
	capacity := h.service.GetMaintainedCapacity()

	c.JSON(http.StatusOK, gin.H{
		"maintained_capacity": capacity,
		"unit":                "kW",
	})
}

// SetTestFrequencyRequest represents test frequency update request
type SetTestFrequencyRequest struct {
	Frequency float64 `json:"frequency" binding:"required,gte=49.0,lte=51.0"` // Hz
}

// SetTestFrequency sets the test frequency for testing purposes
// @Summary Set test frequency
// @Description Set frequency value for test frequency source (testing only)
// @Tags FCR-N
// @Accept json
// @Produce json
// @Param request body SetTestFrequencyRequest true "Test frequency request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/fcrn/test/frequency [put]
func (h *FCRNHandlers) SetTestFrequency(c *gin.Context) {
	var req SetTestFrequencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid test frequency request", logger.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	if err := h.service.SetTestFrequency(req.Frequency); err != nil {
		h.log.Error("Failed to set test frequency", logger.Err(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to set test frequency",
			Message: err.Error(),
		})
		return
	}

	h.log.Info("Test frequency updated successfully",
		logger.Float64("frequency", req.Frequency))
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Test frequency updated successfully",
		Data: map[string]interface{}{
			"frequency": req.Frequency,
		},
	})
}

// Common response structures
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
