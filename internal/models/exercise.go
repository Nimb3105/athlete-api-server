package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exercise struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Type        string             `bson:"type" json:"type"`
	Intensity   string             `bson:"intensity" json:"intensity"`
	Duration    int                `bson:"duration" json:"duration"` // Duration in minutes
	Description string             `bson:"description" json:"description"`
	Equipment   string             `bson:"equipment" json:"equipment"`
	Muscle      string             `bson:"muscle" json:"muscle"`
	MediaURL    string             `bson:"mediaUrl" json:"mediaUrl"`
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}