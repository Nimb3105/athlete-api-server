package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// CoachService cung cấp các phương thức nghiệp vụ cho Coach
type CoachService struct {
	coachRepo *repositories.CoachRepository
}

// NewCoachService tạo một CoachService mới
func NewCoachService(coachRepo *repositories.CoachRepository) *CoachService {
	return &CoachService{coachRepo}
}

// Create tạo một coach mới
func (s *CoachService) Create(ctx context.Context, coach *models.Coach) (*models.Coach, error) {
	if coach.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.coachRepo.Create(ctx, coach)
}

// GetByID lấy coach theo ID
func (s *CoachService) GetByID(ctx context.Context, id string) (*models.Coach, error) {
	return s.coachRepo.GetByID(ctx, id)
}

func (s *CoachService) GetByUserID(ctx context.Context, userID string) (*models.Coach, error) {
	return s.coachRepo.GetByUserID(ctx, userID)
}

// GetAll lấy danh sách tất cả coach với phân trang
func (s *CoachService) GetAll(ctx context.Context, page, limit int64) ([]models.Coach, error) {
	return s.coachRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin coach
func (s *CoachService) Update(ctx context.Context, coach *models.Coach) (*models.Coach, error) {
	if coach.ID.IsZero() {
		return nil, errors.New("invalid coach ID")
	}
	return s.coachRepo.Update(ctx, coach)
}

// Delete xóa coach theo ID
func (s *CoachService) Delete(ctx context.Context, id string) error {
	return s.coachRepo.Delete(ctx, id)
}
