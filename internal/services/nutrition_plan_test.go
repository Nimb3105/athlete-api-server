package services

import (
	"be/internal/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Interfaces for dependency injection
type NutritionPlanRepositoryInterface interface {
	Create(ctx context.Context, nutritionPlan *models.NutritionPlan) (*models.NutritionPlan, error)
	Delete(ctx context.Context, id string) error
}

type FoodRepositoryInterface interface {
	GetByID(ctx context.Context, id string) (*models.Food, error)
}

type PlanFoodServiceInterface interface {
	Create(ctx context.Context, planFood *models.PlanFood) (*models.PlanFood, error)
}

// Mock implementations
type MockNutritionPlanRepository struct {
	mock.Mock
}

func (m *MockNutritionPlanRepository) Create(ctx context.Context, nutritionPlan *models.NutritionPlan) (*models.NutritionPlan, error) {
	args := m.Called(ctx, nutritionPlan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NutritionPlan), args.Error(1)
}

func (m *MockNutritionPlanRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockFoodRepository struct {
	mock.Mock
}

func (m *MockFoodRepository) GetByID(ctx context.Context, id string) (*models.Food, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Food), args.Error(1)
}

type MockPlanFoodService struct {
	mock.Mock
}

func (m *MockPlanFoodService) Create(ctx context.Context, planFood *models.PlanFood) (*models.PlanFood, error) {
	args := m.Called(ctx, planFood)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlanFood), args.Error(1)
}

// TestableNutritionPlanService - version của service có thể test được
type TestableNutritionPlanService struct {
	nutritionPlanRepo NutritionPlanRepositoryInterface
	planFoodService   PlanFoodServiceInterface
	foodRepo          FoodRepositoryInterface
}

func NewTestableNutritionPlanService(
	nutritionPlanRepo NutritionPlanRepositoryInterface,
	planFoodService PlanFoodServiceInterface,
	foodRepo FoodRepositoryInterface,
) *TestableNutritionPlanService {
	return &TestableNutritionPlanService{
		nutritionPlanRepo: nutritionPlanRepo,
		planFoodService:   planFoodService,
		foodRepo:          foodRepo,
	}
}

// Copy logic của hàm Create từ service gốc
func (s *TestableNutritionPlanService) Create(ctx context.Context, nutritionPlan *models.NutritionPlan, foodIDs []string) (*models.NutritionPlan, error) {
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
			return nil, errors.New("không tìm thấy món ăn với ID " + foodIDStr + ": " + err.Error())
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
		return nil, errors.New("không thể tạo kế hoạch dinh dưỡng: " + err.Error())
	}

	// Tạo các bản ghi PlanFood
	for _, foodIDStr := range foodIDs {
		foodID, err := primitive.ObjectIDFromHex(foodIDStr)
		if err != nil {
			// Xóa NutritionPlan nếu tạo PlanFood thất bại để tránh dữ liệu không nhất quán
			if deleteErr := s.nutritionPlanRepo.Delete(ctx, createdPlan.ID.Hex()); deleteErr != nil {
				// Ghi log lỗi xóa nếu cần
			}
			return nil, errors.New("ID món ăn không hợp lệ: " + err.Error())
		}

		planFood := &models.PlanFood{
			FoodID:          foodID,
			NutritionPlanID: createdPlan.ID,
		}

		if _, err := s.planFoodService.Create(ctx, planFood); err != nil {
			// Xóa NutritionPlan nếu tạo PlanFood thất bại
			if deleteErr := s.nutritionPlanRepo.Delete(ctx, createdPlan.ID.Hex()); deleteErr != nil {
				// Log error if needed
			}
			return nil, errors.New("không thể tạo bản ghi món ăn: " + err.Error())
		}
	}

	return createdPlan, nil
}

// Test helper function để tạo Food mock
func createMockFood(id string, calories int) *models.Food {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return &models.Food{
		ID:       objectID,
		Name:     "Test Food",
		Calories: calories,
	}
}

// Test helper function để tạo NutritionPlan
func createTestNutritionPlan() *models.NutritionPlan {
	userID := primitive.NewObjectID()
	createBy := primitive.NewObjectID()
	
	return &models.NutritionPlan{
		Name:        "Test Plan",
		UserID:      userID,
		CreateBy:    createBy,
		MealCount:   1,
		MealType:    "Breakfast",
		MealTime:    time.Now(),
		Description: "Test description",
	}
}

func TestNutritionPlanService_Create_Success(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	foodIDs := []string{"507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012"}

	// Mock foods
	food1 := createMockFood("507f1f77bcf86cd799439011", 200)
	food2 := createMockFood("507f1f77bcf86cd799439012", 300)

	// Mock expectations
	mockFoodRepo.On("GetByID", ctx, foodIDs[0]).Return(food1, nil)
	mockFoodRepo.On("GetByID", ctx, foodIDs[1]).Return(food2, nil)

	createdPlan := *nutritionPlan
	createdPlan.ID = primitive.NewObjectID()
	createdPlan.TotalCalories = 500
	mockNutritionPlanRepo.On("Create", ctx, mock.AnythingOfType("*models.NutritionPlan")).Return(&createdPlan, nil)

	mockPlanFoodService.On("Create", ctx, mock.AnythingOfType("*models.PlanFood")).Return(&models.PlanFood{}, nil).Times(2)

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 500, result.TotalCalories)
	assert.Equal(t, createdPlan.ID, result.ID)

	mockNutritionPlanRepo.AssertExpectations(t)
	mockFoodRepo.AssertExpectations(t)
	mockPlanFoodService.AssertExpectations(t)
}

