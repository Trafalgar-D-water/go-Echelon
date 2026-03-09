package models

import "time"

// ChannelType classifies a channel's purpose.
type ChannelType string

const (
	ChannelTypeSavedMessages ChannelType = "SavedMessages" // personal notes
	ChannelTypeDirectMessage ChannelType = "DirectMessage" // 1-on-1 DM
	ChannelTypeGroup         ChannelType = "Group"         // group DM
	ChannelTypeTextChannel   ChannelType = "TextChannel"   // server text channel
	ChannelTypeVoiceChannel  ChannelType = "VoiceChannel"  // server voice channel
)

// Channel represents any type of communication channel.
type Channel struct {
	ID          ID          `bson:"_id,omitempty" json:"id"`
	ChannelType ChannelType `bson:"channel_type" json:"channel_type"`

	// --- Server channels ---
	ServerID    *ID     `bson:"server_id,omitempty" json:"server_id,omitempty"`
	Name        *string `bson:"name,omitempty" json:"name,omitempty"`
	Description *string `bson:"description,omitempty" json:"description,omitempty"`
	NSFW        bool    `bson:"nsfw" json:"nsfw"`

	// --- DM / Group ---
	OwnerID    *ID  `bson:"owner_id,omitempty" json:"owner_id,omitempty"`     // group owner
	Recipients []ID `bson:"recipients,omitempty" json:"recipients,omitempty"` // DM participants
	IconID     *ID  `bson:"icon_id,omitempty" json:"icon_id,omitempty"`

	// Last message tracking (for unread indicators)
	LastMessageID *ID `bson:"last_message_id,omitempty" json:"last_message_id,omitempty"`

	// Permission overrides keyed by role/user ID string
	RolePermissions map[string]OverrideField `bson:"role_permissions,omitempty" json:"role_permissions,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
