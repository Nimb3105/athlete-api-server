package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NutritionPlan struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	CreateBy      primitive.ObjectID `bson:"createBy" json:"createBy"`
	StartDate     time.Time          `bson:"startDate" json:"startDate"`
	EndDate       time.Time          `bson:"endDate" json:"endDate"`
	TotalCalories int                `bson:"totalCalories" json:"totalCalories"`
	MealsPerDay   int                `bson:"mealsPerDay" json:"mealsPerDay"`
	Notes         string             `bson:"notes" json:"notes"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
