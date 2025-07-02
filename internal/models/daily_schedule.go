package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailySchedule struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	Name      string             `bson:"name" json:"name"`
	Note      string             `bson:"note" json:"note"`
	SportId   primitive.ObjectID `bson:"sportId" json:"sportId"`
	StartDate time.Time          `bson:"startDate" json:"startDate"`
	EndDate   time.Time          `bson:"endDate" json:"endDate"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
