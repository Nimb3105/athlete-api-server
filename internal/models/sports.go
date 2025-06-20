package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sport struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Position string             `bson:"position" json:"position"` // List of positions in the sport
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
