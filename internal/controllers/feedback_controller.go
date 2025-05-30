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

// FeedbackController xử lý các yêu cầu HTTP cho Feedback
type FeedbackController struct {
	feedbackService *services.FeedbackService
}

// NewFeedbackController tạo một FeedbackController mới
func NewFeedbackController(feedbackService *services.FeedbackService) *FeedbackController {
	return &FeedbackController{feedbackService}
}

// CreateFeedback tạo một feedback mới
func (c *FeedbackController) CreateFeedback(ctx *gin.Context) {
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
		"id": true, "userId": true, "scheduleId": true, "content": true,
		"date": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var feedback models.Feedback
	if err := json.Unmarshal(bodyBytes, &feedback); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdFeedback, err := c.feedbackService.Create(ctx, &feedback)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdFeedback})
}

// GetFeedbackByID lấy feedback theo ID
func (c *FeedbackController) GetFeedbackByID(ctx *gin.Context) {
	id := ctx.Param("id")
	feedback, err := c.feedbackService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "feedback not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": feedback})
}

// GetFeedbackByUserID lấy danh sách feedback theo UserID
func (c *FeedbackController) GetFeedbackByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	feedbacks, err := c.feedbackService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(feedbacks) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Feedback{},
			"notes":   "Không có feedback nào cho người dùng này",
			"message": "khong có dữ liệu nào"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": feedbacks})
}

// GetAllFeedbacks lấy danh sách tất cả feedback với phân trang
func (c *FeedbackController) GetAllFeedbacks(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	feedbacks, err := c.feedbackService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(feedbacks) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Feedback{},
			"notes":   "Không có feedback nào",
			"message": "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": feedbacks})
}

// UpdateFeedback cập nhật thông tin feedback
func (c *FeedbackController) UpdateFeedback(ctx *gin.Context) {
	var feedback models.Feedback
	if err := ctx.ShouldBindJSON(&feedback); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedFeedback, err := c.feedbackService.Update(ctx, &feedback)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedFeedback})
}

// DeleteFeedback xóa feedback theo ID
func (c *FeedbackController) DeleteFeedback(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.feedbackService.Delete(ctx, id); err != nil {
		if err.Error() == "feedback not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "feedback deleted"})
}
