package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoachCertification struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId"` // CoachID
	Name       string             `bson:"name" json:"name"`
	DateIssued time.Time          `bson:"dateIssued" json:"dateIssued"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
