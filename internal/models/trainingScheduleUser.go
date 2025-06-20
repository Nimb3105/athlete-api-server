package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingScheduleUser struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ScheduleID primitive.ObjectID `bson:"scheduleId" json:"scheduleId"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
