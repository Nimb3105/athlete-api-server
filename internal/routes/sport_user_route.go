package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// uss
// SetupSportUserRoutes gắn các route cho SportUser vào router group
func SetupSportUserRoutes(router *gin.Engine, SportUserController *controllers.SportUserController) {
	SportUsers := router.Group("/sport-users")
	{
		SportUsers.POST("", SportUserController.CreateSportUser)
		SportUsers.GET(":id", SportUserController.GetSportUserByID)
		SportUsers.GET("user/:userID", SportUserController.GetSportUserByUserID)
		SportUsers.GET("sport/:sportID", SportUserController.GetSportUserBySportID)
		SportUsers.GET("", SportUserController.GetAllSportUsers)
		SportUsers.PUT(":id", SportUserController.UpdateSportUser)
		SportUsers.DELETE(":id", SportUserController.DeleteSportUser)
		SportUsers.GET("/user/:userID/all", SportUserController.GetAllByUserID)
	}
}
