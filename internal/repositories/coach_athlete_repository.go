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

func (r *CoachAthleteRepository) GetAllAssignedAthleteIds(ctx context.Context) ([]primitive.ObjectID, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []primitive.ObjectID
	for cursor.Next(ctx) {
		var rel struct {
			AthleteId primitive.ObjectID `bson:"athleteId"`
		}
		if err := cursor.Decode(&rel); err != nil {
			return nil, err
		}
		ids = append(ids, rel.AthleteId)
	}
	return ids, nil
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

func (r *CoachAthleteRepository) GetByAthleteId(ctx context.Context, athleteId string) (*models.CoachAthlete, error) {
	objectID, err := primitive.ObjectIDFromHex(athleteId)
	if err != nil {
		return nil, fmt.Errorf("athleteId không hợp lệ: %v", err)
	}
	var coachAthlete models.CoachAthlete
	err = r.collection.FindOne(ctx, bson.M{"athleteId": objectID}).Decode(&coachAthlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("coach-athlete relationship not found")
		}
		return nil, err
	}
	return &coachAthlete, nil
}

func (r *CoachAthleteRepository) GetAllByCoachId(ctx context.Context, coachId string, page, limit int64) ([]models.CoachAthlete, int64, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	// Lọc theo athleteId
	objID, err := primitive.ObjectIDFromHex(coachId)
	if err != nil {
		return nil, 0, fmt.Errorf("coachId không hợp lệ: %v", err)
	}
	filter := bson.M{"coachId": objID}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var coachAthletes []models.CoachAthlete
	if err = cursor.All(ctx, &coachAthletes); err != nil {
		return nil, 0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return coachAthletes, totalCount, nil
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
func (r *CoachAthleteRepository) DeleteAllByCoachId(ctx context.Context, coachId string) error {
	objectID, err := primitive.ObjectIDFromHex(coachId)
	if err != nil {
		return fmt.Errorf("invalid coach ID: %w", err)
	}
	filter := bson.M{"coachId": objectID}

	res, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no coach-athlete relationships found")
	}
	return nil
}
