package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

// NewConnection establishes a new connection to RabbitMQ.
func NewConnection() (*amqp.Connection, error) {
	// if amqpURL == "" {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/" // Default to localhost if not set
	}
	// }

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	log.Println("Successfully connected to RabbitMQ")
	return conn, nil
}
