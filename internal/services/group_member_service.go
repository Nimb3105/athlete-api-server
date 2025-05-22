package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// GroupMemberService provides business logic methods for GroupMember
type GroupMemberService struct {
	groupMemberRepo *repositories.GroupMemberRepository
}

// NewGroupMemberService creates a new GroupMemberService
func NewGroupMemberService(groupMemberRepo *repositories.GroupMemberRepository) *GroupMemberService {
	return &GroupMemberService{groupMemberRepo}
}

// Create creates a new group member
func (s *GroupMemberService) Create(ctx context.Context, groupMember *models.GroupMember) (*models.GroupMember, error) {
	if groupMember.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.groupMemberRepo.Create(ctx, groupMember)
}

// GetByID retrieves a group member by ID
func (s *GroupMemberService) GetByID(ctx context.Context, id string) (*models.GroupMember, error) {
	return s.groupMemberRepo.GetByID(ctx, id)
}

// GetByUserID retrieves a group member by user ID
func (s *GroupMemberService) GetByUserID(ctx context.Context, userID string) (*models.GroupMember, error) {
	return s.groupMemberRepo.GetByUserID(ctx, userID)
}

// GetAll retrieves all group members with pagination
func (s *GroupMemberService) GetAll(ctx context.Context, page, limit int64) ([]models.GroupMember, error) {
	return s.groupMemberRepo.GetAll(ctx, page, limit)
}

// Update updates group member information
func (s *GroupMemberService) Update(ctx context.Context, groupMember *models.GroupMember) (*models.GroupMember, error) {
	if groupMember.ID.IsZero() {
		return nil, errors.New("invalid group member ID")
	}
	if groupMember.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	return s.groupMemberRepo.Update(ctx, groupMember)
}

// Delete deletes a group member by ID
func (s *GroupMemberService) Delete(ctx context.Context, id string) error {
	return s.groupMemberRepo.Delete(ctx, id)
}