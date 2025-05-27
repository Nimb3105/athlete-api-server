package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// NutritionMealService provides business logic methods for NutritionMeal
type FoodService struct {
	nutritionMealRepo *repositories.FoodRepository
}

// NewNutritionMealService creates a new NutritionMealService
func NewFoodService(nutritionMealRepo *repositories.FoodRepository) *FoodService {
	return &FoodService{nutritionMealRepo}
}

// Create creates a new nutrition meal
func (s *FoodService) Create(ctx context.Context, nutritionMeal *models.Food) (*models.Food, error) {
	if nutritionMeal.Name == "" {
		return nil, errors.New("food name is required")
	}
	if nutritionMeal.Calories <= 0 {
		return nil, errors.New("calories must be greater than zero")
	}
	return s.nutritionMealRepo.Create(ctx, nutritionMeal)
}

// GetByID retrieves a nutrition meal by ID
func (s *FoodService) GetByID(ctx context.Context, id string) (*models.Food, error) {
	return s.nutritionMealRepo.GetByID(ctx, id)
}

// GetByNutritionPlanID retrieves nutrition meals by nutrition plan ID
func (s *FoodService) GetByNutritionPlanID(ctx context.Context, nutritionPlanID string) ([]models.Food, error) {
	return s.nutritionMealRepo.GetByNutritionPlanID(ctx, nutritionPlanID)
}

// GetAll retrieves all nutrition meals with pagination
func (s *FoodService) GetAll(ctx context.Context, page, limit int64) ([]models.Food, error) {
	return s.nutritionMealRepo.GetAll(ctx, page, limit)
}

// Update updates nutrition meal information
func (s *FoodService) Update(ctx context.Context, nutritionMeal *models.Food) (*models.Food, error) {
	if nutritionMeal.ID.IsZero() {
		return nil, errors.New("invalid nutrition meal ID")
	}
	if nutritionMeal.Name == "" {
		return nil, errors.New("meal type is required")
	}
	if nutritionMeal.Calories <= 0 {
		return nil, errors.New("calories must be greater than zero")
	}
	return s.nutritionMealRepo.Update(ctx, nutritionMeal)
}

// Delete deletes a nutrition meal by ID
func (s *FoodService) Delete(ctx context.Context, id string) error {
	return s.nutritionMealRepo.Delete(ctx, id)
}