package routes

import (
	"be/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupNutritionMealRoutes attaches routes for NutritionMeal to the router group
func SetupFoodRoutes(router *gin.Engine, foodController *controllers.FoodController) {
	foods := router.Group("/foods")
	{
		foods.POST("", foodController.CreateNutritionMeal)
		foods.GET(":id", foodController.GetNutritionMealByID)
		foods.GET("", foodController.GetAllNutritionMeals)
		foods.PUT(":id", foodController.UpdateNutritionMeal)
		foods.DELETE(":id", foodController.DeleteNutritionMeal)
		foods.GET("foodType/:foodType",foodController.GetAllByFoodType)
	}
}