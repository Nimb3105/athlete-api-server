package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes gắn các route cho Notification vào router group
func SetupNotificationRoutes(router *gin.Engine, notificationController *controllers.NotificationController) {
	notifications := router.Group("/notifications")
	{
		notifications.POST("", notificationController.CreateNotification)
		notifications.GET(":id", notificationController.GetNotificationByID)
		notifications.GET("user/:userID", notificationController.GetNotificationsByUserID)
		notifications.PUT(":id", notificationController.UpdateNotification)
		notifications.DELETE(":id", notificationController.DeleteNotification)
	}
}