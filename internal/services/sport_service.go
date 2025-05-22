package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// SportService cung cấp các phương thức nghiệp vụ cho Sport
type SportService struct {
	sportRepo *repositories.SportRepository
}

// NewSportService tạo một SportService mới
func NewSportService(sportRepo *repositories.SportRepository) *SportService {
	return &SportService{sportRepo}
}

// Create tạo một sport mới
func (s *SportService) Create(ctx context.Context, sport *models.Sport) (*models.Sport, error) {
	if sport.Name == "" {
		return nil, errors.New("name is required")
	}
	return s.sportRepo.Create(ctx, sport)
}

// GetByID lấy sport theo ID
func (s *SportService) GetByID(ctx context.Context, id string) (*models.Sport, error) {
	return s.sportRepo.GetByID(ctx, id)
}

// GetAll lấy danh sách tất cả sport với phân trang
func (s *SportService) GetAll(ctx context.Context, page, limit int64) ([]models.Sport, error) {
	return s.sportRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin sport
func (s *SportService) Update(ctx context.Context, sport *models.Sport) (*models.Sport, error) {
	if sport.ID.IsZero() {
		return nil, errors.New("invalid sport ID")
	}
	return s.sportRepo.Update(ctx, sport)
}

// Delete xóa sport theo ID
func (s *SportService) Delete(ctx context.Context, id string) error {
	return s.sportRepo.Delete(ctx, id)
}
