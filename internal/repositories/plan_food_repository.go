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

// PlanFoodRepository provides CRUD methods for PlanFood
type PlanFoodRepository struct {
	collection *mongo.Collection
}

func NewPlanFoodRepository(collection *mongo.Collection) *PlanFoodRepository {
	return &PlanFoodRepository{collection: collection}
}

func (r *PlanFoodRepository) Create(ctx context.Context, planFood *models.PlanFood) (*models.PlanFood, error) {
	planFood.CreatedAt = time.Now()
	planFood.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, planFood)
	if err != nil {
		return nil, err
	}

	planFood.ID = result.InsertedID.(primitive.ObjectID)
	return planFood, nil
}

func (r *PlanFoodRepository) GetByID(ctx context.Context, id string) (*models.PlanFood, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid plan-food ID: %w", err)
	}

	var planFood models.PlanFood
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&planFood)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("plan-food not found")
		}
		return nil, err
	}

	return &planFood, nil
}

func (r *PlanFoodRepository) GetAllByNutritionPlanID(ctx context.Context, nutritionPlanID string) ([]models.PlanFood, error) {
	objectID, err := primitive.ObjectIDFromHex(nutritionPlanID)
	if err != nil {
		return nil,err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"nutritionPlanId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var planFoods []models.PlanFood
	if err = cursor.All(ctx, &planFoods); err != nil {
		return nil, err
	}

	return planFoods, nil
}

func (r *PlanFoodRepository) GetAll(ctx context.Context, page, limit int64) ([]models.PlanFood, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var planFoods []models.PlanFood
	if err = cursor.All(ctx, &planFoods); err != nil {
		return nil, err
	}

	return planFoods, nil
}

func (r *PlanFoodRepository) Update(ctx context.Context, planFood *models.PlanFood) (*models.PlanFood, error) {
	planFood.UpdatedAt = time.Now()

	filter := bson.M{"_id": planFood.ID}
	update := bson.M{"$set": planFood}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return planFood, nil
}

func (r *PlanFoodRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid plan-food ID: %w", err)
	}

	// No foreign key constraints assumed for PlanFood, as it likely doesn't have dependent records
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("error deleting plan-food: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("plan-food not found")
	}

	return nil
}
