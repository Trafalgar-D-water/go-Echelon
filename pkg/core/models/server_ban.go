package models

import "time"

// ServerBan represents a user that has been banned from a server.
type ServerBan struct {
	// Composite key: server + user
	ID struct {
		ServerID ID `bson:"server" json:"server"`
		UserID   ID `bson:"user" json:"user"`
	} `bson:"_id" json:"id"`

	Reason   *string   `bson:"reason,omitempty" json:"reason,omitempty"`
	BannedAt time.Time `bson:"banned_at" json:"banned_at"`
	BannedBy ID        `bson:"banned_by" json:"banned_by"`
}
