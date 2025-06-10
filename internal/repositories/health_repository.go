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

// HealthRepository provides CRUD methods for Health
type HealthRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewHealthRepository(collection *mongo.Collection, db *mongo.Database) *HealthRepository {
	return &HealthRepository{collection: collection, db: db}
}

func (r *HealthRepository) Create(ctx context.Context, health *models.Health) (*models.Health, error) {
	health.CreatedAt = time.Now()
	health.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, health)
	if err != nil {
		return nil, err
	}

	health.ID = result.InsertedID.(primitive.ObjectID)
	return health, nil
}

func (r *HealthRepository) GetByID(ctx context.Context, id string) (*models.Health, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var health models.Health
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&health)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("health record not found")
		}
		return nil, err
	}

	return &health, nil
}

func (r *HealthRepository) GetByUserID(ctx context.Context, userID string) ([]models.Health, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var health []models.Health
	if err = cursor.All(ctx, &health); err != nil {
		return nil, err
	}

	return health, nil
}

func (r *HealthRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Health, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var healthRecords []models.Health
	if err = cursor.All(ctx, &healthRecords); err != nil {
		return nil, err
	}

	return healthRecords, nil
}

func (r *HealthRepository) Update(ctx context.Context, health *models.Health) (*models.Health, error) {
	health.UpdatedAt = time.Now()

	filter := bson.M{"_id": health.ID}
	update := bson.M{"$set": health}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return health, nil
}

func (r *HealthRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID hồ sơ sức khỏe không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("medical_histories"), bson.M{"healthId": objectID}, "lịch sử y tế"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa health
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa hồ sơ sức khỏe: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("hồ sơ sức khỏe không tồn tại")
	}

	return nil
}
