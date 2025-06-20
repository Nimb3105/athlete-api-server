package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// CoachAthleteService provides business logic for CoachAthlete
type CoachAthleteService struct {
	coachAthleteRepo *repositories.CoachAthleteRepository
}

// NewCoachAthleteService creates a new CoachAthleteService
func NewCoachAthleteService(coachAthleteRepo *repositories.CoachAthleteRepository) *CoachAthleteService {
	return &CoachAthleteService{coachAthleteRepo}
}

// Create creates a new coach-athlete relationship
func (s *CoachAthleteService) Create(ctx context.Context, coachAthlete *models.CoachAthlete) (*models.CoachAthlete, error) {
	if coachAthlete.CoachID.IsZero() || coachAthlete.AthleteID.IsZero() {
		return nil, errors.New("coach ID and athlete ID are required")
	}
	return s.coachAthleteRepo.Create(ctx, coachAthlete)
}

// GetByID retrieves a coach-athlete relationship by ID
func (s *CoachAthleteService) GetByID(ctx context.Context, id string) (*models.CoachAthlete, error) {
	return s.coachAthleteRepo.GetByID(ctx, id)
}

// GetByAthleteId retrieves a coach-athlete relationship by athlete ID
func (s *CoachAthleteService) GetByAthleteId(ctx context.Context, athleteId string) (*models.CoachAthlete, error) {
	if athleteId == "" {
		return nil, errors.New("athlete ID is required")
	}
	return s.coachAthleteRepo.GetByAthleteId(ctx, athleteId)
}

func (s *CoachAthleteService) GetAllByCoachId(ctx context.Context, id string, page,limit int64) ([]models.CoachAthlete,int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
	return s.coachAthleteRepo.GetAllByCoachId(ctx, id,page,limit)
}

// GetAll retrieves all coach-athlete relationships with pagination
func (s *CoachAthleteService) GetAll(ctx context.Context, page, limit int64) ([]models.CoachAthlete, error) {
	
	return s.coachAthleteRepo.GetAll(ctx, page, limit)
}

// Update updates a coach-athlete relationship
func (s *CoachAthleteService) Update(ctx context.Context, coachAthlete *models.CoachAthlete) (*models.CoachAthlete, error) {
	if coachAthlete.ID.IsZero() {
		return nil, errors.New("invalid coach-athlete ID")
	}
	return s.coachAthleteRepo.Update(ctx, coachAthlete)
}

// Delete deletes a coach-athlete relationship by ID
func (s *CoachAthleteService) Delete(ctx context.Context, id string) error {
	return s.coachAthleteRepo.Delete(ctx, id)
}
