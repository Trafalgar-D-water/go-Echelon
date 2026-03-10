package models

import "time"

// BotFlags bitfield for special bot attributes.
type BotFlags uint32

const (
	BotFlagVerified        BotFlags = 1 << 0
	BotFlagPublic          BotFlags = 1 << 1 // can be invited by anyone
	BotFlagAnalyticsOptOut BotFlags = 1 << 2
)

// Bot represents an automated user (application bot).
type Bot struct {
	ID           ID       `bson:"_id,omitempty" json:"id"` // same as the bot's User.ID
	OwnerID      ID       `bson:"owner_id" json:"owner_id"`
	Token        string   `bson:"token" json:"-"` // API token — never exposed
	Flags        BotFlags `bson:"flags" json:"flags"`
	Interactions *string  `bson:"interactions_url,omitempty" json:"interactions_url,omitempty"` // webhook URL

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
