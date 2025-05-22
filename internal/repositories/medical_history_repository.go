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

// MedicalHistoryRepository provides CRUD methods for MedicalHistory
type MedicalHistoryRepository struct {
	collection *mongo.Collection
}

func NewMedicalHistoryRepository(collection *mongo.Collection) *MedicalHistoryRepository {
	return &MedicalHistoryRepository{collection}
}

func (r *MedicalHistoryRepository) Create(ctx context.Context, medicalHistory *models.MedicalHistory) (*models.MedicalHistory, error) {
	medicalHistory.CreatedAt = time.Now()
	medicalHistory.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, medicalHistory)
	if err != nil {
		return nil, err
	}

	medicalHistory.ID = result.InsertedID.(primitive.ObjectID)
	return medicalHistory, nil
}

func (r *MedicalHistoryRepository) GetByID(ctx context.Context, id string) (*models.MedicalHistory, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var medicalHistory models.MedicalHistory
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&medicalHistory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("medical history not found")
		}
		return nil, err
	}

	return &medicalHistory, nil
}

func (r *MedicalHistoryRepository) GetByHealthID(ctx context.Context, healthID string) (*models.MedicalHistory, error) {
	objectID, err := primitive.ObjectIDFromHex(healthID)
	if err != nil {
		return nil, err
	}

	var medicalHistory models.MedicalHistory
	err = r.collection.FindOne(ctx, bson.M{"healthId": objectID}).Decode(&medicalHistory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("medical history not found")
		}
		return nil, err
	}

	return &medicalHistory, nil
}

func (r *MedicalHistoryRepository) GetAll(ctx context.Context, page, limit int64) ([]models.MedicalHistory, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var medicalHistories []models.MedicalHistory
	if err = cursor.All(ctx, &medicalHistories); err != nil {
		return nil, err
	}

	return medicalHistories, nil
}

func (r *MedicalHistoryRepository) Update(ctx context.Context, medicalHistory *models.MedicalHistory) (*models.MedicalHistory, error) {
	medicalHistory.UpdatedAt = time.Now()

	filter := bson.M{"_id": medicalHistory.ID}
	update := bson.M{"$set": medicalHistory}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return medicalHistory, nil
}

func (r *MedicalHistoryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("medical history not found")
	}

	return nil
}