func TestNutritionPlanService_Create_MissingUserID(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	nutritionPlan.UserID = primitive.ObjectID{} // Empty ObjectID
	foodIDs := []string{"507f1f77bcf86cd799439011"}

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ID vận động viên là bắt buộc", err.Error())
}

func TestNutritionPlanService_Create_MissingCreateBy(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	nutritionPlan.CreateBy = primitive.ObjectID{} // Empty ObjectID
	foodIDs := []string{"507f1f77bcf86cd799439011"}

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ID huấn luyện viên là bắt buộc", err.Error())
}

func TestNutritionPlanService_Create_FoodNotFound(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	foodIDs := []string{"507f1f77bcf86cd799439011"}

	// Mock expectations
	mockFoodRepo.On("GetByID", ctx, foodIDs[0]).Return(nil, errors.New("food not found"))

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "không tìm thấy món ăn với ID")

	mockFoodRepo.AssertExpectations(t)
}

func TestNutritionPlanService_Create_ZeroTotalCalories(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	foodIDs := []string{"507f1f77bcf86cd799439011"}

	// Mock food with 0 calories
	food := createMockFood("507f1f77bcf86cd799439011", 0)
	mockFoodRepo.On("GetByID", ctx, foodIDs[0]).Return(food, nil)

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "tổng calo phải lớn hơn 0", err.Error())

	mockFoodRepo.AssertExpectations(t)
}

func TestNutritionPlanService_Create_RepositoryCreateError(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	foodIDs := []string{"507f1f77bcf86cd799439011"}

	// Mock food
	food := createMockFood("507f1f77bcf86cd799439011", 200)
	mockFoodRepo.On("GetByID", ctx, foodIDs[0]).Return(food, nil)

	// Mock repository error
	mockNutritionPlanRepo.On("Create", ctx, mock.AnythingOfType("*models.NutritionPlan")).Return(nil, errors.New("database error"))

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "không thể tạo kế hoạch dinh dưỡng")

	mockFoodRepo.AssertExpectations(t)
	mockNutritionPlanRepo.AssertExpectations(t)
}

func TestNutritionPlanService_Create_InvalidFoodID(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	foodIDs := []string{"invalid-id"}

	// Mock food with valid ID first
	food := createMockFood("507f1f77bcf86cd799439011", 200)
	mockFoodRepo.On("GetByID", ctx, "invalid-id").Return(food, nil)

	createdPlan := *nutritionPlan
	createdPlan.ID = primitive.NewObjectID()
	createdPlan.TotalCalories = 200
	mockNutritionPlanRepo.On("Create", ctx, mock.AnythingOfType("*models.NutritionPlan")).Return(&createdPlan, nil)

	// Mock delete call when invalid ID causes error
	mockNutritionPlanRepo.On("Delete", ctx, createdPlan.ID.Hex()).Return(nil)

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "ID món ăn không hợp lệ")

	mockFoodRepo.AssertExpectations(t)
	mockNutritionPlanRepo.AssertExpectations(t)
}

func TestNutritionPlanService_Create_PlanFoodCreateError(t *testing.T) {
	// Arrange
	mockNutritionPlanRepo := new(MockNutritionPlanRepository)
	mockFoodRepo := new(MockFoodRepository)
	mockPlanFoodService := new(MockPlanFoodService)

	service := NewTestableNutritionPlanService(mockNutritionPlanRepo, mockPlanFoodService, mockFoodRepo)

	ctx := context.Background()
	nutritionPlan := createTestNutritionPlan()
	foodIDs := []string{"507f1f77bcf86cd799439011"}

	// Mock food
	food := createMockFood("507f1f77bcf86cd799439011", 200)
	mockFoodRepo.On("GetByID", ctx, foodIDs[0]).Return(food, nil)

	createdPlan := *nutritionPlan
	createdPlan.ID = primitive.NewObjectID()
	createdPlan.TotalCalories = 200
	mockNutritionPlanRepo.On("Create", ctx, mock.AnythingOfType("*models.NutritionPlan")).Return(&createdPlan, nil)

	// Mock PlanFood create error
	mockPlanFoodService.On("Create", ctx, mock.AnythingOfType("*models.PlanFood")).Return(nil, errors.New("planfood create error"))

	// Mock delete call when PlanFood creation fails
	mockNutritionPlanRepo.On("Delete", ctx, createdPlan.ID.Hex()).Return(nil)

	// Act
	result, err := service.Create(ctx, nutritionPlan, foodIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "không thể tạo bản ghi món ăn")

	mockFoodRepo.AssertExpectations(t)
	mockNutritionPlanRepo.AssertExpectations(t)
	mockPlanFoodService.AssertExpectations(t)
}