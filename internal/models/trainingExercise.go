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
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}