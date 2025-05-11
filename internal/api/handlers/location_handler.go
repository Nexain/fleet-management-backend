package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Nexain/fleet-management-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	locationService service.LocationService
}

func NewLocationHandler(locationService service.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

func (h *LocationHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *LocationHandler) GetLastLocation(c *gin.Context) {
	fmt.Println("[LocationHandler] Incoming GetLastLocation")
	vehicleID := c.Param("vehicle_id")
	ctx := c.Request.Context()

	location, err := h.locationService.GetLastLocation(ctx, vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) GetLocationHistory(c *gin.Context) {
	fmt.Println("[LocationHandler] Incoming GetLocationHistory")
	vehicleID := c.Param("vehicle_id")
	start := c.Query("start")
	end := c.Query("end")

	// Parse start and end times
	startTime, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
		return
	}
	endTime, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
		return
	}

	ctx := c.Request.Context()

	history, err := h.locationService.GetLocationHistory(ctx, vehicleID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
