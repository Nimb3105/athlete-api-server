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
	collection        *mongo.Collection
	dailyScheduleRepo *DailyScheduleRepository // Assuming you have this repository for TrainingScheduleUser
	db                *mongo.Database
}

func NewTrainingScheduleRepository(collection *mongo.Collection, dailyScheduleRepo *DailyScheduleRepository, db *mongo.Database) *TrainingScheduleRepository {
	return &TrainingScheduleRepository{collection: collection, dailyScheduleRepo: dailyScheduleRepo, db: db}
}

func (r *TrainingScheduleRepository) Create(ctx context.Context, schedule *models.TrainingSchedule) (*models.TrainingSchedule, error) {
	// Set timestamps for the TrainingSchedule
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	// Bước 4: Tìm tất cả TrainingSchedule trong cùng ngày dựa trên Date và scheduleIDs
	startOfDay := time.Date(schedule.Date.Year(), schedule.Date.Month(), schedule.Date.Day(), 0, 0, 0, 0, schedule.Date.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	filterSchedule := bson.M{
		"dailyScheduleId": schedule.DailyScheduleId,
		"date": bson.M{
			"$gte": startOfDay,
			"$lte": endOfDay,
		},
	}

	cursor, err := r.collection.Find(ctx, filterSchedule)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var existingSchedules []models.TrainingSchedule = make([]models.TrainingSchedule, 0)
	if err = cursor.All(ctx, &existingSchedules); err != nil {
		return nil, err
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

func (r *TrainingScheduleRepository) GetAllByDailyScheduleId(ctx context.Context, dailyScheduleId string, date string) ([]models.TrainingSchedule, error) {
	// Filter by sportName field
	objectId, err := primitive.ObjectIDFromHex(dailyScheduleId)
	if err != nil {
		return nil, fmt.Errorf("invalid dailyScheduleId: %w", err)
	}
	parsedTime, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		return nil, fmt.Errorf("invalid day format (expected RFC3339): %w", err)
	}

	startOfDay := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, parsedTime.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := bson.M{
		"dailyScheduleId": objectId,
		"date": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}
	// Set sort by createdAt descending
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trainingSchedules []models.TrainingSchedule
	if err = cursor.All(ctx, &trainingSchedules); err != nil {
		return nil, err
	}

	return trainingSchedules, nil
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

	updateFiels := bson.M{
		"dailyScheduleId": schedule.DailyScheduleId,
		"date":            schedule.Date,
		"startTime":       schedule.StartTime,
		"endTime":         schedule.EndTime,
		"status":          schedule.Status,
		"location":        schedule.Location,
		"type":            schedule.Type,
		"notes":           schedule.Notes,
		"progress":        schedule.Progress,
		"createdBy":       schedule.CreatedBy,
		"sportId":         schedule.SportId,
		"createdAt":       schedule.CreatedAt,
		"updatedAt":       time.Now(),
	}

	filter := bson.M{"_id": schedule.ID}
	update := bson.M{"$set": updateFiels}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	schedule.UpdatedAt = updateFiels["updatedAt"].(time.Time)
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

// MarkOverdue đánh dấu các lịch tập đã quá hạn là "chưa hoàn thành" và trả về danh sách ID đã cập nhật
func (r *TrainingScheduleRepository) MarkOverdue(ctx context.Context, now time.Time) ([]primitive.ObjectID, int64, error) {
	filter := bson.M{
		"endTime": bson.M{"$lt": now},
		"status":  bson.M{"$ne": "hoàn thành"},
	}

	// Tìm các lịch tập thỏa mãn điều kiện
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var scheduleIDs []primitive.ObjectID
	for cursor.Next(ctx) {
		var result struct {
			ID primitive.ObjectID `bson:"_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, 0, err
		}
		scheduleIDs = append(scheduleIDs, result.ID)
	}

	if len(scheduleIDs) == 0 {
		return nil, 0, nil
	}

	// Cập nhật trạng thái
	update := bson.M{
		"$set": bson.M{
			"status":    "chưa hoàn thành",
			"updatedAt": now,
		},
	}
	updateFilter := bson.M{"_id": bson.M{"$in": scheduleIDs}}

	res, err := r.collection.UpdateMany(ctx, updateFilter, update)
	if err != nil {
		return nil, 0, err
	}
	return scheduleIDs, res.ModifiedCount, nil
}
