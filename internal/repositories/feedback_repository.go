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

// FeedbackRepository cung cấp các phương thức CRUD cho Feedback
type FeedbackRepository struct {
	collection *mongo.Collection
}

func NewFeedbackRepository(collection *mongo.Collection) *FeedbackRepository {
	return &FeedbackRepository{collection}
}

func (r *FeedbackRepository) Create(ctx context.Context, feedback *models.Feedback) (*models.Feedback, error) {
	feedback.CreatedAt = time.Now()
	feedback.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, feedback)
	if err != nil {
		return nil, err
	}

	feedback.ID = result.InsertedID.(primitive.ObjectID)
	return feedback, nil
}

func (r *FeedbackRepository) GetByID(ctx context.Context, id string) (*models.Feedback, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var feedback models.Feedback
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&feedback)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("feedback not found")
		}
		return nil, err
	}

	return &feedback, nil
}

func (r *FeedbackRepository) GetByUserID(ctx context.Context, userID string) ([]models.Feedback, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var feedbacks []models.Feedback
	if err = cursor.All(ctx, &feedbacks); err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (r *FeedbackRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Feedback, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var feedbacks []models.Feedback
	if err = cursor.All(ctx, &feedbacks); err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (r *FeedbackRepository) Update(ctx context.Context, feedback *models.Feedback) (*models.Feedback, error) {
	feedback.UpdatedAt = time.Now()

	filter := bson.M{"_id": feedback.ID}
	update := bson.M{"$set": feedback}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

func (r *FeedbackRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("feedback not found")
	}

	return nil
}