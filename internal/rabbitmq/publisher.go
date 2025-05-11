package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func NewPublisher(channel *amqp.Channel, exchange string) *Publisher {
	return &Publisher{
		Channel:  channel,
		Exchange: exchange,
	}
}

func (p *Publisher) PublishGeofenceEvent(input GeofenceEvent) error {
	// Check if the channel is closed
	if p.Channel == nil {
		log.Println("RabbitMQ channel is closed, recreating...")
		// newChannel, err := p.Channel.Connection.Channel
		// if err != nil {
		// 	return fmt.Errorf("failed to recreate RabbitMQ channel: %w", err)
		// }
		// p.Channel = newChannel
	}

	event := GeofenceEvent{
		VehicleID: input.VehicleID,
		Event:     "geofence_entry",
		Timestamp: input.Timestamp,
	}
	event.Location.Latitude = input.Location.Latitude
	event.Location.Longitude = input.Location.Longitude

	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	err = p.Channel.Publish(
		p.Exchange,
		"", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
		return err
	}

	fmt.Println("[Publisher] Success Publish GeofenceEvent, data:", string(body))
	return nil
}
