package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Exercise struct định nghĩa bài tập
type Exercise struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BodyPart         string             `bson:"bodyPart" json:"bodyPart"`
	Equipment        string             `bson:"equipment" json:"equipment"`
	Name             string             `bson:"name" json:"name"`
	Target           string             `bson:"target" json:"target"`
	SportName        string             `bson:"sportName" json:"sportName"`
	SecondaryMuscles []string           `bson:"secondaryMuscles" json:"secondaryMuscles"`
	Instructions     []string           `bson:"instructions" json:"instructions"`
	GifUrl           string             `bson:"gifUrl" json:"gifUrl"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}
