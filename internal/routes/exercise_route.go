package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupExerciseRoutes(router *gin.Engine, controller *controllers.ExerciseController) {
	exercise := router.Group("/exercises")
	{
		exercise.POST("", controller.CreateExercise)
		exercise.GET("/:id", controller.GetByID)
		exercise.GET("", controller.GetAll)
		exercise.PUT("/:id", controller.Update)
		exercise.DELETE("/:id", controller.Delete)
		//exercise.GET("/bodyPart/:bodyPart", controller.GetAllByBodyPart)
		//exercise.GET("/sport/:sportName", controller.GetAllBySportName)
		exercise.GET("/sport/:sportId", controller.GetAllBySportId)
	}
}
