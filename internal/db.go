package internal

// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Achievement struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID      primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
// 	Title       string             `bson:"title" json:"title"`
// 	Description string             `bson:"description" json:"description"`
// 	Date        time.Time          `bson:"date" json:"date"`
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Athlete struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
// 	AthleteType string             `bson:"athleteType" json:"athleteType"` // Professional, Amateur, Collegiate, etc.
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type AthleteMatch struct {
// 	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	MatchID primitive.ObjectID `bson:"matchId" json:"matchId"`
// 	UserID  primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Coach struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
// 	Experience     string             `bson:"experience" json:"experience"` // Years of experience
// 	Specialization string             `bson:"specialization" json:"specialization"`
// 	Level          string             `bson:"level" json:"level"` // Junior, Senior, Professional, etc.
// 	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type CoachCertification struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID     primitive.ObjectID `bson:"userId" json:"userId"` // CoachID
// 	Name       string             `bson:"name" json:"name"`
// 	DateIssued time.Time          `bson:"dateIssued" json:"dateIssued"`
// 	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Exercise struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Name        string             `bson:"name" json:"name"`
// 	Type        string             `bson:"type" json:"type"`
// 	Intensity   string             `bson:"intensity" json:"intensity"`
// 	Duration    int                `bson:"duration" json:"duration"` // Duration in minutes
// 	Description string             `bson:"description" json:"description"`
// 	Equipment   string             `bson:"equipment" json:"equipment"`
// 	Muscle      string             `bson:"muscle" json:"muscle"`
// 	MediaURL    string             `bson:"mediaUrl" json:"mediaUrl"`
// 	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Feedback struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
// 	ScheduleID primitive.ObjectID `bson:"scheduleId" json:"scheduleId"`
// 	Content    string             `bson:"content" json:"content"`
// 	Date       time.Time          `bson:"date" json:"date"`
// 	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Group struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Name        string             `bson:"name" json:"name"`
// 	Description string             `bson:"description" json:"description"`
// 	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"` // UserID
// 	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type GroupMember struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	GroupID   primitive.ObjectID `bson:"groupId" json:"groupId"`
// 	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Health struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID    primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
// 	Height    float64            `bson:"height" json:"height"` // Height in cm
// 	Weight    float64            `bson:"weight" json:"weight"` // Weight in kg
// 	BMI       float64            `bson:"bmi" json:"bmi"`
// 	BloodType string             `bson:"bloodType" json:"bloodType"`
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Injury struct {
// 	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID          primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
// 	Type            string             `bson:"type" json:"type"`
// 	Date            time.Time          `bson:"date" json:"date"`
// 	Severity        string             `bson:"severity" json:"severity"` // Mild, Moderate, Severe
// 	LocationOnBody  string             `bson:"locationOnBody" json:"locationOnBody"`
// 	CauseOfInjury   string             `bson:"causeOfInjury" json:"causeOfInjury"`
// 	RecoveryStatus  string             `bson:"recoveryStatus" json:"recoveryStatus"` // Recovering, Recovered
// 	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type MatchSchedule struct {
// 	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	TournamentID primitive.ObjectID `bson:"tournamentId" json:"tournamentId"`
// 	Date         time.Time          `bson:"date" json:"date"`
// 	Time         time.Time          `bson:"time" json:"time"`
// 	Location     string             `bson:"location" json:"location"`
// 	Opponent     string             `bson:"opponent" json:"opponent"`
// 	MatchType    string             `bson:"matchType" json:"matchType"`
// 	Status       string             `bson:"status" json:"status"` // Scheduled, Completed, Cancelled
// 	Round        string             `bson:"round" json:"round"`
// 	Score        string             `bson:"score" json:"score"`
// 	Notes        string             `bson:"notes" json:"notes"`
// 	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type MedicalHistory struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	HealthID    primitive.ObjectID `bson:"healthId" json:"healthId"`
// 	Date        time.Time          `bson:"date" json:"date"`
// 	Description string             `bson:"description" json:"description"`
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Message struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	GroupID   primitive.ObjectID `bson:"groupId" json:"groupId"`
// 	SenderID  primitive.ObjectID `bson:"senderId" json:"senderId"` // UserID
// 	SentDate  time.Time          `bson:"sentDate" json:"sentDate"`
// 	Content   string             `bson:"content" json:"content"`
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Notification struct {
// 	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
// 	ScheduleID      primitive.ObjectID `bson:"scheduleId,omitempty" json:"scheduleId"`
// 	NutritionPlanID primitive.ObjectID `bson:"nutritionPlanId,omitempty" json:"nutritionPlanId"`
// 	SentDate        time.Time          `bson:"sentDate" json:"sentDate"`
// 	Status          string             `bson:"status" json:"status"` // Read, Unread
// 	Type            string             `bson:"type" json:"type"`
// 	Content         string             `bson:"content" json:"content"`
// 	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type NutritionMeal struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	NutritionPlanID primitive.ObjectID `bson:"nutritionPlanId" json:"nutritionPlanId"`
// 	MealTime       time.Time          `bson:"mealTime" json:"mealTime"`
// 	MealType       string             `bson:"mealType" json:"mealType"` // Breakfast, Lunch, Dinner, Snack
// 	Description    string             `bson:"description" json:"description"`
// 	Calories       int                `bson:"calories" json:"calories"`
// 	Notes          string             `bson:"notes" json:"notes"`
// 	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type NutritionPlan struct {
// 	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
// 	CreateBy       primitive.ObjectID `bson:"CreateBy" json:"CreateBy"`
// 	StartDate     time.Time          `bson:"startDate" json:"startDate"`
// 	EndDate       time.Time          `bson:"endDate" json:"endDate"`
// 	TotalCalories int                `bson:"totalCalories" json:"totalCalories"`
// 	MealsPerDay   int                `bson:"mealsPerDay" json:"mealsPerDay"`
// 	Notes         string             `bson:"notes" json:"notes"`
// 	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Performance struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID     primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
// 	ScheduleID primitive.ObjectID `bson:"scheduleId" json:"scheduleId"`
// 	Value      float64            `bson:"value" json:"value"`
// 	Date       time.Time          `bson:"date" json:"date"`
// 	MetricType string             `bson:"metricType" json:"metricType"` // Speed, Distance, Weight, etc.
// 	Notes      string             `bson:"notes" json:"notes"`
// 	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Progress struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID     primitive.ObjectID `bson:"userId" json:"userId"`
// 	MetricType string             `bson:"metricType" json:"metricType"`
// 	Value      float64            `bson:"value" json:"value"`
// 	Date       time.Time          `bson:"date" json:"date"`
// 	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Reminder struct {
// 	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	UserID          primitive.ObjectID `bson:"userId,oitempty" json:"userId"`
// 	ScheduleID      primitive.ObjectID `bson:"scheduleId,omitempty" json:"scheduleId"`
// 	NutritionPlanID primitive.ObjectID `bson:"nutritionPlanId,omitempty" json:"nutritionPlanId"`
// 	ReminderTime    time.Time          `bson:"reminderTime" json:"reminderTime"`
// 	ReminderDate    time.Time          `bson:"reminderDate" json:"reminderDate"`
// 	Content         string             `bson:"content" json:"content"`
// 	Status          string             `bson:"status" json:"status"` // Pending, Sent, Dismissed
// 	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Sport struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Name      string             `bson:"name" json:"name"`
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type SportAthlete struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	SportID   primitive.ObjectID `bson:"sportId" json:"sportId"`
// 	UserID    primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
// 	Position  string             `bson:"position" json:"position"`
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Team struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Name        string             `bson:"name" json:"name"`
// 	SportID     primitive.ObjectID `bson:"sportId" json:"sportId"`
// 	Description string             `bson:"description" json:"description"`
// 	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"` // UserID
// 	CreatedDate time.Time          `bson:"createdDate" json:"createdDate"`
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type TeamMember struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	TeamID    primitive.ObjectID `bson:"teamId" json:"teamId"`
// 	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Tournament struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Name        string             `bson:"name" json:"name"`
// 	Location    string             `bson:"location" json:"location"`
// 	StartDate   time.Time          `bson:"startDate" json:"startDate"`
// 	EndDate     time.Time          `bson:"endDate" json:"endDate"`
// 	Level       string             `bson:"level" json:"level"` // Local, Regional, National, International
// 	Organizer   string             `bson:"organizer" json:"organizer"`
// 	Description string             `bson:"description" json:"description"`
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type TrainingExercise struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	ScheduleID primitive.ObjectID `bson:"scheduleId" json:"scheduleId"`
// 	ExerciseID primitive.ObjectID `bson:"exerciseId" json:"exerciseId"`
// 	Order      int                `bson:"order" json:"order"`
// 	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type TrainingSchedule struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Date      time.Time          `bson:"date" json:"date"`
// 	StartTime time.Time          `bson:"startTime" json:"startTime"`
// 	EndTime   time.Time          `bson:"endTime" json:"endTime"`
// 	Status    string             `bson:"status" json:"status"` // Scheduled, Completed, Cancelled
// 	Location  string             `bson:"location" json:"location"`
// 	Type      string             `bson:"type" json:"type"`
// 	Notes     string             `bson:"notes" json:"notes"`
// 	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"` // CoachID
// 	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type TrainingScheduleUser struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	ScheduleID primitive.ObjectID `bson:"scheduleId" json:"scheduleId"`
// 	UserID     primitive.ObjectID `bson:"userId" json:"userId"` // AthleteID
// 	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
// package models

// import (
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )
// const(
// 	RoleAdmin     = "admin"
//     RoleUser      = "athlete"
//     RoleModerator = "coach"
// )

// // User represents the base user model
// type User struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Gender      string             `bson:"gender" json:"gender"`
// 	FullName    string             `bson:"fullName" json:"fullName"`
// 	HashPass    string             `bson:"hashPass" json:"hashPass"`
// 	Email       string             `bson:"email" json:"email"`
// 	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
// 	DateOfBirth time.Time          `bson:"dateOfBirth" json:"dateOfBirth"`
// 	Role        string             `bson:"role" json:"role"`
// 	Status      string             `bson:"status" json:"status"`
// 	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
// 	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
// }
