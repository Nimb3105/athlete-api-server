package routes

// import (
// 	"be/internal/controllers"
// 	"github.com/gin-gonic/gin"
// )

// // SetupPerformanceRoutes attaches routes for Performance to the router group
// func SetupPerformanceRoutes(router *gin.Engine, performanceController *controllers.PerformanceController) {
// 	performances := router.Group("/performances")
// 	{
// 		performances.POST("", performanceController.CreatePerformance)
// 		performances.GET(":id", performanceController.GetPerformanceByID)
// 		performances.GET("user/:userID", performanceController.GetPerformancesByUserID)
// 		performances.GET("", performanceController.GetAllPerformances)
// 		performances.PUT(":id", performanceController.UpdatePerformance)
// 		performances.DELETE(":id", performanceController.DeletePerformance)
// 	}
// }