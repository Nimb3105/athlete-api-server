package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// AchievementService cung cấp các phương thức nghiệp vụ cho Achievement
type AchievementService struct {
	achievementRepo *repositories.AchievementRepository
}

// NewAchievementService tạo một AchievementService mới
func NewAchievementService(achievementRepo *repositories.AchievementRepository) *AchievementService {
	return &AchievementService{achievementRepo}
}

// Create tạo một achievement mới
func (s *AchievementService) Create(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error) {
	if achievement.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if achievement.Title == "" {
		return nil, errors.New("achievement title is required")
	}
	return s.achievementRepo.Create(ctx, achievement)
}

// GetByID lấy achievement theo ID
func (s *AchievementService) GetByID(ctx context.Context, id string) (*models.Achievement, error) {
	return s.achievementRepo.GetByID(ctx, id)
}

// GetByUserID lấy danh sách achievement theo UserID
func (s *AchievementService) GetByUserID(ctx context.Context, userID string) ([]models.Achievement, error) {
	return s.achievementRepo.GetByUserID(ctx, userID)
}

// GetAll lấy danh sách tất cả achievement với phân trang
func (s *AchievementService) GetAll(ctx context.Context, page, limit int64) ([]models.Achievement, error) {
	return s.achievementRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin achievement
func (s *AchievementService) Update(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error) {
	if achievement.ID.IsZero() {
		return nil, errors.New("invalid achievement ID")
	}
	if achievement.Title == "" {
		return nil, errors.New("achievement title is required")
	}
	return s.achievementRepo.Update(ctx, achievement)
}

// Delete xóa achievement theo ID
func (s *AchievementService) Delete(ctx context.Context, id string) error {
	return s.achievementRepo.Delete(ctx, id)
}
