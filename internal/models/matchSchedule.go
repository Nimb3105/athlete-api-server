package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchSchedule struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TournamentID primitive.ObjectID `bson:"tournamentId" json:"tournamentId"`
	Date         time.Time          `bson:"date" json:"date"`
	Location     string             `bson:"location" json:"location"`
	Opponent     string             `bson:"opponent" json:"opponent"`
	MatchType    string             `bson:"matchType" json:"matchType"`
	Status       string             `bson:"status" json:"status"` // Scheduled, Completed, Cancelled
	Round        string             `bson:"round" json:"round"`
	Score        string             `bson:"score" json:"score"`
	Notes        string             `bson:"notes" json:"notes"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}
