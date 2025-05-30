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

// TeamRepository provides CRUD methods for Team
type TeamRepository struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewTeamRepository(collection *mongo.Collection,db         *mongo.Database) *TeamRepository {
	return &TeamRepository{collection:collection,db:db,}
}

func (r *TeamRepository) Create(ctx context.Context, team *models.Team) (*models.Team, error) {
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, team)
	if err != nil {
		return nil, err
	}

	team.ID = result.InsertedID.(primitive.ObjectID)
	return team, nil
}

func (r *TeamRepository) GetByID(ctx context.Context, id string) (*models.Team, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var team models.Team
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&team)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("team not found")
		}
		return nil, err
	}

	return &team, nil
}

func (r *TeamRepository) GetBySportID(ctx context.Context, sportID string) ([]models.Team, error) {
	objectID, err := primitive.ObjectIDFromHex(sportID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"sportId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var teams []models.Team
	if err = cursor.All(ctx, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *TeamRepository) GetAll(ctx context.Context, page, limit int64) ([]models.Team, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var teams []models.Team
	if err = cursor.All(ctx, &teams); err != nil {
		return nil, err
	}

	return teams, nil
}

func (r *TeamRepository) Update(ctx context.Context, team *models.Team) (*models.Team, error) {
	team.UpdatedAt = time.Now()

	filter := bson.M{"_id": team.ID}
	update := bson.M{"$set": team}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID đội không hợp lệ: %w", err)
	}

	// Kiểm tra ràng buộc khóa ngoại
	configs := []ForeignKeyCheckConfig{
		{r.db.Collection("team_members"), bson.M{"teamId": objectID}, "thành viên đội"},
	}
	if err := CheckForeignKeyConstraints(ctx, configs); err != nil {
		return err
	}

	// Xóa team
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("lỗi khi xóa đội: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("đội không tồn tại")
	}

	return nil
}