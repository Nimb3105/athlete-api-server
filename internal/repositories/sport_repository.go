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

// SportRepository cung cấp các phương thức CRUD cho Sport
type SportRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

// NewSportRepository tạo một SportRepository mới
func NewSportRepository(collection *mongo.Collection,db         *mongo.Database) *SportRepository {
	return &SportRepository{collection:collection,db:db,}
}

// Create tạo một Sport mới
func (r *SportRepository) Create(ctx context.Context, Sport *models.Sport) (*models.Sport, error) {
	Sport.CreatedAt = time.Now()
	Sport.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, Sport)
	if err != nil {
		return nil, err
	}

	Sport.ID = result.InsertedID.(primitive.ObjectID)
	return Sport, nil
}

// GetByID lấy Sport theo ID
func (r *SportRepository) GetByID(ctx context.Context, id string) (*models.Sport, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var Sport models.Sport
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&Sport)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("sport not found")
		}
		return nil, err
	}

	return &Sport, nil
}

// GetAll lấy danh sách tất cả Sport với phân trang
func (r *SportRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Sport, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Sports []models.Sport
	if err = cursor.All(ctx, &Sports); err != nil {
		return nil, err
	}

	return Sports, nil
}

// Update cập nhật thông tin Sport
func (r *SportRepository) Update(ctx context.Context, Sport *models.Sport) (*models.Sport, error) {
	Sport.UpdatedAt = time.Now()

	filter := bson.M{"_id": Sport.ID}
	update := bson.M{"$set": Sport}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return Sport, nil
}

// Delete xóa Sport theo ID
func (r *SportRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID môn thể thao không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("sport_athletes"), bson.M{"sportId": objectID}, "vận động viên môn thể thao"},
		{r.db.Collection("teams"), bson.M{"sportId": objectID}, "đội"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa sport
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa môn thể thao: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("môn thể thao không tồn tại")
	}

	return nil
}
