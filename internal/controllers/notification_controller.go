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

// NotificationController xử lý các yêu cầu HTTP cho Notification
type NotificationController struct {
	notificationService *services.NotificationService
}

// NewNotificationController tạo một NotificationController mới
func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{notificationService}
}

// GetAllNotifications lấy tất cả notifications với phân trang
func (c *NotificationController) GetAllNotifications(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)
	notifications, err := c.notificationService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": notifications})
}

// CreateNotification tạo một notification mới
func (c *NotificationController) CreateNotification(ctx *gin.Context) {
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
		"sentDate": true, "status": true, "type": true, "content": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var notification models.Notification
	if err := json.Unmarshal(bodyBytes, &notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdNotification, err := c.notificationService.Create(ctx, &notification)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdNotification})
}

// GetNotificationByID lấy notification theo ID
func (c *NotificationController) GetNotificationByID(ctx *gin.Context) {
	id := ctx.Param("id")
	notification, err := c.notificationService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "notification not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": notification})
}

// GetNotificationsByUserID lấy danh sách notification theo UserID với phân trang
func (c *NotificationController) GetNotificationsByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	notifications, err := c.notificationService.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(notifications) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Notification{},
			"message": "Không có dữ liệu nào",
			"notes":   "Không có thông báo nào cho người dùng này",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": notifications})
}

// UpdateNotification cập nhật thông tin notification
func (c *NotificationController) UpdateNotification(ctx *gin.Context) {
	var notification models.Notification
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNotification, err := c.notificationService.Update(ctx, &notification)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedNotification})
}

// DeleteNotification xóa notification theo ID
func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.notificationService.Delete(ctx, id); err != nil {
		if err.Error() == "notification not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "notification deleted"})
}
