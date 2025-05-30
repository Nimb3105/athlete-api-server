package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// PlanFoodService provides business logic for PlanFood
type PlanFoodService struct {
	planFoodRepo *repositories.PlanFoodRepository
}

// NewPlanFoodService creates a new PlanFoodService
func NewPlanFoodService(planFoodRepo *repositories.PlanFoodRepository) *PlanFoodService {
	return &PlanFoodService{planFoodRepo}
}

// Create creates a new plan-food association
func (s *PlanFoodService) Create(ctx context.Context, planFood *models.PlanFood) (*models.PlanFood, error) {
	if planFood.FoodID.IsZero() || planFood.NutritionPlanID.IsZero() {
		return nil, errors.New("food ID and nutrition plan ID are required")
	}
	return s.planFoodRepo.Create(ctx, planFood)
}

// GetByID retrieves a plan-food association by ID
func (s *PlanFoodService) GetByID(ctx context.Context, id string) (*models.PlanFood, error) {
	return s.planFoodRepo.GetByID(ctx, id)
}

// GetAll retrieves all plan-food associations with pagination
func (s *PlanFoodService) GetAll(ctx context.Context, page, limit int64) ([]models.PlanFood, error) {
	return s.planFoodRepo.GetAll(ctx, page, limit)
}

// Update updates a plan-food association
func (s *PlanFoodService) Update(ctx context.Context, planFood *models.PlanFood) (*models.PlanFood, error) {
	if planFood.ID.IsZero() {
		return nil, errors.New("invalid plan-food ID")
	}
	return s.planFoodRepo.Update(ctx, planFood)
}

// Delete deletes a plan-food association by ID
func (s *PlanFoodService) Delete(ctx context.Context, id string) error {
	return s.planFoodRepo.Delete(ctx, id)
}

func (s* PlanFoodService) GetAllByNutritionPlanID(ctx context.Context, nutritionPlanID string) ([]models.PlanFood, error) {
	if nutritionPlanID == "" {
		return nil, errors.New("nutrition plan ID is required")
	}
	return s.planFoodRepo.GetAllByNutritionPlanID(ctx, nutritionPlanID)
}
