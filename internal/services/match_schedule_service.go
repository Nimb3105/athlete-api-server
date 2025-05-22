package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// MatchScheduleService provides business logic methods for MatchSchedule
type MatchScheduleService struct {
	matchScheduleRepo *repositories.MatchScheduleRepository
}

// NewMatchScheduleService creates a new MatchScheduleService
func NewMatchScheduleService(matchScheduleRepo *repositories.MatchScheduleRepository) *MatchScheduleService {
	return &MatchScheduleService{matchScheduleRepo}
}

// Create creates a new match schedule
func (s *MatchScheduleService) Create(ctx context.Context, matchSchedule *models.MatchSchedule) (*models.MatchSchedule, error) {
	if matchSchedule.TournamentID.IsZero() {
		return nil, errors.New("tournament ID is required")
	}
	return s.matchScheduleRepo.Create(ctx, matchSchedule)
}

// GetByID retrieves a match schedule by ID
func (s *MatchScheduleService) GetByID(ctx context.Context, id string) (*models.MatchSchedule, error) {
	return s.matchScheduleRepo.GetByID(ctx, id)
}

// GetByTournamentID retrieves a match schedule by tournament ID
func (s *MatchScheduleService) GetByTournamentID(ctx context.Context, tournamentID string) (*models.MatchSchedule, error) {
	return s.matchScheduleRepo.GetByTournamentID(ctx, tournamentID)
}

// GetAll retrieves all match schedules with pagination
func (s *MatchScheduleService) GetAll(ctx context.Context, page, limit int64) ([]models.MatchSchedule, error) {
	return s.matchScheduleRepo.GetAll(ctx, page, limit)
}

// Update updates match schedule information
func (s *MatchScheduleService) Update(ctx context.Context, matchSchedule *models.MatchSchedule) (*models.MatchSchedule, error) {
	if matchSchedule.ID.IsZero() {
		return nil, errors.New("invalid match schedule ID")
	}
	if matchSchedule.TournamentID.IsZero() {
		return nil, errors.New("tournament ID is required")
	}
	return s.matchScheduleRepo.Update(ctx, matchSchedule)
}

// Delete deletes a match schedule by ID
func (s *MatchScheduleService) Delete(ctx context.Context, id string) error {
	return s.matchScheduleRepo.Delete(ctx, id)
}