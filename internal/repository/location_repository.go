package repository

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
    "fleet-management-backend/internal/models"
)

type LocationRepository struct {
    db *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
    return &LocationRepository{db: db}
}

func (r *LocationRepository) SaveLocation(location *models.Location) error {
    query := `
        INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp)
        VALUES ($1, $2, $3, $4)
    `
    _, err := r.db.Exec(query, location.VehicleID, location.Latitude, location.Longitude, location.Timestamp)
    if err != nil {
        log.Printf("Error saving location: %v", err)
        return fmt.Errorf("could not save location: %w", err)
    }
    return nil
}

func (r *LocationRepository) GetLastLocation(vehicleID string) (*models.Location, error) {
    query := `
        SELECT vehicle_id, latitude, longitude, timestamp
        FROM vehicle_locations
        WHERE vehicle_id = $1
        ORDER BY timestamp DESC
        LIMIT 1
    `
    location := &models.Location{}
    err := r.db.QueryRow(query, vehicleID).Scan(&location.VehicleID, &location.Latitude, &location.Longitude, &location.Timestamp)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        log.Printf("Error getting last location: %v", err)
        return nil, fmt.Errorf("could not get last location: %w", err)
    }
    return location, nil
}

func (r *LocationRepository) GetLocationHistory(vehicleID string, start, end int64) ([]models.Location, error) {
    query := `
        SELECT vehicle_id, latitude, longitude, timestamp
        FROM vehicle_locations
        WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3
        ORDER BY timestamp
    `
    rows, err := r.db.Query(query, vehicleID, start, end)
    if err != nil {
        log.Printf("Error getting location history: %v", err)
        return nil, fmt.Errorf("could not get location history: %w", err)
    }
    defer rows.Close()

    var locations []models.Location
    for rows.Next() {
        var location models.Location
        if err := rows.Scan(&location.VehicleID, &location.Latitude, &location.Longitude, &location.Timestamp); err != nil {
            log.Printf("Error scanning location: %v", err)
            return nil, fmt.Errorf("could not scan location: %w", err)
        }
        locations = append(locations, location)
    }
    return locations, nil
}