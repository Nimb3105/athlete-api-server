package services

import (
	"context"
	"errors"

	"be/internal/models"
	"be/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingExerciseService struct {
	repo *repositories.TrainingExerciseRepository
}

func NewTrainingExerciseService(repo *repositories.TrainingExerciseRepository) *TrainingExerciseService {
	return &TrainingExerciseService{repo}
}

func (s *TrainingExerciseService) Create(ctx context.Context, trainingExercise *models.TrainingExercise) (*models.TrainingExercise, error) {
	if trainingExercise.ScheduleID.IsZero() {
		return nil, errors.New("schedule ID is required")
	}
	if trainingExercise.ExerciseID.IsZero() {
		return nil, errors.New("exercise ID is required")
	}
	if trainingExercise.Order < 0 {
		return nil, errors.New("order cannot be negative")
	}

	return s.repo.Create(ctx, trainingExercise)
}

func (s *TrainingExerciseService) GetByID(ctx context.Context, id string) (*models.TrainingExercise, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.New("invalid training exercise ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TrainingExerciseService) GetByScheduleID(ctx context.Context, scheduleID string) ([]models.TrainingExercise, error) {
	if _, err := primitive.ObjectIDFromHex(scheduleID); err != nil {
		return nil, errors.New("invalid schedule ID")
	}
	return s.repo.GetByScheduleID(ctx, scheduleID)
}

func (s *TrainingExerciseService) GetAll(ctx context.Context, page, limit int64) ([]models.TrainingExercise, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("invalid page or limit")
	}
	return s.repo.GetAll(ctx, page, limit)
}

func (s *TrainingExerciseService) Update(ctx context.Context, trainingExercise *models.TrainingExercise) (*models.TrainingExercise, error) {
	if trainingExercise.ID.IsZero() {
		return nil, errors.New("training exercise ID is required")
	}
	if trainingExercise.ExerciseID.IsZero() {
		return nil, errors.New("exercise ID is required")
	}
	if trainingExercise.Order < 0 {
		return nil, errors.New("order cannot be negative")
	}

	return s.repo.Update(ctx, trainingExercise)
}

func (s *TrainingExerciseService) Delete(ctx context.Context, id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return errors.New("invalid training exercise ID")
	}
	return s.repo.Delete(ctx, id)
}

func (s *TrainingExerciseService) GetAllByUserID(ctx context.Context, scheduleId string) ([]models.TrainingExercise, error) {
	return s.repo.GetAllByScheduleID(ctx, scheduleId)
}
