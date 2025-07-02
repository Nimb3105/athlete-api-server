package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyScheduleService struct {
	repo                    *repositories.DailyScheduleRepository
	TrainingScheduleService *TrainingScheduleService
}

func NewDailyScheduleService(repo *repositories.DailyScheduleRepository, trainingScheduleService *TrainingScheduleService) *DailyScheduleService {

	return &DailyScheduleService{repo: repo, TrainingScheduleService: trainingScheduleService}
}

func (s *DailyScheduleService) Create(ctx context.Context, dailySchedule *models.DailySchedule, trainingSchedules []models.CreateTrainingScheduleRequest, trainingExercises []*models.TrainingExercise) (*models.DailySchedule, error) {
	createDaily, err := s.repo.Create(ctx, dailySchedule)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo daily: %v", err)
	}

	for _, ts := range trainingSchedules {
		ts.DailyScheduleId = createDaily.ID
		_, err := s.TrainingScheduleService.Create(ctx, &ts.TrainingSchedule, ts.TrainingExercise)
		if err != nil {
			// Xóa DailySchedule đã tạo nếu có lỗi
			if deleteErr := s.repo.Delete(ctx, createDaily.ID.Hex()); deleteErr != nil {
				fmt.Printf("không thể xóa daily schedule đã tạo %s: %v\n", createDaily.ID.Hex(), deleteErr)
			}
			return nil, fmt.Errorf("lỗi tạo training schedule: %v", err)
		}
	}

	return createDaily, nil
}

func (s *DailyScheduleService) GetById(ctx context.Context, id string) (*models.DailySchedule, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.New("invalid dailyschedule id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DailyScheduleService) GetByUserID(ctx context.Context, day string, userId string) (*models.DailySchedule, error) {
	if _, err := primitive.ObjectIDFromHex(userId); err != nil {
		return nil, errors.New("invalid dailyschedule id")
	}

	return s.repo.GetByUserID(ctx, day, userId)
}

func (s *DailyScheduleService) GetAll(ctx context.Context, page, limit int64) ([]models.DailySchedule, int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}

	return s.repo.GetAll(ctx, page, limit)
}

func (s *DailyScheduleService) Update(ctx context.Context, dailySchedule *models.DailySchedule) (*models.DailySchedule, error) {
	return s.repo.Update(ctx, dailySchedule)
}

func (s *DailyScheduleService) Delete(ctx context.Context, id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return errors.New("invalid dailyschedule id")
	}
	return s.repo.Delete(ctx, id)
}
