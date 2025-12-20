package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

// NewID generates a new MongoDB ObjectID
func NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// StringToID converts a hex string to an ObjectID
func StringToID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
