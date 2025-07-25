package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanFood struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FoodID          primitive.ObjectID `bson:"foodId" json:"foodId"`                   // ID of the food item
	NutritionPlanID primitive.ObjectID `bson:"nutritionPlanId" json:"nutritionPlanId"` // ID of the nutrition plan              // e.g., Breakfast, Lunch, Dinner
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
