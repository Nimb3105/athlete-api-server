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

// CoachAthleteRepository provides CRUD methods for CoachAthlete
type CoachAthleteRepository struct {
	collection *mongo.Collection
}

func NewCoachAthleteRepository(collection *mongo.Collection) *CoachAthleteRepository {
	return &CoachAthleteRepository{collection: collection}
}

func (r *CoachAthleteRepository) Create(ctx context.Context, coachAthlete *models.CoachAthlete) (*models.CoachAthlete, error) {
	coachAthlete.CreatedAt = time.Now()
	coachAthlete.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, coachAthlete)
	if err != nil {
		return nil, err
	}

	coachAthlete.ID = result.InsertedID.(primitive.ObjectID)
	return coachAthlete, nil
}

func (r *CoachAthleteRepository) GetByID(ctx context.Context, id string) (*models.CoachAthlete, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid coach-athlete ID: %w", err)
	}

	var coachAthlete models.CoachAthlete
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&coachAthlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("coach-athlete relationship not found")
		}
		return nil, err
	}

	return &coachAthlete, nil
}

func (r *CoachAthleteRepository) GetAll(ctx context.Context, page, limit int64) ([]models.CoachAthlete, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var coachAthletes []models.CoachAthlete
	if err = cursor.All(ctx, &coachAthletes); err != nil {
		return nil, err
	}

	return coachAthletes, nil
}

func (r *CoachAthleteRepository) Update(ctx context.Context, coachAthlete *models.CoachAthlete) (*models.CoachAthlete, error) {
	coachAthlete.UpdatedAt = time.Now()

	filter := bson.M{"_id": coachAthlete.ID}
	update := bson.M{"$set": coachAthlete}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return coachAthlete, nil
}

func (r *CoachAthleteRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid coach-athlete ID: %w", err)
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("error deleting coach-athlete relationship: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("coach-athlete relationship not found")
	}

	return nil
}
