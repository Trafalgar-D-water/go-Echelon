package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ID is an alias for MongoDB ObjectID — used across all models.
type ID = primitive.ObjectID

// NewID generates a new unique ObjectID.
func NewID() ID {
	return primitive.NewObjectID()
}

// IDFromHex parses a hex string into an ObjectID.
func IDFromHex(s string) (ID, error) {
	return primitive.ObjectIDFromHex(s)
}
