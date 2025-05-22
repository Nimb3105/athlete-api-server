package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MedicalHistory struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HealthID    primitive.ObjectID `bson:"healthId" json:"healthId"`
	Date        time.Time          `bson:"date" json:"date"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
