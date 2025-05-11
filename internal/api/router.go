package api

import (
	"github.com/Nexain/fleet-management-backend/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(locationHandler handlers.LocationHandler) *gin.Engine {
	router := gin.Default()

	// locationHandler := handlers.LocationHandler{}

	router.GET("/ping", locationHandler.Ping)
	router.GET("/vehicles/:vehicle_id/location", locationHandler.GetLastLocation)
	router.GET("/vehicles/:vehicle_id/history", locationHandler.GetLocationHistory)

	return router
}
