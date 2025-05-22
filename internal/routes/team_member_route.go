package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupTeamMemberRoutes attaches routes for TeamMember to the router group
func SetupTeamMemberRoutes(router *gin.Engine, teamMemberController *controllers.TeamMemberController) {
	teamMembers := router.Group("/team-members")
	{
		teamMembers.POST("", teamMemberController.CreateTeamMember)
		teamMembers.GET(":id", teamMemberController.GetTeamMemberByID)
		teamMembers.GET("team/:teamID", teamMemberController.GetTeamMembersByTeamID)
		teamMembers.GET("user/:userID", teamMemberController.GetTeamMembersByUserID)
		teamMembers.GET("", teamMemberController.GetAllTeamMembers)
		teamMembers.PUT(":id", teamMemberController.UpdateTeamMember)
		teamMembers.DELETE(":id", teamMemberController.DeleteTeamMember)
	}
}