package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RelFriend   int = 1
	RelBlocked  int = 2
	RelIncoming int = 3 // Target received the request
	RelOutgoing int = 4 // Actor sent the request
)

type Relationship struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActorID   string             `bson:"actor_id" json:"actor_id"`
	TargetID  string             `bson:"target_id" json:"target_id"`
	Type      int                `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
