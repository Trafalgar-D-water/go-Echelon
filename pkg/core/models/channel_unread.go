package models

// ChannelUnread tracks which messages a user has already seen in a channel.
type ChannelUnread struct {
	// Composite key: channel + user
	ID struct {
		ChannelID ID `bson:"channel" json:"channel"`
		UserID    ID `bson:"user" json:"user"`
	} `bson:"_id" json:"id"`

	// ID of the last message the user has read
	LastID *ID `bson:"last_id,omitempty" json:"last_id,omitempty"`

	// IDs of messages that explicitly mentioned this user (unread)
	Mentions []ID `bson:"mentions" json:"mentions"`
}
