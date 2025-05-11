package rabbitmq

import (
    "encoding/json"
    "log"

    "github.com/streadway/amqp"
)

type Publisher struct {
    Channel *amqp.Channel
    Exchange string
}

type GeofenceEvent struct {
    VehicleID string `json:"vehicle_id"`
    Event     string `json:"event"`
    Location  struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    } `json:"location"`
    Timestamp int64 `json:"timestamp"`
}

func NewPublisher(channel *amqp.Channel, exchange string) *Publisher {
    return &Publisher{
        Channel: channel,
        Exchange: exchange,
    }
}

func (p *Publisher) PublishGeofenceEvent(vehicleID string, latitude, longitude float64, timestamp int64) error {
    event := GeofenceEvent{
        VehicleID: vehicleID,
        Event:     "geofence_entry",
        Timestamp: timestamp,
    }
    event.Location.Latitude = latitude
    event.Location.Longitude = longitude

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

    return nil
}