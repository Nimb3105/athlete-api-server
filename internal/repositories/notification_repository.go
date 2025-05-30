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

// NotificationRepository cung cấp các phương thức CRUD cho Notification
type NotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(collection *mongo.Collection) *NotificationRepository {
	return &NotificationRepository{collection}
}

// GetAll lấy tất cả thông báo với phân trang
func (r *NotificationRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Notification, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var notifications []models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) (*models.Notification, error) {
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		return nil, err
	}

	notification.ID = result.InsertedID.(primitive.ObjectID)
	return notification, nil
}

// GetByID lấy thông báo theo ID
func (r *NotificationRepository) GetByID(ctx context.Context, id string) (*models.Notification, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var notification models.Notification
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}

	return &notification, nil
}

// GetByUserID lấy danh sách thông báo theo UserID
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID string, page, limit int64) ([]models.Notification, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

// Update cập nhật thông tin thông báo
func (r *NotificationRepository) Update(ctx context.Context, notification *models.Notification) (*models.Notification, error) {
	notification.UpdatedAt = time.Now()

	filter := bson.M{"_id": notification.ID}
	update := bson.M{"$set": notification}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

// Delete xóa thông báo theo ID
func (r *NotificationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("notification not found")
	}

	return nil
}
