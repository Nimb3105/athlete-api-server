package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAthleteMatchRoutes attaches routes for AthleteMatch to the router group
func SetupUserMatchRoutes(router *gin.Engine, userMatchController *controllers.UserMatchController) {
	userMatches := router.Group("/user-matchs")
	{
		userMatches.POST("", userMatchController.CreateUserMatch)
		userMatches.GET(":id", userMatchController.GetUserMatchByID)
		userMatches.GET("user/:userID", userMatchController.GetUserMatchByUserID)
		userMatches.GET("", userMatchController.GetAllUserMatches)
		userMatches.PUT(":id", userMatchController.UpdateUserMatch)
		userMatches.DELETE(":id", userMatchController.DeleteUserMatch)
	}
}
