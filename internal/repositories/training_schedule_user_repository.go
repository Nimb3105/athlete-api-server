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

type TrainingScheduleUserRepository struct {
	collection *mongo.Collection
}

func NewTrainingScheduleUserRepository(collection *mongo.Collection) *TrainingScheduleUserRepository {
	return &TrainingScheduleUserRepository{collection}
}

func (r *TrainingScheduleUserRepository) Create(ctx context.Context, scheduleAthlete *models.TrainingScheduleUser) (*models.TrainingScheduleUser, error) {
	scheduleAthlete.CreatedAt = time.Now()
	scheduleAthlete.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, scheduleAthlete)
	if err != nil {
		return nil, err
	}

	scheduleAthlete.ID = result.InsertedID.(primitive.ObjectID)
	return scheduleAthlete, nil
}

func (r *TrainingScheduleUserRepository) GetByID(ctx context.Context, id string) (*models.TrainingScheduleUser, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var scheduleAthlete models.TrainingScheduleUser
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&scheduleAthlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("training schedule athlete not found")
		}
		return nil, err
	}

	return &scheduleAthlete, nil
}

func (r *TrainingScheduleUserRepository) GetByScheduleID(ctx context.Context, scheduleID string) ([]models.TrainingScheduleUser, error) {
	objectID, err := primitive.ObjectIDFromHex(scheduleID)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"scheduleId": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scheduleAthletes []models.TrainingScheduleUser
	if err = cursor.All(ctx, &scheduleAthletes); err != nil {
		return nil, err
	}

	return scheduleAthletes, nil
}

func (r *TrainingScheduleUserRepository) GetByUserID(ctx context.Context, userID string) ([]models.TrainingScheduleUser, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scheduleAthletes []models.TrainingScheduleUser
	if err = cursor.All(ctx, &scheduleAthletes); err != nil {
		return nil, err
	}

	return scheduleAthletes, nil
}

func (r *TrainingScheduleUserRepository) GetAll(ctx context.Context, page, limit int64) ([]models.TrainingScheduleUser, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scheduleAthletes []models.TrainingScheduleUser
	if err = cursor.All(ctx, &scheduleAthletes); err != nil {
		return nil, err
	}

	return scheduleAthletes, nil
}

func (r *TrainingScheduleUserRepository) Update(ctx context.Context, scheduleAthlete *models.TrainingScheduleUser) (*models.TrainingScheduleUser, error) {
	scheduleAthlete.UpdatedAt = time.Now()

	filter := bson.M{"_id": scheduleAthlete.ID}
	update := bson.M{"$set": scheduleAthlete}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return scheduleAthlete, nil
}

func (r *TrainingScheduleUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("training schedule athlete not found")
	}

	return nil
}

// GetAllByUserID lấy danh sách tất cả SportAthlete theo userID
func (r *TrainingScheduleUserRepository) GetAllByUserID(ctx context.Context, userID string) ([]models.TrainingScheduleUser, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var SportAthletes []models.TrainingScheduleUser
	if err = cursor.All(ctx, &SportAthletes); err != nil {
		return nil, err
	}

	return SportAthletes, nil
}
