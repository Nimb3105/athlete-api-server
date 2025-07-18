package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	StoragePath string
}

func NewImageController() *ImageController {
	return &ImageController{
		StoragePath: "./images",
	}
}

// UploadImage handles image and video uploads to the ./be/images directory
func (ic *ImageController) UploadImage(ctx *gin.Context) {
	// Limit file size to 50MB
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Yêu cầu tệp hình ảnh hoặc video"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".gif": true,
		".jpg": true,
		".png": true,
		".mp4": true,
		".mov": true,
		".avi": true,
	}
	if !allowedExtensions[ext] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng tệp không hợp lệ (chỉ chấp nhận .gif, .jpg, .png, .mp4, .mov, .avi)"})
		return
	}

	// Sanitize filename
	safeFilename := regexp.MustCompile(`[^a-zA-Z0-9._-]`).ReplaceAllString(file.Filename, "_")
	if safeFilename == "" || safeFilename == "." || safeFilename == ".." {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tên tệp không hợp lệ"})
		return
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(ic.StoragePath, 0755); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo thư mục lưu trữ"})
		return
	}

	filePath := filepath.Join(ic.StoragePath, safeFilename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lưu tệp"})
		return
	}

	// Tạo URL đầy đủ
	fileURL := fmt.Sprintf("/images/%s", safeFilename)

	// Trả về JSON chứa URL
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tải tệp thành công",
		"url":     fileURL, // Trả về URL cho frontend
	})
}

// DeleteImage handles image and video deletion from the ./be/images directory
func (ic *ImageController) DeleteImage(ctx *gin.Context) {
	filename := ctx.Param("name")

	// Prevent path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tên tệp không hợp lệ"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExtensions := map[string]bool{
		".gif": true,
		".jpg": true,
		".png": true,
		".mp4": true,
		".mov": true,
		".avi": true,
	}
	if !allowedExtensions[ext] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Chỉ được xóa tệp hình ảnh hoặc video (.gif, .jpg, .png, .mp4, .mov, .avi)"})
		return
	}

	// File path
	filePath := filepath.Join(ic.StoragePath, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tệp không tồn tại"})
		return
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể xóa tệp"})
		return
	}

	// Return success message
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Xóa tệp thành công",
	})
}
