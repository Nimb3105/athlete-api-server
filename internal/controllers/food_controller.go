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
type FoodController struct {
	foodService *services.FoodService
}

// NewNutritionMealController creates a new NutritionMealController
func NewFoodController(nutritionMealService *services.FoodService) *FoodController {
	return &FoodController{nutritionMealService}
}

// CreateNutritionMeal creates a new nutrition meal
func (c *FoodController) CreateNutritionMeal(ctx *gin.Context) {
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
		"id": true, "name": true, "foodType": true, "foodImage": true,
		"description": true, "calories": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var nutritionMeal models.Food
	if err := json.Unmarshal(bodyBytes, &nutritionMeal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdNutritionMeal, err := c.foodService.Create(ctx, &nutritionMeal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdNutritionMeal})
}

// GetNutritionMealByID retrieves a nutrition meal by ID
func (c *FoodController) GetNutritionMealByID(ctx *gin.Context) {
	id := ctx.Param("id")
	nutritionMeal, err := c.foodService.GetByID(ctx, id)
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

// GetAllNutritionMeals retrieves all nutrition meals with pagination
func (c *FoodController) GetAllNutritionMeals(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	nutritionMeals, err := c.foodService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": nutritionMeals})
}

// UpdateNutritionMeal updates a nutrition meal
func (c *FoodController) UpdateNutritionMeal(ctx *gin.Context) {
	ctx.Param("id")
	var nutritionMeal models.Food
	if err := ctx.ShouldBindJSON(&nutritionMeal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNutritionMeal, err := c.foodService.Update(ctx, &nutritionMeal)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedNutritionMeal})
}

// DeleteNutritionMeal deletes a nutrition meal by ID
func (c *FoodController) DeleteNutritionMeal(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.foodService.Delete(ctx, id); err != nil {
		if err.Error() == "nutrition meal not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "nutrition meal deleted"})
}
