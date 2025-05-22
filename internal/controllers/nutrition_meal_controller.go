package controllers

import (
	"be/internal/models"
	"be/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NutritionMealController handles HTTP requests for NutritionMeal
type NutritionMealController struct {
	nutritionMealService *services.NutritionMealService
}

// NewNutritionMealController creates a new NutritionMealController
func NewNutritionMealController(nutritionMealService *services.NutritionMealService) *NutritionMealController {
	return &NutritionMealController{nutritionMealService}
}

// CreateNutritionMeal creates a new nutrition meal
func (c *NutritionMealController) CreateNutritionMeal(ctx *gin.Context) {
	var bodyBytes []byte
	if rawData, err := ctx.GetRawData(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc dữ liệu"})
		return
	} else {
		bodyBytes = rawData
	}

	var tempMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không đúng định dạng"})
		return
	}

	validFields := map[string]bool{
		"id": true, "nutritionPlanId": true, "mealTime": true, "mealType": true,
		"description": true, "calories": true, "notes": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var nutritionMeal models.NutritionMeal
	if err := json.Unmarshal(bodyBytes, &nutritionMeal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdNutritionMeal, err := c.nutritionMealService.Create(ctx, &nutritionMeal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdNutritionMeal})
}

// GetNutritionMealByID retrieves a nutrition meal by ID
func (c *NutritionMealController) GetNutritionMealByID(ctx *gin.Context) {
	id := ctx.Param("id")
	nutritionMeal, err := c.nutritionMealService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "nutrition meal not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": nutritionMeal})
}

// GetNutritionMealsByNutritionPlanID retrieves nutrition meals by nutrition plan ID
func (c *NutritionMealController) GetNutritionMealsByNutritionPlanID(ctx *gin.Context) {
	nutritionPlanID := ctx.Param("nutritionPlanID")
	nutritionMeals, err := c.nutritionMealService.GetByNutritionPlanID(ctx, nutritionPlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": nutritionMeals})
}

// GetAllNutritionMeals retrieves all nutrition meals with pagination
func (c *NutritionMealController) GetAllNutritionMeals(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	nutritionMeals, err := c.nutritionMealService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": nutritionMeals})
}

// UpdateNutritionMeal updates a nutrition meal
func (c *NutritionMealController) UpdateNutritionMeal(ctx *gin.Context) {
	ctx.Param("id")
	var nutritionMeal models.NutritionMeal
	if err := ctx.ShouldBindJSON(&nutritionMeal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNutritionMeal, err := c.nutritionMealService.Update(ctx, &nutritionMeal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedNutritionMeal})
}

// DeleteNutritionMeal deletes a nutrition meal by ID
func (c *NutritionMealController) DeleteNutritionMeal(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.nutritionMealService.Delete(ctx, id); err != nil {
		if err.Error() == "nutrition meal not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "nutrition meal deleted"})
}