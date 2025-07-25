package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupCoachAthleteRoutes attaches routes for CoachAthlete to the router group
func SetupCoachAthleteRoutes(router *gin.Engine, coachAthleteController *controllers.CoachAthleteController) {
	coachAthletes := router.Group("/coach-athletes")
	{
		coachAthletes.POST("", coachAthleteController.CreateCoachAthlete)
		coachAthletes.GET(":id", coachAthleteController.GetCoachAthleteByID)
		coachAthletes.GET("/athlete/:athleteId", coachAthleteController.GetCoachAthleteByAthleteId)
		coachAthletes.GET("", coachAthleteController.GetAllCoachAthletes)
		coachAthletes.PUT(":id", coachAthleteController.UpdateCoachAthlete)
		coachAthletes.DELETE(":id", coachAthleteController.DeleteCoachAthlete)
		coachAthletes.GET("/user/:userId", coachAthleteController.GetAllByCoachId)
		coachAthletes.DELETE("/coach/:coachId", coachAthleteController.DeleteAllByCoachId)
	}
}
