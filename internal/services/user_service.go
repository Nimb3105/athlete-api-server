package services

import (
	"be/config"
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService cung cấp các phương thức nghiệp vụ cho User
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService tạo một UserService mới
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo}
}

// Create tạo một user mới
func (s *UserService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if user.Email == "" {
		return nil, errors.New("email is required")
	}
	// Kiểm tra email đã tồn tại
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(ctx, user)
}

// GetByID lấy user theo ID
func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetByEmail lấy user theo email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// GetAll lấy danh sách tất cả user với phân trang
func (s *UserService) GetAll(ctx context.Context, page, limit int64) ([]models.User, error) {
	return s.userRepo.GetAll(ctx, page, limit)
}

// Update cập nhật thông tin user
func (s *UserService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if user.ID.IsZero() {
		return nil, errors.New("invalid user ID")
	}
	return s.userRepo.Update(ctx, user)
}

// Delete xóa user theo ID
func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// Login xác thực người dùng và trả về JWT token
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {

	//lấy config
	var cfg = config.LoadConfig()

	// Lấy người dùng từ repository
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Kiểm tra mật khẩu
	password = strings.TrimSpace(password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Tạo JWT token với Role
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"role":    user.Role, // Thêm Role vào claims
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret)) // Lưu key trong config
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
