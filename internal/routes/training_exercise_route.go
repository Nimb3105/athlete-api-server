package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupTrainingExerciseRoutes(r *gin.Engine, controller *controllers.TrainingExerciseController) {
	trainingExercise := r.Group("/training-exercises")
	{
		trainingExercise.POST("", controller.CreateTrainingExercise)
		trainingExercise.GET("/:id", controller.GetByID)
		trainingExercise.GET("/schedule/:scheduleId", controller.GetByScheduleID)
		trainingExercise.GET("", controller.GetAll)
		trainingExercise.PUT("", controller.Update)
		trainingExercise.DELETE("/:id", controller.Delete)
		 trainingExercise.GET("/schedule/:scheduleId/all", controller.GetAllTrainingExerciseByScheduleId)
	}
}
