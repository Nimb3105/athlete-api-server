package services

import (
	"context"
	"errors"

	"be/internal/models"
	"be/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseService struct {
	repo *repositories.ExerciseRepository
}

func NewExerciseService(repo *repositories.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo}
}

func (s *ExerciseService) Create(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error) {
	if exercise.Name == "" {
		return nil, errors.New("exercise name is required")
	}
	if exercise.Type == "" {
		return nil, errors.New("exercise type is required")
	}
	if exercise.Duration < 0 {
		return nil, errors.New("duration cannot be negative")
	}

	return s.repo.Create(ctx, exercise)
}

func (s *ExerciseService) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.New("invalid exercise ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ExerciseService) GetAll(ctx context.Context, page, limit int64) ([]models.Exercise, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("invalid page or limit")
	}
	return s.repo.GetAll(ctx, page, limit)
}

func (s *ExerciseService) Update(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error) {
	if exercise.ID.IsZero() {
		return nil, errors.New("exercise ID is required")
	}
	if exercise.Name == "" {
		return nil, errors.New("exercise name is required")
	}
	if exercise.Type == "" {
		return nil, errors.New("exercise type is required")
	}
	if exercise.Duration < 0 {
		return nil, errors.New("duration cannot be negative")
	}

	return s.repo.Update(ctx, exercise)
}

func (s *ExerciseService) Delete(ctx context.Context, id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return errors.New("invalid exercise ID")
	}
	return s.repo.Delete(ctx, id)
}
