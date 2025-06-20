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

// AthleteRepository cung cấp các phương thức CRUD cho Athlete
type AthleteRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

// NewAthleteRepository tạo một AthleteRepository mới
func NewAthleteRepository(athleteCollection *mongo.Collection, db *mongo.Database) *AthleteRepository {
	return &AthleteRepository{collection: athleteCollection, db: db}
}

// Create tạo một athlete mới
func (r *AthleteRepository) Create(ctx context.Context, athlete *models.Athlete) (*models.Athlete, error) {
	athlete.CreatedAt = time.Now()
	athlete.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, athlete)
	if err != nil {
		return nil, err
	}

	athlete.ID = result.InsertedID.(primitive.ObjectID)
	return athlete, nil
}

// GetByID lấy athlete theo ID
func (r *AthleteRepository) GetByID(ctx context.Context, id string) (*models.Athlete, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var athlete models.Athlete
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&athlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("athlete not found")
		}
		return nil, err
	}

	return &athlete, nil
}

// GetByUserID lấy athlete theo UserID
func (r *AthleteRepository) GetByUserID(ctx context.Context, userID string) (*models.Athlete, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var athlete models.Athlete
	err = r.collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&athlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("athlete not found")
		}
		return nil, err
	}

	return &athlete, nil
}

// GetAll lấy danh sách tất cả athlete với phân trang
func (r *AthleteRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Athlete, int64, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var athletes []models.Athlete
	if err = cursor.All(ctx, &athletes); err != nil {
		return nil, 0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return athletes, totalCount, nil
}

// Update cập nhật thông tin athlete
func (r *AthleteRepository) Update(ctx context.Context, athlete *models.Athlete) (*models.Athlete, error) {
	athlete.UpdatedAt = time.Now()

	filter := bson.M{"_id": athlete.ID}
	update := bson.M{"$set": athlete}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return athlete, nil
}

// Delete xóa athlete theo ID
func (r *AthleteRepository) Delete(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("ID vận động viên không hợp lệ: %w", err)
	}
	// Xóa athlete
	result, err := r.collection.DeleteOne(ctx, bson.M{"userId": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa vận động viên: %w", err)
	}

	if result.DeletedCount == 0 {
		return nil
	}

	return nil
}

func (r *AthleteRepository) Exists(ctx context.Context, userID string) (bool, error) {
	// Kiểm tra xem userId có tồn tại không
	count, err := r.collection.CountDocuments(ctx, bson.M{"userId": userID})
	if err != nil {
		return false, fmt.Errorf("lỗi khi kiểm tra vận động viên: %w", err)
	}

	// Nếu số lượng lớn hơn 0, tức là tồn tại
	return count > 0, nil
}
