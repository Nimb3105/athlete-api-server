package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupDailyScheduleRoutes(router *gin.Engine, controller *controllers.DailyScheduleController) {
	dailySchedule := router.Group("/dailySchedules")
	{
		dailySchedule.POST("", controller.Create)
		dailySchedule.GET("/:id", controller.GetByID)
		dailySchedule.GET("/user/:userID/:day", controller.GetByUserID)
		dailySchedule.GET("", controller.GetAll)
		dailySchedule.PUT(":id", controller.Update)
		dailySchedule.DELETE("/:id", controller.Delete)
		dailySchedule.GET("/creator/:creatorId", controller.GetByCreatorId)
		dailySchedule.GET("/user/:userID/all", controller.GetAllDailySchedulesByUserId)
	}
}
