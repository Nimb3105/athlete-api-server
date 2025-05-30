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

// ReminderController xử lý các yêu cầu HTTP cho Reminder
type ReminderController struct {
	reminderService *services.ReminderService
}

// NewReminderController tạo một ReminderController mới
func NewReminderController(reminderService *services.ReminderService) *ReminderController {
	return &ReminderController{reminderService}
}

func (c *ReminderController) GetAllReminders(ctx *gin.Context) {
	// Lấy danh sách tất cả reminders với phân trang
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	reminders, err := c.reminderService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reminders})
}

// CreateReminder tạo một reminder mới
func (c *ReminderController) CreateReminder(ctx *gin.Context) {
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
		"id": true, "userId": true, "scheduleId": true, "nutritionPlanId": true,
		"reminderTime": true, "reminderDate": true, "content": true, "status": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var reminder models.Reminder
	if err := json.Unmarshal(bodyBytes, &reminder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdReminder, err := c.reminderService.Create(ctx, &reminder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdReminder})
}

// GetReminderByID lấy reminder theo ID
func (c *ReminderController) GetReminderByID(ctx *gin.Context) {
	id := ctx.Param("id")
	reminder, err := c.reminderService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "reminder not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reminder})
}

// GetRemindersByUserID lấy danh sách reminder theo UserID với phân trang
func (c *ReminderController) GetRemindersByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	reminders, err := c.reminderService.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(reminders) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":  []models.Reminder{},
			"notes": "Không có reminder nào cho người dùng này",
			"message": "Không có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": reminders})
}

// UpdateReminder cập nhật thông tin reminder
func (c *ReminderController) UpdateReminder(ctx *gin.Context) {
	var reminder models.Reminder
	if err := ctx.ShouldBindJSON(&reminder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedReminder, err := c.reminderService.Update(ctx, &reminder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedReminder})
}

// DeleteReminder xóa reminder theo ID
func (c *ReminderController) DeleteReminder(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.reminderService.Delete(ctx, id); err != nil {
		if err.Error() == "reminder not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "reminder deleted"})
}
