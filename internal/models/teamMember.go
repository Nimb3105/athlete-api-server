package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamMember struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TeamID    primitive.ObjectID `bson:"teamId" json:"teamId"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	DateJoined time.Time          `bson:"dateJoined" json:"dateJoined"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
