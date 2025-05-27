package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	FoodType    string             `bson:"foodType" json:"foodType"`
	FoodImage   string             `bson:"foodImage" json:"foodImage"`
	Description string             `bson:"description" json:"description"`
	Calories    int                `bson:"calories" json:"calories"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
