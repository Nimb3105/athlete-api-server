package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// CoachCertificationService cung cấp các phương thức nghiệp vụ cho CoachCertification
type CoachCertificationService struct {
	certificationRepo *repositories.CoachCertificationRepository
}

// NewCoachCertificationService tạo một CoachCertificationService mới
func NewCoachCertificationService(certificationRepo *repositories.CoachCertificationRepository) *CoachCertificationService {
	return &CoachCertificationService{certificationRepo}
}

// Create tạo một coach certification mới
func (s *CoachCertificationService) Create(ctx context.Context, certification *models.CoachCertification) (*models.CoachCertification, error) {
	if certification.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if certification.Name == "" {
		return nil, errors.New("certification name is required")
	}
	return s.certificationRepo.Create(ctx, certification)
}

// GetByID lấy coach certification theo ID
func (s *CoachCertificationService) GetByID(ctx context.Context, id string) (*models.CoachCertification, error) {
	return s.certificationRepo.GetByID(ctx, id)
}

// GetByUserID lấy danh sách coach certification theo UserID
func (s *CoachCertificationService) GetByUserID(ctx context.Context, userID string) ([]models.CoachCertification, error) {
	return s.certificationRepo.GetByUserID(ctx, userID)
}

// GetAll lấy danh sách tất cả coach certification với phân trang
func (s *CoachCertificationService) GetAll(ctx context.Context, page, limit int64) ([]models.CoachCertification, error) {
	return s.certificationRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin coach certification
func (s *CoachCertificationService) Update(ctx context.Context, certification *models.CoachCertification) (*models.CoachCertification, error) {
	if certification.ID.IsZero() {
		return nil, errors.New("invalid certification ID")
	}
	if certification.Name == "" {
		return nil, errors.New("certification name is required")
	}
	return s.certificationRepo.Update(ctx, certification)
}

// Delete xóa coach certification theo ID
func (s *CoachCertificationService) Delete(ctx context.Context, id string) error {
	return s.certificationRepo.Delete(ctx, id)
}