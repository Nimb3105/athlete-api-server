package internal

// import (
// 	"context"
// 	"errors"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// func (r *AthleteRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}
// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("athlete not found")
// 	}

// 	return nil
// }

// func (r *CoachRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("coach not found")
// 	}

// 	return nil
// }

// func (r *ExerciseRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("exercise not found")
// 	}

// 	return nil
// }

// func (r *GroupRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("group not found")
// 	}

// 	return nil
// }

// func (r *HealthRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("health record not found")
// 	}

// 	return nil
// }

// func (r *MatchScheduleRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("match schedule not found")
// 	}

// 	return nil
// }

// func (r *NutritionPlanRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("nutrition plan not found")
// 	}

// 	return nil
// }

// func (r *SportRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("sport not found")
// 	}

// 	return nil
// }

// func (r *TeamRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("team not found")
// 	}

// 	return nil
// }

// func (r *TournamentRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	// Nếu không có trận đấu liên quan, tiến hành xóa giải đấu
// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("giải đấu không tồn tại")
// 	}

// 	return nil
// }

// func (r *TrainingScheduleRepository) Delete(ctx context.Context, id string) error {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		return err
// 	}

// 	if result.DeletedCount == 0 {
// 		return errors.New("training schedule not found")
// 	}

// 	return nil
// }
