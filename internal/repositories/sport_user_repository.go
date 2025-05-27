package repositories

import (
	"be/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SportUserRepository cung cấp các phương thức CRUD cho SportUser
type SportUserRepository struct {
	collection *mongo.Collection
}

// NewSportUserRepository tạo một SportUserRepository mới
func NewSportUserRepository(collection *mongo.Collection) *SportUserRepository {
	return &SportUserRepository{collection}
}

// Create tạo một SportUser mới
func (r *SportUserRepository) Create(ctx context.Context, SportUser *models.SportUser) (*models.SportUser, error) {
	SportUser.CreatedAt = time.Now()
	SportUser.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, SportUser)
	if err != nil {
		return nil, err
	}

	SportUser.ID = result.InsertedID.(primitive.ObjectID)
	return SportUser, nil
}

// GetByID lấy SportUser theo ID
func (r *SportUserRepository) GetByID(ctx context.Context, id string) (*models.SportUser, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var SportUser models.SportUser
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&SportUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("SportUser not found")
		}
		return nil, err
	}

	return &SportUser, nil
}

// GetByUserID lấy SportUser theo UserID
func (r *SportUserRepository) GetByUserID(ctx context.Context, userID string) (*models.SportUser, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var SportUser models.SportUser
	err = r.collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&SportUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("SportUser not found")
		}
		return nil, err
	}

	return &SportUser, nil
}

func (r *SportUserRepository) GetBySportID(ctx context.Context, SportID string) (*models.SportUser, error) {
	objectID, err := primitive.ObjectIDFromHex(SportID)
	if err != nil {
		return nil, err
	}

	var SportUser models.SportUser
	err = r.collection.FindOne(ctx, bson.M{"sportId": objectID}).Decode(&SportUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("SportUser not found")
		}
		return nil, err
	}

	return &SportUser, nil
}

// GetAll lấy danh sách tất cả SportUser với phân trang
func (r *SportUserRepository) GetAll(ctx context.Context, page, limit int64) ([]models.SportUser, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var SportUsers []models.SportUser
	if err = cursor.All(ctx, &SportUsers); err != nil {
		return nil, err
	}

	return SportUsers, nil
}

// Update cập nhật thông tin SportUser
func (r *SportUserRepository) Update(ctx context.Context, SportUser *models.SportUser) (*models.SportUser, error) {
	SportUser.UpdatedAt = time.Now()

	filter := bson.M{"_id": SportUser.ID}
	update := bson.M{"$set": SportUser}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return SportUser, nil
}

// Delete xóa SportUser theo ID
func (r *SportUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("SportUser not found")
	}

	return nil
}

// GetAllByUserID lấy danh sách tất cả SportUser theo userID
func (r *SportUserRepository) GetAllByUserID(ctx context.Context, userID string) ([]models.SportUser, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var SportUsers []models.SportUser
	if err = cursor.All(ctx, &SportUsers); err != nil {
		return nil, err
	}

	return SportUsers, nil
}
