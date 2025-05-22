package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// TournamentService provides business logic methods for Tournament
type TournamentService struct {
	tournamentRepo *repositories.TournamentRepository
}

// NewTournamentService creates a new TournamentService
func NewTournamentService(tournamentRepo *repositories.TournamentRepository) *TournamentService {
	return &TournamentService{tournamentRepo}
}

// Create creates a new tournament
func (s *TournamentService) Create(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error) {
	if tournament.Name == "" {
		return nil, errors.New("tournament name is required")
	}
	if tournament.StartDate.IsZero() {
		return nil, errors.New("start date is required")
	}
	if tournament.EndDate.Before(tournament.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}
	return s.tournamentRepo.Create(ctx, tournament)
}

// GetByID retrieves a tournament by ID
func (s *TournamentService) GetByID(ctx context.Context, id string) (*models.Tournament, error) {
	return s.tournamentRepo.GetByID(ctx, id)
}

// GetAll retrieves all tournaments with pagination
func (s *TournamentService) GetAll(ctx context.Context, page, limit int64) ([]models.Tournament, error) {
	return s.tournamentRepo.GetAll(ctx, page, limit)
}

// Update updates tournament information
func (s *TournamentService) Update(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error) {
	if tournament.ID.IsZero() {
		return nil, errors.New("invalid tournament ID")
	}
	if tournament.Name == "" {
		return nil, errors.New("tournament name is required")
	}
	if tournament.StartDate.IsZero() {
		return nil, errors.New("start date is required")
	}
	if tournament.EndDate.Before(tournament.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}
	return s.tournamentRepo.Update(ctx, tournament)
}

// Delete deletes a tournament by ID
func (s *TournamentService) Delete(ctx context.Context, id string) error {
	return s.tournamentRepo.Delete(ctx, id)
}