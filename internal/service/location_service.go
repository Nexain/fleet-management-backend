package service

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/fleet-management-backend/internal/models"
	"github.com/yourusername/fleet-management-backend/internal/repository"
)

type LocationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) SaveLocation(ctx context.Context, location *models.Location) error {
	if location == nil {
		return errors.New("location cannot be nil")
	}
	return s.repo.Save(ctx, location)
}

func (s *LocationService) GetLastLocation(ctx context.Context, vehicleID string) (*models.Location, error) {
	return s.repo.FindLastByVehicleID(ctx, vehicleID)
}

func (s *LocationService) GetLocationHistory(ctx context.Context, vehicleID string, startTime, endTime time.Time) ([]models.Location, error) {
	return s.repo.FindByVehicleIDAndTimeRange(ctx, vehicleID, startTime, endTime)
}