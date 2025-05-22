package repositories

import (
	"be/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CoachRepository cung cấp các phương thức CRUD cho Coach
type CoachRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewCoachRepository(coachCollection *mongo.Collection, db *mongo.Database) *CoachRepository {
	return &CoachRepository{collection: coachCollection, db: db}
}

func (r *CoachRepository) Create(ctx context.Context, Coach *models.Coach) (*models.Coach, error) {
	Coach.CreatedAt = time.Now()
	Coach.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, Coach)
	if err != nil {
		return nil, err
	}

	Coach.ID = result.InsertedID.(primitive.ObjectID)
	return Coach, nil
}

// GetByID lấy người dùng theo ID
func (r *CoachRepository) GetByID(ctx context.Context, id string) (*models.Coach, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var Coach models.Coach
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&Coach)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("coach not found")
		}
		return nil, err
	}

	return &Coach, nil
}

func (r *CoachRepository) GetByUserID(ctx context.Context, userID string) (*models.Coach, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var coach models.Coach
	err = r.collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&coach)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("coach not found")
		}
		return nil, err
	}

	return &coach, nil
}

// GetAll lấy danh sách tất cả người dùng với phân trang
func (r *CoachRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Coach, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Coachs []models.Coach
	if err = cursor.All(ctx, &Coachs); err != nil {
		return nil, err
	}

	return Coachs, nil
}

// Update cập nhật thông tin người dùng
func (r *CoachRepository) Update(ctx context.Context, Coach *models.Coach) (*models.Coach, error) {
	Coach.UpdatedAt = time.Now()

	filter := bson.M{"_id": Coach.ID}
	update := bson.M{"$set": Coach}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return Coach, nil
}

// Delete xóa người dùng theo ID
func (r *CoachRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID huấn luyện viên không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("nutrition_plans"), bson.M{"coachId": objectID}, "kế hoạch dinh dưỡng"},
		{r.db.Collection("training_schedules"), bson.M{"createdBy": objectID}, "lịch tập luyện"},
		{r.db.Collection("coach_certifications"), bson.M{"userId": objectID}, "chứng chỉ huấn luyện viên"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa coach
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa huấn luyện viên: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("huấn luyện viên không tồn tại")
	}

	return nil
}
