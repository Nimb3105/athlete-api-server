package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupPlanFoodRoutes attaches routes for PlanFood to the router group
func SetupPlanFoodRoutes(router *gin.Engine, planFoodController *controllers.PlanFoodController) {
	planFoods := router.Group("/plan-foods")
	{
		planFoods.POST("", planFoodController.CreatePlanFood)
		planFoods.GET(":id", planFoodController.GetPlanFoodByID)
		planFoods.GET("", planFoodController.GetAllPlanFoods)
		planFoods.PUT(":id", planFoodController.UpdatePlanFood)
		planFoods.DELETE(":id", planFoodController.DeletePlanFood)
		planFoods.GET("nutrition/:nutritionPlanId/all", planFoodController.GetAllByNutritionPlanID)
	}
}
