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

// MatchScheduleController handles HTTP requests for MatchSchedule
type MatchScheduleController struct {
	matchScheduleService *services.MatchScheduleService
}

// NewMatchScheduleController creates a new MatchScheduleController
func NewMatchScheduleController(matchScheduleService *services.MatchScheduleService) *MatchScheduleController {
	return &MatchScheduleController{matchScheduleService}
}

// CreateMatchSchedule creates a new match schedule
func (c *MatchScheduleController) CreateMatchSchedule(ctx *gin.Context) {
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
		"id": true, "tournamentId": true, "date": true, "location": true,
		"opponent": true, "matchType": true, "status": true, "round": true, "score": true,
		"notes": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var matchSchedule models.MatchSchedule
	if err := json.Unmarshal(bodyBytes, &matchSchedule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdMatchSchedule, err := c.matchScheduleService.Create(ctx, &matchSchedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdMatchSchedule})
}

// GetMatchScheduleByID retrieves a match schedule by ID
func (c *MatchScheduleController) GetMatchScheduleByID(ctx *gin.Context) {
	id := ctx.Param("id")
	matchSchedule, err := c.matchScheduleService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "match schedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": matchSchedule})
}

// GetMatchScheduleByTournamentID retrieves a match schedule by tournament ID
func (c *MatchScheduleController) GetMatchScheduleByTournamentID(ctx *gin.Context) {
	tournamentID := ctx.Param("tournamentID")
	matchSchedule, err := c.matchScheduleService.GetByTournamentID(ctx, tournamentID)
	if err != nil {
		if err.Error() == "match schedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": matchSchedule})
}

// GetAllMatchSchedules retrieves all match schedules with pagination
func (c *MatchScheduleController) GetAllMatchSchedules(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	matchSchedules, err := c.matchScheduleService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(matchSchedules) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.MatchSchedule{},
			"totalCount": 0,
			"notes":      "Không có lịch thi đấu nào",
			"message":    "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": matchSchedules})
}

// UpdateMatchSchedule updates a match schedule
func (c *MatchScheduleController) UpdateMatchSchedule(ctx *gin.Context) {
	ctx.Param("id")
	var matchSchedule models.MatchSchedule
	if err := ctx.ShouldBindJSON(&matchSchedule); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMatchSchedule, err := c.matchScheduleService.Update(ctx, &matchSchedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedMatchSchedule})
}

// DeleteMatchSchedule deletes a match schedule by ID
func (c *MatchScheduleController) DeleteMatchSchedule(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.matchScheduleService.Delete(ctx, id); err != nil {
		if err.Error() == "match schedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "match schedule deleted"})
}
