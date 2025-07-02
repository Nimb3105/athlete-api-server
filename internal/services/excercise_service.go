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
	if exercise.BodyPart == "" {
		return nil, errors.New("exercise bodyPart is required")
	}
	if exercise.Equipment == "" {
		return nil, errors.New("exercise Equipment is negative")
	}
	if exercise.Target == "" {
		return nil, errors.New("exercise Target is negative")
	}
	if len(exercise.Instructions) == 0 {
		return nil, errors.New("exercise instructions are required")
	}

	return s.repo.Create(ctx, exercise)
}

func (s *ExerciseService) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.New("invalid exercise ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ExerciseService) GetAllBySportId(ctx context.Context, sportId string) ([]models.Exercise, error) {
	if _, err := primitive.ObjectIDFromHex(sportId); err != nil {
		return nil, errors.New("invalid sport Id")
	}
	return s.repo.GetAllBySportId(ctx, sportId)
}

func (s *ExerciseService) GetAllBySportName(ctx context.Context, sportName string, page, limit int64) ([]models.Exercise, int64, error) {
	if sportName == "" {
		return nil, 0, errors.New("sportName is required")
	}
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
	return s.repo.GetAllBySportName(ctx, sportName, page, limit)
}

func (s *ExerciseService) GetAllByBodyPart(ctx context.Context, bodyPart string, page, limit int64) ([]models.Exercise, int64, error) {
	if bodyPart == "" {
		return nil, 0, errors.New("bodyPart is required")
	}
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
	return s.repo.GetAllByBodyPart(ctx, bodyPart, page, limit)
}

func (s *ExerciseService) GetAll(ctx context.Context, page, limit int64) ([]models.Exercise, int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
	return s.repo.GetAll(ctx, page, limit)
}

func (s *ExerciseService) Update(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error) {
	if exercise.Name == "" {
		return nil, errors.New("exercise name is required")
	}
	if exercise.BodyPart == "" {
		return nil, errors.New("exercise bodyPart is required")
	}
	if exercise.Equipment == "" {
		return nil, errors.New("exercise Equipment is negative")
	}
	if exercise.Target == "" {
		return nil, errors.New("exercise Target is negative")
	}
	if len(exercise.Instructions) == 0 {
		return nil, errors.New("exercise instructions are required")
	}

	return s.repo.Update(ctx, exercise)
}

func (s *ExerciseService) Delete(ctx context.Context, id string) error {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return errors.New("invalid exercise ID")
	}
	return s.repo.Delete(ctx, id)
}
