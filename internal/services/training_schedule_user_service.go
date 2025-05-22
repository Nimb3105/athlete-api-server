package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"be/internal/models"
	"be/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingScheduleUserService struct {
	repo                *repositories.TrainingScheduleUserRepository
	notificationService *NotificationService
	reminderService     *ReminderService
	scheduleRepo        *repositories.TrainingScheduleRepository
}

// NewTrainingScheduleUserService tạo một TrainingScheduleUserService mới
func NewTrainingScheduleUserService(
	repo *repositories.TrainingScheduleUserRepository,
	notificationService *NotificationService,
	reminderService *ReminderService,
	scheduleRepo *repositories.TrainingScheduleRepository,
) *TrainingScheduleUserService {
	return &TrainingScheduleUserService{
		repo:                repo,
		notificationService: notificationService,
		reminderService:     reminderService,
		scheduleRepo:        scheduleRepo,
	}
}

// Create tạo một lịch tập cho vận động viên, gửi thông báo và tạo lời nhắc
func (s *TrainingScheduleUserService) Create(ctx context.Context, scheduleAthlete *models.TrainingScheduleUser) (*models.TrainingScheduleUser, error) {
	if scheduleAthlete.ScheduleID.IsZero() {
		return nil, errors.New("schedule ID is required")
	}
	if scheduleAthlete.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}

	// Lấy thông tin lịch tập để tạo nội dung thông báo và lời nhắc
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleAthlete.ScheduleID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %v", err)
	}

	// Kiểm tra nếu lịch tập không ở trạng thái "Scheduled"
	if schedule.Status != "Scheduled" {
		return nil, errors.New("cannot assign athlete to a non-scheduled training schedule")
	}

	// Tạo lịch tập cho vận động viên
	createdScheduleAthlete, err := s.repo.Create(ctx, scheduleAthlete)
	if err != nil {
		return nil, err
	}

	// Tạo thông báo cho vận động viên
	notification := &models.Notification{
		UserID:     scheduleAthlete.UserID,
		ScheduleID: scheduleAthlete.ScheduleID,
		Content: fmt.Sprintf(
			"Bạn đã được gán lịch tập mới tại %s, ngày %s, bắt đầu lúc %s",
			schedule.Location,
			schedule.Date.Format("02/01/2006"),
			schedule.StartTime.Format("15:04"),
		),
		Type:      "new_schedule",
		SentDate:  time.Now(),
		Status:    "Unread",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.notificationService.Create(ctx, notification)
	if err != nil {
		// Log lỗi nhưng không làm gián đoạn nghiệp vụ chính
		fmt.Printf("Failed to create notification: %v\n", err)
	}

	// Tạo lời nhắc cho vận động viên (30 phút trước giờ bắt đầu)
	reminderTime := schedule.StartTime.Add(-30 * time.Minute)
	reminder := &models.Reminder{
		UserID:       scheduleAthlete.UserID,
		ScheduleID:   scheduleAthlete.ScheduleID,
		ReminderDate:         schedule.Date, // Lấy từ TrainingSchedule
		ReminderTime: reminderTime,
		Content: fmt.Sprintf(
			"Nhắc nhở: Lịch tập tại %s, ngày %s, sắp bắt đầu lúc %s",
			schedule.Location,
			schedule.Date.Format("02/01/2006"),
			schedule.StartTime.Format("15:04"),
		),
		Status:    "Pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.reminderService.Create(ctx, reminder)
	if err != nil {
		// Log lỗi nhưng không làm gián đoạn nghiệp vụ chính
		fmt.Printf("Failed to create reminder: %v\n", err)
	}

	return createdScheduleAthlete, nil
}

func (s *TrainingScheduleUserService) GetByID(ctx context.Context, id string) (*models.TrainingScheduleUser, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.New("invalid training schedule athlete ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TrainingScheduleUserService) GetByScheduleID(ctx context.Context, scheduleID string) ([]models.TrainingScheduleUser, error) {
	if _, err := primitive.ObjectIDFromHex(scheduleID); err != nil {
		return nil, errors.New("invalid schedule ID")
	}
	return s.repo.GetByScheduleID(ctx, scheduleID)
}

func (s *TrainingScheduleUserService) GetByUserID(ctx context.Context, userID string) ([]models.TrainingScheduleUser, error) {
	if _, err := primitive.ObjectIDFromHex(userID); err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetByUserID(ctx, userID)
}

func (s *TrainingScheduleUserService) GetAll(ctx context.Context, page, limit int64) ([]models.TrainingScheduleUser, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("invalid page or limit")
	}
	return s.repo.GetAll(ctx, page, limit)
}

func (s *TrainingScheduleUserService) Update(ctx context.Context, scheduleAthlete *models.TrainingScheduleUser) (*models.TrainingScheduleUser, error) {
	if scheduleAthlete.ID.IsZero() {
		return nil, errors.New("training schedule athlete ID is required")
	}
	if scheduleAthlete.ScheduleID.IsZero() {
		return nil, errors.New("schedule ID is required")
	}
	if scheduleAthlete.UserID.IsZero() {
		return nil, errors.New("user ID is required")
	}

	return s.repo.Update(ctx, scheduleAthlete)
}

func (s *TrainingScheduleUserService) Delete(ctx context.Context, id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return errors.New("invalid training schedule athlete ID")
	}
	return s.repo.Delete(ctx, id)
}
func (s *TrainingScheduleUserService) GetAllByUserID(ctx context.Context, userID string) ([]models.TrainingScheduleUser, error) {
	return s.repo.GetAllByUserID(ctx, userID)
}
