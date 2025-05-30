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

// UserMatchController xử lý các yêu cầu HTTP cho AthleteMatch
type UserMatchController struct {
	userMatchService *services.UserMatchService
}

// NewUserMatchController tạo một UserMatchController mới
func NewUserMatchController(userMatchService *services.UserMatchService) *UserMatchController {
	return &UserMatchController{userMatchService: userMatchService}
}

// CreateAthleteMatch tạo một athlete match mới
func (c *UserMatchController) CreateUserMatch(ctx *gin.Context) {
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

	var userMatch models.UserMatch
	if err := json.Unmarshal(bodyBytes, &userMatch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdUserMatch, err := c.userMatchService.Create(ctx, &userMatch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdUserMatch})
}

// GetAthleteMatchByID lấy athlete match theo ID
func (c *UserMatchController) GetUserMatchByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userMatch, err := c.userMatchService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "user match not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": userMatch})
}

// GetAthleteMatchByUserID lấy danh sách athlete match theo UserID
func (c *UserMatchController) GetUserMatchByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	userMatches, err := c.userMatchService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(userMatches) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.UserMatch{},
			"message": "không có dữ liệu nào",
			"notes":   "bạn có thể tạo user match mới nếu bạn muốn",
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": userMatches})
}

// GetAllAthleteMatches lấy danh sách tất cả athlete match với phân trang
func (c *UserMatchController) GetAllUserMatches(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	userMatches, err := c.userMatchService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": userMatches})
}

// UpdateAthleteMatch cập nhật thông tin athlete match
func (c *UserMatchController) UpdateUserMatch(ctx *gin.Context) {
	var userMatch models.UserMatch
	if err := ctx.ShouldBindJSON(&userMatch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUserMatch, err := c.userMatchService.Update(ctx, &userMatch)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedUserMatch})
}

// DeleteAthleteMatch xóa athlete match theo ID
func (c *UserMatchController) DeleteUserMatch(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.userMatchService.Delete(ctx, id); err != nil {
		if err.Error() == "athlete match not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "athlete match deleted"})
}
