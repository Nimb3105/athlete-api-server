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

// NutritionMealRepository provides CRUD methods for NutritionMeal
type FoodRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewFoodRepository(collection *mongo.Collection, db *mongo.Database) *FoodRepository {
	return &FoodRepository{collection: collection, db: db}
}

func (r *FoodRepository) Create(ctx context.Context, nutritionMeal *models.Food) (*models.Food, error) {
	nutritionMeal.CreatedAt = time.Now()
	nutritionMeal.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, nutritionMeal)
	if err != nil {
		return nil, err
	}

	nutritionMeal.ID = result.InsertedID.(primitive.ObjectID)
	return nutritionMeal, nil
}

func (r *FoodRepository) GetByID(ctx context.Context, id string) (*models.Food, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var nutritionMeal models.Food
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&nutritionMeal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("nutrition meal not found")
		}
		return nil, err
	}

	return &nutritionMeal, nil
}

func (r *FoodRepository) GetByNutritionPlanID(ctx context.Context, nutritionPlanID string) ([]models.Food, error) {
	objectID, err := primitive.ObjectIDFromHex(nutritionPlanID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"nutritionPlanId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var nutritionMeals []models.Food
	if err = cursor.All(ctx, &nutritionMeals); err != nil {
		return nil, err
	}

	return nutritionMeals, nil
}

func (r *FoodRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Food, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var nutritionMeals []models.Food
	if err = cursor.All(ctx, &nutritionMeals); err != nil {
		return nil, err
	}

	return nutritionMeals, nil
}

func (r *FoodRepository) Update(ctx context.Context, nutritionMeal *models.Food) (*models.Food, error) {
	nutritionMeal.UpdatedAt = time.Now()

	filter := bson.M{"_id": nutritionMeal.ID}
	update := bson.M{"$set": nutritionMeal}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return nutritionMeal, nil
}

func (r *FoodRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("plan_foods"), bson.M{"foodId": objectID}, "kế hoạc dinh dưỡng"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("nutrition meal not found")
	}

	return nil
}
