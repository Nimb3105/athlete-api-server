package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupTournamentRoutes attaches routes for Tournament to the router group
func SetupTournamentRoutes(router *gin.Engine, tournamentController *controllers.TournamentController) {
	tournaments := router.Group("/tournaments")
	{
		tournaments.POST("", tournamentController.CreateTournament)
		tournaments.GET(":id", tournamentController.GetTournamentByID)
		tournaments.GET("", tournamentController.GetAllTournaments)
		tournaments.PUT(":id", tournamentController.UpdateTournament)
		tournaments.DELETE(":id", tournamentController.DeleteTournament)
	}
}