package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// ReminderService cung cấp các phương thức nghiệp vụ cho Reminder
type ReminderService struct {
	reminderRepo *repositories.ReminderRepository
}

// NewReminderService tạo một ReminderService mới
func NewReminderService(reminderRepo *repositories.ReminderRepository) *ReminderService {
	return &ReminderService{reminderRepo}
}

func (s *ReminderService) GetAll(ctx context.Context, page, limit int64) ([]models.Reminder, error) {
	// Lấy danh sách tất cả reminders với phân trang
	return s.reminderRepo.GetAll(ctx, page, limit)
}

// Create tạo một reminder mới
func (s *ReminderService) Create(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
	if reminder.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}
	if reminder.Content == "" {
		return nil, errors.New("content is required")
	}
	if reminder.ReminderTime.IsZero() {
		return nil, errors.New("reminder time is required")
	}

	return s.reminderRepo.Create(ctx, reminder)
}

// GetByID lấy reminder theo ID
func (s *ReminderService) GetByID(ctx context.Context, id string) (*models.Reminder, error) {
	return s.reminderRepo.GetByID(ctx, id)
}

// GetByUserID lấy danh sách reminder theo UserID với phân trang
func (s *ReminderService) GetByUserID(ctx context.Context, userID string, page, limit int64) ([]models.Reminder, error) {
	return s.reminderRepo.GetByUserID(ctx, userID, page, limit)
}

// Update cập nhật thông tin reminder
func (s *ReminderService) Update(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
	if reminder.ID.IsZero() {
		return nil, errors.New("invalid reminder ID")
	}
	if reminder.Content == "" {
		return nil, errors.New("content is required")
	}
	return s.reminderRepo.Update(ctx, reminder)
}

// Delete xóa reminder theo ID
func (s *ReminderService) Delete(ctx context.Context, id string) error {
	return s.reminderRepo.Delete(ctx, id)
}
