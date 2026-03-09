package models

// Permission is a bitmask of server/channel permissions.
type Permission uint64

const (
	// General
	PermManageChannel     Permission = 1 << 0
	PermManageServer      Permission = 1 << 1
	PermManagePermissions Permission = 1 << 2
	PermManageRoles       Permission = 1 << 3
	PermManageCustomEmoji Permission = 1 << 4
	PermManageWebhooks    Permission = 1 << 5

	// Membership
	PermKickMembers    Permission = 1 << 6
	PermBanMembers     Permission = 1 << 7
	PermTimeoutMembers Permission = 1 << 8
	PermAssignRoles    Permission = 1 << 9

	// Messaging
	PermSendMessages       Permission = 1 << 22
	PermSendEmbeds         Permission = 1 << 23
	PermSendAttachments    Permission = 1 << 24
	PermSendMassMentions   Permission = 1 << 25
	PermSendVoiceMessages  Permission = 1 << 26
	PermReadMessageHistory Permission = 1 << 27
	PermManageMessages     Permission = 1 << 28
	PermReactToMessages    Permission = 1 << 29

	// Voice
	PermConnect     Permission = 1 << 33
	PermSpeak       Permission = 1 << 34
	PermVideo       Permission = 1 << 35
	PermMuteUsers   Permission = 1 << 36
	PermDeafenUsers Permission = 1 << 37
	PermMoveUsers   Permission = 1 << 38

	// Special
	PermGrantAllSafe Permission = 0x000F_FFFF_FFFF_FFFF
	PermGrantAll     Permission = ^Permission(0)
)

// Has reports whether the permission set includes the given permission.
func (p Permission) Has(flag Permission) bool {
	return p&flag == flag
}

// OverrideField represents a permission override for a role or user on a channel.
type OverrideField struct {
	Allow Permission `bson:"allow" json:"allow"`
	Deny  Permission `bson:"deny" json:"deny"`
}
