package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleAdmin     = "admin"
	RoleUser      = "athlete"
	RoleModerator = "coach"
)

// User represents the base user model
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Gender      string             `bson:"gender" json:"gender"`
	ImageURL    string             `bson:"imageUrl" json:"imageUrl"`
	FullName    string             `bson:"fullName" json:"fullName"`
	Password    string             `bson:"password" json:"password"`
	Email       string             `bson:"email" json:"email"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
	DateOfBirth time.Time          `bson:"dateOfBirth" json:"dateOfBirth"`
	Role        string             `bson:"role" json:"role"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
