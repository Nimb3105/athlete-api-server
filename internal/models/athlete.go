package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Athlete struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	AthleteType string             `bson:"athleteType" json:"athleteType"` // Professional, Amateur, Collegiate, etc.
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
