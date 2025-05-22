package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// SportAthleteService cung cấp các phương thức nghiệp vụ cho SportAthlete
type SportAthleteService struct {
	sportAthleteRepo *repositories.SportAthleteRepository
}

// NewSportAthleteService tạo một SportAthleteService mới
func NewSportAthleteService(sportAthleteRepo *repositories.SportAthleteRepository) *SportAthleteService {
	return &SportAthleteService{sportAthleteRepo}
}

// Create tạo một sport athlete mới
func (s *SportAthleteService) Create(ctx context.Context, sportAthlete *models.SportAthlete) (*models.SportAthlete, error) {
	if sportAthlete.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.sportAthleteRepo.Create(ctx, sportAthlete)
}

// GetByID lấy sport athlete theo ID
func (s *SportAthleteService) GetByID(ctx context.Context, id string) (*models.SportAthlete, error) {
	return s.sportAthleteRepo.GetByID(ctx, id)
}

// GetByUserID lấy sport athlete theo UserID
func (s *SportAthleteService) GetByUserID(ctx context.Context, userID string) (*models.SportAthlete, error) {
	return s.sportAthleteRepo.GetByUserID(ctx, userID)
}

func (s *SportAthleteService) GetBySportID(ctx context.Context, sportID string) (*models.SportAthlete, error) {
	return s.sportAthleteRepo.GetBySportID(ctx, sportID)
}

// GetAll lấy danh sách tất cả sport athlete với phân trang
func (s *SportAthleteService) GetAll(ctx context.Context, page, limit int64) ([]models.SportAthlete, error) {
	return s.sportAthleteRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin sport athlete
func (s *SportAthleteService) Update(ctx context.Context, sportAthlete *models.SportAthlete) (*models.SportAthlete, error) {
	if sportAthlete.ID.IsZero() {
		return nil, errors.New("invalid sport athlete ID")
	}
	return s.sportAthleteRepo.Update(ctx, sportAthlete)
}

// Delete xóa sport athlete theo ID
func (s *SportAthleteService) Delete(ctx context.Context, id string) error {
	return s.sportAthleteRepo.Delete(ctx, id)
}

func (s *SportAthleteService) GetAllByUserID(ctx context.Context, userID string) ([]models.SportAthlete, error) {
	return s.sportAthleteRepo.GetAllByUserID(ctx, userID)
}
