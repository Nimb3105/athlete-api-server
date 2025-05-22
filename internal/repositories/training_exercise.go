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

type TrainingExerciseRepository struct {
	collection *mongo.Collection
}

func NewTrainingExerciseRepository(collection *mongo.Collection) *TrainingExerciseRepository {
	return &TrainingExerciseRepository{collection}
}

func (r *TrainingExerciseRepository) Create(ctx context.Context, trainingExercise *models.TrainingExercise) (*models.TrainingExercise, error) {
	trainingExercise.CreatedAt = time.Now()
	trainingExercise.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, trainingExercise)
	if err != nil {
		return nil, err
	}

	trainingExercise.ID = result.InsertedID.(primitive.ObjectID)
	return trainingExercise, nil
}

func (r *TrainingExerciseRepository) GetByID(ctx context.Context, id string) (*models.TrainingExercise, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var trainingExercise models.TrainingExercise
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&trainingExercise)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("training exercise not found")
		}
		return nil, err
	}

	return &trainingExercise, nil
}

func (r *TrainingExerciseRepository) GetByScheduleID(ctx context.Context, scheduleID string) ([]models.TrainingExercise, error) {
	objectID, err := primitive.ObjectIDFromHex(scheduleID)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "order", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{"scheduleId": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trainingExercises []models.TrainingExercise
	if err = cursor.All(ctx, &trainingExercises); err != nil {
		return nil, err
	}

	return trainingExercises, nil
}

func (r *TrainingExerciseRepository) GetAll(ctx context.Context, page, limit int64) ([]models.TrainingExercise, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trainingExercises []models.TrainingExercise
	if err = cursor.All(ctx, &trainingExercises); err != nil {
		return nil, err
	}

	return trainingExercises, nil
}

func (r *TrainingExerciseRepository) Update(ctx context.Context, trainingExercise *models.TrainingExercise) (*models.TrainingExercise, error) {
	trainingExercise.UpdatedAt = time.Now()

	filter := bson.M{"_id": trainingExercise.ID}
	update := bson.M{"$set": trainingExercise}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return trainingExercise, nil
}

func (r *TrainingExerciseRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("training exercise not found")
	}

	return nil
}

func (r *TrainingExerciseRepository) GetAllByScheduleID(ctx context.Context, scheduleID string) ([]models.TrainingExercise, error) {
	// Convert scheduleID from hex string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(scheduleID)
	if err != nil {
		return nil, err
	}

	// Find all documents matching the filter
	cursor, err := r.collection.Find(ctx, bson.M{"scheduleId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Slice to store the results
	var exercises []models.TrainingExercise

	// Decode all documents into the exercises slice
	if err = cursor.All(ctx, &exercises); err != nil {
		return nil, err
	}

	return exercises, nil
}
