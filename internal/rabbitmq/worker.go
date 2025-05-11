package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func StartWorker(conn *amqp.Connection) {
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"geofence_alerts",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	log.Println("Waiting for messages. To exit press CTRL+C")

	for msg := range msgs {
		var event GeofenceEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Error decoding message: %s", err)
			continue
		}
		processGeofenceEvent(event)
	}
}

func processGeofenceEvent(event GeofenceEvent) {
	log.Printf("Received geofence event: %+v", event)
}
