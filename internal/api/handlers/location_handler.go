package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fleet-management-backend/internal/service"
)

type LocationHandler struct {
	locationService service.LocationService
}

func NewLocationHandler(locationService service.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

func (h *LocationHandler) GetLastLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	location, err := h.locationService.GetLastLocation(vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

func (h *LocationHandler) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	start := c.Query("start")
	end := c.Query("end")
	history, err := h.locationService.GetLocationHistory(vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}