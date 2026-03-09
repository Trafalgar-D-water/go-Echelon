package models

import "time"

// PolicyChange records when a user acknowledged a platform policy update.
type PolicyChange struct {
	ID        ID        `bson:"_id,omitempty" json:"id"`
	UserID    ID        `bson:"user_id" json:"user_id"`
	PolicyKey string    `bson:"policy_key" json:"policy_key"` // e.g. "tos_v2", "privacy_v3"
	AckedAt   time.Time `bson:"acked_at" json:"acked_at"`
}
