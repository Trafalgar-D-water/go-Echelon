package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username      string             `bson:"username" json:"username"`
	Discriminator string             `bson:"discriminator" json:"discriminator"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"-"`
	Online        bool               `bson:"online" json:"online"`
	IsVerified    bool               `bson:"is_verified" json:"is_verified"`
	OTP           string             `bson:"otp,omitempty" json:"-"`
	DOB           time.Time          `bson:"dob" json:"dob"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
