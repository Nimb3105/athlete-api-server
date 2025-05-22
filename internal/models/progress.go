package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Progress struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
	MetricType string             `bson:"metricType" json:"metricType"`
	Value      float64            `bson:"value" json:"value"`
	Date       time.Time          `bson:"date" json:"date"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}