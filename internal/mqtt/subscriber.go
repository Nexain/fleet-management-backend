package mqtt

import (
	"encoding/json"
	"log"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/yourusername/fleet-management-backend/internal/models"
	"github.com/yourusername/fleet-management-backend/internal/repository"
)

type Subscriber struct {
	repo *repository.LocationRepository
}

func NewSubscriber(repo *repository.LocationRepository) *Subscriber {
	return &Subscriber{repo: repo}
}

func (s *Subscriber) Start() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("fleet_management_subscriber")
	opts.OnConnect = func(client mqtt.Client) {
		if token := client.Subscribe("/fleet/vehicle/#/location", 0, s.messageHandler); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func (s *Subscriber) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var location models.Location
	if err := json.Unmarshal(msg.Payload(), &location); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	if err := s.repo.SaveLocation(location); err != nil {
		log.Printf("Failed to save location: %v", err)
	}
}