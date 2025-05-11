package service

import (
	"fmt"
	"log"
	"math"

	"github.com/Nexain/fleet-management-backend/internal/rabbitmq"
)

type GeofenceService struct {
	radius    float64
	center    Location
	publisher *rabbitmq.Publisher
}

type Location struct {
	Latitude  float64
	Longitude float64
}

func NewGeofenceService(radius float64, center Location, publisher *rabbitmq.Publisher) *GeofenceService {
	return &GeofenceService{
		radius:    radius,
		center:    center,
		publisher: publisher,
	}
}

// func (g *GeofenceService) CheckGeofence(vehicleLat, vehicleLon float64) bool {
// 	distance := g.calculateDistance(g.Lat, g.Lon, vehicleLat, vehicleLon)
// 	log.Printf("Distance to geofence: %.2f meters", distance)
// 	return distance <= g.Radius
// }

func (g *GeofenceService) IsInsideGeofence(location Location) bool {
	distance := g.calculateDistance(g.center, location)
	log.Printf("Distance to geofence: %.2f meters", distance)
	return distance <= g.radius
}

func (g *GeofenceService) calculateDistance(loc1, loc2 Location) float64 {
	const earthRadius = 6371000 // meters

	lat1 := toRadians(loc1.Latitude)
	lat2 := toRadians(loc2.Latitude)
	deltaLat := toRadians(loc2.Latitude - loc1.Latitude)
	deltaLon := toRadians(loc2.Longitude - loc1.Longitude)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (g *GeofenceService) triggerGeofenceEvent(vehicleID string, location Location, timestamp int64) {
	message := rabbitmq.GeofenceEvent{
		VehicleID: vehicleID,
		Event:     "geofence_entry",
		Location: rabbitmq.Location{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
		Timestamp: timestamp,
	}

	err := g.publisher.PublishGeofenceEvent(message)
	if err != nil {
		fmt.Printf("Failed to publish geofence event: %v\n", err)
	}
}

func toRadians(degree float64) float64 {
	return degree * math.Pi / 180
}
