package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// NotificationService cung cấp các phương thức nghiệp vụ cho Notification
type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
}

// NewNotificationService tạo một NotificationService mới
func NewNotificationService(notificationRepo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{notificationRepo}
}

// Create tạo một notification mới
func (s *NotificationService) Create(ctx context.Context, notification *models.Notification) (*models.Notification, error) {
	if notification.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if notification.Content == "" {
		return nil, errors.New("content is required")
	}
	if notification.Type == "" {
		return nil, errors.New("notification type is required")
	}
	if notification.SentDate.IsZero() {
		return nil, errors.New("sent date is required")
	}

	return s.notificationRepo.Create(ctx, notification)
}

// GetByID lấy notification theo ID
func (s *NotificationService) GetByID(ctx context.Context, id string) (*models.Notification, error) {
	return s.notificationRepo.GetByID(ctx, id)
}

// GetByUserID lấy danh sách notification theo UserID với phân trang
func (s *NotificationService) GetByUserID(ctx context.Context, userID string, page, limit int64) ([]models.Notification, error) {
	return s.notificationRepo.GetByUserID(ctx, userID, page, limit)
}

// Update cập nhật thông tin notification
func (s *NotificationService) Update(ctx context.Context, notification *models.Notification) (*models.Notification, error) {
	if notification.ID.IsZero() {
		return nil, errors.New("invalid notification ID")
	}
	if notification.Content == "" {
		return nil, errors.New("content is required")
	}
	return s.notificationRepo.Update(ctx, notification)
}

// Delete xóa notification theo ID
func (s *NotificationService) Delete(ctx context.Context, id string) error {
	return s.notificationRepo.Delete(ctx, id)
}
