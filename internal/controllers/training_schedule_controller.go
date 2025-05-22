package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"be/internal/models"
	"be/internal/services"

	"github.com/gin-gonic/gin"
)

type TrainingScheduleController struct {
	service *services.TrainingScheduleService
}

func NewTrainingScheduleController(service *services.TrainingScheduleService) *TrainingScheduleController {
	return &TrainingScheduleController{service}
}

func (c *TrainingScheduleController) CreateTrainingSchedule(ctx *gin.Context) {
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
		"id": true, "date": true, "startTime": true, "endTime": true,
		"status": true, "location": true, "type": true, "notes": true,
		"createdBy": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var trainingSchedule models.TrainingSchedule
	if err := json.Unmarshal(bodyBytes, &trainingSchedule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdTrainingSchedule, err := c.service.Create(ctx, &trainingSchedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdTrainingSchedule})
}

func (c *TrainingScheduleController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	schedule, err := c.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "training schedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": schedule})
}

func (c *TrainingScheduleController) GetAll(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	schedules, err := c.service.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": schedules})
}

func (c *TrainingScheduleController) Update(ctx *gin.Context) {
	ctx.Param("id")
	var schedule models.TrainingSchedule
	if err := ctx.ShouldBindJSON(&schedule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSchedule, err := c.service.Update(ctx, &schedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedSchedule})
}

func (c *TrainingScheduleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx, id); err != nil {
		if err.Error() == "training schedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "training schedule deleted"})
}
