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

// AthleteController xử lý các yêu cầu HTTP cho Athlete
type AthleteController struct {
	athleteService *services.AthleteService
}

// NewAthleteController tạo một AthleteController mới
func NewAthleteController(athleteService *services.AthleteService) *AthleteController {
	return &AthleteController{athleteService}
}

// CreateAthlete tạo một athlete mới
func (c *AthleteController) CreateAthlete(ctx *gin.Context) {
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
		"id": true, "userId": true, "athleteType": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var athlete models.Athlete
	if err := json.Unmarshal(bodyBytes, &athlete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdAthlete, err := c.athleteService.Create(ctx, &athlete)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdAthlete})
}

// GetAthleteByID lấy athlete theo ID
func (c *AthleteController) GetAthleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	athlete, err := c.athleteService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "athlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": athlete})
}

// GetAthleteByUserID lấy athlete theo UserID
func (c *AthleteController) GetAthleteByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	athlete, err := c.athleteService.GetByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "athlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": athlete})
}

// GetAllAthletes lấy danh sách tất cả athlete với phân trang
func (c *AthleteController) GetAllAthletes(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	athletes, totalCount, err := c.athleteService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(athletes) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Athlete{},
			"message": "Không có dữ liệu nào",
			"note":    "Chưa có vận động viên nào được ghi nhận",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": athletes, "totalCount": totalCount})
}

// UpdateAthlete cập nhật thông tin athlete
func (c *AthleteController) UpdateAthlete(ctx *gin.Context) {
	ctx.Param("id")
	var athlete models.Athlete
	if err := ctx.ShouldBindJSON(&athlete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAthlete, err := c.athleteService.Update(ctx, &athlete)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedAthlete})
}

// DeleteAthlete xóa athlete theo ID
func (c *AthleteController) DeleteAthlete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.athleteService.Delete(ctx, id); err != nil {
		if err.Error() == "athlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "athlete deleted"})
}
