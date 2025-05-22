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

// SportAthleteController xử lý các yêu cầu HTTP cho SportAthlete
type SportAthleteController struct {
	sportAthleteService *services.SportAthleteService
}

// NewSportAthleteController tạo một SportAthleteController mới
func NewSportAthleteController(sportAthleteService *services.SportAthleteService) *SportAthleteController {
	return &SportAthleteController{sportAthleteService}
}

// CreateSportAthlete tạo một sport athlete mới
func (c *SportAthleteController) CreateSportAthlete(ctx *gin.Context) {
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
		"id": true, "sportId": true, "userId": true, "position": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var sportAthlete models.SportAthlete
	if err := json.Unmarshal(bodyBytes, &sportAthlete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdSportAthlete, err := c.sportAthleteService.Create(ctx, &sportAthlete)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdSportAthlete})
}

// GetSportAthleteByID lấy sport athlete theo ID
func (c *SportAthleteController) GetSportAthleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	sportAthlete, err := c.sportAthleteService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "sportAthlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sportAthlete})
}

// GetSportAthleteByUserID lấy sport athlete theo UserID
func (c *SportAthleteController) GetSportAthleteByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	sportAthlete, err := c.sportAthleteService.GetByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "sportAthlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sportAthlete})
}

func (c *SportAthleteController) GetSportAthleteBySportID(ctx *gin.Context) {
	sportID := ctx.Param("sportID")
	sportAthlete, err := c.sportAthleteService.GetByUserID(ctx, sportID)
	if err != nil {
		if err.Error() == "sportAthlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sportAthlete})
}

// GetAllSportAthletes lấy danh sách tất cả sport athlete với phân trang
func (c *SportAthleteController) GetAllSportAthletes(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	sportAthletes, err := c.sportAthleteService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": sportAthletes})
}

// UpdateSportAthlete cập nhật thông tin sport athlete
func (c *SportAthleteController) UpdateSportAthlete(ctx *gin.Context) {
	ctx.Param("id")
	var sportAthlete models.SportAthlete
	if err := ctx.ShouldBindJSON(&sportAthlete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSportAthlete, err := c.sportAthleteService.Update(ctx, &sportAthlete)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedSportAthlete})
}

// DeleteSportAthlete xóa sport athlete theo ID
func (c *SportAthleteController) DeleteSportAthlete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.sportAthleteService.Delete(ctx, id); err != nil {
		if err.Error() == "sportAthlete not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "sportAthlete deleted"})
}

func (c *SportAthleteController) GetAllByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	SportAthletes, err := c.sportAthleteService.GetAllByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": SportAthletes})
}
