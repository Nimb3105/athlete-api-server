package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NutritionMeal struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NutritionPlanID primitive.ObjectID `bson:"nutritionPlanId" json:"nutritionPlanId"`
	MealTime       time.Time          `bson:"mealTime" json:"mealTime"`
	MealType       string             `bson:"mealType" json:"mealType"` // Breakfast, Lunch, Dinner, Snack
	Description    string             `bson:"description" json:"description"`
	Calories       int                `bson:"calories" json:"calories"`
	Notes          string             `bson:"notes" json:"notes"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}