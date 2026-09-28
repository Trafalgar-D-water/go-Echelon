package permissions

type RolePermissionInput struct {
	ID       string
	Rank     int
	Override Override
}

type ServerPermissionInput struct {
	IsPrivileged       bool
	IsOwner            bool
	IsMember           bool
	DefaultPermissions Mask

	Roles []RolePermissionInput

	TimedOut   bool
	CanPublish bool
	CanReceive bool
}

type ChannelPermissionInput struct {
	ServerPermissions PermissionValue

	DefaultOverride Override

	RoleOverrides []RolePermissionInput

	MemberOverride Override

	CanPublish bool
	CanReceive bool
	TimedOut   bool
}

type PermissionDelegationInput struct {
	ActorPermissions Mask

	ActorRank  int
	TargetRank int

	GrantedPermissions Mask

	IsOwner bool
}
