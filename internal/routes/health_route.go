package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupHealthRoutes attaches routes for Health to the router group
func SetupHealthRoutes(router *gin.Engine, healthController *controllers.HealthController) {
	healthRecords := router.Group("/health-records")
	{
		healthRecords.POST("", healthController.CreateHealth)
		healthRecords.GET(":id", healthController.GetHealthByID)
		healthRecords.GET("user/:userID", healthController.GetHealthByUserID)
		healthRecords.GET("", healthController.GetAllHealthRecords)
		healthRecords.PUT(":id", healthController.UpdateHealth)
		healthRecords.DELETE(":id", healthController.DeleteHealth)
	}
}
