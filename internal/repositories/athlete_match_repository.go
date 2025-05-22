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
type AthleteMatchRepository struct {
	collection *mongo.Collection
}

func NewAthleteMatchRepository(collection *mongo.Collection) *AthleteMatchRepository {
	return &AthleteMatchRepository{collection}
}

func (r *AthleteMatchRepository) Create(ctx context.Context, athleteMatch *models.AthleteMatch) (*models.AthleteMatch, error) {
	athleteMatch.CreatedAt = time.Now()
	athleteMatch.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, athleteMatch)
	if err != nil {
		return nil, err
	}

	athleteMatch.ID = result.InsertedID.(primitive.ObjectID)
	return athleteMatch, nil
}

func (r *AthleteMatchRepository) GetByID(ctx context.Context, id string) (*models.AthleteMatch, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var athleteMatch models.AthleteMatch
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&athleteMatch)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("athlete match not found")
		}
		return nil, err
	}

	return &athleteMatch, nil
}

func (r *AthleteMatchRepository) GetByUserID(ctx context.Context, userID string) ([]models.AthleteMatch, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var athleteMatches []models.AthleteMatch
	if err = cursor.All(ctx, &athleteMatches); err != nil {
		return nil, err
	}

	return athleteMatches, nil
}

func (r *AthleteMatchRepository) GetAll(ctx context.Context, page, limit int64) ([]models.AthleteMatch, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var athleteMatches []models.AthleteMatch
	if err = cursor.All(ctx, &athleteMatches); err != nil {
		return nil, err
	}

	return athleteMatches, nil
}

func (r *AthleteMatchRepository) Update(ctx context.Context, athleteMatch *models.AthleteMatch) (*models.AthleteMatch, error) {
	athleteMatch.UpdatedAt = time.Now()

	filter := bson.M{"_id": athleteMatch.ID}
	update := bson.M{"$set": athleteMatch}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return athleteMatch, nil
}

func (r *AthleteMatchRepository) Delete(ctx context.Context, id string) error {
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
