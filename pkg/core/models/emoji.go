package models

import "time"

// EmojiParent indicates where this emoji belongs (a Server or is Detached).
type EmojiParent struct {
	Type     string `bson:"type" json:"type"` // "Server" | "Detached"
	ServerID *ID    `bson:"server_id,omitempty" json:"server_id,omitempty"`
}

// Emoji is a custom emoji uploaded to a server.
type Emoji struct {
	ID        ID          `bson:"_id,omitempty" json:"id"`
	Parent    EmojiParent `bson:"parent" json:"parent"`
	CreatorID ID          `bson:"creator_id" json:"creator_id"`
	Name      string      `bson:"name" json:"name"`
	Animated  bool        `bson:"animated" json:"animated"`
	NSFW      bool        `bson:"nsfw" json:"nsfw"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
}
