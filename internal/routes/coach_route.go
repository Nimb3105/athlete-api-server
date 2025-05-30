package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupCoachRoutes gắn các route cho Coach vào router group
func SetupCoachRoutes(router *gin.Engine, coachController *controllers.CoachController) {
	coaches := router.Group("/coachs")
	{
		coaches.POST("", coachController.CreateCoach)
		coaches.GET(":id", coachController.GetCoachByID)
		coaches.GET("user/:userID", coachController.GetCoachByUserID)
		coaches.GET("", coachController.GetAllCoaches)
		coaches.PUT(":id", coachController.UpdateCoach)
		coaches.DELETE(":id", coachController.DeleteCoach)
	}
}
