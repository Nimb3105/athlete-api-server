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

// NutritionPlanRepository provides CRUD methods for NutritionPlan
type NutritionPlanRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewNutritionPlanRepository(collection *mongo.Collection,db         *mongo.Database) *NutritionPlanRepository {
	return &NutritionPlanRepository{collection:collection,db:db,}
}

func (r *NutritionPlanRepository) Create(ctx context.Context, nutritionPlan *models.NutritionPlan) (*models.NutritionPlan, error) {
	nutritionPlan.CreatedAt = time.Now()
	nutritionPlan.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, nutritionPlan)
	if err != nil {
		return nil, err
	}

	nutritionPlan.ID = result.InsertedID.(primitive.ObjectID)
	return nutritionPlan, nil
}

func (r *NutritionPlanRepository) GetByID(ctx context.Context, id string) (*models.NutritionPlan, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var nutritionPlan models.NutritionPlan
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&nutritionPlan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("nutrition plan not found")
		}
		return nil, err
	}

	return &nutritionPlan, nil
}

func (r *NutritionPlanRepository) GetByAthleteID(ctx context.Context, athleteID string) ([]models.NutritionPlan, error) {
	objectID, err := primitive.ObjectIDFromHex(athleteID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"athleteId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var nutritionPlans []models.NutritionPlan
	if err = cursor.All(ctx, &nutritionPlans); err != nil {
		return nil, err
	}

	return nutritionPlans, nil
}

func (r *NutritionPlanRepository) GetAll(ctx context.Context, page, limit int64) ([]models.NutritionPlan, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var nutritionPlans []models.NutritionPlan
	if err = cursor.All(ctx, &nutritionPlans); err != nil {
		return nil, err
	}

	return nutritionPlans, nil
}

func (r *NutritionPlanRepository) Update(ctx context.Context, nutritionPlan *models.NutritionPlan) (*models.NutritionPlan, error) {
	nutritionPlan.UpdatedAt = time.Now()

	filter := bson.M{"_id": nutritionPlan.ID}
	update := bson.M{"$set": nutritionPlan}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return nutritionPlan, nil
}

func (r *NutritionPlanRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID kế hoạch dinh dưỡng không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("nutrition_meals"), bson.M{"nutritionPlanId": objectID}, "bữa ăn dinh dưỡng"},
		{r.db.Collection("notifications"), bson.M{"nutritionPlanId": objectID}, "thông báo"},
		{r.db.Collection("reminders"), bson.M{"nutritionPlanId": objectID}, "lời nhắc"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa nutrition plan
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa kế hoạch dinh dưỡng: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("kế hoạch dinh dưỡng không tồn tại")
	}

	return nil
}