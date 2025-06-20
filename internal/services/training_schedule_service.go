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
}

func NewTrainingScheduleService(repo *repositories.TrainingScheduleRepository,
	TrainingExerciseService *TrainingExerciseService) *TrainingScheduleService {
	return &TrainingScheduleService{repo: repo, TrainingExerciseService: TrainingExerciseService}
}

func (s *TrainingScheduleService) Create(ctx context.Context, schedule *models.TrainingSchedule, trainingExercises []*models.TrainingExercise, athleteId string) (*models.TrainingSchedule, error) {
	if schedule.Date.IsZero() || schedule.StartTime.IsZero() || schedule.EndTime.IsZero() {
		return nil, errors.New("date, start time, and end time are required")
	}
	if schedule.EndTime.Before(schedule.StartTime) {
		return nil, errors.New("end time cannot be before start time")
	}
	if schedule.Status == "" {
		schedule.Status = "Scheduled"
	}
	if schedule.CreatedBy.IsZero() {
		return nil, errors.New("created by is required")
	}

	createSchedule, err := s.repo.Create(ctx, schedule, athleteId)
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

// training_schedule_service.go
func (s *TrainingScheduleService) AutoMarkOverdue(ctx context.Context) (int64, error) {
	return s.repo.MarkOverdue(ctx, time.Now())
}
