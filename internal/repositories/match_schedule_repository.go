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

// MatchScheduleRepository provides CRUD methods for MatchSchedule
type MatchScheduleRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewMatchScheduleRepository(collection *mongo.Collection,db         *mongo.Database) *MatchScheduleRepository {
	return &MatchScheduleRepository{collection:collection,db:db,}
}

func (r *MatchScheduleRepository) Create(ctx context.Context, matchSchedule *models.MatchSchedule) (*models.MatchSchedule, error) {
	matchSchedule.CreatedAt = time.Now()
	matchSchedule.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, matchSchedule)
	if err != nil {
		return nil, err
	}

	matchSchedule.ID = result.InsertedID.(primitive.ObjectID)
	return matchSchedule, nil
}

func (r *MatchScheduleRepository) GetByID(ctx context.Context, id string) (*models.MatchSchedule, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var matchSchedule models.MatchSchedule
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&matchSchedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("match schedule not found")
		}
		return nil, err
	}

	return &matchSchedule, nil
}

func (r *MatchScheduleRepository) GetByTournamentID(ctx context.Context, tournamentID string) (*models.MatchSchedule, error) {
	objectID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return nil, err
	}

	var matchSchedule models.MatchSchedule
	err = r.collection.FindOne(ctx, bson.M{"tournamentId": objectID}).Decode(&matchSchedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("match schedule not found")
		}
		return nil, err
	}

	return &matchSchedule, nil
}

func (r *MatchScheduleRepository) GetAll(ctx context.Context, page, limit int64) ([]models.MatchSchedule, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matchSchedules []models.MatchSchedule
	if err = cursor.All(ctx, &matchSchedules); err != nil {
		return nil, err
	}

	return matchSchedules, nil
}

func (r *MatchScheduleRepository) Update(ctx context.Context, matchSchedule *models.MatchSchedule) (*models.MatchSchedule, error) {
	matchSchedule.UpdatedAt = time.Now()

	filter := bson.M{"_id": matchSchedule.ID}
	update := bson.M{"$set": matchSchedule}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return matchSchedule, nil
}

func (r *MatchScheduleRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID lịch thi đấu không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("athlete_matches"), bson.M{"matchId": objectID}, "trận đấu của vận động viên"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa match schedule
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa lịch thi đấu: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("lịch thi đấu không tồn tại")
	}

	return nil
}