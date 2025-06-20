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

type TrainingScheduleRepository struct {
	collection       *mongo.Collection
	scheduleUserRepo *TrainingScheduleUserRepository // Assuming you have this repository for TrainingScheduleUser
	db               *mongo.Database
}

func NewTrainingScheduleRepository(collection *mongo.Collection, scheduleUserRepo *TrainingScheduleUserRepository, db *mongo.Database) *TrainingScheduleRepository {
	return &TrainingScheduleRepository{collection: collection, scheduleUserRepo: scheduleUserRepo, db: db}
}

func (r *TrainingScheduleRepository) Create(ctx context.Context, schedule *models.TrainingSchedule, athleteId string) (*models.TrainingSchedule, error) {
	// Set timestamps for the TrainingSchedule
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	// Bước 1: Convert athleteId to primitive.ObjectID
	athleteObjectID, err := primitive.ObjectIDFromHex(athleteId)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi chuyển athleteId sang ObjectID: %v", err)
	}

	// Bước 2: Tìm tất cả TrainingScheduleUser bằng athleteId
	filterUser := bson.M{
		"userId": athleteObjectID,
	}
	var scheduleUsers []models.TrainingScheduleUser
	cursor, err := r.scheduleUserRepo.collection.Find(ctx, filterUser)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn TrainingScheduleUser: %v", err)
	}
	if err := cursor.All(ctx, &scheduleUsers); err != nil {
		return nil, fmt.Errorf("lỗi khi giải mã TrainingScheduleUser: %v", err)
	}

	// Bước 3: Lấy tất cả scheduleId từ TrainingScheduleUser
	scheduleIDs := make([]primitive.ObjectID, 0, len(scheduleUsers))
	for _, su := range scheduleUsers {
		scheduleIDs = append(scheduleIDs, su.ScheduleID)
	}

	// Bước 4: Tìm tất cả TrainingSchedule trong cùng ngày dựa trên Date và scheduleIDs
	startOfDay := time.Date(schedule.Date.Year(), schedule.Date.Month(), schedule.Date.Day(), 0, 0, 0, 0, schedule.Date.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	filterSchedule := bson.M{
		"_id": bson.M{
			"$in": scheduleIDs,
		},
		"date": bson.M{
			"$gte": startOfDay,
			"$lte": endOfDay,
		},
	}
	var existingSchedules []models.TrainingSchedule
	cursor, err = r.collection.Find(ctx, filterSchedule)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi truy vấn lịch trình cùng ngày: %v", err)
	}
	if err := cursor.All(ctx, &existingSchedules); err != nil {
		return nil, fmt.Errorf("lỗi khi giải mã lịch trình: %v", err)
	}

	// Bước 5: Kiểm tra điều kiện giờ bắt đầu
	for _, existing := range existingSchedules {
		// So sánh chỉ giờ, phút, giây trong ngày
		newStartTime := schedule.StartTime
		existingEndTime := existing.EndTime

		// Nếu StartTime và EndTime chứa ngày khác, chỉ lấy giờ/phút/giây
		newStartHour, newStartMin, newStartSec := newStartTime.Hour(), newStartTime.Minute(), newStartTime.Second()
		existingEndHour, existingEndMin, existingEndSec := existingEndTime.Hour(), existingEndTime.Minute(), existingEndTime.Second()

		newStartTotalSec := newStartHour*3600 + newStartMin*60 + newStartSec
		existingEndTotalSec := existingEndHour*3600 + existingEndMin*60 + existingEndSec

		if newStartTotalSec <= existingEndTotalSec {
			return nil, fmt.Errorf("giờ bắt đầu (%v) phải sau giờ kết thúc (%v) của lịch trình ID %v",
				newStartTime.Format("15:04:05"), existingEndTime.Format("15:04:05"), existing.ID.Hex())
		}
	}

	// Bước 6: Chèn TrainingSchedule vào collection
	result, err := r.collection.InsertOne(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi chèn lịch trình: %v", err)
	}

	// Set ID cho TrainingSchedule
	schedule.ID = result.InsertedID.(primitive.ObjectID)

	// Bước 7: Tạo và chèn TrainingScheduleUser
	scheduleUser := &models.TrainingScheduleUser{
		ScheduleID: schedule.ID,
		UserID:     athleteObjectID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = r.scheduleUserRepo.Create(ctx, scheduleUser)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi chèn TrainingScheduleUser: %v", err)
	}

	return schedule, nil
}

func (r *TrainingScheduleRepository) GetByID(ctx context.Context, id string) (*models.TrainingSchedule, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var schedule models.TrainingSchedule
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("training schedule not found")
		}
		return nil, err
	}

	return &schedule, nil
}

func (r *TrainingScheduleRepository) GetAll(ctx context.Context, page, limit int64) ([]models.TrainingSchedule, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []models.TrainingSchedule
	if err = cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *TrainingScheduleRepository) Update(ctx context.Context, schedule *models.TrainingSchedule) (*models.TrainingSchedule, error) {
	schedule.UpdatedAt = time.Now()

	filter := bson.M{"_id": schedule.ID}
	update := bson.M{"$set": schedule}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

func (r *TrainingScheduleRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID lịch tập luyện không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("feedbacks"), bson.M{"scheduleId": objectID}, "phản hồi"},
		{r.db.Collection("performances"), bson.M{"scheduleId": objectID}, "hiệu suất"},
		{r.db.Collection("notifications"), bson.M{"scheduleId": objectID}, "thông báo"},
		{r.db.Collection("reminders"), bson.M{"scheduleId": objectID}, "lời nhắc"},
		{r.db.Collection("training_exercises"), bson.M{"scheduleId": objectID}, "bài tập trong lịch tập luyện"},
		{r.db.Collection("training_schedule_users"), bson.M{"scheduleId": objectID}, "người dùng lịch tập luyện"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa training schedule
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa lịch tập luyện: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("lịch tập luyện không tồn tại")
	}

	return nil
}

// training_schedule_repository.go
func (r *TrainingScheduleRepository) MarkOverdue(ctx context.Context, now time.Time) (int64, error) {
	filter := bson.M{
		"endTime": bson.M{"$lt": now},
		"status":  bson.M{"$ne": "hoàn thành"},
	}
	update := bson.M{
		"$set": bson.M{
			"status":    "chưa hoàn thành",
			"updatedAt": now,
		},
	}

	res, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
