package repositories

// import (
// 	"be/internal/models"
// 	"context"
// 	"errors"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // PerformanceRepository provides CRUD methods for Performance
// type PerformanceRepository struct {
// 	collection *mongo.Collection
// }

// func NewPerformanceRepository(collection *mongo.Collection) *PerformanceRepository {
// 	return &PerformanceRepository{collection}
// }

// func (r *PerformanceRepository) Create(ctx context.Context, performance *models.Performance) (*models.Performance, error) {
// 	performance.CreatedAt = time.Now()
// 	performance.UpdatedAt = time.Now()

// 	result, err := r.collection.InsertOne(ctx, performance)
// 	if err != nil {
// 		return nil, err
// 	}

// 	performance.ID = result.InsertedID.(primitive.ObjectID)
// 	return performance, nil
// }

// func (r *PerformanceRepository) GetByID(ctx context.Context, id string) (*models.Performance, error) {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var performance models.Performance
// 	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&performance)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, errors.New("performance not found")
// 		}
// 		return nil, err
// 	}

// 	return &performance, nil
// }

// func (r *PerformanceRepository) GetByUserID(ctx context.Context, userID string) ([]models.Performance, error) {
// 	objectID, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var performances []models.Performance
// 	if err = cursor.All(ctx, &performances); err != nil {
// 		return nil, err
// 	}

// 	return performances, nil
// }

// func (r *PerformanceRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Performance, error) {
// 	opts := options.Find()
// 	opts.SetSkip((page - 1) * limit)
// 	opts.SetLimit(limit)
// 	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

// 	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var performances []models.Performance
// 	if err = cursor.All(ctx, &performances); err != nil {
// 		return nil, err
// 	}

// 	return performances, nil
// }

// func (r *PerformanceRepository) Update(ctx context.Context, performance *models.Performance) (*models.Performance, error) {
// 	performance.UpdatedAt = time.Now()

// 	filter := bson.M{"_id": performance.ID}
// 	update := bson.M{"$set": performance}

// 	_, err := r.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return performance, nil
// }

// func (r *PerformanceRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("performance not found")
// 	}

// 	return nil
// }