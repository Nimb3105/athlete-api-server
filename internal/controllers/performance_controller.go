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

// PerformanceController handles HTTP requests for Performance
type PerformanceController struct {
	performanceService *services.PerformanceService
}

// NewPerformanceController creates a new PerformanceController
func NewPerformanceController(performanceService *services.PerformanceService) *PerformanceController {
	return &PerformanceController{performanceService}
}

// CreatePerformance creates a new performance record
func (c *PerformanceController) CreatePerformance(ctx *gin.Context) {
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
		"id": true, "userId": true, "scheduleId": true, "value": true,
		"date": true, "metricType": true, "notes": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var performance models.Performance
	if err := json.Unmarshal(bodyBytes, &performance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdPerformance, err := c.performanceService.Create(ctx, &performance)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdPerformance})
}

// GetPerformanceByID retrieves a performance record by ID
func (c *PerformanceController) GetPerformanceByID(ctx *gin.Context) {
	id := ctx.Param("id")
	performance, err := c.performanceService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "performance not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": performance})
}

// GetPerformancesByUserID retrieves performance records by user ID
func (c *PerformanceController) GetPerformancesByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	performances, err := c.performanceService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(performances) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Performance{},
			"message": "không có dữ liệu nào",
			"notes":   "Không có hiệu suất nào cho người dùng này"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": performances})
}

// GetAllPerformances retrieves all performance records with pagination
func (c *PerformanceController) GetAllPerformances(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	performances, err := c.performanceService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(performances) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Performance{},
			"totalCount": 0,
			"notes":      "Không có hiệu suất nào",
			"message":    "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": performances})
}

// UpdatePerformance updates a performance record
func (c *PerformanceController) UpdatePerformance(ctx *gin.Context) {
	ctx.Param("id")
	var performance models.Performance
	if err := ctx.ShouldBindJSON(&performance); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPerformance, err := c.performanceService.Update(ctx, &performance)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPerformance})
}

// DeletePerformance deletes a performance record by ID
func (c *PerformanceController) DeletePerformance(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.performanceService.Delete(ctx, id); err != nil {
		if err.Error() == "performance not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "performance deleted"})
}
