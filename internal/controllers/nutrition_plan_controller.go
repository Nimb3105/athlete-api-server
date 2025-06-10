package controllers

import (
	"be/internal/models"
	"be/internal/services"
	"context"
	"time"

	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateNutritionPlanRequest struct {
	models.NutritionPlan
	FoodIDs []string `json:"foodIds" validate:"required,dive,hex24"`
}

func hex24Validator(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	_, err := primitive.ObjectIDFromHex(str)
	return err == nil
}

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// NutritionPlanController handles HTTP requests for NutritionPlan
type NutritionPlanController struct {
	nutritionPlanService *services.NutritionPlanService
}

// NewNutritionPlanController creates a new NutritionPlanController
func NewNutritionPlanController(nutritionPlanService *services.NutritionPlanService) *NutritionPlanController {
	return &NutritionPlanController{nutritionPlanService}
}

// CreateNutritionPlan creates a new nutrition plan
func (c *NutritionPlanController) CreateNutritionPlan(ctx *gin.Context) {
	// Khởi tạo trình xác thực
	validate := validator.New()
	validate.RegisterValidation("hex24", hex24Validator)

	var req CreateNutritionPlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	// Xác thực yêu cầu
	if err := validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{Error: fmt.Sprintf("Xác thực thất bại: %v", err)})
		return
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Tạo kế hoạch dinh dưỡng
	createdNutritionPlan, err := c.nutritionPlanService.Create(ctxTimeout, &req.NutritionPlan, req.FoodIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, APIResponse{Data: createdNutritionPlan})
}

// GetNutritionPlanByID retrieves a nutrition plan by ID
func (c *NutritionPlanController) GetNutritionPlanByID(ctx *gin.Context) {
	id := ctx.Param("id")
	nutritionPlan, err := c.nutritionPlanService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "nutrition plan not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": nutritionPlan})
}

// GetNutritionPlansByAthleteID retrieves nutrition plans by athlete ID
func (c *NutritionPlanController) GetNutritionPlansByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	nutritionPlans, err := c.nutritionPlanService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(nutritionPlans) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.NutritionPlan{},
			"note":    "Chưa có kế hoạch dinh dưỡng nào cho vận động viên này",
			"message": "Không có dữ liệu nào"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": nutritionPlans})
}

// GetAllNutritionPlans retrieves all nutrition plans with pagination
func (c *NutritionPlanController) GetAllNutritionPlans(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	nutritionPlans, err := c.nutritionPlanService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(nutritionPlans) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.NutritionPlan{},
			"message": "Không có dữ liệu nào",
			"note":    "Chưa có kế hoạch dinh dưỡng nào được ghi nhận"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": nutritionPlans})
}

// UpdateNutritionPlan updates a nutrition plan
func (c *NutritionPlanController) UpdateNutritionPlan(ctx *gin.Context) {
	ctx.Param("id")
	var nutritionPlan models.NutritionPlan
	if err := ctx.ShouldBindJSON(&nutritionPlan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNutritionPlan, err := c.nutritionPlanService.Update(ctx, &nutritionPlan)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedNutritionPlan})
}

// DeleteNutritionPlan deletes a nutrition plan by ID
func (c *NutritionPlanController) DeleteNutritionPlan(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.nutritionPlanService.Delete(ctx, id); err != nil {
		if err.Error() == "nutrition plan not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "nutrition plan deleted"})
}
