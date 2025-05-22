package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupMessageRoutes attaches routes for Message to the router group
func SetupMessageRoutes(router *gin.Engine, messageController *controllers.MessageController) {
	messages := router.Group("/messages")
	{
		messages.POST("", messageController.CreateMessage)
		messages.GET(":id", messageController.GetMessageByID)
		messages.GET("group/:groupID", messageController.GetMessagesByGroupID)
		messages.GET("", messageController.GetAllMessages)
		messages.PUT(":id", messageController.UpdateMessage)
		messages.DELETE(":id", messageController.DeleteMessage)
	}
}