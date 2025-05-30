package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NutritionPlan struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"` // Name of the nutrition plan
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	CreateBy      primitive.ObjectID `bson:"createBy" json:"createBy"`
	TotalCalories int                `bson:"totalCalories" json:"totalCalories"`
	MealCount     int                `bson:"mealCount" json:"mealCount"`
	MealType      string             `bson:"mealType" json:"mealType"`       // e.g., "Breakfast", "Lunch", "Dinner", "Snack"
	MealTime      time.Time          `bson:"mealTime" json:"mealTime"`       // Time of the meal
	Description   string             `bson:"description" json:"description"` // Description of the nutrition plan
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
