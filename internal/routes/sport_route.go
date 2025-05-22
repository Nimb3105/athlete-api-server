package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupSportRoutes gắn các route cho Sport vào router group
func SetupSportRoutes(router *gin.Engine, sportController *controllers.SportController) {
	sports := router.Group("/sports")
	{
		sports.POST("", sportController.CreateSport)
		sports.GET(":id", sportController.GetSportByID)
		sports.GET("", sportController.GetAllSports)
		sports.PUT(":id", sportController.UpdateSport)
		sports.DELETE(":id", sportController.DeleteSport)
	}
}
