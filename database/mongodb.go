package database

import (
	"be/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB chứa client và các collection
type MongoDB struct {
	Client                         *mongo.Client
	Database                       *mongo.Database
	UserCollection                 *mongo.Collection
	CoachCollection                *mongo.Collection
	AthleteCollection              *mongo.Collection
	SportUserCollection            *mongo.Collection
	SportCollection                *mongo.Collection
	ExerciseCollection             *mongo.Collection
	TrainingScheduleCollection     *mongo.Collection
	TrainingExerciseCollection     *mongo.Collection
	TrainingScheduleUserCollection *mongo.Collection
	NotificationCollection         *mongo.Collection
	ReminderCollection             *mongo.Collection
	AchivementCollection           *mongo.Collection
	UserMatchCollection            *mongo.Collection
	CoachCertificationCollection   *mongo.Collection
	FeedbackCollection             *mongo.Collection
	GroupCollection                *mongo.Collection
	GroupMemberCollection          *mongo.Collection
	HealthCollection               *mongo.Collection
	InjuryCollection               *mongo.Collection
	MatchScheduleCollection        *mongo.Collection
	MedicalHistoryCollection       *mongo.Collection
	MessageCollection              *mongo.Collection
	NutritionPlanCollection        *mongo.Collection
	FoodCollection                 *mongo.Collection
	PerformanceCollection          *mongo.Collection
	ProgressCollection             *mongo.Collection
	TournamentCollection           *mongo.Collection
	TeamMemberCollection           *mongo.Collection
	TeamCollection                 *mongo.Collection
	PlanFoodCollection             *mongo.Collection
	CoachAthleteCollection         *mongo.Collection
	DailyScheduleCollection        *mongo.Collection
}

// ConnectMongoDB khởi tạo kết nối tới MongoDB
func ConnectMongoDB(config *config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Kiểm tra kết nối
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(config.MongoDatabase)

	return &MongoDB{
		Client:                         client,
		Database:                       db,
		UserCollection:                 db.Collection("users"),
		CoachCollection:                db.Collection("coaches"),
		AthleteCollection:              db.Collection("athletes"),
		SportUserCollection:            db.Collection("sport_users"),
		SportCollection:                db.Collection("sports"),
		ExerciseCollection:             db.Collection("exercises"),
		TrainingScheduleCollection:     db.Collection("training_schedules"),
		TrainingScheduleUserCollection: db.Collection("training_schedule_users"),
		TrainingExerciseCollection:     db.Collection("training_exercises"),
		NotificationCollection:         db.Collection("notifications"),
		ReminderCollection:             db.Collection("reminders"),
		AchivementCollection:           db.Collection("achivements"),
		UserMatchCollection:            db.Collection("user_matches"),
		CoachCertificationCollection:   db.Collection("coach_certifications"),
		FeedbackCollection:             db.Collection("feedbacks"),
		GroupCollection:                db.Collection("groups"),
		GroupMemberCollection:          db.Collection("group_members"),
		HealthCollection:               db.Collection("healths"),
		InjuryCollection:               db.Collection("injuries"),
		MatchScheduleCollection:        db.Collection("match_schedules"),
		MedicalHistoryCollection:       db.Collection("medical_histories"),
		MessageCollection:              db.Collection("messages"),
		NutritionPlanCollection:        db.Collection("nutrition_plans"),
		FoodCollection:                 db.Collection("foods"),
		PerformanceCollection:          db.Collection("performances"),
		ProgressCollection:             db.Collection("progresses"),
		TournamentCollection:           db.Collection("tournaments"),
		TeamMemberCollection:           db.Collection("team_members"),
		TeamCollection:                 db.Collection("teams"),
		PlanFoodCollection:             db.Collection("plan_foods"),
		CoachAthleteCollection:         db.Collection("coach_athletes"),
		DailyScheduleCollection:        db.Collection("daily_schedules"),
	}, nil
}

// Close đóng kết nối MongoDB
func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
