package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupFeedbackRoutes attaches routes for Feedback to the router group
func SetupFeedbackRoutes(router *gin.Engine, feedbackController *controllers.FeedbackController) {
	feedbacks := router.Group("/feedbacks")
	{
		feedbacks.POST("", feedbackController.CreateFeedback)
		feedbacks.GET(":id", feedbackController.GetFeedbackByID)
		feedbacks.GET("user/:userID", feedbackController.GetFeedbackByUserID)
		feedbacks.GET("", feedbackController.GetAllFeedbacks)
		feedbacks.PUT(":id", feedbackController.UpdateFeedback)
		feedbacks.DELETE(":id", feedbackController.DeleteFeedback)
	}
}