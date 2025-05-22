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

// SportController xử lý các yêu cầu HTTP cho Sport
type SportController struct {
	sportService *services.SportService
}

// NewSportController tạo một SportController mới
func NewSportController(sportService *services.SportService) *SportController {
	return &SportController{sportService}
}

// CreateSport tạo một sport mới
func (c *SportController) CreateSport(ctx *gin.Context) {
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
		"id": true, "name": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var sport models.Sport
	if err := json.Unmarshal(bodyBytes, &sport); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdSport, err := c.sportService.Create(ctx, &sport)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdSport})
}

// GetSportByID lấy sport theo ID
func (c *SportController) GetSportByID(ctx *gin.Context) {
	id := ctx.Param("id")
	sport, err := c.sportService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "sport not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sport})
}

// GetAllSports lấy danh sách tất cả sport với phân trang
func (c *SportController) GetAllSports(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	sports, err := c.sportService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sports})
}

// UpdateSport cập nhật thông tin sport
func (c *SportController) UpdateSport(ctx *gin.Context) {
	ctx.Param("id")
	var sport models.Sport
	if err := ctx.ShouldBindJSON(&sport); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSport, err := c.sportService.Update(ctx, &sport)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedSport})
}

// DeleteSport xóa sport theo ID
func (c *SportController) DeleteSport(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.sportService.Delete(ctx, id); err != nil {
		if err.Error() == "sport not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "sport deleted"})
}
