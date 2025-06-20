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

// UserController xử lý các yêu cầu HTTP cho User
type UserController struct {
	userService *services.UserService
}

// NewUserController tạo một UserController mới
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService}
}

func (c *UserController) GetAllUserCoachBySportId(ctx *gin.Context) {
	sportId := ctx.Param("sportId")
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	users, totalCount, err := c.userService.GetAllUserCoachBySportId(ctx, page, limit, sportId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.Athlete{},
			"message": "Không có dữ liệu nào",
			"note":    "Chưa có vận động viên nào được ghi nhận",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users, "totalCount": totalCount})
}

// CreateUser tạo một user mới
func (c *UserController) CreateUser(ctx *gin.Context) {
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
		"email": true, "password": true, "fullName": true, "gender": true,
		"phoneNumber": true, "dateOfBirth": true, "role": true, "status": true,
		"createdAt": true, "updatedAt": true, "imageUrl": true, "sportId": true,
	}
	for key := range tempMap {
		if !validFields[key] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Trường không hợp lệ: %s", key)})
			return
		}
	}

	var user models.User
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không thể ánh xạ dữ liệu vào User model"})
		return
	}

	createdUser, err := c.userService.Create(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdUser})
}

// GetUserByID lấy user theo ID
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.userService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// GetUserByEmail lấy user theo email
func (c *UserController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	user, err := c.userService.GetByEmail(ctx, email)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

// GetAllUsers lấy danh sách tất cả user với phân trang
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "10"), 10, 64)

	users, totalCount, err := c.userService.GetAll(ctx, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"data":    []models.User{},
			"message": "Không có dữ liệu nào",
			"notes":   "Không có người dùng nào được tìm thấy",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users, "totalCount": totalCount})
}

// UpdateUser cập nhật thông tin user
func (c *UserController) UpdateUser(ctx *gin.Context) {
	ctx.Param("id")
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := c.userService.Update(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedUser})
}

// DeleteUser xóa user theo ID
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.userService.DeleteUser(ctx, id); err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "user deleted"})
}

// Login xử lý đăng nhập người dùng
func (c *UserController) Login(ctx *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	fmt.Println("Mật khẩu đăng nhập:", req.Password)
	token, err := c.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
