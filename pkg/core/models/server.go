package models

import "time"

// ServerFlags contains bitfield flags for a server.
type ServerFlags uint32

const (
	ServerFlagVerified ServerFlags = 1 << 0
	ServerFlagOfficial ServerFlags = 1 << 1
)

// Category groups channels inside a server.
type Category struct {
	ID       ID     `bson:"id" json:"id"`
	Title    string `bson:"title" json:"title"`
	Channels []ID   `bson:"channels" json:"channels"`
}

// SystemMessages holds IDs of special system message channels.
type SystemMessages struct {
	UserJoinedID *ID `bson:"user_joined,omitempty" json:"user_joined,omitempty"`
	UserLeftID   *ID `bson:"user_left,omitempty" json:"user_left,omitempty"`
	UserBannedID *ID `bson:"user_banned,omitempty" json:"user_banned,omitempty"`
	UserKickedID *ID `bson:"user_kicked,omitempty" json:"user_kicked,omitempty"`
}

// Role represents a permission group within a server.
type Role struct {
	ID          ID         `bson:"_id,omitempty" json:"id"`
	Name        string     `bson:"name" json:"name"`
	Colour      *string    `bson:"colour,omitempty" json:"colour,omitempty"` // hex e.g. "#FF5733"
	Hoist       bool       `bson:"hoist" json:"hoist"`                       // shown separately in member list
	Rank        int        `bson:"rank" json:"rank"`                         // lower = higher priority
	Permissions Permission `bson:"permissions" json:"permissions"`
}

// Server is the top-level community/guild.
type Server struct {
	ID             ID             `bson:"_id,omitempty" json:"id"`
	OwnerID        ID             `bson:"owner_id" json:"owner_id"`
	Name           string         `bson:"name" json:"name"`
	Description    *string        `bson:"description,omitempty" json:"description,omitempty"`
	IconID         *ID            `bson:"icon_id,omitempty" json:"icon_id,omitempty"`
	BannerID       *ID            `bson:"banner_id,omitempty" json:"banner_id,omitempty"`
	Categories     []Category     `bson:"categories" json:"categories"`
	SystemMessages SystemMessages `bson:"system_messages" json:"system_messages"`
	Roles          []Role         `bson:"roles" json:"roles"`
	Flags          ServerFlags    `bson:"flags" json:"flags"`
	NSFW           bool           `bson:"nsfw" json:"nsfw"`
	CreatedAt      time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `bson:"updated_at" json:"updated_at"`
}
