package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupMatchScheduleRoutes attaches routes for MatchSchedule to the router group
func SetupMatchScheduleRoutes(router *gin.Engine, matchScheduleController *controllers.MatchScheduleController) {
	matchSchedules := router.Group("/match-schedules")
	{
		matchSchedules.POST("", matchScheduleController.CreateMatchSchedule)
		matchSchedules.GET(":id", matchScheduleController.GetMatchScheduleByID)
		matchSchedules.GET("tournament/:tournamentID", matchScheduleController.GetMatchScheduleByTournamentID)
		matchSchedules.GET("", matchScheduleController.GetAllMatchSchedules)
		matchSchedules.PUT(":id", matchScheduleController.UpdateMatchSchedule)
		matchSchedules.DELETE(":id", matchScheduleController.DeleteMatchSchedule)
	}
}