 package services

// import (
// 	"be/internal/models"
// 	"be/internal/repositories"
// 	"context"
// 	"errors"
// )

// // SportUserService cung cấp các phương thức nghiệp vụ cho SportUser
// type SportUserService struct {
// 	SportUserRepo *repositories.SportUserRepository
// }

// // NewSportUserService tạo một SportUserService mới
// func NewSportUserService(SportUserRepo *repositories.SportUserRepository) *SportUserService {
// 	return &SportUserService{SportUserRepo}
// }

// // Create tạo một sport athlete mới
// func (s *SportUserService) Create(ctx context.Context, SportUser *models.SportUser) (*models.SportUser, error) {
// 	if SportUser.UserID.IsZero() {
// 		return nil, errors.New("user ID is required")
// 	}
// 	return s.SportUserRepo.Create(ctx, SportUser)
// }

// // GetByID lấy sport athlete theo ID
// func (s *SportUserService) GetByID(ctx context.Context, id string) (*models.SportUser, error) {
// 	return s.SportUserRepo.GetByID(ctx, id)
// }

// // GetByUserID lấy sport athlete theo UserID
// func (s *SportUserService) GetByUserID(ctx context.Context, userID string) (*models.SportUser, error) {
// 	return s.SportUserRepo.GetByUserID(ctx, userID)
// }

// func (s *SportUserService) GetBySportID(ctx context.Context, sportID string) (*models.SportUser, error) {
// 	return s.SportUserRepo.GetBySportID(ctx, sportID)
// }

// // GetAll lấy danh sách tất cả sport athlete với phân trang
// func (s *SportUserService) GetAll(ctx context.Context, page, limit int64) ([]models.SportUser, error) {
// 	return s.SportUserRepo.GetAll(ctx, page, limit)
// }

// // Update cập nhật thông tin sport athlete
// func (s *SportUserService) Update(ctx context.Context, SportUser *models.SportUser) (*models.SportUser, error) {
// 	if SportUser.ID.IsZero() {
// 		return nil, errors.New("invalid sport athlete ID")
// 	}
// 	return s.SportUserRepo.Update(ctx, SportUser)
// }

// // Delete xóa sport athlete theo ID
// func (s *SportUserService) Delete(ctx context.Context, id string) error {
// 	return s.SportUserRepo.Delete(ctx, id)
// }

// func (s *SportUserService) GetAllByUserID(ctx context.Context, userID string) ([]models.SportUser, error) {
// 	return s.SportUserRepo.GetAllByUserID(ctx, userID)
// }
