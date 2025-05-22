package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// NutritionMealService provides business logic methods for NutritionMeal
type NutritionMealService struct {
	nutritionMealRepo *repositories.NutritionMealRepository
}

// NewNutritionMealService creates a new NutritionMealService
func NewNutritionMealService(nutritionMealRepo *repositories.NutritionMealRepository) *NutritionMealService {
	return &NutritionMealService{nutritionMealRepo}
}

// Create creates a new nutrition meal
func (s *NutritionMealService) Create(ctx context.Context, nutritionMeal *models.NutritionMeal) (*models.NutritionMeal, error) {
	if nutritionMeal.NutritionPlanID.IsZero() {
		return nil, errors.New("nutrition plan ID is required")
	}
	if nutritionMeal.MealType == "" {
		return nil, errors.New("meal type is required")
	}
	if nutritionMeal.Calories <= 0 {
		return nil, errors.New("calories must be greater than zero")
	}
	return s.nutritionMealRepo.Create(ctx, nutritionMeal)
}

// GetByID retrieves a nutrition meal by ID
func (s *NutritionMealService) GetByID(ctx context.Context, id string) (*models.NutritionMeal, error) {
	return s.nutritionMealRepo.GetByID(ctx, id)
}

// GetByNutritionPlanID retrieves nutrition meals by nutrition plan ID
func (s *NutritionMealService) GetByNutritionPlanID(ctx context.Context, nutritionPlanID string) ([]models.NutritionMeal, error) {
	return s.nutritionMealRepo.GetByNutritionPlanID(ctx, nutritionPlanID)
}

// GetAll retrieves all nutrition meals with pagination
func (s *NutritionMealService) GetAll(ctx context.Context, page, limit int64) ([]models.NutritionMeal, error) {
	return s.nutritionMealRepo.GetAll(ctx, page, limit)
}

// Update updates nutrition meal information
func (s *NutritionMealService) Update(ctx context.Context, nutritionMeal *models.NutritionMeal) (*models.NutritionMeal, error) {
	if nutritionMeal.ID.IsZero() {
		return nil, errors.New("invalid nutrition meal ID")
	}
	if nutritionMeal.NutritionPlanID.IsZero() {
		return nil, errors.New("nutrition plan ID is required")
	}
	if nutritionMeal.MealType == "" {
		return nil, errors.New("meal type is required")
	}
	if nutritionMeal.Calories <= 0 {
		return nil, errors.New("calories must be greater than zero")
	}
	return s.nutritionMealRepo.Update(ctx, nutritionMeal)
}

// Delete deletes a nutrition meal by ID
func (s *NutritionMealService) Delete(ctx context.Context, id string) error {
	return s.nutritionMealRepo.Delete(ctx, id)
}