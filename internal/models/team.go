package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	SportID     primitive.ObjectID `bson:"sportId" json:"sportId"`
	Description string             `bson:"description" json:"description"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"` // UserID
	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}