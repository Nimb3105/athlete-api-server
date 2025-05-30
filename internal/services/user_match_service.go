package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// UserMatchService cung cấp các phương thức nghiệp vụ cho AthleteMatch
type UserMatchService struct {
	userMatchRepo *repositories.UserMatchRepository
}

// NewUserMatchService tạo một UserMatchService mới
func NewUserMatchService(userMatchRepo *repositories.UserMatchRepository) *UserMatchService {
	return &UserMatchService{userMatchRepo}
}

// Create tạo một athlete match mới
func (s *UserMatchService) Create(ctx context.Context, userMatch *models.UserMatch) (*models.UserMatch, error) {
	if userMatch.UserID.IsZero() || userMatch.MatchID.IsZero() {
		return nil, errors.New("user ID and match ID are required")
	}
	return s.userMatchRepo.Create(ctx, userMatch)
}

// GetByID lấy athlete match theo ID
func (s *UserMatchService) GetByID(ctx context.Context, id string) (*models.UserMatch, error) {
	return s.userMatchRepo.GetByID(ctx, id)
}

// GetByUserID lấy danh sách athlete match theo UserID
func (s *UserMatchService) GetByUserID(ctx context.Context, userID string) ([]models.UserMatch, error) {
	return s.userMatchRepo.GetByUserID(ctx, userID)
}

// GetAll lấy danh sách tất cả athlete match với phân trang
func (s *UserMatchService) GetAll(ctx context.Context, page, limit int64) ([]models.UserMatch, error) {
	return s.userMatchRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin athlete match
func (s *UserMatchService) Update(ctx context.Context, userMatch *models.UserMatch) (*models.UserMatch, error) {
	if userMatch.ID.IsZero() {
		return nil, errors.New("invalid athlete match ID")
	}
	return s.userMatchRepo.Update(ctx, userMatch)
}

// Delete xóa athlete match theo ID
func (s *UserMatchService) Delete(ctx context.Context, id string) error {
	return s.userMatchRepo.Delete(ctx, id)
}