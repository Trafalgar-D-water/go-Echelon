package permissions

func ResolveChannel(input ChannelPermissionInput) PermissionValue {
	// 1. Start with the already calculated server permissions.
	value := input.ServerPermissions

	// 2. Apply channel @everyone/default override.
	value = input.DefaultOverride.Apply(value)

	// 3. Apply channel role overrides.
	roles, err := SortRoles(input.RoleOverrides)
	if err != nil {
		return FromRaw(0)
	}

	for _, role := range roles {
		value = role.Override.Apply(value)
	}

	// 4. Apply member-specific override.
	value = input.MemberOverride.Apply(value)

	// 5. Apply publish restriction.
	if !input.CanPublish {
		value.Revoke(Speak)
		value.Revoke(Video)
	}

	// 6. Apply receive restriction.
	if !input.CanReceive {
		value.Revoke(Listen)
	}

	// 7. Apply timeout restriction.
	if input.TimedOut {
		value.Restrict(ALLOW_IN_TIMEOUT)
	}

	// 8. If the user cannot view the channel,
	// they have no permissions in the channel.
	if !value.Has(ViewChannel) {
		return FromRaw(0)
	}

	// 9. Return final permissions.
	return value
}
