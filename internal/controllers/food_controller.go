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

	nutritionMeals,totalCount, err := c.foodService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(nutritionMeals) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Food{},
			"totalCount": 0,
			"notes":      "khồn có bữa ăn dinh dưỡng nào",
			"message":    "Không có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": nutritionMeals,"totalCount":totalCount})
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

func (h *FoodController) FindFoodsByFilter(c *gin.Context) {
	// Nhận foodType: ưu tiên path param, sau đó tới query
	foodType := c.Param("foodType")
	if foodType == "" {
		foodType = c.DefaultQuery("foodType", "")
	}

	parseInt := func(key string) int {
		valStr := c.DefaultQuery(key, "-1")
		val, _ := strconv.Atoi(valStr)
		return val
	}

	caloriesMin := parseInt("caloriesMin")
	caloriesMax := parseInt("caloriesMax")
	proteinMin := parseInt("proteinMin")
	proteinMax := parseInt("proteinMax")
	carbsMin := parseInt("carbsMin")
	carbsMax := parseInt("carbsMax")
	fatMin := parseInt("fatMin")
	fatMax := parseInt("fatMax")

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)

	foods, total, err := h.foodService.FindFoodsByFilter(
		c.Request.Context(),
		foodType,
		caloriesMin, caloriesMax,
		proteinMin, proteinMax,
		carbsMin, carbsMax,
		fatMin, fatMax,
		page, limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(foods) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":       []models.Food{},
			"totalCount": 0,
			"notes":      "không có dữ liệu",
			"message":    "chưa có món ăn nào phù hợp",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       foods,
		"totalCount": total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + limit - 1) / limit,
	})
}
