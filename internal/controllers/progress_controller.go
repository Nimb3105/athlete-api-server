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

// ProgressController handles HTTP requests for Progress
type ProgressController struct {
	progressService *services.ProgressService
}

// NewProgressController creates a new ProgressController
func NewProgressController(progressService *services.ProgressService) *ProgressController {
	return &ProgressController{progressService}
}

// CreateProgress creates a new progress record
func (c *ProgressController) CreateProgress(ctx *gin.Context) {
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
		"id": true, "userId": true, "metricType": true, "value": true,
		"date": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var progress models.Progress
	if err := json.Unmarshal(bodyBytes, &progress); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdProgress, err := c.progressService.Create(ctx, &progress)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdProgress})
}


// GetProgressByID retrieves a progress record by ID
func (c *ProgressController) GetProgressByID(ctx *gin.Context) {
	id := ctx.Param("id")
	progress, err := c.progressService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "progress not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": progress})
}

// GetProgressesByUserID retrieves progress records by user ID
func (c *ProgressController) GetProgressesByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	progresses, err := c.progressService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": progresses})
}

// GetAllProgresses retrieves all progress records with pagination
func (c *ProgressController) GetAllProgresses(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	progresses, err := c.progressService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": progresses})
}

// UpdateProgress updates a progress record
func (c *ProgressController) UpdateProgress(ctx *gin.Context) {
	ctx.Param("id")
	var progress models.Progress
	if err := ctx.ShouldBindJSON(&progress); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProgress, err := c.progressService.Update(ctx, &progress)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedProgress})
}

// DeleteProgress deletes a progress record by ID
func (c *ProgressController) DeleteProgress(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.progressService.Delete(ctx, id); err != nil {
		if err.Error() == "progress not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "progress deleted"})
}