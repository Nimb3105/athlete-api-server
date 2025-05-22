package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
)

// NutritionPlanService provides business logic methods for NutritionPlan
type NutritionPlanService struct {
	nutritionPlanRepo *repositories.NutritionPlanRepository
}

// NewNutritionPlanService creates a new NutritionPlanService
func NewNutritionPlanService(nutritionPlanRepo *repositories.NutritionPlanRepository) *NutritionPlanService {
	return &NutritionPlanService{nutritionPlanRepo}
}

// Create creates a new nutrition plan
func (s *NutritionPlanService) Create(ctx context.Context, nutritionPlan *models.NutritionPlan) (*models.NutritionPlan, error) {
	if nutritionPlan.UserID.IsZero() {
		return nil, errors.New("athlete ID is required")
	}
	if nutritionPlan.CreateBy.IsZero() {
		return nil, errors.New("coach ID is required")
	}
	if nutritionPlan.StartDate.IsZero() {
		return nil, errors.New("start date is required")
	}
	if nutritionPlan.EndDate.Before(nutritionPlan.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}
	if nutritionPlan.TotalCalories <= 0 {
		return nil, errors.New("total calories must be greater than zero")
	}
	return s.nutritionPlanRepo.Create(ctx, nutritionPlan)
}

// GetByID retrieves a nutrition plan by ID
func (s *NutritionPlanService) GetByID(ctx context.Context, id string) (*models.NutritionPlan, error) {
	return s.nutritionPlanRepo.GetByID(ctx, id)
}

// GetByAthleteID retrieves nutrition plans by athlete ID
func (s *NutritionPlanService) GetByAthleteID(ctx context.Context, athleteID string) ([]models.NutritionPlan, error) {
	return s.nutritionPlanRepo.GetByAthleteID(ctx, athleteID)
}

// GetAll retrieves all nutrition plans with pagination
func (s *NutritionPlanService) GetAll(ctx context.Context, page, limit int64) ([]models.NutritionPlan, error) {
	return s.nutritionPlanRepo.GetAll(ctx, page, limit)
}

// Update updates nutrition plan information
func (s *NutritionPlanService) Update(ctx context.Context, nutritionPlan *models.NutritionPlan) (*models.NutritionPlan, error) {
	if nutritionPlan.ID.IsZero() {
		return nil, errors.New("invalid nutrition plan ID")
	}
	if nutritionPlan.UserID.IsZero() {
		return nil, errors.New("athlete ID is required")
	}
	if nutritionPlan.CreateBy.IsZero() {
		return nil, errors.New("coach ID is required")
	}
	if nutritionPlan.StartDate.IsZero() {
		return nil, errors.New("start date is required")
	}
	if nutritionPlan.EndDate.Before(nutritionPlan.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}
	if nutritionPlan.TotalCalories <= 0 {
		return nil, errors.New("total calories must be greater than zero")
	}
	return s.nutritionPlanRepo.Update(ctx, nutritionPlan)
}

// Delete deletes a nutrition plan by ID
func (s *NutritionPlanService) Delete(ctx context.Context, id string) error {
	return s.nutritionPlanRepo.Delete(ctx, id)
}