package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"fleet-management-backend/internal/service"
)

type HistoryHandler struct {
	locationService *service.LocationService
}

func NewHistoryHandler(locationService *service.LocationService) *HistoryHandler {
	return &HistoryHandler{locationService: locationService}
}

func (h *HistoryHandler) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	start := c.Query("start")
	end := c.Query("end")

	startTime, err := time.ParseInt(start, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
		return
	}

	endTime, err := time.ParseInt(end, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
		return
	}

	history, err := h.locationService.GetLocationHistory(vehicleID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve location history"})
		return
	}

	c.JSON(http.StatusOK, history)
}