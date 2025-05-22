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

// InjuryRepository provides CRUD methods for Injury
type InjuryRepository struct {
	collection *mongo.Collection
}

func NewInjuryRepository(collection *mongo.Collection) *InjuryRepository {
	return &InjuryRepository{collection}
}

func (r *InjuryRepository) Create(ctx context.Context, injury *models.Injury) (*models.Injury, error) {
	injury.CreatedAt = time.Now()
	injury.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, injury)
	if err != nil {
		return nil, err
	}

	injury.ID = result.InsertedID.(primitive.ObjectID)
	return injury, nil
}

func (r *InjuryRepository) GetByID(ctx context.Context, id string) (*models.Injury, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var injury models.Injury
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&injury)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("injury not found")
		}
		return nil, err
	}

	return &injury, nil
}

func (r *InjuryRepository) GetByUserID(ctx context.Context, userID string) (*models.Injury, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var injury models.Injury
	err = r.collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&injury)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("injury not found")
		}
		return nil, err
	}

	return &injury, nil
}

func (r *InjuryRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Injury, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var injuries []models.Injury
	if err = cursor.All(ctx, &injuries); err != nil {
		return nil, err
	}

	return injuries, nil
}

func (r *InjuryRepository) Update(ctx context.Context, injury *models.Injury) (*models.Injury, error) {
	injury.UpdatedAt = time.Now()

	filter := bson.M{"_id": injury.ID}
	update := bson.M{"$set": injury}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return injury, nil
}

func (r *InjuryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("injury not found")
	}

	return nil
}