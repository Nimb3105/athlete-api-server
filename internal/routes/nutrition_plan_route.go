package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupNutritionPlanRoutes attaches routes for NutritionPlan to the router group
func SetupNutritionPlanRoutes(router *gin.Engine, nutritionPlanController *controllers.NutritionPlanController) {
	nutritionPlans := router.Group("/nutrition-plans")
	{
		nutritionPlans.POST("", nutritionPlanController.CreateNutritionPlan)
		nutritionPlans.GET(":id", nutritionPlanController.GetNutritionPlanByID)
		nutritionPlans.GET("user/:userID", nutritionPlanController.GetNutritionPlansByUserID)
		nutritionPlans.GET("", nutritionPlanController.GetAllNutritionPlans)
		nutritionPlans.PUT(":id", nutritionPlanController.UpdateNutritionPlan)
		nutritionPlans.DELETE(":id", nutritionPlanController.DeleteNutritionPlan)
	}
}
