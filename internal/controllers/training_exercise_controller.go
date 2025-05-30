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

type TrainingExerciseController struct {
	service *services.TrainingExerciseService
}

func NewTrainingExerciseController(service *services.TrainingExerciseService) *TrainingExerciseController {
	return &TrainingExerciseController{service}
}

func (c *TrainingExerciseController) CreateTrainingExercise(ctx *gin.Context) {
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
		"id": true, "scheduleId": true, "exerciseId": true, "order": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var trainingExercise models.TrainingExercise
	if err := json.Unmarshal(bodyBytes, &trainingExercise); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdTrainingExercise, err := c.service.Create(ctx, &trainingExercise)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdTrainingExercise})
}

func (c *TrainingExerciseController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	trainingExercise, err := c.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "training exercise not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": trainingExercise})
}

func (c *TrainingExerciseController) GetByScheduleID(ctx *gin.Context) {
	scheduleID := ctx.Param("scheduleId")
	trainingExercises, err := c.service.GetByScheduleID(ctx, scheduleID)
	if err != nil {
		if err.Error() == "invalid schedule ID" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(trainingExercises) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.TrainingExercise{},
			"message": "không có dữ liệu nào",
			"notes":   "Không có bài tập nào được tìm thấy cho lịch tập này",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": trainingExercises})
}

func (c *TrainingExerciseController) GetAll(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	trainingExercises, err := c.service.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(trainingExercises) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"data": []models.TrainingExercise{}, "message": "không có dữ liệu nào", " notes": "Bạn có thể tạo bài tập mới bằng cách sử dụng API tạo bài tập"})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": trainingExercises})
}

func (c *TrainingExerciseController) Update(ctx *gin.Context) {
	ctx.Param("id")
	var trainingExercise models.TrainingExercise
	if err := ctx.ShouldBindJSON(&trainingExercise); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTrainingExercise, err := c.service.Update(ctx, &trainingExercise)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedTrainingExercise})
}

func (c *TrainingExerciseController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx, id); err != nil {
		if err.Error() == "training exercise not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "training exercise deleted"})
}

func (c *TrainingExerciseController) GetAllTrainingExerciseByScheduleId(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	TrainingExercise, err := c.service.GetByScheduleID(ctx, scheduleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": TrainingExercise})
}
