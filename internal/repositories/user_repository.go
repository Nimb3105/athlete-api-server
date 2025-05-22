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
func (r *UserRepository) GetAll(ctx context.Context, page, limit int64) ([]models.User, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
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

// Delete xóa người dùng theo ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID người dùng không hợp lệ: %w", err)
	}

	// Danh sách các collection cần kiểm tra
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("achievements"), bson.M{"userId": objectID}, "thành tích"},
		{r.db.Collection("athletes"), bson.M{"userId": objectID}, "vận động viên"},
		{r.db.Collection("athlete_matches"), bson.M{"userId": objectID}, "trận đấu của vận động viên"},
		{r.db.Collection("coaches"), bson.M{"userId": objectID}, "huấn luyện viên"},
		{r.db.Collection("coach_certifications"), bson.M{"userId": objectID}, "chứng chỉ huấn luyện viên"},
		{r.db.Collection("feedbacks"), bson.M{"userId": objectID}, "phản hồi"},
		{r.db.Collection("groups"), bson.M{"createdBy": objectID}, "nhóm"},
		{r.db.Collection("group_members"), bson.M{"userId": objectID}, "thành viên nhóm"},
		{r.db.Collection("healths"), bson.M{"userId": objectID}, "sức khỏe"},
		{r.db.Collection("injuries"), bson.M{"userId": objectID}, "chấn thương"},
		{r.db.Collection("messages"), bson.M{"senderId": objectID}, "tin nhắn"},
		{r.db.Collection("notifications"), bson.M{"userId": objectID}, "thông báo"},
		{r.db.Collection("performances"), bson.M{"userId": objectID}, "hiệu suất"},
		{r.db.Collection("progresses"), bson.M{"userId": objectID}, "tiến độ"},
		{r.db.Collection("reminders"), bson.M{"userId": objectID}, "lời nhắc"},
		{r.db.Collection("sport_athletes"), bson.M{"userId": objectID}, "vận động viên môn thể thao"},
		{r.db.Collection("teams"), bson.M{"createdBy": objectID}, "đội"},
		{r.db.Collection("team_members"), bson.M{"userId": objectID}, "thành viên đội"},
		{r.db.Collection("training_schedule_users"), bson.M{"userId": objectID}, "người dùng lịch tập luyện"},
		{r.db.Collection("nutrition_plans"), bson.M{"$or": []bson.M{{"userId": objectID}, {"createby": objectID}}}, "kế hoạch dinh dưỡng"},
		{r.db.Collection("training_schedules"), bson.M{"createdBy": objectID}, "lịch tập luyện"},
	}

	// Kiểm tra ràng buộc khóa ngoại
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa user nếu không có ràng buộc
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa người dùng: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("người dùng không tồn tại")
	}

	return nil
}
