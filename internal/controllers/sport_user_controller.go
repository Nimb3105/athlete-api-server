package controllers

// import (
// 	"be/internal/models"
// 	"be/internal/services"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // SportUserController xử lý các yêu cầu HTTP cho SportUser
// type SportUserController struct {
// 	SportUserService *services.SportUserService
// }

// // NewSportUserController tạo một SportUserController mới
// func NewSportUserController(SportUserService *services.SportUserService) *SportUserController {
// 	return &SportUserController{SportUserService}
// }

// // CreateSportUser tạo một sport athlete mới
// func (c *SportUserController) CreateSportUser(ctx *gin.Context) {
// 	var bodyBytes []byte
// 	if rawData, err := ctx.GetRawData(); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc dữ liệu"})
// 		return
// 	} else {
// 		bodyBytes = rawData
// 	}

// 	var tempMap map[string]interface{}
// 	if err := json.Unmarshal(bodyBytes, &tempMap); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không đúng định dạng"})
// 		return
// 	}

// 	validFields := map[string]bool{
// 		"id": true, "sportId": true, "userId": true, "position": true,
// 		"createdAt": true, "updatedAt": true,
// 	}
// 	for key := range tempMap {
// 		if !validFields[key] {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
// 			return
// 		}
// 	}

// 	var SportUser models.SportUser
// 	if err := json.Unmarshal(bodyBytes, &SportUser); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào model"})
// 		return
// 	}

// 	createdSportUser, err := c.SportUserService.Create(ctx, &SportUser)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{"data": createdSportUser})
// }

// // GetSportUserByID lấy sport athlete theo ID
// func (c *SportUserController) GetSportUserByID(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	SportUser, err := c.SportUserService.GetByID(ctx, id)
// 	if err != nil {
// 		if err.Error() == "SportUser not found" {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": SportUser})
// }

// // GetSportUserByUserID lấy sport athlete theo UserID
// func (c *SportUserController) GetSportUserByUserID(ctx *gin.Context) {
// 	userID := ctx.Param("userID")
// 	SportUser, err := c.SportUserService.GetByUserID(ctx, userID)
// 	if err != nil {
// 		if err.Error() == "SportUser not found" {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": SportUser})
// }

// func (c *SportUserController) GetSportUserBySportID(ctx *gin.Context) {
// 	sportID := ctx.Param("sportID")
// 	SportUser, err := c.SportUserService.GetByUserID(ctx, sportID)
// 	if err != nil {
// 		if err.Error() == "SportUser not found" {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": SportUser})
// }

// // GetAllSportUsers lấy danh sách tất cả sport athlete với phân trang
// func (c *SportUserController) GetAllSportUsers(ctx *gin.Context) {
// 	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
// 	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

// 	SportUsers, err := c.SportUserService.GetAll(ctx, page, limit)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if len(SportUsers) == 0 {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"data":    []models.SportUser{},
// 			"message": "Không có dữ liệu nào",
// 			"note":    "Chưa có vận động viên thể thao nào được ghi nhận",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": SportUsers})
// }

// // UpdateSportUser cập nhật thông tin sport athlete
// func (c *SportUserController) UpdateSportUser(ctx *gin.Context) {
// 	ctx.Param("id")
// 	var SportUser models.SportUser
// 	if err := ctx.ShouldBindJSON(&SportUser); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updatedSportUser, err := c.SportUserService.Update(ctx, &SportUser)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": updatedSportUser})
// }

// // DeleteSportUser xóa sport athlete theo ID
// func (c *SportUserController) DeleteSportUser(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	if err := c.SportUserService.Delete(ctx, id); err != nil {
// 		if err.Error() == "SportUser not found" {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": "SportUser deleted"})
// }

// func (c *SportUserController) GetAllByUserID(ctx *gin.Context) {
// 	userID := ctx.Param("userID")
// 	SportUsers, err := c.SportUserService.GetAllByUserID(ctx, userID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if len(SportUsers) == 0 {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"data":    []models.SportUser{},
// 			"note":    "Người dùng chưa có môn thể thao nào",
// 			"message": "Không có dữ liệu nào",
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": SportUsers})
// }
