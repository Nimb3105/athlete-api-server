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

// GroupMemberController handles HTTP requests for GroupMember
type GroupMemberController struct {
	groupMemberService *services.GroupMemberService
}

// NewGroupMemberController creates a new GroupMemberController
func NewGroupMemberController(groupMemberService *services.GroupMemberService) *GroupMemberController {
	return &GroupMemberController{groupMemberService}
}

// CreateGroupMember creates a new group member
func (c *GroupMemberController) CreateGroupMember(ctx *gin.Context) {
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
		"id": true, "groupId": true, "userId": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var groupMember models.GroupMember
	if err := json.Unmarshal(bodyBytes, &groupMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdGroupMember, err := c.groupMemberService.Create(ctx, &groupMember)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdGroupMember})
}

// GetGroupMemberByID retrieves a group member by ID
func (c *GroupMemberController) GetGroupMemberByID(ctx *gin.Context) {
	id := ctx.Param("id")
	groupMember, err := c.groupMemberService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "group member not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": groupMember})
}

// GetGroupMemberByUserID retrieves a group member by user ID
func (c *GroupMemberController) GetGroupMemberByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	groupMembers, err := c.groupMemberService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(groupMembers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data": []models.TeamMember{},
			"message": "không có dữ liệu nào",
			"note": "Bạn có thể tạo group member mới bằng cách gửi yêu cầu POST tới /team_members",
		})
		return
	}


	ctx.JSON(http.StatusOK, gin.H{"data": groupMembers})
}

// GetAllGroupMembers retrieves all group members with pagination
func (c *GroupMemberController) GetAllGroupMembers(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	groupMembers, err := c.groupMemberService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(groupMembers) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.GroupMember{},
			"totalCount": 0,
			"notes":      "Không có thành viên nhóm nào",
			"message":    "Không có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": groupMembers})
}

// UpdateGroupMember updates a group member
func (c *GroupMemberController) UpdateGroupMember(ctx *gin.Context) {
	ctx.Param("id")
	var groupMember models.GroupMember
	if err := ctx.ShouldBindJSON(&groupMember); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedGroupMember, err := c.groupMemberService.Update(ctx, &groupMember)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedGroupMember})
}

// DeleteGroupMember deletes a group member by ID
func (c *GroupMemberController) DeleteGroupMember(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.groupMemberService.Delete(ctx, id); err != nil {
		if err.Error() == "group member not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "group member deleted"})
}
