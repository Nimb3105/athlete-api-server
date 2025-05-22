package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AthleteMatch struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchID primitive.ObjectID `bson:"matchId" json:"matchId"`
	UserID  primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}