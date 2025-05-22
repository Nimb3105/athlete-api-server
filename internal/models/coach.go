package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coach struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Experience     string             `bson:"experience" json:"experience"` // Years of experience
	Specialization string             `bson:"specialization" json:"specialization"`
	Level          string             `bson:"level" json:"level"` // Junior, Senior, Professional, etc.
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}
