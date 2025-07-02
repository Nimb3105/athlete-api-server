package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTrainingScheduleRequest struct {
	TrainingSchedule
	TrainingExercise []*TrainingExercise `json:"trainingExercises" validate:"required,dive"`
}

type TrainingSchedule struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	DailyScheduleId primitive.ObjectID `bson:"dailyScheduleId" json:"dailyScheduleId"`
	Date            time.Time           `bson:"date" json:"date"`
	StartTime       time.Time           `bson:"startTime" json:"startTime"`
	EndTime         time.Time           `bson:"endTime" json:"endTime"`
	Status          string              `bson:"status" json:"status"` // Scheduled, Completed, Cancelled
	Location        string              `bson:"location" json:"location"`
	Type            string              `bson:"type" json:"type"`
	Notes           string              `bson:"notes" json:"notes"`
	Progress        float64             `bson:"progress" json:"progress"`   // Percentage of completion
	CreatedBy       primitive.ObjectID  `bson:"createdBy" json:"createdBy"` // CoachID
	SportId         primitive.ObjectID  `bson:"sportId" json:"sportId"`
	CreatedAt       time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time           `bson:"updatedAt" json:"updatedAt"`
}
