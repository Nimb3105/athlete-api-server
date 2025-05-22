package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// TeamService provides business logic methods for Team
type TeamService struct {
	teamRepo *repositories.TeamRepository
}

// NewTeamService creates a new TeamService
func NewTeamService(teamRepo *repositories.TeamRepository) *TeamService {
	return &TeamService{teamRepo}
}

// Create creates a new team
func (s *TeamService) Create(ctx context.Context, team *models.Team) (*models.Team, error) {
	if team.SportID.IsZero() {
		return nil, errors.New("sport ID is required")
	}
	if team.CreatedBy.IsZero() {
		return nil, errors.New("created by user ID is required")
	}
	if team.Name == "" {
		return nil, errors.New("team name is required")
	}
	return s.teamRepo.Create(ctx, team)
}

// GetByID retrieves a team by ID
func (s *TeamService) GetByID(ctx context.Context, id string) (*models.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}

// GetBySportID retrieves teams by sport ID
func (s *TeamService) GetBySportID(ctx context.Context, sportID string) ([]models.Team, error) {
	return s.teamRepo.GetBySportID(ctx, sportID)
}

// GetAll retrieves all teams with pagination
func (s *TeamService) GetAll(ctx context.Context, page, limit int64) ([]models.Team, error) {
	return s.teamRepo.GetAll(ctx, page, limit)
}

// Update updates team information
func (s *TeamService) Update(ctx context.Context, team *models.Team) (*models.Team, error) {
	if team.ID.IsZero() {
		return nil, errors.New("invalid team ID")
	}
	if team.Name == "" {
		return nil, errors.New("team name is required")
	}
	return s.teamRepo.Update(ctx, team)
}

// Delete deletes a team by ID
func (s *TeamService) Delete(ctx context.Context, id string) error {
	return s.teamRepo.Delete(ctx, id)
}