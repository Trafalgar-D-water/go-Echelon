package models

import "time"

// Invite is a single-use or unlimited-use invite link to a channel/server.
type Invite struct {
	Code      string     `bson:"_id" json:"code"` // the invite code (primary key)
	CreatorID ID         `bson:"creator_id" json:"creator_id"`
	ChannelID ID         `bson:"channel_id" json:"channel_id"`
	ServerID  *ID        `bson:"server_id,omitempty" json:"server_id,omitempty"`
	Uses      int        `bson:"uses" json:"uses"`
	MaxUses   *int       `bson:"max_uses,omitempty" json:"max_uses,omitempty"` // nil = unlimited
	ExpiresAt *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
}
