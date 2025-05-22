package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupAchievementRoutes attaches routes for Achievement to the router group
func SetupAchievementRoutes(router *gin.Engine, achievementController *controllers.AchievementController) {
	achievements := router.Group("/achievements")
	{
		achievements.POST("", achievementController.CreateAchievement)
		achievements.GET(":id", achievementController.GetAchievementByID)
		achievements.GET("user/:userID", achievementController.GetAchievementByUserID)
		achievements.GET("", achievementController.GetAllAchievements)
		achievements.PUT(":id", achievementController.UpdateAchievement)
		achievements.DELETE(":id", achievementController.DeleteAchievement)
	}
}