package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupAthleteMatchRoutes attaches routes for AthleteMatch to the router group
func SetupAthleteMatchRoutes(router *gin.Engine, athleteMatchController *controllers.AthleteMatchController) {
	athleteMatches := router.Group("/athlete-matches")
	{
		athleteMatches.POST("", athleteMatchController.CreateAthleteMatch)
		athleteMatches.GET(":id", athleteMatchController.GetAthleteMatchByID)
		athleteMatches.GET("user/:userID", athleteMatchController.GetAthleteMatchByUserID)
		athleteMatches.GET("", athleteMatchController.GetAllAthleteMatches)
		athleteMatches.PUT(":id", athleteMatchController.UpdateAthleteMatch)
		athleteMatches.DELETE(":id", athleteMatchController.DeleteAthleteMatch)
	}
}