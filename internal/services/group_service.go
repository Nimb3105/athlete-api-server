package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// GroupService cung cấp các phương thức nghiệp vụ cho Group
type GroupService struct {
	groupRepo *repositories.GroupRepository
}

// NewGroupService tạo một GroupService mới
func NewGroupService(groupRepo *repositories.GroupRepository) *GroupService {
	return &GroupService{groupRepo}
}

// Create tạo một group mới
func (s *GroupService) Create(ctx context.Context, group *models.Group) (*models.Group, error) {
	if group.CreatedBy.IsZero() {
		return nil, errors.New("created by ID is required")
	}
	if group.Name == "" {
		return nil, errors.New("group name is required")
	}
	return s.groupRepo.Create(ctx, group)
}

// GetByID lấy group theo ID
func (s *GroupService) GetByID(ctx context.Context, id string) (*models.Group, error) {
	return s.groupRepo.GetByID(ctx, id)
}

// GetByCreatedBy lấy danh sách group theo CreatedBy
func (s *GroupService) GetByCreatedBy(ctx context.Context, createdBy string) ([]models.Group, error) {
	return s.groupRepo.GetByCreatedBy(ctx, createdBy)
}

// GetAll lấy danh sách tất cả group với phân trang
func (s *GroupService) GetAll(ctx context.Context, page, limit int64) ([]models.Group, error) {
	return s.groupRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin group
func (s *GroupService) Update(ctx context.Context, group *models.Group) (*models.Group, error) {
	if group.ID.IsZero() {
		return nil, errors.New("invalid group ID")
	}
	if group.Name == "" {
		return nil, errors.New("group name is required")
	}
	return s.groupRepo.Update(ctx, group)
}

// Delete xóa group theo ID
func (s *GroupService) Delete(ctx context.Context, id string) error {
	return s.groupRepo.Delete(ctx, id)
}