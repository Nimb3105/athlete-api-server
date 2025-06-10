package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Health struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
	Height    float64            `bson:"height" json:"height"` // Height in cm
	Weight    float64            `bson:"weight" json:"weight"` // Weight in kg
	BMI       float64            `bson:"bmi" json:"bmi"`
	BloodType string             `bson:"bloodType" json:"bloodType"`
	Date      time.Time          `bson:"date" json:"date"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
