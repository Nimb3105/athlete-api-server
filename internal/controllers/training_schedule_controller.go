package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"be/internal/models"
	"be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TrainingScheduleController struct {
	service *services.TrainingScheduleService
}

func NewTrainingScheduleController(service *services.TrainingScheduleService) *TrainingScheduleController {
	return &TrainingScheduleController{service}
}

func (c *TrainingScheduleController) CreateTrainingSchedule(ctx *gin.Context) {

	validate := validator.New()
	validate.RegisterValidation("hex24", hex24Validator)

	var req models.CreateTrainingScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{Error: "Dữ liệu không hợp lệ"})
		return
	}

	if err := validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{Error: fmt.Sprintf("Xác thực thất bại: %v", err)})
		return
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	createdTrainingSchedule, err := c.service.Create(ctxTimeout, &req.TrainingSchedule, req.TrainingExercise)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, APIResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, APIResponse{Data: createdTrainingSchedule})
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

func (c *TrainingScheduleController) GetAllByDailyScheduleId(ctx *gin.Context) {
	dailyScheduleId := ctx.Param("dailyScheduleId")
	date := ctx.Param("date")

	TrainingSchedules, err := c.service.GetAllByDailyScheduleId(ctx, dailyScheduleId, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(TrainingSchedules) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.TrainingSchedule{},
			"notes":   "không có lịch tập nào",
			"message": "chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": TrainingSchedules,
	})
}

func (c *TrainingScheduleController) GetAll(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	schedules, err := c.service.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(schedules) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.TrainingSchedule{},
			"totalCount": 0,
			"message":    "Không có lịch tập nào",
			"notes":      "không có dữ liệu nào",
		})
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
