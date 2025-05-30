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

// AchievementController xử lý các yêu cầu HTTP cho Achievement
type AchievementController struct {
	achievementService *services.AchievementService
}

// NewAchievementController tạo một AchievementController mới
func NewAchievementController(achievementService *services.AchievementService) *AchievementController {
	return &AchievementController{achievementService}
}

// CreateAchievement tạo một achievement mới
func (c *AchievementController) CreateAchievement(ctx *gin.Context) {
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
		"id": true, "userId": true, "title": true, "description": true,
		"date": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var achievement models.Achievement
	if err := json.Unmarshal(bodyBytes, &achievement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdAchievement, err := c.achievementService.Create(ctx, &achievement)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdAchievement})
}

// GetAchievementByID lấy achievement theo ID
func (c *AchievementController) GetAchievementByID(ctx *gin.Context) {
	id := ctx.Param("id")
	achievement, err := c.achievementService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "achievement not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": achievement})
}

// GetAchievementByUserID lấy danh sách achievement theo UserID
func (c *AchievementController) GetAchievementByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	achievements, err := c.achievementService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(achievements) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Achievement{},
			"message": "Không có dữ liệu nào",
			"note":    "Chưa có thành tích nào được ghi nhận cho người dùng này",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": achievements})
}

// GetAllAchievements lấy danh sách tất cả achievement với phân trang
func (c *AchievementController) GetAllAchievements(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	achievements, err := c.achievementService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(achievements) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Achievement{},
			"message": "Không có dữ liệu nào",
			"note":    "Chưa có thành tích nào được ghi nhận",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": achievements})
}

// UpdateAchievement cập nhật thông tin achievement
func (c *AchievementController) UpdateAchievement(ctx *gin.Context) {
	var achievement models.Achievement
	if err := ctx.ShouldBindJSON(&achievement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAchievement, err := c.achievementService.Update(ctx, &achievement)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedAchievement})
}

// DeleteAchievement xóa achievement theo ID
func (c *AchievementController) DeleteAchievement(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.achievementService.Delete(ctx, id); err != nil {
		if err.Error() == "achievement not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "achievement deleted"})
}
