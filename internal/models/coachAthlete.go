package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoachAthlete struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CoachID   primitive.ObjectID `bson:"coachId" json:"coachId"`
	AthleteID primitive.ObjectID `bson:"athleteId" json:"athleteId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
