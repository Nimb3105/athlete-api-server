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

// CoachCertificationRepository cung cấp các phương thức CRUD cho CoachCertification
type CoachCertificationRepository struct {
	collection *mongo.Collection
}

func NewCoachCertificationRepository(collection *mongo.Collection) *CoachCertificationRepository {
	return &CoachCertificationRepository{collection}
}

func (r *CoachCertificationRepository) Create(ctx context.Context, certification *models.CoachCertification) (*models.CoachCertification, error) {
	certification.CreatedAt = time.Now()
	certification.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, certification)
	if err != nil {
		return nil, err
	}

	certification.ID = result.InsertedID.(primitive.ObjectID)
	return certification, nil
}

func (r *CoachCertificationRepository) GetByID(ctx context.Context, id string) (*models.CoachCertification, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var certification models.CoachCertification
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&certification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("coach certification not found")
		}
		return nil, err
	}

	return &certification, nil
}

func (r *CoachCertificationRepository) GetByUserID(ctx context.Context, userID string) ([]models.CoachCertification, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certifications []models.CoachCertification
	if err = cursor.All(ctx, &certifications); err != nil {
		return nil, err
	}

	return certifications, nil
}

func (r *CoachCertificationRepository) GetAll(ctx context.Context, page, limit int64) ([]models.CoachCertification, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certifications []models.CoachCertification
	if err = cursor.All(ctx, &certifications); err != nil {
		return nil, err
	}

	return certifications, nil
}

func (r *CoachCertificationRepository) Update(ctx context.Context, certification *models.CoachCertification) (*models.CoachCertification, error) {
	certification.UpdatedAt = time.Now()

	filter := bson.M{"_id": certification.ID}
	update := bson.M{"$set": certification}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return certification, nil
}

func (r *CoachCertificationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("coach certification not found")
	}

	return nil
}