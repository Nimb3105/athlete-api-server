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
		"id": true, "userId": true, "CreateBy": true, "startDate": true,
		"endDate": true, "totalCalories": true, "mealsPerDay": true, "notes": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var nutritionPlan models.NutritionPlan
	if err := json.Unmarshal(bodyBytes, &nutritionPlan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdNutritionPlan, err := c.nutritionPlanService.Create(ctx, &nutritionPlan)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdNutritionPlan})
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
func (c *NutritionPlanController) GetNutritionPlansByAthleteID(ctx *gin.Context) {
	athleteID := ctx.Param("athleteID")
	nutritionPlans, err := c.nutritionPlanService.GetByAthleteID(ctx, athleteID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
