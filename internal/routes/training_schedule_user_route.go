package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTrainingScheduleUserRoutes(r *gin.Engine, controller *controllers.TrainingScheduleUserController) {
	scheduleAthlete := r.Group("/training-schedule-users")
	{
		scheduleAthlete.POST("", controller.CreateTrainingScheduleUser)

		scheduleAthlete.GET("/user/:userID/all", controller.GetAllTrainingScheduleUserByUserID)
		scheduleAthlete.GET("/:id", controller.GetByID)
		scheduleAthlete.GET("/schedule/:scheduleID", controller.GetByScheduleID)
		scheduleAthlete.GET("/user/:userID", controller.GetByUserID)
		scheduleAthlete.GET("", controller.GetAll)
		scheduleAthlete.PUT("", controller.Update)
		scheduleAthlete.DELETE("/:id", controller.Delete)
	}
}
