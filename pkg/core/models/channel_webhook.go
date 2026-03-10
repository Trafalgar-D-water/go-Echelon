package models

import "time"

// Webhook allows external services to post messages into a channel.
type Webhook struct {
	ID        ID        `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	AvatarID  *ID       `bson:"avatar_id,omitempty" json:"avatar_id,omitempty"`
	ChannelID ID        `bson:"channel_id" json:"channel_id"`
	ServerID  ID        `bson:"server_id" json:"server_id"`
	Token     string    `bson:"token" json:"-"` // never exposed in API by default
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
