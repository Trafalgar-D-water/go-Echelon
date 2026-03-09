package models

import "time"

// MemberCompositeKey is the composite primary key for a server member.
// MongoDB stores this as the _id.
type MemberCompositeKey struct {
	ServerID ID `bson:"server" json:"server"`
	UserID   ID `bson:"user" json:"user"`
}

// Member represents a user's membership within a server.
type Member struct {
	ID       MemberCompositeKey `bson:"_id" json:"id"`
	Nickname *string            `bson:"nickname,omitempty" json:"nickname,omitempty"`
	Roles    []ID               `bson:"roles" json:"roles"`
	Avatar   *ID                `bson:"avatar,omitempty" json:"avatar,omitempty"`
	JoinedAt time.Time          `bson:"joined_at" json:"joined_at"`
	// Timeout — member cannot send messages until this time
	Timeout *time.Time `bson:"timeout,omitempty" json:"timeout,omitempty"`
}
