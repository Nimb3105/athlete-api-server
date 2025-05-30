package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// AthleteService cung cấp các phương thức nghiệp vụ cho Athlete
type AthleteService struct {
	athleteRepo *repositories.AthleteRepository
}

// NewAthleteService tạo một AthleteService mới
func NewAthleteService(athleteRepo *repositories.AthleteRepository) *AthleteService {
	return &AthleteService{athleteRepo}
}

// Create tạo một athlete mới
func (s *AthleteService) Create(ctx context.Context, athlete *models.Athlete) (*models.Athlete, error) {
	if athlete.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.athleteRepo.Create(ctx, athlete)
}

// GetByID lấy athlete theo ID
func (s *AthleteService) GetByID(ctx context.Context, id string) (*models.Athlete, error) {
	return s.athleteRepo.GetByID(ctx, id)
}

// GetByUserID lấy athlete theo UserID
func (s *AthleteService) GetByUserID(ctx context.Context, userID string) (*models.Athlete, error) {
	return s.athleteRepo.GetByUserID(ctx, userID)
}

// GetAll lấy danh sách tất cả athlete với phân trang
func (s *AthleteService) GetAll(ctx context.Context, page, limit int64) ([]models.Athlete, int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
	return s.athleteRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin athlete
func (s *AthleteService) Update(ctx context.Context, athlete *models.Athlete) (*models.Athlete, error) {
	if athlete.ID.IsZero() {
		return nil, errors.New("invalid athlete ID")
	}
	return s.athleteRepo.Update(ctx, athlete)
}

// Delete xóa athlete theo ID
func (s *AthleteService) Delete(ctx context.Context, id string) error {
	return s.athleteRepo.Delete(ctx, id)
}
