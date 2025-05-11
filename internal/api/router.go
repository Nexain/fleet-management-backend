package api

import (
    "github.com/gin-gonic/gin"
    "fleet-management-backend/internal/api/handlers"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    locationHandler := handlers.LocationHandler{}

    router.GET("/vehicles/:vehicle_id/location", locationHandler.GetLastLocation)
    router.GET("/vehicles/:vehicle_id/history", locationHandler.GetLocationHistory)

    return router
}