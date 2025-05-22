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

// CoachController xử lý các yêu cầu HTTP cho Coach
type CoachController struct {
	coachService *services.CoachService
}

// NewCoachController tạo một CoachController mới
func NewCoachController(coachService *services.CoachService) *CoachController {
	return &CoachController{coachService}
}

// CreateCoach tạo một coach mới
func (c *CoachController) CreateCoach(ctx *gin.Context) {
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
		"id": true, "userId": true, "experience": true, "specialization": true,
		"level": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var coach models.Coach
	if err := json.Unmarshal(bodyBytes, &coach); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdCoach, err := c.coachService.Create(ctx, &coach)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdCoach})
}

// GetCoachByID lấy coach theo ID
func (c *CoachController) GetCoachByID(ctx *gin.Context) {
	id := ctx.Param("id")
	coach, err := c.coachService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "coach not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": coach})
}

func (c *CoachController) GetCoachByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	coach, err := c.coachService.GetByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "coach not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": coach})
}

// GetAllCoaches lấy danh sách tất cả coach với phân trang
func (c *CoachController) GetAllCoaches(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	coaches, err := c.coachService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": coaches})
}

// UpdateCoach cập nhật thông tin coach
func (c *CoachController) UpdateCoach(ctx *gin.Context) {
	ctx.Param("id")
	var coach models.Coach
	if err := ctx.ShouldBindJSON(&coach); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCoach, err := c.coachService.Update(ctx, &coach)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedCoach})
}

// DeleteCoach xóa coach theo ID
func (c *CoachController) DeleteCoach(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.coachService.Delete(ctx, id); err != nil {
		if err.Error() == "coach not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "coach deleted"})
}
