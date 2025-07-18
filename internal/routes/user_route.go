package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes gắn các route cho User vào router group
func SetupUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	router.POST("/login", userController.Login)
	users := router.Group("/users")
	{
		users.POST("", userController.CreateUser)
		users.GET(":id", userController.GetUserByID)
		users.GET("email/:email", userController.GetUserByEmail)
		users.GET("", userController.GetAllUsers)
		users.PUT(":id", userController.UpdateUser)
		users.DELETE(":id", userController.DeleteUser)
		users.GET("coach/sport/:sportId", userController.GetAllUserCoachBySportId)
		users.GET("role/:role", userController.GetUsersByRoleWithPagination)
		users.GET("unassigned-athletes/:sportId",userController.GetUnassignedAthletes)
	}
}
