package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// AthleteMatchService cung cấp các phương thức nghiệp vụ cho AthleteMatch
type AthleteMatchService struct {
	athleteMatchRepo *repositories.AthleteMatchRepository
}

// NewAthleteMatchService tạo một AthleteMatchService mới
func NewAthleteMatchService(athleteMatchRepo *repositories.AthleteMatchRepository) *AthleteMatchService {
	return &AthleteMatchService{athleteMatchRepo}
}

// Create tạo một athlete match mới
func (s *AthleteMatchService) Create(ctx context.Context, athleteMatch *models.AthleteMatch) (*models.AthleteMatch, error) {
	if athleteMatch.UserID.IsZero() || athleteMatch.MatchID.IsZero() {
		return nil, errors.New("user ID and match ID are required")
	}
	return s.athleteMatchRepo.Create(ctx, athleteMatch)
}

// GetByID lấy athlete match theo ID
func (s *AthleteMatchService) GetByID(ctx context.Context, id string) (*models.AthleteMatch, error) {
	return s.athleteMatchRepo.GetByID(ctx, id)
}

// GetByUserID lấy danh sách athlete match theo UserID
func (s *AthleteMatchService) GetByUserID(ctx context.Context, userID string) ([]models.AthleteMatch, error) {
	return s.athleteMatchRepo.GetByUserID(ctx, userID)
}

// GetAll lấy danh sách tất cả athlete match với phân trang
func (s *AthleteMatchService) GetAll(ctx context.Context, page, limit int64) ([]models.AthleteMatch, error) {
	return s.athleteMatchRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin athlete match
func (s *AthleteMatchService) Update(ctx context.Context, athleteMatch *models.AthleteMatch) (*models.AthleteMatch, error) {
	if athleteMatch.ID.IsZero() {
		return nil, errors.New("invalid athlete match ID")
	}
	return s.athleteMatchRepo.Update(ctx, athleteMatch)
}

// Delete xóa athlete match theo ID
func (s *AthleteMatchService) Delete(ctx context.Context, id string) error {
	return s.athleteMatchRepo.Delete(ctx, id)
}