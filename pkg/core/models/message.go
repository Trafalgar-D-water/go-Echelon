package models

import "time"

// Reply holds a reference to the original message being replied to.
type Reply struct {
	MessageID ID   `bson:"message_id" json:"message_id"`
	Mention   bool `bson:"mention" json:"mention"` // whether the reply pings the original author
}

// MessageContent holds the textual body + attachments.
type MessageContent struct {
	// Plain text / markdown body
	Text *string `bson:"text,omitempty" json:"text,omitempty"`
	// Attachment file IDs (references to autumn/file storage)
	Attachments []ID `bson:"attachments,omitempty" json:"attachments,omitempty"`
}

// MessageFlag bitfield for system messages.
type MessageFlag uint32

const (
	MessageFlagSuppressNotifications MessageFlag = 1 << 0
	MessageFlagSystem                MessageFlag = 1 << 1
)

// Message is a single message in a channel.
type Message struct {
	ID        ID             `bson:"_id,omitempty" json:"id"`
	ChannelID ID             `bson:"channel_id" json:"channel_id"`
	AuthorID  ID             `bson:"author_id" json:"author_id"`
	Content   MessageContent `bson:"content" json:"content"`
	Replies   []Reply        `bson:"replies,omitempty" json:"replies,omitempty"`
	Mentions  []ID           `bson:"mentions,omitempty" json:"mentions,omitempty"`
	Embeds    []Embed        `bson:"embeds,omitempty" json:"embeds,omitempty"`
	Reactions []Reaction     `bson:"reactions,omitempty" json:"reactions,omitempty"`
	Flags     MessageFlag    `bson:"flags" json:"flags"`
	Pinned    bool           `bson:"pinned" json:"pinned"`
	Edited    *time.Time     `bson:"edited,omitempty" json:"edited,omitempty"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
}

// Reaction groups users who reacted with the same emoji.
type Reaction struct {
	EmojiID *ID    `bson:"emoji_id,omitempty" json:"emoji_id,omitempty"` // nil = unicode emoji
	Emoji   string `bson:"emoji" json:"emoji"`                           // unicode char or custom emoji name
	Count   int    `bson:"count" json:"count"`
	UserIDs []ID   `bson:"user_ids" json:"user_ids"`
}
