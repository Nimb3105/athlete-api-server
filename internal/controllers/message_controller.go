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

// MessageController handles HTTP requests for Message
type MessageController struct {
	messageService *services.MessageService
}

// NewMessageController creates a new MessageController
func NewMessageController(messageService *services.MessageService) *MessageController {
	return &MessageController{messageService}
}

// CreateMessage creates a new message
func (c *MessageController) CreateMessage(ctx *gin.Context) {
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
		"id": true, "groupId": true, "senderId": true, "sentDate": true,
		"content": true, "createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var message models.Message
	if err := json.Unmarshal(bodyBytes, &message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdMessage, err := c.messageService.Create(ctx, &message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdMessage})
}

// GetMessageByID retrieves a message by ID
func (c *MessageController) GetMessageByID(ctx *gin.Context) {
	id := ctx.Param("id")
	message, err := c.messageService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "message not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": message})
}

// GetMessagesByGroupID retrieves messages by group ID
func (c *MessageController) GetMessagesByGroupID(ctx *gin.Context) {
	groupID := ctx.Param("groupID")
	messages, err := c.messageService.GetByGroupID(ctx, groupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(messages) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Message{},
			"message": "không có dữ liệu nào",
			"notes":   "Không có tin nhắn nào được tìm thấy cho GroupID: " + groupID,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": messages})
}

// GetAllMessages retrieves all messages with pagination
func (c *MessageController) GetAllMessages(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	messages, err := c.messageService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(messages) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":       []models.Message{},
			"totalCount": 0,
			"notes":      "Không có tin nhắn nào",
			"message":    "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": messages})
}

// UpdateMessage updates a message
func (c *MessageController) UpdateMessage(ctx *gin.Context) {
	ctx.Param("id")
	var message models.Message
	if err := ctx.ShouldBindJSON(&message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMessage, err := c.messageService.Update(ctx, &message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedMessage})
}

// DeleteMessage deletes a message by ID
func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.messageService.Delete(ctx, id); err != nil {
		if err.Error() == "message not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "message deleted"})
}
