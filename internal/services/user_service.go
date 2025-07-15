package services

import (
	"be/config"
	"be/internal/models"
	"be/internal/repositories"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserService cung cấp các phương thức nghiệp vụ cho User
type UserService struct {
	client      *mongo.Client
	userRepo    *repositories.UserRepository
	athleteRepo *repositories.AthleteRepository
	coachRepo   *repositories.CoachRepository
	DB          *mongo.Database
	coachAthleteRepo  *repositories.CoachAthleteRepository
}


// NewUserService tạo một UserService mới
func NewUserService(coachAthleteRepo  *repositories.CoachAthleteRepository,DB *mongo.Database, client *mongo.Client, userRepo *repositories.UserRepository, athleteRepo *repositories.AthleteRepository,
	coachRepo *repositories.CoachRepository) *UserService {
	return &UserService{client: client, userRepo: userRepo, athleteRepo: athleteRepo, coachRepo: coachRepo, DB: DB,coachAthleteRepo: coachAthleteRepo}
}


func (s *UserService) GetUnassignedAthletes(ctx context.Context, sportId string) ([]models.User, error) {
	assignedAthleteIds, err := s.coachAthleteRepo.GetAllAssignedAthleteIds(ctx)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindUnassignedAthletesBySport(ctx, sportId, assignedAthleteIds)
}

func (s *UserService) GetUsersByRoleWithPagination(ctx context.Context, page, limit int64, role string) ([]models.User, int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
	return s.userRepo.GetUsersByRoleWithPagination(ctx, page, limit, role)
}

func (s *UserService) GetAllUserCoachBySportId(ctx context.Context, page, limit int64, sportId string) ([]models.User, int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
	return s.userRepo.GetAllUserCoachBySportId(ctx, page, limit, sportId)
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
func (s *UserService) GetAll(ctx context.Context, page, limit int64) ([]models.User, int64, error) {
	if page < 1 || limit < 1 {
		return nil, 0, errors.New("invalid page or limit")
	}
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
// func (s *UserService) Delete(ctx context.Context, id string) error {
// 	// Kiểm tra xem User có trong bảng Athlete hoặc Coach không
// 	hasAthlete, err := s.athleteRepo.Exists(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	hasCoach, err := s.coachRepo.Exists(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	// Nếu User có trong bảng Athlete, xóa luôn Athlete
// 	if hasAthlete {
// 		if err := s.athleteRepo.Delete(ctx, id); err != nil {
// 			return err
// 		}
// 	}

// 	// Nếu User có trong bảng Coach, xóa luôn Coach
// 	if hasCoach {
// 		if err := s.coachRepo.Delete(ctx, id); err != nil {
// 			return err
// 		}
// 	}

// 	// Xóa User
// 	if err := s.userRepo.Delete(ctx, id); err != nil {
// 		return err
// 	}

// 	return nil
// }

// internal/services/user_service.go

// internal/services/user_service.go

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID người dùng không hợp lệ: %w", err)
	}

	// === KIỂM TRA KHÓA NGOẠI ===
	configs := []repositories.ForeignKeyCheckConfig{
		// ... (giữ nguyên phần kiểm tra khóa ngoại của bạn)
		{Collection: s.DB.Collection("achievements"), Filter: bson.M{"userId": objectID}, Name: "thành tích"},
		{Collection: s.DB.Collection("coach_certifications"), Filter: bson.M{"userId": objectID}, Name: "chứng chỉ huấn luyện viên"},
		{Collection: s.DB.Collection("feedbacks"), Filter: bson.M{"userId": objectID}, Name: "phản hồi"},
		{Collection: s.DB.Collection("groups"), Filter: bson.M{"createdBy": objectID}, Name: "nhóm"},
		{Collection: s.DB.Collection("group_members"), Filter: bson.M{"userId": objectID}, Name: "thành viên nhóm"},
		{Collection: s.DB.Collection("healths"), Filter: bson.M{"userId": objectID}, Name: "sức khỏe"},
		{Collection: s.DB.Collection("injuries"), Filter: bson.M{"userId": objectID}, Name: "chấn thương"},
		{Collection: s.DB.Collection("messages"), Filter: bson.M{"senderId": objectID}, Name: "tin nhắn"},
		{Collection: s.DB.Collection("notifications"), Filter: bson.M{"userId": objectID}, Name: "thông báo"},
		{Collection: s.DB.Collection("reminders"), Filter: bson.M{"userId": objectID}, Name: "lời nhắc"},
		{Collection: s.DB.Collection("team_members"), Filter: bson.M{"userId": objectID}, Name: "thành viên đội"},
		{Collection: s.DB.Collection("training_schedule_users"), Filter: bson.M{"userId": objectID}, Name: "người dùng lịch tập luyện"},
		{Collection: s.DB.Collection("nutrition_plans"), Filter: bson.M{"$or": []bson.M{{"userId": objectID}, {"createBy": objectID}}}, Name: "kế hoạch dinh dưỡng"},
		{Collection: s.DB.Collection("training_schedules"), Filter: bson.M{"createdBy": objectID}, Name: "lịch tập luyện"},
		{Collection: s.DB.Collection("user_matches"), Filter: bson.M{"userId": objectID}, Name: "trận đấu của vận động viên"},
		{Collection: s.DB.Collection("coach_athletes"), Filter: bson.M{"$or": []bson.M{{"athleteId": objectID}, {"coachId": objectID}}}, Name: "mối quan hệ huấn luyện viên - vận động viên"},
	}

	if err := repositories.CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// === TIẾN HÀNH XÓA ===
	// Xóa bản ghi Athlete (nếu có)
	if err := s.athleteRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("không thể xóa vận động viên: %w", err)
	}

	// Xóa bản ghi Coach (nếu có)
	if err := s.coachRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("không thể xóa huấn luyện viên: %w", err)
	}

	// Cuối cùng, xóa user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("không thể xóa người dùng: %w", err)
	}

	return nil
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
