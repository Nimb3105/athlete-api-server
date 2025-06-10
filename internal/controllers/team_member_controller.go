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

// TeamMemberController handles HTTP requests for TeamMember
type TeamMemberController struct {
	teamMemberService *services.TeamMemberService
}

// NewTeamMemberController creates a new TeamMemberController
func NewTeamMemberController(teamMemberService *services.TeamMemberService) *TeamMemberController {
	return &TeamMemberController{teamMemberService}
}

// CreateTeamMember creates a new team member
func (c *TeamMemberController) CreateTeamMember(ctx *gin.Context) {
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
		"id": true, "teamId": true, "userId": true,"dateJoined":true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var teamMember models.TeamMember
	if err := json.Unmarshal(bodyBytes, &teamMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdTeamMember, err := c.teamMemberService.Create(ctx, &teamMember)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdTeamMember})
}
// GetTeamMemberByID retrieves a team member by ID
func (c *TeamMemberController) GetTeamMemberByID(ctx *gin.Context) {
	id := ctx.Param("id")
	teamMember, err := c.teamMemberService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "team member not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": teamMember})
}

// GetTeamMembersByTeamID retrieves team members by team ID
func (c *TeamMemberController) GetTeamMembersByTeamID(ctx *gin.Context) {
	teamID := ctx.Param("teamID")
	teamMembers, err := c.teamMemberService.GetByTeamID(ctx, teamID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(teamMembers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data": []models.TeamMember{},
			"message":"không có dữ liệu nào",
			"note": "Bạn có thể tạo team member mới bằng cách gửi yêu cầu POST tới /team_members",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": teamMembers})
}

// GetTeamMembersByUserID retrieves team members by user ID
func (c *TeamMemberController) GetTeamMembersByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	teamMembers, err := c.teamMemberService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(teamMembers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data": []models.TeamMember{},
			"message": "không có dữ liệu nào",
			"note": "Bạn có thể tạo team member mới bằng cách gửi yêu cầu POST tới /team_members",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": teamMembers})
}

// GetAllTeamMembers retrieves all team members with pagination
func (c *TeamMemberController) GetAllTeamMembers(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	teamMembers, err := c.teamMemberService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(teamMembers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.TeamMember{},
			"totalCount": 0,
			"notes":      "Không có thành viên đội nào",
			"message":    "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": teamMembers})
}

// UpdateTeamMember updates a team member
func (c *TeamMemberController) UpdateTeamMember(ctx *gin.Context) {
	ctx.Param("id")
	var teamMember models.TeamMember
	if err := ctx.ShouldBindJSON(&teamMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTeamMember, err := c.teamMemberService.Update(ctx, &teamMember)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedTeamMember})
}

// DeleteTeamMember deletes a team member by ID
func (c *TeamMemberController) DeleteTeamMember(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.teamMemberService.Delete(ctx, id); err != nil {
		if err.Error() == "team member not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "team member deleted"})
}	