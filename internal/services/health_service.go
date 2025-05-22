package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// HealthService provides business logic methods for Health
type HealthService struct {
	healthRepo *repositories.HealthRepository
}

// NewHealthService creates a new HealthService
func NewHealthService(healthRepo *repositories.HealthRepository) *HealthService {
	return &HealthService{healthRepo}
}

// Create creates a new health record
func (s *HealthService) Create(ctx context.Context, health *models.Health) (*models.Health, error) {
	if health.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.healthRepo.Create(ctx, health)
}

// GetByID retrieves a health record by ID
func (s *HealthService) GetByID(ctx context.Context, id string) (*models.Health, error) {
	return s.healthRepo.GetByID(ctx, id)
}

// GetByUserID retrieves a health record by user ID
func (s *HealthService) GetByUserID(ctx context.Context, userID string) (*models.Health, error) {
	return s.healthRepo.GetByUserID(ctx, userID)
}

// GetAll retrieves all health records with pagination
func (s *HealthService) GetAll(ctx context.Context, page, limit int64) ([]models.Health, error) {
	return s.healthRepo.GetAll(ctx, page, limit)
}

// Update updates health record information
func (s *HealthService) Update(ctx context.Context, health *models.Health) (*models.Health, error) {
	if health.ID.IsZero() {
		return nil, errors.New("invalid health record ID")
	}
	if health.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.healthRepo.Update(ctx, health)
}

// Delete deletes a health record by ID
func (s *HealthService) Delete(ctx context.Context, id string) error {
	return s.healthRepo.Delete(ctx, id)
}