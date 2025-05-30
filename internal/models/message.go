package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID   primitive.ObjectID `bson:"groupId" json:"groupId"`
	SenderID  primitive.ObjectID `bson:"senderId" json:"senderId"` // UserID
	Content   string             `bson:"content" json:"content"`
	Date      time.Time          `bson:"date" json:"date"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
