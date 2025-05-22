package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// InjuryService provides business logic methods for Injury
type InjuryService struct {
	injuryRepo *repositories.InjuryRepository
}

// NewInjuryService creates a new InjuryService
func NewInjuryService(injuryRepo *repositories.InjuryRepository) *InjuryService {
	return &InjuryService{injuryRepo}
}

// Create creates a new injury record
func (s *InjuryService) Create(ctx context.Context, injury *models.Injury) (*models.Injury, error) {
	if injury.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.injuryRepo.Create(ctx, injury)
}

// GetByID retrieves an injury by ID
func (s *InjuryService) GetByID(ctx context.Context, id string) (*models.Injury, error) {
	return s.injuryRepo.GetByID(ctx, id)
}

// GetByUserID retrieves an injury by user ID
func (s *InjuryService) GetByUserID(ctx context.Context, userID string) (*models.Injury, error) {
	return s.injuryRepo.GetByUserID(ctx, userID)
}

// GetAll retrieves all injury records with pagination
func (s *InjuryService) GetAll(ctx context.Context, page, limit int64) ([]models.Injury, error) {
	return s.injuryRepo.GetAll(ctx, page, limit)
}

// Update updates injury information
func (s *InjuryService) Update(ctx context.Context, injury *models.Injury) (*models.Injury, error) {
	if injury.ID.IsZero() {
		return nil, errors.New("invalid injury ID")
	}
	if injury.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.injuryRepo.Update(ctx, injury)
}

// Delete deletes an injury by ID
func (s *InjuryService) Delete(ctx context.Context, id string) error {
	return s.injuryRepo.Delete(ctx, id)
}