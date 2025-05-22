package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupInjuryRoutes attaches routes for Injury to the router group
func SetupInjuryRoutes(router *gin.Engine, injuryController *controllers.InjuryController) {
	injuries := router.Group("/injuries")
	{
		injuries.POST("", injuryController.CreateInjury)
		injuries.GET(":id", injuryController.GetInjuryByID)
		injuries.GET("user/:userID", injuryController.GetInjuryByUserID)
		injuries.GET("", injuryController.GetAllInjuries)
		injuries.PUT(":id", injuryController.UpdateInjury)
		injuries.DELETE(":id", injuryController.DeleteInjury)
	}
}