package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       string             `bson:"userId" json:"userID"`
	RefreshToken string             `bson:"refreshToken" json:"refreshToken"`
	UserAgent    string             `bson:"userAgent" json:"userAgent"`
	IP           string             `bson:"ip" json:"ip"`
	ExpiresAt    time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
