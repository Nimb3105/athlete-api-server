package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingExercise struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ScheduleID primitive.ObjectID `bson:"scheduleId" json:"scheduleId"`
	ExerciseID primitive.ObjectID `bson:"exerciseId" json:"exerciseId"`
	Order      int                `bson:"order" json:"order"`
	Reps       int                `bson:"reps" json:"reps"`
	Sets       int                `bson:"sets" json:"sets"`
	Weight     float64            `bson:"weight" json:"weight"`
	Duration   int                `bson:"duration" json:"duration"` // Duration in seconds
	Distance   float64            `bson:"distance" json:"distance"` // Distance in meters
	ActualReps  int                `bson:"actualReps" json:"actualReps"`
	ActualSets  int                `bson:"actualSets" json:"actualSets"`
	ActualWeight float64           `bson:"actualWeight" json:"actualWeight"`
	ActualDuration int           `bson:"actualDuration" json:"actualDuration"` // Actual duration in seconds
	ActualDistance float64        `bson:"actualDistance" json:"actualDistance"` // Actual distance in meters
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}