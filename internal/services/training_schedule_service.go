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

type TrainingScheduleService struct {
	repo                    *repositories.TrainingScheduleRepository
	TrainingExerciseService *TrainingExerciseService
	TrainingExerciseRepo    *repositories.TrainingExerciseRepository
}

func NewTrainingScheduleService(repo *repositories.TrainingScheduleRepository,
	TrainingExerciseService *TrainingExerciseService, TrainingExerciseRepo *repositories.TrainingExerciseRepository) *TrainingScheduleService {
	return &TrainingScheduleService{repo: repo, TrainingExerciseService: TrainingExerciseService, TrainingExerciseRepo: TrainingExerciseRepo}
}

func (s *TrainingScheduleService) Create(ctx context.Context, schedule *models.TrainingSchedule, trainingExercises []*models.TrainingExercise) (*models.TrainingSchedule, error) {

	createSchedule, err := s.repo.Create(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo lịch tập: %v", err)
	}

	for _, trainingExercise := range trainingExercises {
		// Gán ScheduleID cho TrainingExercise
		trainingExercise.ScheduleID = createSchedule.ID

		if _, err := s.TrainingExerciseService.Create(ctx, trainingExercise); err != nil {
			if deleteErr := s.repo.Delete(ctx, createSchedule.ID.Hex()); deleteErr != nil {
				fmt.Printf("không thể xóa lịch tập đã tạo %s: %v\n", createSchedule.ID.Hex(), deleteErr)
			}
			return nil, fmt.Errorf("không thể tạo lịch tập với ID: %v", err)
		}
	}

	return createSchedule, nil
}

func (s *TrainingScheduleService) GetByID(ctx context.Context, id string) (*models.TrainingSchedule, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.New("invalid schedule ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TrainingScheduleService) GetAllByDailyScheduleId(ctx context.Context, dailyScheduleId string, date string) ([]models.TrainingSchedule, error) {
	if _, err := primitive.ObjectIDFromHex(dailyScheduleId); err != nil {
		return nil, errors.New("invalid daily Schedule Id")
	}
	return s.repo.GetAllByDailyScheduleId(ctx, dailyScheduleId, date)
}

func (s *TrainingScheduleService) GetAll(ctx context.Context, page, limit int64) ([]models.TrainingSchedule, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("invalid page or limit")
	}
	return s.repo.GetAll(ctx, page, limit)
}

func (s *TrainingScheduleService) Update(ctx context.Context, schedule *models.TrainingSchedule) (*models.TrainingSchedule, error) {
	if schedule.ID.IsZero() {
		return nil, errors.New("schedule ID is required")
	}
	if schedule.Date.IsZero() || schedule.StartTime.IsZero() || schedule.EndTime.IsZero() {
		return nil, errors.New("date, start time, and end time are required")
	}
	if schedule.EndTime.Before(schedule.StartTime) {
		return nil, errors.New("end time cannot be before start time")
	}
	if schedule.CreatedBy.IsZero() {
		return nil, errors.New("created by is required")
	}

	return s.repo.Update(ctx, schedule)
}

func (s *TrainingScheduleService) Delete(ctx context.Context, id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return errors.New("invalid schedule ID")
	}
	return s.repo.Delete(ctx, id)
}

// AutoMarkOverdue tự động đánh dấu các lịch tập và bài tập đã quá hạn
func (s *TrainingScheduleService) AutoMarkOverdue(ctx context.Context) (int64, int64, error) {
	scheduleIDs, updatedSchedules, err := s.repo.MarkOverdue(ctx, time.Now())
	if err != nil {
		return 0, 0, err
	}

	if len(scheduleIDs) > 0 {
		updatedExercises, err := s.TrainingExerciseRepo.UpdateStatusByScheduleIds(ctx, scheduleIDs, "chưa hoàn thành")
		if err != nil {
			return updatedSchedules, 0, err
		}
		return updatedSchedules, updatedExercises, nil
	}

	return 0, 0, nil
}
