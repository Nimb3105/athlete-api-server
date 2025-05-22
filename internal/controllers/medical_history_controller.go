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

// MedicalHistoryController handles HTTP requests for MedicalHistory
type MedicalHistoryController struct {
	medicalHistoryService *services.MedicalHistoryService
}

// NewMedicalHistoryController creates a new MedicalHistoryController
func NewMedicalHistoryController(medicalHistoryService *services.MedicalHistoryService) *MedicalHistoryController {
	return &MedicalHistoryController{medicalHistoryService}
}

// CreateMedicalHistory creates a new medical history record
func (c *MedicalHistoryController) CreateMedicalHistory(ctx *gin.Context) {
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
		"id": true, "healthId": true, "date": true, "description": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var medicalHistory models.MedicalHistory
	if err := json.Unmarshal(bodyBytes, &medicalHistory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdMedicalHistory, err := c.medicalHistoryService.Create(ctx, &medicalHistory)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdMedicalHistory})
}

// GetMedicalHistoryByID retrieves a medical history by ID
func (c *MedicalHistoryController) GetMedicalHistoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	medicalHistory, err := c.medicalHistoryService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "medical history not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": medicalHistory})
}

// GetMedicalHistoryByHealthID retrieves a medical history by health ID
func (c *MedicalHistoryController) GetMedicalHistoryByHealthID(ctx *gin.Context) {
	healthID := ctx.Param("healthID")
	medicalHistory, err := c.medicalHistoryService.GetByHealthID(ctx, healthID)
	if err != nil {
		if err.Error() == "medical history not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": medicalHistory})
}

// GetAllMedicalHistories retrieves all medical history records with pagination
func (c *MedicalHistoryController) GetAllMedicalHistories(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	medicalHistories, err := c.medicalHistoryService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": medicalHistories})
}

// UpdateMedicalHistory updates a medical history record
func (c *MedicalHistoryController) UpdateMedicalHistory(ctx *gin.Context) {
	ctx.Param("id")
	var medicalHistory models.MedicalHistory
	if err := ctx.ShouldBindJSON(&medicalHistory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMedicalHistory, err := c.medicalHistoryService.Update(ctx, &medicalHistory)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedMedicalHistory})
}

// DeleteMedicalHistory deletes a medical history by ID
func (c *MedicalHistoryController) DeleteMedicalHistory(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.medicalHistoryService.Delete(ctx, id); err != nil {
		if err.Error() == "medical history not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "medical history deleted"})
}