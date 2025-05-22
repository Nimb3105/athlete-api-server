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

// ProgressRepository provides CRUD methods for Progress
type ProgressRepository struct {
	collection *mongo.Collection
}

func NewProgressRepository(collection *mongo.Collection) *ProgressRepository {
	return &ProgressRepository{collection}
}

func (r *ProgressRepository) Create(ctx context.Context, progress *models.Progress) (*models.Progress, error) {
	progress.CreatedAt = time.Now()
	progress.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, progress)
	if err != nil {
		return nil, err
	}

	progress.ID = result.InsertedID.(primitive.ObjectID)
	return progress, nil
}

func (r *ProgressRepository) GetByID(ctx context.Context, id string) (*models.Progress, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var progress models.Progress
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&progress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("progress not found")
		}
		return nil, err
	}

	return &progress, nil
}

func (r *ProgressRepository) GetByUserID(ctx context.Context, userID string) ([]models.Progress, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var progresses []models.Progress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, err
	}

	return progresses, nil
}

func (r *ProgressRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Progress, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var progresses []models.Progress
	if err = cursor.All(ctx, &progresses); err != nil {
		return nil, err
	}

	return progresses, nil
}

func (r *ProgressRepository) Update(ctx context.Context, progress *models.Progress) (*models.Progress, error) {
	progress.UpdatedAt = time.Now()

	filter := bson.M{"_id": progress.ID}
	update := bson.M{"$set": progress}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *ProgressRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("progress not found")
	}

	return nil
}