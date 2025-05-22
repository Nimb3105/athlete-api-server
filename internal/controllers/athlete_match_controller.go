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

// AthleteMatchController xử lý các yêu cầu HTTP cho AthleteMatch
type AthleteMatchController struct {
	athleteMatchService *services.AthleteMatchService
}

// NewAthleteMatchController tạo một AthleteMatchController mới
func NewAthleteMatchController(athleteMatchService *services.AthleteMatchService) *AthleteMatchController {
	return &AthleteMatchController{athleteMatchService}
}

// CreateAthleteMatch tạo một athlete match mới
func (c *AthleteMatchController) CreateAthleteMatch(ctx *gin.Context) {
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
		"id": true, "matchId": true, "userId": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var athleteMatch models.AthleteMatch
	if err := json.Unmarshal(bodyBytes, &athleteMatch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdAthleteMatch, err := c.athleteMatchService.Create(ctx, &athleteMatch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdAthleteMatch})
}

// GetAthleteMatchByID lấy athlete match theo ID
func (c *AthleteMatchController) GetAthleteMatchByID(ctx *gin.Context) {
	id := ctx.Param("id")
	athleteMatch, err := c.athleteMatchService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "athlete match not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": athleteMatch})
}

// GetAthleteMatchByUserID lấy danh sách athlete match theo UserID
func (c *AthleteMatchController) GetAthleteMatchByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	athleteMatches, err := c.athleteMatchService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": athleteMatches})
}

// GetAllAthleteMatches lấy danh sách tất cả athlete match với phân trang
func (c *AthleteMatchController) GetAllAthleteMatches(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	athleteMatches, err := c.athleteMatchService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": athleteMatches})
}

// UpdateAthleteMatch cập nhật thông tin athlete match
func (c *AthleteMatchController) UpdateAthleteMatch(ctx *gin.Context) {
	var athleteMatch models.AthleteMatch
	if err := ctx.ShouldBindJSON(&athleteMatch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAthleteMatch, err := c.athleteMatchService.Update(ctx, &athleteMatch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedAthleteMatch})
}

// DeleteAthleteMatch xóa athlete match theo ID
func (c *AthleteMatchController) DeleteAthleteMatch(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.athleteMatchService.Delete(ctx, id); err != nil {
		if err.Error() == "athlete match not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "athlete match deleted"})
}