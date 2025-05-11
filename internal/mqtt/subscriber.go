package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Nexain/fleet-management-backend/internal/models"
	"github.com/Nexain/fleet-management-backend/internal/rabbitmq"
	"github.com/Nexain/fleet-management-backend/internal/repository"
	"github.com/Nexain/fleet-management-backend/internal/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Subscriber struct {
	repo              *repository.LocationRepository
	rabbitmqPublisher *rabbitmq.Publisher
	ctx               context.Context
}

func NewSubscriber(ctx context.Context, repo *repository.LocationRepository, rabbitmq *rabbitmq.Publisher) *Subscriber {
	return &Subscriber{repo: repo, ctx: ctx, rabbitmqPublisher: rabbitmq}
}

func (s *Subscriber) Start() {
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "localhost" // Default to localhost if the environment variable is not set
	}

	opts := mqtt.NewClientOptions().AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID("fleet_management_subscriber")
	opts.OnConnect = func(client mqtt.Client) {
		if token := client.Subscribe("/fleet/vehicle/+/location", 0, s.messageHandler); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

	fmt.Println("Connecting to MQTT broker...")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func (s *Subscriber) messageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Println("[MQTT Subscriber] Message received")
	var location models.Location

	if err := json.Unmarshal(msg.Payload(), &location); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// Check if the vehicle is within the geofence
	geofenceService := service.NewGeofenceService(50, service.Location{
		Latitude:  -6.2088,  // Example geofence center latitude
		Longitude: 106.8456, // Example geofence center longitude
	}, s.rabbitmqPublisher) // 50 meters radius
	if geofenceService.IsInsideGeofence(service.Location{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}) {
		log.Printf("Vehicle %s entered the geofence", location.VehicleID)

		// Publish geofence event to RabbitMQ
		s.publishGeofenceEvent(rabbitmq.GeofenceEvent{
			VehicleID: location.VehicleID,
			Event:     "geofence_alert",
			Location: rabbitmq.Location{
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
			},
			Timestamp: location.Timestamp,
		})
	}

	if err := s.repo.SaveLocation(s.ctx, &location); err != nil {
		log.Printf("Failed to save location: %v", err)
	}
}

func (s *Subscriber) publishGeofenceEvent(event rabbitmq.GeofenceEvent) {
	err := s.rabbitmqPublisher.PublishGeofenceEvent(event)
	if err != nil {
		log.Printf("Failed to publish geofence event: %v", err)
	}
}
