package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// FeedbackService cung cấp các phương thức nghiệp vụ cho Feedback
type FeedbackService struct {
	feedbackRepo *repositories.FeedbackRepository
}

// NewFeedbackService tạo một FeedbackService mới
func NewFeedbackService(feedbackRepo *repositories.FeedbackRepository) *FeedbackService {
	return &FeedbackService{feedbackRepo}
}

// Create tạo một feedback mới
func (s *FeedbackService) Create(ctx context.Context, feedback *models.Feedback) (*models.Feedback, error) {
	if feedback.UserID.IsZero() || feedback.ScheduleID.IsZero() {
		return nil, errors.New("user ID and schedule ID are required")
	}
	if feedback.Content == "" {
		return nil, errors.New("feedback content is required")
	}
	return s.feedbackRepo.Create(ctx, feedback)
}

// GetByID lấy feedback theo ID
func (s *FeedbackService) GetByID(ctx context.Context, id string) (*models.Feedback, error) {
	return s.feedbackRepo.GetByID(ctx, id)
}

// GetByUserID lấy danh sách feedback theo UserID
func (s *FeedbackService) GetByUserID(ctx context.Context, userID string) ([]models.Feedback, error) {
	return s.feedbackRepo.GetByUserID(ctx, userID)
}

// GetAll lấy danh sách tất cả feedback với phân trang
func (s *FeedbackService) GetAll(ctx context.Context, page, limit int64) ([]models.Feedback, error) {
	return s.feedbackRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin feedback
func (s *FeedbackService) Update(ctx context.Context, feedback *models.Feedback) (*models.Feedback, error) {
	if feedback.ID.IsZero() {
		return nil, errors.New("invalid feedback ID")
	}
	if feedback.Content == "" {
		return nil, errors.New("feedback content is required")
	}
	return s.feedbackRepo.Update(ctx, feedback)
}

// Delete xóa feedback theo ID
func (s *FeedbackService) Delete(ctx context.Context, id string) error {
	return s.feedbackRepo.Delete(ctx, id)
}