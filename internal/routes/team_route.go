package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupTeamRoutes attaches routes for Team to the router group
func SetupTeamRoutes(router *gin.Engine, teamController *controllers.TeamController) {
	teams := router.Group("/teams")
	{
		teams.POST("", teamController.CreateTeam)
		teams.GET(":id", teamController.GetTeamByID)
		teams.GET("sport/:sportID", teamController.GetTeamsBySportID)
		teams.GET("", teamController.GetAllTeams)
		teams.PUT(":id", teamController.UpdateTeam)
		teams.DELETE(":id", teamController.DeleteTeam)
	}
}