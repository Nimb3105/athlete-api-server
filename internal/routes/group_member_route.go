package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupGroupMemberRoutes attaches routes for GroupMember to the router group
func SetupGroupMemberRoutes(router *gin.Engine, groupMemberController *controllers.GroupMemberController) {
	groupMembers := router.Group("/group-members")
	{
		groupMembers.POST("", groupMemberController.CreateGroupMember)
		groupMembers.GET(":id", groupMemberController.GetGroupMemberByID)
		groupMembers.GET("user/:userID", groupMemberController.GetGroupMemberByUserID)
		groupMembers.GET("", groupMemberController.GetAllGroupMembers)
		groupMembers.PUT(":id", groupMemberController.UpdateGroupMember)
		groupMembers.DELETE(":id", groupMemberController.DeleteGroupMember)
	}
}
