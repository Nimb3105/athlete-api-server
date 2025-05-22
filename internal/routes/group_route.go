package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupGroupRoutes attaches routes for Group to the router group
func SetupGroupRoutes(router *gin.Engine, groupController *controllers.GroupController) {
	groups := router.Group("/groups")
	{
		groups.POST("", groupController.CreateGroup)
		groups.GET(":id", groupController.GetGroupByID)
		groups.GET("created-by/:createdBy", groupController.GetGroupByCreatedBy)
		groups.GET("", groupController.GetAllGroups)
		groups.PUT(":id", groupController.UpdateGroup)
		groups.DELETE(":id", groupController.DeleteGroup)
	}
}