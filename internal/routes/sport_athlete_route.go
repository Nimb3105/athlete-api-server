package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// uss
// SetupSportAthleteRoutes gắn các route cho SportAthlete vào router group
func SetupSportAthleteRoutes(router *gin.Engine, sportAthleteController *controllers.SportAthleteController) {
	sportAthletes := router.Group("/sport-athletes")
	{
		sportAthletes.POST("", sportAthleteController.CreateSportAthlete)
		sportAthletes.GET(":id", sportAthleteController.GetSportAthleteByID)
		sportAthletes.GET("user/:userID", sportAthleteController.GetSportAthleteByUserID)
		sportAthletes.GET("sport/:sportID", sportAthleteController.GetSportAthleteBySportID)
		sportAthletes.GET("", sportAthleteController.GetAllSportAthletes)
		sportAthletes.PUT(":id", sportAthleteController.UpdateSportAthlete)
		sportAthletes.DELETE(":id", sportAthleteController.DeleteSportAthlete)
		sportAthletes.GET("/user/:userID/all", sportAthleteController.GetAllByUserID)
	}
}
