package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reminder struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId,oitempty" json:"userId"`
	ScheduleID      primitive.ObjectID `bson:"scheduleId,omitempty" json:"scheduleId"`
	NutritionPlanID primitive.ObjectID `bson:"nutritionPlanId,omitempty" json:"nutritionPlanId"`
	ReminderTime    time.Time          `bson:"reminderTime" json:"reminderTime"`
	ReminderDate    time.Time          `bson:"reminderDate" json:"reminderDate"`
	Content         string             `bson:"content" json:"content"`
	Status          string             `bson:"status" json:"status"` // Pending, Sent, Dismissed
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
