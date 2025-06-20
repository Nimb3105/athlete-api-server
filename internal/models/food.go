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
	Protein     int                `bson:"protein" json:"protein"`
	Carbs       int                `bson:"carbs" json:"carbs"`
	Fat         int                `bson:"fat" json:"fat"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
