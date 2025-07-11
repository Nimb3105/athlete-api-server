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

type DailyScheduleRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewDailyScheduleRepository(collection *mongo.Collection, db *mongo.Database) *DailyScheduleRepository {
	return &DailyScheduleRepository{collection: collection, db: db}
}

func (r *DailyScheduleRepository) Create(ctx context.Context, dailySchedule *models.DailySchedule) (*models.DailySchedule, error) {
	dailySchedule.CreatedAt = time.Now()
	dailySchedule.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, dailySchedule)
	if err != nil {
		return nil, err
	}

	dailySchedule.ID = result.InsertedID.(primitive.ObjectID)

	return dailySchedule, nil
}

func (r *DailyScheduleRepository) GetByID(ctx context.Context, id string) (*models.DailySchedule, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var dailySchedule models.DailySchedule
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&dailySchedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("daily schedule not found")
		}
		return nil, err
	}

	return &dailySchedule, nil
}

func (r *DailyScheduleRepository) GetByUserID(ctx context.Context, day string, userId string) (*models.DailySchedule, error) {

	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse(time.RFC3339, day)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"userId":    objectID,
		"startDate": bson.M{"$lte": parsedDate},
		"endDate":   bson.M{"$gte": parsedDate},
	}

	var dailySchedule models.DailySchedule
	err = r.collection.FindOne(ctx, filter).Decode(&dailySchedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("daily schedule not found")
		}
		return nil, err
	}

	return &dailySchedule, nil
}

func (r *DailyScheduleRepository) GetAll(ctx context.Context, page, limit int64) ([]models.DailySchedule, int64, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var dailySchedule []models.DailySchedule
	if err = cursor.All(ctx, &dailySchedule); err != nil {
		return nil, 0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	return dailySchedule, totalCount, nil
}

func (r *DailyScheduleRepository) Update(ctx context.Context, dailySchedule *models.DailySchedule) (*models.DailySchedule, error) {
	dailySchedule.UpdatedAt = time.Now()

	filter := bson.M{"_id": dailySchedule.ID}
	update := bson.M{"$set": dailySchedule}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return dailySchedule, nil
}

func (r *DailyScheduleRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID daily không hợp lệ: %w", err)
	}

	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("training_schedules"), bson.M{"dailyScheduleId": objectID}, "Daily lịch tập nằm trong lịch tập"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa daily lịch tập: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("daily lịch tập không tồn tại")
	}

	return nil

}
