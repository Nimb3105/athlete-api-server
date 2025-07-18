package routes

import (
	"be/internal/controllers"

	"github.com/gin-gonic/gin" // Thay bằng đường dẫn thực tế đến package controllers
)

// SetupImageRoutes registers routes for image-related operations
func SetupImageRoutes(router *gin.Engine, imageController *controllers.ImageController) {
	// Nhóm các route liên quan đến hình ảnh dưới /images
	imageGroup := router.Group("/images")
	{
		// Route để tải lên hình ảnh
		imageGroup.POST("/upload", imageController.UploadImage)
		// Route để xóa hình ảnh
		imageGroup.DELETE("/:name", imageController.DeleteImage)
	}

	// Phục vụ file tĩnh từ thư mục ./be/images
	router.Static("/images", "./images")
}
