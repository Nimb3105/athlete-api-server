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

// ReminderRepository cung cấp các phương thức CRUD cho Reminder
type ReminderRepository struct {
	collection *mongo.Collection
}

func NewReminderRepository(collection *mongo.Collection) *ReminderRepository {
	return &ReminderRepository{collection}
}

func (r *ReminderRepository) Create(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
	reminder.CreatedAt = time.Now()
	reminder.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, reminder)
	if err != nil {
		return nil, err
	}

	reminder.ID = result.InsertedID.(primitive.ObjectID)
	return reminder, nil
}

// GetByID lấy lời nhắc theo ID
func (r *ReminderRepository) GetByID(ctx context.Context, id string) (*models.Reminder, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var reminder models.Reminder
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&reminder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("reminder not found")
		}
		return nil, err
	}

	return &reminder, nil
}

// GetByUserID lấy danh sách lời nhắc theo UserID
func (r *ReminderRepository) GetByUserID(ctx context.Context, userID string, page, limit int64) ([]models.Reminder, error) {
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

	var reminders []models.Reminder
	if err = cursor.All(ctx, &reminders); err != nil {
		return nil, err
	}

	return reminders, nil
}

// Update cập nhật thông tin lời nhắc
func (r *ReminderRepository) Update(ctx context.Context, reminder *models.Reminder) (*models.Reminder, error) {
	reminder.UpdatedAt = time.Now()

	filter := bson.M{"_id": reminder.ID}
	update := bson.M{"$set": reminder}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

// Delete xóa lời nhắc theo ID
func (r *ReminderRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("reminder not found")
	}

	return nil
}
