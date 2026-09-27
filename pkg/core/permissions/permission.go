package permissions

type Permission uint64

const (
	// General
	ManageChannel       Permission = 1 << 0
	ManageServer        Permission = 1 << 1
	ManagePermissions   Permission = 1 << 2
	ManageRole          Permission = 1 << 3
	ManageCustomisation Permission = 1 << 4

	// Bit 5 is reserved.

	KickMembers    Permission = 1 << 6
	BanMembers     Permission = 1 << 7
	TimeoutMembers Permission = 1 << 8
	AssignRoles    Permission = 1 << 9

	ChangeNickname  Permission = 1 << 10
	ManageNicknames Permission = 1 << 11
	ChangeAvatar    Permission = 1 << 12
	RemoveAvatars   Permission = 1 << 13

	// Bits 14-19 are reserved.

	ViewChannel        Permission = 1 << 20
	ReadMessageHistory Permission = 1 << 21
	SendMessage        Permission = 1 << 22
	ManageMessages     Permission = 1 << 23
	ManageWebhooks     Permission = 1 << 24
	InviteOthers       Permission = 1 << 25
	SendEmbeds         Permission = 1 << 26
	UploadFiles        Permission = 1 << 27
	Masquerade         Permission = 1 << 28
	React              Permission = 1 << 29

	Connect Permission = 1 << 30
	Speak   Permission = 1 << 31
	Video   Permission = 1 << 32

	MuteMembers   Permission = 1 << 33
	DeafenMembers Permission = 1 << 34
	MoveMembers   Permission = 1 << 35
	Listen        Permission = 1 << 36

	MentionEveryone Permission = 1 << 37
	MentionRoles    Permission = 1 << 38
	BypassSlowmode  Permission = 1 << 39
	ViewAuditLogs   Permission = 1 << 40
)

// GrantAllSafe contains the permission range that is safe for
// normal permission assignment.
//
// Bits 0-51 are part of the safe permission range.
// Higher bits are reserved/prohibited for normal permissions.
const GrantAllSafe Permission = (1 << 52) - 1
