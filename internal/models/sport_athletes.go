package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SportAthlete struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SportID   primitive.ObjectID `bson:"sportId" json:"sportId"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
	Position  string             `bson:"position" json:"position"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
