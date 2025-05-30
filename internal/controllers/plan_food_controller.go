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

// PlanFoodController handles HTTP requests for PlanFood
type PlanFoodController struct {
	planFoodService *services.PlanFoodService
}

// NewPlanFoodController creates a new PlanFoodController
func NewPlanFoodController(planFoodService *services.PlanFoodService) *PlanFoodController {
	return &PlanFoodController{planFoodService}
}

// CreatePlanFood creates a new plan-food association
func (c *PlanFoodController) CreatePlanFood(ctx *gin.Context) {
	var bodyBytes []byte
	if rawData, err := ctx.GetRawData(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read request body"})
		return
	} else {
		bodyBytes = rawData
	}

	var tempMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	validFields := map[string]bool{
		"id": true, "foodId": true, "nutritionPlanId": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
			return
		}
	}

	var planFood models.PlanFood
	if err := json.Unmarshal(bodyBytes, &planFood); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot map data to model"})
		return
	}

	createdPlanFood, err := c.planFoodService.Create(ctx, &planFood)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdPlanFood})
}

// GetPlanFoodByID retrieves a plan-food association by ID
func (c *PlanFoodController) GetPlanFoodByID(ctx *gin.Context) {
	id := ctx.Param("id")
	planFood, err := c.planFoodService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "plan-food not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": planFood})
}

// GetAllPlanFoods retrieves all plan-food associations with pagination
func (c *PlanFoodController) GetAllPlanFoods(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	planFoods, err := c.planFoodService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(planFoods) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"data": []models.PlanFood{}, "message": "không có dữ liệu nào"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": planFoods})
}

// UpdatePlanFood updates a plan-food association
func (c *PlanFoodController) UpdatePlanFood(ctx *gin.Context) {
	ctx.Param("id")
	var planFood models.PlanFood
	if err := ctx.ShouldBindJSON(&planFood); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlanFood, err := c.planFoodService.Update(ctx, &planFood)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPlanFood})
}

// DeletePlanFood deletes a plan-food association by ID
func (c *PlanFoodController) DeletePlanFood(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.planFoodService.Delete(ctx, id); err != nil {
		if err.Error() == "plan-food not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "plan-food deleted"})
}

func (c *PlanFoodController) GetAllByNutritionPlanID(ctx *gin.Context) {
	nutritionPlanID := ctx.Param("nutritionPlanId")
	planFoods, err := c.planFoodService.GetAllByNutritionPlanID(ctx, nutritionPlanID)
	if err != nil {
		if err.Error() == "plan-food not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": planFoods})
}