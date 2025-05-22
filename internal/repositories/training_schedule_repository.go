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

type TrainingScheduleRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewTrainingScheduleRepository(collection *mongo.Collection,db         *mongo.Database) *TrainingScheduleRepository {
	return &TrainingScheduleRepository{collection:collection,db:db}
}

func (r *TrainingScheduleRepository) Create(ctx context.Context, schedule *models.TrainingSchedule) (*models.TrainingSchedule, error) {
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, schedule)
	if err != nil {
		return nil, err
	}

	schedule.ID = result.InsertedID.(primitive.ObjectID)
	return schedule, nil
}

func (r *TrainingScheduleRepository) GetByID(ctx context.Context, id string) (*models.TrainingSchedule, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var schedule models.TrainingSchedule
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("training schedule not found")
		}
		return nil, err
	}

	return &schedule, nil
}

func (r *TrainingScheduleRepository) GetAll(ctx context.Context, page, limit int64) ([]models.TrainingSchedule, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []models.TrainingSchedule
	if err = cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *TrainingScheduleRepository) Update(ctx context.Context, schedule *models.TrainingSchedule) (*models.TrainingSchedule, error) {
	schedule.UpdatedAt = time.Now()

	filter := bson.M{"_id": schedule.ID}
	update := bson.M{"$set": schedule}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *TrainingScheduleRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID lịch tập luyện không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("feedbacks"), bson.M{"scheduleId": objectID}, "phản hồi"},
		{r.db.Collection("performances"), bson.M{"scheduleId": objectID}, "hiệu suất"},
		{r.db.Collection("notifications"), bson.M{"scheduleId": objectID}, "thông báo"},
		{r.db.Collection("reminders"), bson.M{"scheduleId": objectID}, "lời nhắc"},
		{r.db.Collection("training_exercises"), bson.M{"scheduleId": objectID}, "bài tập trong lịch tập luyện"},
		{r.db.Collection("training_schedule_users"), bson.M{"scheduleId": objectID}, "người dùng lịch tập luyện"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa training schedule
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa lịch tập luyện: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("lịch tập luyện không tồn tại")
	}

	return nil
}