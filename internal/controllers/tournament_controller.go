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

// TournamentController handles HTTP requests for Tournament
type TournamentController struct {
	tournamentService *services.TournamentService
}

// NewTournamentController creates a new TournamentController
func NewTournamentController(tournamentService *services.TournamentService) *TournamentController {
	return &TournamentController{tournamentService}
}

// CreateTournament creates a new tournament
func (c *TournamentController) CreateTournament(ctx *gin.Context) {
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
		"id": true, "name": true, "location": true, "startDate": true,
		"endDate": true, "level": true, "organizer": true, "description": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var tournament models.Tournament
	if err := json.Unmarshal(bodyBytes, &tournament); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdTournament, err := c.tournamentService.Create(ctx, &tournament)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdTournament})
}

// GetTournamentByID retrieves a tournament by ID
func (c *TournamentController) GetTournamentByID(ctx *gin.Context) {
	id := ctx.Param("id")
	tournament, err := c.tournamentService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "tournament not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tournament})
}

// GetAllTournaments retrieves all tournaments with pagination
func (c *TournamentController) GetAllTournaments(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	tournaments, err := c.tournamentService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(tournaments) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Tournament{},
			"totalCount": 0,
			"message": "không có dữ liệu nào",
			"notes": "bạn có thể tạo một giải đấu mới bằng cách sử dụng API tạo giải đấu",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tournaments})
}

// UpdateTournament updates a tournament
func (c *TournamentController) UpdateTournament(ctx *gin.Context) {
	ctx.Param("id")
	var tournament models.Tournament
	if err := ctx.ShouldBindJSON(&tournament); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTournament, err := c.tournamentService.Update(ctx, &tournament)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedTournament})
}

// DeleteTournament deletes a tournament by ID
func (c *TournamentController) DeleteTournament(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.tournamentService.Delete(ctx, id); err != nil {
		if err.Error() == "tournament not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "tournament deleted"})
}