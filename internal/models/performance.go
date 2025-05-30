package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Performance struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
	ScheduleID primitive.ObjectID `bson:"scheduleId" json:"scheduleId"`
	Value      float64            `bson:"value" json:"value"`
	MetricType string             `bson:"metricType" json:"metricType"` // Speed, Distance, Weight, etc.
	Notes      string             `bson:"notes" json:"notes"`
	Date       time.Time          `bson:"date" json:"date"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}
