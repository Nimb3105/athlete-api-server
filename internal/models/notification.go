package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
	ScheduleID      primitive.ObjectID `bson:"scheduleId,omitempty" json:"scheduleId"`
	NutritionPlanID primitive.ObjectID `bson:"nutritionPlanId,omitempty" json:"nutritionPlanId"`
	SentDate        time.Time          `bson:"sentDate" json:"sentDate"`
	Status          string             `bson:"status" json:"status"` // Read, Unread
	Type            string             `bson:"type" json:"type"`
	Content         string             `bson:"content" json:"content"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
