package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAthleteRoutes gắn các route cho Athlete vào router group
func SetupAthleteRoutes(router *gin.Engine, athleteController *controllers.AthleteController) {
	athletes := router.Group("/athletes")
	{
		athletes.POST("", athleteController.CreateAthlete)
		athletes.GET(":id", athleteController.GetAthleteByID)
		athletes.GET("user/:userID", athleteController.GetAthleteByUserID)
		athletes.GET("", athleteController.GetAllAthletes)
		athletes.PUT(":id", athleteController.UpdateAthlete)
		athletes.DELETE(":id", athleteController.DeleteAthlete)
	}
}
