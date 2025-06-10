
package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// ProgressService provides business logic methods for Progress
type ProgressService struct {
	progressRepo *repositories.ProgressRepository
}

// NewProgressService creates a new ProgressService
func NewProgressService(progressRepo *repositories.ProgressRepository) *ProgressService {
	return &ProgressService{progressRepo}
}

// Create creates a new progress record
func (s *ProgressService) Create(ctx context.Context, progress *models.Progress) (*models.Progress, error) {
	if progress.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if progress.MetricType == "" {
		return nil, errors.New("metric type is required")
	}
	if progress.Value <= 0 {
		return nil, errors.New("value must be greater than zero")
	}
	return s.progressRepo.Create(ctx, progress)
}

// GetByID retrieves a progress record by ID
func (s *ProgressService) GetByID(ctx context.Context, id string) (*models.Progress, error) {
	return s.progressRepo.GetByID(ctx, id)
}

// GetByUserID retrieves progress records by user ID
func (s *ProgressService) GetByUserID(ctx context.Context, userID string) ([]models.Progress, error) {
	return s.progressRepo.GetByUserID(ctx, userID)
}

// GetAll retrieves all progress records with pagination
func (s *ProgressService) GetAll(ctx context.Context, page, limit int64) ([]models.Progress, error) {
	return s.progressRepo.GetAll(ctx, page, limit)
}

// Update updates progress record information
func (s *ProgressService) Update(ctx context.Context, progress *models.Progress) (*models.Progress, error) {
	if progress.ID.IsZero() {
		return nil, errors.New("invalid progress ID")
	}
	if progress.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if progress.MetricType == "" {
		return nil, errors.New("metric type is required")
	}
	if progress.Value <= 0 {
		return nil, errors.New("value must be greater than zero")
	}
	return s.progressRepo.Update(ctx, progress)
}

// Delete deletes a progress record by ID
func (s *ProgressService) Delete(ctx context.Context, id string) error {
	return s.progressRepo.Delete(ctx, id)
}