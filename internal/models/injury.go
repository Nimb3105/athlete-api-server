package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Injury struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
	Type            string             `bson:"type" json:"type"`
	Date            time.Time          `bson:"date" json:"date"`
	Severity        string             `bson:"severity" json:"severity"` // Mild, Moderate, Severe
	LocationOnBody  string             `bson:"locationOnBody" json:"locationOnBody"`
	CauseOfInjury   string             `bson:"causeOfInjury" json:"causeOfInjury"`
	RecoveryStatus  string             `bson:"recoveryStatus" json:"recoveryStatus"` // Recovering, Recovered
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}