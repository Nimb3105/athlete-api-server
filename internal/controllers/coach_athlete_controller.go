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

func (c *CoachAthleteController) GetAllByAthleteId(ctx *gin.Context) {

	athleteId := ctx.Param("userId")
	coachAthletes, err := c.coachAthleteService.GetAllByAthleteID(ctx, athleteId)
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
