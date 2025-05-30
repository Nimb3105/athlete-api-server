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

type ExerciseRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewExerciseRepository(collection *mongo.Collection, db *mongo.Database) *ExerciseRepository {
	return &ExerciseRepository{collection: collection, db: db}
}

func (r *ExerciseRepository) Create(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error) {
	exercise.CreatedAt = time.Now()
	exercise.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, exercise)
	if err != nil {
		return nil, err
	}

	exercise.ID = result.InsertedID.(primitive.ObjectID)
	return exercise, nil
}

func (r *ExerciseRepository) GetByID(ctx context.Context, id string) (*models.Exercise, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var exercise models.Exercise
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&exercise)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("exercise not found")
		}
		return nil, err
	}

	return &exercise, nil
}

func (r *ExerciseRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Exercise,int64, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil,0, err
	}
	defer cursor.Close(ctx)

	var exercises []models.Exercise
	if err = cursor.All(ctx, &exercises); err != nil {
		return nil,0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return exercises, totalCount, nil
}

func (r *ExerciseRepository) Update(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error) {
	exercise.UpdatedAt = time.Now()

	filter := bson.M{"_id": exercise.ID}
	update := bson.M{"$set": exercise}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return exercise, nil
}

func (r *ExerciseRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID bài tập không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("training_exercises"), bson.M{"exerciseId": objectID}, "bài tập trong lịch tập luyện"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa exercise
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa bài tập: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("bài tập không tồn tại")
	}

	return nil
}
