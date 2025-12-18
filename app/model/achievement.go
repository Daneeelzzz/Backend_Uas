package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// [SRS 3.2.1] Collection achievements (MongoDB)
type AchievementMongo struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	StudentRefID    string                 `bson:"student_ref_id" json:"student_ref_id"` // User ID dari Postgres
	AchievementType string                 `bson:"achievement_type" json:"achievement_type"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	
	// Field Dinamis (Inti NoSQL)
	Details         map[string]interface{} `bson:"details" json:"details"`
	
	Tags            []string               `bson:"tags" json:"tags"`
	Points          int                    `bson:"points" json:"points"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
}

// [SRS 3.1.7] Tabel achievement_references (Postgres)
type AchievementRef struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	MongoAchievementID string    `json:"mongo_achievement_id"`
	Status             string    `json:"status"`
	VerifiedBy         *string   `json:"verified_by"` // Pointer karena bisa null
	RejectionNote      string    `json:"rejection_note"`
	CreatedAt          time.Time `json:"created_at"`
}