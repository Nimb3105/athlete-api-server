package controllers

import (
	"be/internal/models"
	"be/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CoachAthleteController handles HTTP requests for CoachAthlete
type CoachAthleteController struct {
	coachAthleteService *services.CoachAthleteService
}

// NewCoachAthleteController creates a new CoachAthleteController
func NewCoachAthleteController(coachAthleteService *services.CoachAthleteService) *CoachAthleteController {
	return &CoachAthleteController{coachAthleteService}
}

// CreateCoachAthlete creates a new coach-athlete relationship
func (c *CoachAthleteController) CreateCoachAthlete(ctx *gin.Context) {
	var bodyBytes []byte
	if rawData, err := ctx.GetRawData(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read request body"})
		return
	} else {
		bodyBytes = rawData
	}

	var tempMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	validFields := map[string]bool{
		"id": true, "coachId": true, "athleteId": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid field: %s", key)})
			return
		}
	}

	var coachAthlete models.CoachAthlete
	if err := json.Unmarshal(bodyBytes, &coachAthlete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot map data to model"})
		return
	}

	createdCoachAthlete, err := c.coachAthleteService.Create(ctx, &coachAthlete)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdCoachAthlete})
}

// GetCoachAthleteByID retrieves a coach-athlete relationship by ID
func (c *CoachAthleteController) GetCoachAthleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	coachAthlete, err := c.coachAthleteService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "coach-athlete relationship not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": coachAthlete})
}

// GetCoachAthleteByAthleteId retrieves a coach-athlete relationship by athlete ID
func (c *CoachAthleteController) GetCoachAthleteByAthleteId(ctx *gin.Context) {
	athleteId := ctx.Param("athleteId")
	if athleteId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Yêu cầu cung cấp athleteId"})
		return
	}
	if _, err := primitive.ObjectIDFromHex(athleteId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "athleteId không hợp lệ"})
		return
	}

	coachAthlete, err := c.coachAthleteService.GetByAthleteId(ctx, athleteId)
	if err != nil {
		if err.Error() == "coach-athlete relationship not found" {
			ctx.JSON(http.StatusOK, gin.H{
				"data":    nil, // hoặc []interface{}{} nếu bạn muốn mảng rỗng
				"message": "Không tìm thấy dữ liệu",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": coachAthlete, "message": "Dữ liệu đã được tìm thấy"})
}

func (c *CoachAthleteController) GetAllByCoachId(ctx *gin.Context) {

	coachId := ctx.Param("userId")
	if coachId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Yêu cầu cung cấp coachId"})
		return
	}
	if _, err := primitive.ObjectIDFromHex(coachId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "coachId không hợp lệ"})
		return
	}

	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	coachAthletes, totalCount, err := c.coachAthleteService.GetAllByCoachId(ctx, coachId, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(coachAthletes) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.CoachAthlete{},
			"note":    "không có mối quan hệ huấn luyện viên - vận động viên nào",
			"message": "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": coachAthletes, "totalCount": totalCount})
}

// GetAllCoachAthletes retrieves all coach-athlete relationships with pagination
func (c *CoachAthleteController) GetAllCoachAthletes(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	coachAthletes, err := c.coachAthleteService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(coachAthletes) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.CoachAthlete{},
			"note":    "không có mối quan hệ huấn luyện viên - vận động viên nào",
			"message": "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": coachAthletes})
}

// UpdateCoachAthlete updates a coach-athlete relationship
func (c *CoachAthleteController) UpdateCoachAthlete(ctx *gin.Context) {
	ctx.Param("id")
	var coachAthlete models.CoachAthlete
	if err := ctx.ShouldBindJSON(&coachAthlete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCoachAthlete, err := c.coachAthleteService.Update(ctx, &coachAthlete)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedCoachAthlete})
}

// DeleteCoachAthlete deletes a coach-athlete relationship by ID
func (c *CoachAthleteController) DeleteCoachAthlete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.coachAthleteService.Delete(ctx, id); err != nil {
		if err.Error() == "coach-athlete relationship not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "coach-athlete relationship deleted"})
}
func (c *CoachAthleteController) DeleteAllByCoachId(ctx *gin.Context) {
	coachId := ctx.Param("coachId")
	err := c.coachAthleteService.DeleteAllByCoachId(ctx, coachId)
	if err != nil {
		if err.Error() == "no coach-athlete relationships found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "all coach-athlete relationships deleted for coachId: " + coachId})
}
