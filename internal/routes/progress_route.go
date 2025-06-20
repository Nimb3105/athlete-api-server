package routes

// import (
// 	"be/internal/controllers"
// 	"github.com/gin-gonic/gin"
// )

// // SetupProgressRoutes attaches routes for Progress to the router group
// func SetupProgressRoutes(router *gin.Engine, progressController *controllers.ProgressController) {
// 	progresses := router.Group("/progresses")
// 	{
// 		progresses.POST("", progressController.CreateProgress)
// 		progresses.GET(":id", progressController.GetProgressByID)
// 		progresses.GET("user/:userID", progressController.GetProgressesByUserID)
// 		progresses.GET("", progressController.GetAllProgresses)
// 		progresses.PUT(":id", progressController.UpdateProgress)
// 		progresses.DELETE(":id", progressController.DeleteProgress)
// 	}
// }