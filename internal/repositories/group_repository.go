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

// GroupRepository cung cấp các phương thức CRUD cho Group
type GroupRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewGroupRepository(collection *mongo.Collection,db         *mongo.Database) *GroupRepository {
	return &GroupRepository{collection:collection,db:db,}
}

func (r *GroupRepository) Create(ctx context.Context, group *models.Group) (*models.Group, error) {
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	group.CreatedDate = time.Now()

	result, err := r.collection.InsertOne(ctx, group)
	if err != nil {
		return nil, err
	}

	group.ID = result.InsertedID.(primitive.ObjectID)
	return group, nil
}

func (r *GroupRepository) GetByID(ctx context.Context, id string) (*models.Group, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var group models.Group
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("group not found")
		}
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepository) GetByCreatedBy(ctx context.Context, createdBy string) ([]models.Group, error) {
	objectID, err := primitive.ObjectIDFromHex(createdBy)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"createdBy": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []models.Group
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Group, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []models.Group
	if err = cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepository) Update(ctx context.Context, group *models.Group) (*models.Group, error) {
	group.UpdatedAt = time.Now()

	filter := bson.M{"_id": group.ID}
	update := bson.M{"$set": group}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (r *GroupRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID nhóm không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("group_members"), bson.M{"groupId": objectID}, "thành viên nhóm"},
		{r.db.Collection("messages"), bson.M{"groupId": objectID}, "tin nhắn"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa group
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa nhóm: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("nhóm không tồn tại")
	}

	return nil
}