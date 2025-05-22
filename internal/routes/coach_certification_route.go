package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupCoachCertificationRoutes attaches routes for CoachCertification to the router group
func SetupCoachCertificationRoutes(router *gin.Engine, certificationController *controllers.CoachCertificationController) {
	certifications := router.Group("/coach-certifications")
	{
		certifications.POST("", certificationController.CreateCoachCertification)
		certifications.GET(":id", certificationController.GetCoachCertificationByID)
		certifications.GET("user/:userID", certificationController.GetCoachCertificationByUserID)
		certifications.GET("", certificationController.GetAllCoachCertifications)
		certifications.PUT(":id", certificationController.UpdateCoachCertification)
		certifications.DELETE(":id", certificationController.DeleteCoachCertification)
	}
}