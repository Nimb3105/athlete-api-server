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

type TrainingScheduleUserController struct {
	service *services.TrainingScheduleUserService
}

func NewTrainingScheduleUserController(service *services.TrainingScheduleUserService) *TrainingScheduleUserController {
	return &TrainingScheduleUserController{service}
}

func (c *TrainingScheduleUserController) CreateTrainingScheduleUser(ctx *gin.Context) {
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
		"id": true, "scheduleId": true, "userId": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var trainingScheduleUser models.TrainingScheduleUser
	if err := json.Unmarshal(bodyBytes, &trainingScheduleUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdTrainingScheduleUser, err := c.service.Create(ctx, &trainingScheduleUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdTrainingScheduleUser})
}

func (c *TrainingScheduleUserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	scheduleAthlete, err := c.service.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "training schedule athlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": scheduleAthlete})
}

func (c *TrainingScheduleUserController) GetByScheduleID(ctx *gin.Context) {
	scheduleID := ctx.Param("scheduleId")
	scheduleAthletes, err := c.service.GetByScheduleID(ctx, scheduleID)
	if err != nil {
		if err.Error() == "invalid schedule ID" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(scheduleAthletes) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.TrainingScheduleUser{},
			"message": "không có dữ liệu nào",
			"notes":   "bạn có thể tạo mới một lịch tập cho người dùng này",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": scheduleAthletes})
}

func (c *TrainingScheduleUserController) GetByUserID(ctx *gin.Context) {
	userID := ctx.Param("userId")
	scheduleAthletes, err := c.service.GetByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "invalid user ID" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(scheduleAthletes) == 0 {

		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.TrainingScheduleUser{},
			"message": "không có dữ liệu nào",
			"notes":   "bạn có thể tạo mới một lịch tập cho người dùng này",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": scheduleAthletes})
}

func (c *TrainingScheduleUserController) GetAll(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	scheduleAthletes, err := c.service.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(scheduleAthletes) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.TrainingScheduleUser{},
			"totalCount": 0,
			"message":    "không có lịch tập nào",
			"notes":      "không có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": scheduleAthletes})
}

func (c *TrainingScheduleUserController) Update(ctx *gin.Context) {
	ctx.Param("id")
	var scheduleAthlete models.TrainingScheduleUser
	if err := ctx.ShouldBindJSON(&scheduleAthlete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedScheduleAthlete, err := c.service.Update(ctx, &scheduleAthlete)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedScheduleAthlete})
}

func (c *TrainingScheduleUserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(ctx, id); err != nil {
		if err.Error() == "training schedule athlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "training schedule athlete deleted"})
}

func (c *TrainingScheduleUserController) GetAllTrainingScheduleUserByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	TrainingScheduleUsers, err := c.service.GetAllByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(TrainingScheduleUsers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.TrainingScheduleUser{},
			"message": "không có lịch tập nào cho người dùng này",
			"notes":   "bạn có thể tạo mới một lịch tập cho người dùng này",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": TrainingScheduleUsers})
}
