package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupReminderRoutes gắn các route cho Reminder vào router group
func SetupReminderRoutes(router *gin.Engine, reminderController *controllers.ReminderController) {
	reminders := router.Group("/reminders")
	{
		reminders.POST("", reminderController.CreateReminder)
		reminders.GET(":id", reminderController.GetReminderByID)
		reminders.GET("user/:userID", reminderController.GetRemindersByUserID)
		reminders.PUT(":id", reminderController.UpdateReminder)
		reminders.DELETE(":id", reminderController.DeleteReminder)
		reminders.GET("", reminderController.GetAllReminders)
	}
}
