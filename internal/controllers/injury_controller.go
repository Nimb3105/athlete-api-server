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

// InjuryController handles HTTP requests for Injury
type InjuryController struct {
	injuryService *services.InjuryService
}

// NewInjuryController creates a new InjuryController
func NewInjuryController(injuryService *services.InjuryService) *InjuryController {
	return &InjuryController{injuryService}
}

// CreateInjury creates a new injury record
func (c *InjuryController) CreateInjury(ctx *gin.Context) {
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
		"id": true, "userId": true, "type": true, "date": true, "severity": true,
		"locationOnBody": true, "causeOfInjury": true, "recoveryStatus": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var injury models.Injury
	if err := json.Unmarshal(bodyBytes, &injury); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdInjury, err := c.injuryService.Create(ctx, &injury)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdInjury})
}

// GetInjuryByID retrieves an injury by ID
func (c *InjuryController) GetInjuryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	injury, err := c.injuryService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "injury not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": injury})
}

// GetInjuryByUserID retrieves an injury by user ID
func (c *InjuryController) GetInjuryByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	injury, err := c.injuryService.GetByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "injury not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": injury})
}

// GetAllInjuries retrieves all injury records with pagination
func (c *InjuryController) GetAllInjuries(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	injuries, err := c.injuryService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": injuries})
}

// UpdateInjury updates an injury record
func (c *InjuryController) UpdateInjury(ctx *gin.Context) {
	ctx.Param("id")
	var injury models.Injury
	if err := ctx.ShouldBindJSON(&injury); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedInjury, err := c.injuryService.Update(ctx, &injury)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedInjury})
}

// DeleteInjury deletes an injury by ID
func (c *InjuryController) DeleteInjury(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.injuryService.Delete(ctx, id); err != nil {
		if err.Error() == "injury not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "injury deleted"})
}