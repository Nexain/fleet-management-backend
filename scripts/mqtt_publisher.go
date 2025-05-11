package main

import (
	"encoding/json"
	"math/rand"
	"time"
	"fmt"
	"os"
	"github.com/eclipse/paho.mqtt.golang"
)

type Location struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("mqtt_publisher")

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	defer client.Disconnect(250)

	vehicleID := "B1234XYZ"

	for {
		location := Location{
			VehicleID: vehicleID,
			Latitude:  -6.2088 + (rand.Float64()-0.5)*0.01, // Randomize latitude slightly
			Longitude: 106.8456 + (rand.Float64()-0.5)*0.01, // Randomize longitude slightly
			Timestamp: time.Now().Unix(),
		}

		payload, err := json.Marshal(location)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			continue
		}

		token := client.Publish(fmt.Sprintf("/fleet/vehicle/%s/location", vehicleID), 0, false, payload)
		token.Wait()

		fmt.Printf("Published: %s\n", payload)

		time.Sleep(2 * time.Second)
	}
}