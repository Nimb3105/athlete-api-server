package services

import (
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NutritionPlanService provides business logic methods for NutritionPlan
type NutritionPlanService struct {
	nutritionPlanRepo *repositories.NutritionPlanRepository
	planFoodService   *PlanFoodService
	foodRepo          *repositories.FoodRepository
}

// NewNutritionPlanService creates a new NutritionPlanService
func NewNutritionPlanService(nutritionPlanRepo *repositories.NutritionPlanRepository, planFoodService *PlanFoodService, foodRepo *repositories.FoodRepository) *NutritionPlanService {
	return &NutritionPlanService{nutritionPlanRepo: nutritionPlanRepo, planFoodService: planFoodService, foodRepo: foodRepo}
}

// Create creates a new nutrition plan
func (s *NutritionPlanService) Create(ctx context.Context, nutritionPlan *models.NutritionPlan, foodIDs []string) (*models.NutritionPlan, error) {
	// Xác thực các trường bắt buộc
	if nutritionPlan.UserID.IsZero() {
		return nil, errors.New("ID vận động viên là bắt buộc")
	}
	if nutritionPlan.CreateBy.IsZero() {
		return nil, errors.New("ID huấn luyện viên là bắt buộc")
	}

	// Tính TotalCalories từ Food
	var totalCalories int
	for _, foodIDStr := range foodIDs {

		// Truy vấn Food từ cơ sở dữ liệu
		food, err := s.foodRepo.GetByID(ctx, foodIDStr)
		if err != nil {
			return nil, fmt.Errorf("không tìm thấy món ăn với ID %s: %v", foodIDStr, err)
		}

		totalCalories += food.Calories
	}

	// Cập nhật TotalCalories trong NutritionPlan
	if totalCalories <= 0 {
		return nil, errors.New("tổng calo phải lớn hơn 0")
	}
	nutritionPlan.TotalCalories = totalCalories

	// Tạo NutritionPlan
	createdPlan, err := s.nutritionPlanRepo.Create(ctx, nutritionPlan)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo kế hoạch dinh dưỡng: %v", err)
	}

	// Tạo các bản ghi PlanFood
	for _, foodIDStr := range foodIDs {
		foodID, err := primitive.ObjectIDFromHex(foodIDStr)
		if err != nil {
			// Xóa NutritionPlan nếu tạo PlanFood thất bại để tránh dữ liệu không nhất quán
			if deleteErr := s.nutritionPlanRepo.Delete(ctx, createdPlan.ID.Hex()); deleteErr != nil {
                // Ghi log lỗi xóa nếu cần
                fmt.Printf("Không thể xóa NutritionPlan %s: %v\n", createdPlan.ID.Hex(), deleteErr)
            }
            return nil, fmt.Errorf("ID món ăn không hợp lệ: %v", err)
		}

		planFood := &models.PlanFood{
			FoodID:          foodID,
			NutritionPlanID: createdPlan.ID,
		}

		if _, err := s.planFoodService.Create(ctx, planFood); err != nil {
            // Xóa NutritionPlan nếu tạo PlanFood thất bại
            if deleteErr := s.nutritionPlanRepo.Delete(ctx, createdPlan.ID.Hex()); deleteErr != nil {
                fmt.Printf("Không thể xóa NutritionPlan %s: %v\n", createdPlan.ID.Hex(), deleteErr)
            }
            return nil, fmt.Errorf("không thể tạo bản ghi món ăn: %v", err)
        }
	}

	return createdPlan, nil
}

// GetByID retrieves a nutrition plan by ID
func (s *NutritionPlanService) GetByID(ctx context.Context, id string) (*models.NutritionPlan, error) {
	return s.nutritionPlanRepo.GetByID(ctx, id)
}

// GetByAthleteID retrieves nutrition plans by athlete ID
func (s *NutritionPlanService) GetByUserID(ctx context.Context, userID string) ([]models.NutritionPlan, error) {
	return s.nutritionPlanRepo.GetByUserID(ctx, userID)
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
	if nutritionPlan.TotalCalories <= 0 {
		return nil, errors.New("total calories must be greater than zero")
	}

	return s.nutritionPlanRepo.Update(ctx, nutritionPlan)
}

// Delete deletes a nutrition plan by ID
func (s *NutritionPlanService) Delete(ctx context.Context, id string) error {
	return s.nutritionPlanRepo.Delete(ctx, id)
}
