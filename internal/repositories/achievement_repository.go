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

// AchievementRepository cung cấp các phương thức CRUD cho Achievement
type AchievementRepository struct {
	collection     *mongo.Collection
}

func NewAchievementRepository(achievementCollection *mongo.Collection) *AchievementRepository {
	return &AchievementRepository{
		collection:     achievementCollection,
	}
}

func (r *AchievementRepository) Create(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error) {
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return achievement, nil
}

func (r *AchievementRepository) GetByID(ctx context.Context, id string) (*models.Achievement, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var achievement models.Achievement
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&achievement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("achievement not found")
		}
		return nil, err
	}

	return &achievement, nil
}

func (r *AchievementRepository) GetByUserID(ctx context.Context, userID string) ([]models.Achievement, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []models.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *AchievementRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Achievement, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []models.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *AchievementRepository) Update(ctx context.Context, achievement *models.Achievement) (*models.Achievement, error) {
	achievement.UpdatedAt = time.Now()

	filter := bson.M{"_id": achievement.ID}
	update := bson.M{"$set": achievement}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return achievement, nil
}

func (r *AchievementRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("achievement not found")
	}

	return nil
}
