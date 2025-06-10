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

// HealthController handles HTTP requests for Health
type HealthController struct {
	healthService *services.HealthService
}

// NewHealthController creates a new HealthController
func NewHealthController(healthService *services.HealthService) *HealthController {
	return &HealthController{healthService}
}

// CreateHealth creates a new health record
func (c *HealthController) CreateHealth(ctx *gin.Context) {
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
		"id": true, "userId": true, "height": true, "weight": true, "date": true,
		"bmi": true, "bloodType": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var health models.Health
	if err := json.Unmarshal(bodyBytes, &health); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdHealth, err := c.healthService.Create(ctx, &health)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdHealth})
}

// GetHealthByID retrieves a health record by ID
func (c *HealthController) GetHealthByID(ctx *gin.Context) {
	id := ctx.Param("id")
	health, err := c.healthService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "health record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": health})
}

// GetHealthByUserID retrieves a health record by user ID
func (c *HealthController) GetHealthByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	health, err := c.healthService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(health) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Health{},
			"message": "không có dữ liệu nào",
			"note":    "chưa có thông tin sức khỏe nào được ghi nhận",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": health})
}

// GetAllHealthRecords retrieves all health records with pagination
func (c *HealthController) GetAllHealthRecords(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	healthRecords, err := c.healthService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(healthRecords) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Health{},
			"totalCount": 0,
			"notes":      "Không có hồ sơ sức khỏe nào",
			"message":    "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": healthRecords})
}

// UpdateHealth updates a health record
func (c *HealthController) UpdateHealth(ctx *gin.Context) {
	ctx.Param("id")
	var health models.Health
	if err := ctx.ShouldBindJSON(&health); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedHealth, err := c.healthService.Update(ctx, &health)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedHealth})
}

// DeleteHealth deletes a health record by ID
func (c *HealthController) DeleteHealth(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.healthService.Delete(ctx, id); err != nil {
		if err.Error() == "health record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "health record deleted"})
}
