package controllers

import (
	"be/internal/models"
	"be/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CoachCertificationController xử lý các yêu cầu HTTP cho CoachCertification
type CoachCertificationController struct {
	certificationService *services.CoachCertificationService
}

// NewCoachCertificationController tạo một CoachCertificationController mới
func NewCoachCertificationController(certificationService *services.CoachCertificationService) *CoachCertificationController {
	return &CoachCertificationController{certificationService}
}

func (c *CoachCertificationController) CreateCoachCertification(ctx *gin.Context) {
	var bodyBytes []byte
	if rawData, err := ctx.GetRawData(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc dữ liệu"})
		return
	} else {
		bodyBytes = rawData
	}

	var tempMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không đúng định dạng"})
		return
	}

	validFields := map[string]bool{
		"id": true, "userId": true, "name": true, "dateIssued": true,
		"createdAt": true, "updatedAt": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var coachCertification models.CoachCertification
	if err := json.Unmarshal(bodyBytes, &coachCertification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
		return
	}

	createdCoachCertification, err := c.certificationService.Create(ctx, &coachCertification)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdCoachCertification})
}

// GetCoachCertificationByID lấy coach certification theo ID
func (c *CoachCertificationController) GetCoachCertificationByID(ctx *gin.Context) {
	id := ctx.Param("id")
	certification, err := c.certificationService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "coach certification not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": certification})
}

// GetCoachCertificationByUserID lấy danh sách coach certification theo UserID
func (c *CoachCertificationController) GetCoachCertificationByUserID(ctx *gin.Context) {
	userID := ctx.Param("userID")
	certifications, err := c.certificationService.GetByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(certifications) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.CoachCertification{},
			"note":    "Không có chứng chỉ huấn luyện viên nào",
			"message": "Chưa có dữ liệu nào",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": certifications})
}

// GetAllCoachCertifications lấy danh sách tất cả coach certification với phân trang
func (c *CoachCertificationController) GetAllCoachCertifications(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	certifications, err := c.certificationService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(certifications) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.CoachCertification{},
			"message": "Không có dữ liệu nào",
			"note":    "Chưa có chứng chỉ huấn luyện viên nào được ghi nhận",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": certifications})
}

// UpdateCoachCertification cập nhật thông tin coach certification
func (c *CoachCertificationController) UpdateCoachCertification(ctx *gin.Context) {
	var certification models.CoachCertification
	if err := ctx.ShouldBindJSON(&certification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCertification, err := c.certificationService.Update(ctx, &certification)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedCertification})
}

// DeleteCoachCertification xóa coach certification theo ID
func (c *CoachCertificationController) DeleteCoachCertification(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.certificationService.Delete(ctx, id); err != nil {
		if err.Error() == "coach certification not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "coach certification deleted"})
}
