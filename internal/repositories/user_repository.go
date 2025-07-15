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

// UserRepository cung cấp các phương thức CRUD cho User
type UserRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewUserRepository(collection *mongo.Collection, db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: collection,
		db:         db,
	}
}

func (r *UserRepository) FindUnassignedAthletesBySport(ctx context.Context, sportId string, excludeIds []primitive.ObjectID) ([]models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(sportId)
    if err != nil {
        return nil, fmt.Errorf("invalid coach ID: %w", err)
    }

	filter := bson.M{
		"role":    "Vận động viên",
		"sportId": objectID,
	}

	if len(excludeIds) > 0 {
		filter["_id"] = bson.M{"$nin": excludeIds}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}


func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// GetByID lấy người dùng theo ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail lấy người dùng theo email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAll lấy danh sách tất cả người dùng với phân trang
func (r *UserRepository) GetAll(ctx context.Context, page, limit int64) ([]models.User,int64, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0,err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0,err
	}

	totalCount,err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil{
		return nil,0,err
	}

	return users,totalCount, nil
}

func (r *UserRepository) GetUsersByRoleWithPagination(ctx context.Context, page, limit int64, role string) ([]models.User, int64, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{"role": role}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{"role": role})
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (r *UserRepository) GetAllUserCoachBySportId(ctx context.Context, page, limit int64, sportId string) ([]models.User, int64, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	objectID, err := primitive.ObjectIDFromHex(sportId)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"sportId": objectID, "role": "Huấn luyện viên"}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// Update cập nhật thông tin người dùng
func (r *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete xóa người dùng theo ID (chỉ thực hiện xóa)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("ID người dùng không hợp lệ: %w", err)
    }

    // Chỉ xóa user, không kiểm tra khóa ngoại ở đây nữa
    result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        return fmt.Errorf("lỗi khi xóa người dùng: %w", err)
    }

    if result.DeletedCount == 0 {
        return errors.New("người dùng không tồn tại")
    }

    return nil
}
