package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// TeamMemberService provides business logic methods for TeamMember
type TeamMemberService struct {
	teamMemberRepo *repositories.TeamMemberRepository
}

// NewTeamMemberService creates a new TeamMemberService
func NewTeamMemberService(teamMemberRepo *repositories.TeamMemberRepository) *TeamMemberService {
	return &TeamMemberService{teamMemberRepo}
}

// Create creates a new team member
func (s *TeamMemberService) Create(ctx context.Context, teamMember *models.TeamMember) (*models.TeamMember, error) {
	if teamMember.TeamID.IsZero() {
		return nil, errors.New("team ID is required")
	}
	if teamMember.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.teamMemberRepo.Create(ctx, teamMember)
}

// GetByID retrieves a team member by ID
func (s *TeamMemberService) GetByID(ctx context.Context, id string) (*models.TeamMember, error) {
	return s.teamMemberRepo.GetByID(ctx, id)
}

// GetByTeamID retrieves team members by team ID
func (s *TeamMemberService) GetByTeamID(ctx context.Context, teamID string) ([]models.TeamMember, error) {
	return s.teamMemberRepo.GetByTeamID(ctx, teamID)
}

// GetByUserID retrieves team members by user ID
func (s *TeamMemberService) GetByUserID(ctx context.Context, userID string) ([]models.TeamMember, error) {
	return s.teamMemberRepo.GetByUserID(ctx, userID)
}

// GetAll retrieves all team members with pagination
func (s *TeamMemberService) GetAll(ctx context.Context, page, limit int64) ([]models.TeamMember, error) {
	return s.teamMemberRepo.GetAll(ctx, page, limit)
}

// Update updates team member information
func (s *TeamMemberService) Update(ctx context.Context, teamMember *models.TeamMember) (*models.TeamMember, error) {
	if teamMember.ID.IsZero() {
		return nil, errors.New("invalid team member ID")
	}
	if teamMember.TeamID.IsZero() {
		return nil, errors.New("team ID is required")
	}
	if teamMember.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.teamMemberRepo.Update(ctx, teamMember)
}

// Delete deletes a team member by ID
func (s *TeamMemberService) Delete(ctx context.Context, id string) error {
	return s.teamMemberRepo.Delete(ctx, id)
}