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

// GroupController xử lý các yêu cầu HTTP cho Group
type GroupController struct {
	groupService *services.GroupService
}

// NewGroupController tạo một GroupController mới
func NewGroupController(groupService *services.GroupService) *GroupController {
	return &GroupController{groupService}
}

// CreateGroup tạo một group mới
func (c *GroupController) CreateGroup(ctx *gin.Context) {
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
		"id": true, "name": true, "description": true, "createdBy": true,
		"createdDate": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var group models.Group
	if err := json.Unmarshal(bodyBytes, &group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdGroup, err := c.groupService.Create(ctx, &group)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdGroup})
}

// GetGroupByID lấy group theo ID
func (c *GroupController) GetGroupByID(ctx *gin.Context) {
	id := ctx.Param("id")
	group, err := c.groupService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "group not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": group})
}

// GetGroupByCreatedBy lấy danh sách group theo CreatedBy
func (c *GroupController) GetGroupByCreatedBy(ctx *gin.Context) {
	createdBy := ctx.Param("createdBy")
	groups, err := c.groupService.GetByCreatedBy(ctx, createdBy)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(groups) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Group{},
			"notes":   "Không có group nào được tìm thấy cho CreatedBy: " + createdBy,
			"message": "Không có dữ liệu nào"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": groups})
}

// GetAllGroups lấy danh sách tất cả group với phân trang
func (c *GroupController) GetAllGroups(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	groups, err := c.groupService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(groups) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Group{},
			"message": "Không có dữ liệu nào",
			"note":    "Chưa có group nào được tạo",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": groups})
}

// UpdateGroup cập nhật thông tin group
func (c *GroupController) UpdateGroup(ctx *gin.Context) {
	var group models.Group
	if err := ctx.ShouldBindJSON(&group); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedGroup, err := c.groupService.Update(ctx, &group)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedGroup})
}

// DeleteGroup xóa group theo ID
func (c *GroupController) DeleteGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.groupService.Delete(ctx, id); err != nil {
		if err.Error() == "group not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "group deleted"})
}
