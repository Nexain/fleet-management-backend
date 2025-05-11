package service

import (
	"context"
	"errors"

	"github.com/Nexain/fleet-management-backend/internal/models"
	repository "github.com/Nexain/fleet-management-backend/internal/repository"
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
	// return nil
	return s.repo.SaveLocation(ctx, location)
}

func (s *LocationService) GetLastLocation(ctx context.Context, vehicleID string) (*models.Location, error) {
	return s.repo.GetLastLocation(ctx, vehicleID)
}

func (s *LocationService) GetLocationHistory(ctx context.Context, vehicleID string, startTime, endTime int64) ([]models.Location, error) {
	return s.repo.GetLocationHistory(ctx, vehicleID, startTime, endTime)
}
