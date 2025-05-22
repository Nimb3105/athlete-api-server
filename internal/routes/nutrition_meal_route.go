package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupNutritionMealRoutes attaches routes for NutritionMeal to the router group
func SetupNutritionMealRoutes(router *gin.Engine, nutritionMealController *controllers.NutritionMealController) {
	nutritionMeals := router.Group("/nutrition-meals")
	{
		nutritionMeals.POST("", nutritionMealController.CreateNutritionMeal)
		nutritionMeals.GET(":id", nutritionMealController.GetNutritionMealByID)
		nutritionMeals.GET("plan/:nutritionPlanID", nutritionMealController.GetNutritionMealsByNutritionPlanID)
		nutritionMeals.GET("", nutritionMealController.GetAllNutritionMeals)
		nutritionMeals.PUT(":id", nutritionMealController.UpdateNutritionMeal)
		nutritionMeals.DELETE(":id", nutritionMealController.DeleteNutritionMeal)
	}
}