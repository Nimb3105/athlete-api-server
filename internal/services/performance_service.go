package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// PerformanceService provides business logic methods for Performance
type PerformanceService struct {
	performanceRepo *repositories.PerformanceRepository
}

// NewPerformanceService creates a new PerformanceService
func NewPerformanceService(performanceRepo *repositories.PerformanceRepository) *PerformanceService {
	return &PerformanceService{performanceRepo}
}

// Create creates a new performance record
func (s *PerformanceService) Create(ctx context.Context, performance *models.Performance) (*models.Performance, error) {
	if performance.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if performance.ScheduleID.IsZero() {
		return nil, errors.New("schedule ID is required")
	}
	if performance.MetricType == "" {
		return nil, errors.New("metric type is required")
	}
	if performance.Value <= 0 {
		return nil, errors.New("value must be greater than zero")
	}
	return s.performanceRepo.Create(ctx, performance)
}

// GetByID retrieves a performance record by ID
func (s *PerformanceService) GetByID(ctx context.Context, id string) (*models.Performance, error) {
	return s.performanceRepo.GetByID(ctx, id)
}

// GetByUserID retrieves performance records by user ID
func (s *PerformanceService) GetByUserID(ctx context.Context, userID string) ([]models.Performance, error) {
	return s.performanceRepo.GetByUserID(ctx, userID)
}

// GetAll retrieves all performance records with pagination
func (s *PerformanceService) GetAll(ctx context.Context, page, limit int64) ([]models.Performance, error) {
	return s.performanceRepo.GetAll(ctx, page, limit)
}

// Update updates performance record information
func (s *PerformanceService) Update(ctx context.Context, performance *models.Performance) (*models.Performance, error) {
	if performance.ID.IsZero() {
		return nil, errors.New("invalid performance ID")
	}
	if performance.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if performance.ScheduleID.IsZero() {
		return nil, errors.New("schedule ID is required")
	}
	if performance.MetricType == "" {
		return nil, errors.New("metric type is required")
	}
	if performance.Value <= 0 {
		return nil, errors.New("value must be greater than zero")
	}
	return s.performanceRepo.Update(ctx, performance)
}

// Delete deletes a performance record by ID
func (s *PerformanceService) Delete(ctx context.Context, id string) error {
	return s.performanceRepo.Delete(ctx, id)
}