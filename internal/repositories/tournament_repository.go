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

// TournamentRepository provides CRUD methods for Tournament
type TournamentRepository struct {
	collection              *mongo.Collection
	db         *mongo.Database
}

func NewTournamentRepository(collection *mongo.Collection, db         *mongo.Database) *TournamentRepository {
	return &TournamentRepository{
		collection:collection,db:db,
	}
}

func (r *TournamentRepository) Create(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error) {
	tournament.CreatedAt = time.Now()
	tournament.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, tournament)
	if err != nil {
		return nil, err
	}

	tournament.ID = result.InsertedID.(primitive.ObjectID)
	return tournament, nil
}

func (r *TournamentRepository) GetByID(ctx context.Context, id string) (*models.Tournament, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var tournament models.Tournament
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&tournament)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("tournament not found")
		}
		return nil, err
	}

	return &tournament, nil
}

func (r *TournamentRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Tournament, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tournaments []models.Tournament
	if err = cursor.All(ctx, &tournaments); err != nil {
		return nil, err
	}

	return tournaments, nil
}

func (r *TournamentRepository) Update(ctx context.Context, tournament *models.Tournament) (*models.Tournament, error) {
	tournament.UpdatedAt = time.Now()

	filter := bson.M{"_id": tournament.ID}
	update := bson.M{"$set": tournament}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return tournament, nil
}

func (r *TournamentRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID giải đấu không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("match_schedules"), bson.M{"tournamentId": objectID}, "lịch thi đấu"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa tournament
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa giải đấu: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("giải đấu không tồn tại")
	}

	return nil
}
