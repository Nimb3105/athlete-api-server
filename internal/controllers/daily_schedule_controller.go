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

type CreateDailyScheduleRequest struct {
	models.DailySchedule
	TrainingSchedules []models.CreateTrainingScheduleRequest `json:"trainingSchedules" validate:"required,dive"`
}

type DailyScheduleController struct {
	service *services.DailyScheduleService
}

func NewDailyScheduleController(service *services.DailyScheduleService) *DailyScheduleController {
	return &DailyScheduleController{service: service}
}

func (c *DailyScheduleController) GetAllDailySchedulesByUserId(ctx *gin.Context) {
	userId := ctx.Param("userID")

	schedules, err := c.service.GetAllDailySchedulesByUserId(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": schedules})
}

func (c *DailyScheduleController) GetByCreatorId(ctx *gin.Context) {
	creatorId := ctx.Param("creatorId")

	schedules, err := c.service.GetByCreatorId(ctx, creatorId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": schedules})
}

func (c *DailyScheduleController) Create(ctx *gin.Context) {
	var bodyBytes []byte
	if rawData, err := ctx.GetRawData(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc dữ liệu"})
		return
	} else {
		bodyBytes = rawData
	}

	// Unmarshal sang map để kiểm tra field
	var tempMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không đúng định dạng"})
		return
	}

	// Danh sách field hợp lệ của DailySchedule
	validFields := map[string]bool{
		"id": true, "userId": true, "name": true, "note": true, "createdBy": true,
		"startDate": true, "endDate": true, "createdAt": true, "updatedAt": true,
		"trainingSchedules": true, "sportId": true, // thêm nếu có trong request
	}

	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Trường không hợp lệ: %s", key),
			})
			return
		}
	}

	// Ánh xạ dữ liệu
	var req CreateDailyScheduleRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Lỗi ánh xạ JSON: %v", err),
		})
	}

	// Kiểm tra userId có hợp lệ không
	if req.UserId.IsZero() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Trường 'userId' không hợp lệ hoặc thiếu"})
		return
	}

	// Gọi service để tạo
	created, err := c.service.Create(ctx, &req.DailySchedule, req.TrainingSchedules, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": created})
}

func (c *DailyScheduleController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	dailySchedule, err := c.service.GetById(ctx, id)
	if err != nil {
		if err.Error() == "dailySchedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": dailySchedule})
}

func (c *DailyScheduleController) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	day := ctx.Param("day")

	dailySchedule, err := c.service.GetByUserID(ctx, day, userID)
	if err != nil {
		if err.Error() == "dailySchedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": dailySchedule})
}

func (c *DailyScheduleController) GetAll(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	dailySchedule, totalCount, err := c.service.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(dailySchedule) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.DailySchedule{},
			"totalCount": 0,
			"notes":      "không có dailyschedule nao",
			"message":    "chưa có dữ liệu nào",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":       dailySchedule,
		"totalCount": totalCount,
	})
}

func (c *DailyScheduleController) Update(ctx *gin.Context) {
	ctx.Param("id")
	var dailySchedule models.DailySchedule
	if err := ctx.ShouldBindJSON(&dailySchedule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateDailySchedule, err := c.service.Update(ctx, &dailySchedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updateDailySchedule})
}

func (c *DailyScheduleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx, id); err != nil {
		if err.Error() == "dailySchedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "daily schedule deleted"})
}
