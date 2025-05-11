package main

import (
    "log"
    "net/http"

    "fleet-management-backend/internal/api"
    "fleet-management-backend/internal/mqtt"
    "fleet-management-backend/internal/rabbitmq"
    "fleet-management-backend/internal/repository"
    "fleet-management-backend/internal/service"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize the router
    router := gin.Default()

    // Set up the database repository
    locationRepo := repository.NewLocationRepository()
    
    // Initialize services
    locationService := service.NewLocationService(locationRepo)
    geofenceService := service.NewGeofenceService(locationRepo)

    // Set up the API handlers
    locationHandler := api.NewLocationHandler(locationService)

    // Set up routes
    api.SetupRoutes(router, locationHandler)

    // Start MQTT subscriber
    mqtt.StartSubscriber()

    // Start RabbitMQ worker
    rabbitmq.StartWorker()

    // Start the server
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}