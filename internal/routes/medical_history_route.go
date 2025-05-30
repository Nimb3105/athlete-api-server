package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupMedicalHistoryRoutes attaches routes for MedicalHistory to the router group
func SetupMedicalHistoryRoutes(router *gin.Engine, medicalHistoryController *controllers.MedicalHistoryController) {
	medicalHistories := router.Group("/medical-historys")
	{
		medicalHistories.POST("", medicalHistoryController.CreateMedicalHistory)
		medicalHistories.GET(":id", medicalHistoryController.GetMedicalHistoryByID)
		medicalHistories.GET("health/:healthID", medicalHistoryController.GetMedicalHistoryByHealthID)
		medicalHistories.GET("", medicalHistoryController.GetAllMedicalHistories)
		medicalHistories.PUT(":id", medicalHistoryController.UpdateMedicalHistory)
		medicalHistories.DELETE(":id", medicalHistoryController.DeleteMedicalHistory)
	}
}