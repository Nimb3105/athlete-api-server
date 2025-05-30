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

// AthleteMatchRepository cung cấp các phương thức CRUD cho AthleteMatch
type UserMatchRepository struct {
	collection *mongo.Collection
}

func NewUserMatchRepository(collection *mongo.Collection) *UserMatchRepository {
	return &UserMatchRepository{collection}
}

func (r *UserMatchRepository) Create(ctx context.Context, userMatch *models.UserMatch) (*models.UserMatch, error) {
	userMatch.CreatedAt = time.Now()
	userMatch.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, userMatch)
	if err != nil {
		return nil, err
	}

	userMatch.ID = result.InsertedID.(primitive.ObjectID)
	return userMatch, nil
}

func (r *UserMatchRepository) GetByID(ctx context.Context, id string) (*models.UserMatch, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var userMatch models.UserMatch
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&userMatch)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("athlete match not found")
		}
		return nil, err
	}

	return &userMatch, nil
}

func (r *UserMatchRepository) GetByUserID(ctx context.Context, userID string) ([]models.UserMatch, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userMatches []models.UserMatch
	if err = cursor.All(ctx, &userMatches); err != nil {
		return nil, err
	}

	return userMatches, nil
}

func (r *UserMatchRepository) GetAll(ctx context.Context, page, limit int64) ([]models.UserMatch, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userMatches []models.UserMatch
	if err = cursor.All(ctx, &userMatches); err != nil {
		return nil, err
	}

	return userMatches, nil
}

func (r *UserMatchRepository) Update(ctx context.Context, userMatch *models.UserMatch) (*models.UserMatch, error) {
	userMatch.UpdatedAt = time.Now()

	filter := bson.M{"_id": userMatch.ID}
	update := bson.M{"$set": userMatch}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return userMatch, nil
}

func (r *UserMatchRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("athlete match not found")
	}

	return nil
}
