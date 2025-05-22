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

// GroupMemberRepository provides CRUD methods for GroupMember
type GroupMemberRepository struct {
	collection *mongo.Collection
}

func NewGroupMemberRepository(collection *mongo.Collection) *GroupMemberRepository {
	return &GroupMemberRepository{collection}
}

func (r *GroupMemberRepository) Create(ctx context.Context, groupMember *models.GroupMember) (*models.GroupMember, error) {
	groupMember.CreatedAt = time.Now()
	groupMember.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, groupMember)
	if err != nil {
		return nil, err
	}

	groupMember.ID = result.InsertedID.(primitive.ObjectID)
	return groupMember, nil
}

func (r *GroupMemberRepository) GetByID(ctx context.Context, id string) (*models.GroupMember, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var groupMember models.GroupMember
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&groupMember)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("group member not found")
		}
		return nil, err
	}

	return &groupMember, nil
}

func (r *GroupMemberRepository) GetByUserID(ctx context.Context, userID string) (*models.GroupMember, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var groupMember models.GroupMember
	err = r.collection.FindOne(ctx, bson.M{"userId": objectID}).Decode(&groupMember)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("group member not found")
		}
		return nil, err
	}

	return &groupMember, nil
}

func (r *GroupMemberRepository) GetAll(ctx context.Context, page, limit int64) ([]models.GroupMember, error) {
	opts := options.Find()
	opts.SetSkip((page - 1) * limit)
	opts.SetLimit(limit)
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var groupMembers []models.GroupMember
	if err = cursor.All(ctx, &groupMembers); err != nil {
		return nil, err
	}

	return groupMembers, nil
}

func (r *GroupMemberRepository) Update(ctx context.Context, groupMember *models.GroupMember) (*models.GroupMember, error) {
	groupMember.UpdatedAt = time.Now()

	filter := bson.M{"_id": groupMember.ID}
	update := bson.M{"$set": groupMember}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return groupMember, nil
}

func (r *GroupMemberRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("group member not found")
	}

	return nil
}