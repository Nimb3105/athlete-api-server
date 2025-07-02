package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTrainingScheduleRoutes(r *gin.Engine, controller *controllers.TrainingScheduleController) {
	schedule := r.Group("/training-schedules")
	{
		schedule.POST("", controller.CreateTrainingSchedule)
		schedule.GET("/:id", controller.GetByID)
		schedule.GET("", controller.GetAll)
		schedule.PUT("", controller.Update)
		schedule.DELETE("/:id", controller.Delete)
		schedule.GET("/daily/:dailyScheduleId/:date", controller.GetAllByDailyScheduleId)
	}
}
