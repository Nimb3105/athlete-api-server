package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainingSchedule struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date      time.Time          `bson:"date" json:"date"`
	StartTime time.Time          `bson:"startTime" json:"startTime"`
	EndTime   time.Time          `bson:"endTime" json:"endTime"`
	Status    string             `bson:"status" json:"status"` // Scheduled, Completed, Cancelled
	Location  string             `bson:"location" json:"location"`
	Type      string             `bson:"type" json:"type"`
	Notes     string             `bson:"notes" json:"notes"`
	Progress  float64            `bson:"progress" json:"progress"`   // Percentage of completion
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"` // CoachID
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
