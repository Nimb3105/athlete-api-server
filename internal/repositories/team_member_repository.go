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

// TeamMemberRepository provides CRUD methods for TeamMember
type TeamMemberRepository struct {
	collection *mongo.Collection
}

func NewTeamMemberRepository(collection *mongo.Collection) *TeamMemberRepository {
	return &TeamMemberRepository{collection}
}

func (r *TeamMemberRepository) Create(ctx context.Context, teamMember *models.TeamMember) (*models.TeamMember, error) {
	teamMember.CreatedAt = time.Now()
	teamMember.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, teamMember)
	if err != nil {
		return nil, err
	}

	teamMember.ID = result.InsertedID.(primitive.ObjectID)
	return teamMember, nil
}

func (r *TeamMemberRepository) GetByID(ctx context.Context, id string) (*models.TeamMember, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var teamMember models.TeamMember
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&teamMember)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("team member not found")
		}
		return nil, err
	}

	return &teamMember, nil
}

func (r *TeamMemberRepository) GetByTeamID(ctx context.Context, teamID string) ([]models.TeamMember, error) {
	objectID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"teamId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var teamMembers []models.TeamMember
	if err = cursor.All(ctx, &teamMembers); err != nil {
		return nil, err
	}

	return teamMembers, nil
}

func (r *TeamMemberRepository) GetByUserID(ctx context.Context, userID string) ([]models.TeamMember, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"userId": objectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var teamMembers []models.TeamMember
	if err = cursor.All(ctx, &teamMembers); err != nil {
		return nil, err
	}

	return teamMembers, nil
}

func (r *TeamMemberRepository) GetAll(ctx context.Context, page, limit int64) ([]models.TeamMember, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var teamMembers []models.TeamMember
	if err = cursor.All(ctx, &teamMembers); err != nil {
		return nil, err
	}

	return teamMembers, nil
}

func (r *TeamMemberRepository) Update(ctx context.Context, teamMember *models.TeamMember) (*models.TeamMember, error) {
	teamMember.UpdatedAt = time.Now()

	filter := bson.M{"_id": teamMember.ID}
	update := bson.M{"$set": teamMember}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return teamMember, nil
}

func (r *TeamMemberRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("team member not found")
	}

	return nil
}