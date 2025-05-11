package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Nexain/fleet-management-backend/internal/api"
	"github.com/Nexain/fleet-management-backend/internal/api/handlers"
	"github.com/Nexain/fleet-management-backend/internal/mqtt"
	"github.com/Nexain/fleet-management-backend/internal/rabbitmq"
	"github.com/Nexain/fleet-management-backend/internal/repository"
	"github.com/Nexain/fleet-management-backend/internal/service"
)

func main() {
	// Set up the database connection
	fmt.Println("Connecting to database...")
	dbDsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbDsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	fmt.Println("Connected to database")

	// Set up the database repository
	locationRepo := repository.NewLocationRepository(db)

	// Initialize services
	fmt.Println("Setting up routes...")
	locationService := service.NewLocationService(*locationRepo)

	// Set up the API handlers
	locationHandler := handlers.NewLocationHandler(*locationService)
	router := api.SetupRouter(*locationHandler)
	fmt.Println("Routes set up")

	// Start RabbitMQ worker
	time.Sleep(10 * time.Second) // Wait for RabbitMQ to be ready
	fmt.Println("Starting RabbitMQ worker connection...")
	workerConn, err := rabbitmq.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ for worker: %v", err)
	}
	defer workerConn.Close()

	go func() {
		fmt.Println("Starting RabbitMQ worker...")
		rabbitmq.StartWorker(workerConn)
		fmt.Println("RabbitMQ worker started")
	}()

	// Create a separate RabbitMQ connection for the publisher
	fmt.Println("Starting RabbitMQ publisher connection...")
	publisherConn, err := rabbitmq.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ for publisher: %v", err)
	}
	defer publisherConn.Close()

	channel, err := publisherConn.Channel()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ channel: %v", err)
	}
	defer channel.Close()

	rabbitMQPub := rabbitmq.NewPublisher(channel, "fleet.events")
	fmt.Println("Connected to RabbitMQ for publishing")

	// Start MQTT subscriber
	fmt.Println("Starting MQTT subscriber...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mqttSubscriber := mqtt.NewSubscriber(ctx, locationRepo, rabbitMQPub)
	mqttSubscriber.Start()
	fmt.Println("MQTT subscriber started")

	// Start the server
	fmt.Println("Starting server...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
	fmt.Println("Server started on :8080")
}
