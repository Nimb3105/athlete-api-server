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

// SportAthleteRepository cung cấp các phương thức CRUD cho SportAthlete
type SportAthleteRepository struct {
	collection *mongo.Collection
}

// NewSportAthleteRepository tạo một SportAthleteRepository mới
func NewSportAthleteRepository(collection *mongo.Collection) *SportAthleteRepository {
	return &SportAthleteRepository{collection}
}

// Create tạo một SportAthlete mới
func (r *SportAthleteRepository) Create(ctx context.Context, SportAthlete *models.SportAthlete) (*models.SportAthlete, error) {
	SportAthlete.CreatedAt = time.Now()
	SportAthlete.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, SportAthlete)
	if err != nil {
		return nil, err
	}

	SportAthlete.ID = result.InsertedID.(primitive.ObjectID)
	return SportAthlete, nil
}

// GetByID lấy SportAthlete theo ID
func (r *SportAthleteRepository) GetByID(ctx context.Context, id string) (*models.SportAthlete, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var SportAthlete models.SportAthlete
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&SportAthlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("SportAthlete not found")
		}
		return nil, err
	}

	return &SportAthlete, nil
}

// GetByUserID lấy SportAthlete theo UserID
func (r *SportAthleteRepository) GetByUserID(ctx context.Context, userID string) (*models.SportAthlete, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var SportAthlete models.SportAthlete
	err = r.collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&SportAthlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("SportAthlete not found")
		}
		return nil, err
	}

	return &SportAthlete, nil
}

func (r *SportAthleteRepository) GetBySportID(ctx context.Context, SportID string) (*models.SportAthlete, error) {
	objectID, err := primitive.ObjectIDFromHex(SportID)
	if err != nil {
		return nil, err
	}

	var SportAthlete models.SportAthlete
	err = r.collection.FindOne(ctx, bson.M{"sportId": objectID}).Decode(&SportAthlete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("SportAthlete not found")
		}
		return nil, err
	}

	return &SportAthlete, nil
}

// GetAll lấy danh sách tất cả SportAthlete với phân trang
func (r *SportAthleteRepository) GetAll(ctx context.Context, page, limit int64) ([]models.SportAthlete, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var SportAthletes []models.SportAthlete
	if err = cursor.All(ctx, &SportAthletes); err != nil {
		return nil, err
	}

	return SportAthletes, nil
}

// Update cập nhật thông tin SportAthlete
func (r *SportAthleteRepository) Update(ctx context.Context, SportAthlete *models.SportAthlete) (*models.SportAthlete, error) {
	SportAthlete.UpdatedAt = time.Now()

	filter := bson.M{"_id": SportAthlete.ID}
	update := bson.M{"$set": SportAthlete}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return SportAthlete, nil
}

// Delete xóa SportAthlete theo ID
func (r *SportAthleteRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("SportAthlete not found")
	}

	return nil
}

// GetAllByUserID lấy danh sách tất cả SportAthlete theo userID
func (r *SportAthleteRepository) GetAllByUserID(ctx context.Context, userID string) ([]models.SportAthlete, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var SportAthletes []models.SportAthlete
	if err = cursor.All(ctx, &SportAthletes); err != nil {
		return nil, err
	}

	return SportAthletes, nil
}
