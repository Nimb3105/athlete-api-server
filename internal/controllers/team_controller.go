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

// TeamController handles HTTP requests for Team
type TeamController struct {
	teamService *services.TeamService
}

// NewTeamController creates a new TeamController
func NewTeamController(teamService *services.TeamService) *TeamController {
	return &TeamController{teamService}
}

// CreateTeam creates a new team
func (c *TeamController) CreateTeam(ctx *gin.Context) {
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
		"id": true, "name": true, "sportId": true, "description": true,
		"createdBy": true, "createdDate": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var team models.Team
	if err := json.Unmarshal(bodyBytes, &team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdTeam, err := c.teamService.Create(ctx, &team)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdTeam})
}

// GetTeamByID retrieves a team by ID
func (c *TeamController) GetTeamByID(ctx *gin.Context) {
	id := ctx.Param("id")
	team, err := c.teamService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "team not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": team})
}

// GetTeamsBySportID retrieves teams by sport ID
func (c *TeamController) GetTeamsBySportID(ctx *gin.Context) {
	sportID := ctx.Param("sportID")
	teams, err := c.teamService.GetBySportID(ctx, sportID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(teams) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Team{},
			"message": "ko có đội nào cho môn thể thao này",
			"note":    "Chưa có đội nào được ghi nhận cho môn thể thao này",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": teams})
}

// GetAllTeams retrieves all teams with pagination
func (c *TeamController) GetAllTeams(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	teams, err := c.teamService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(teams) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Team{},
			"totalCount": 0,
			"notes":      "Không có đội nào",
			"message":    "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": teams})
}

// UpdateTeam updates a team
func (c *TeamController) UpdateTeam(ctx *gin.Context) {
	ctx.Param("id")
	var team models.Team
	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTeam, err := c.teamService.Update(ctx, &team)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedTeam})
}

// DeleteTeam deletes a team by ID
func (c *TeamController) DeleteTeam(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.teamService.Delete(ctx, id); err != nil {
		if err.Error() == "team not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "team deleted"})
}
