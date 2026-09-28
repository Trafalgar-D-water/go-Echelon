package permissions

func ResolveServer(input ServerPermissionInput) PermissionValue {
	// 1. Privileged user
	if input.IsPrivileged {
		return FromRaw(uint64(GrantAllSafe))
	}

	// 2. Server owner
	if input.IsOwner {
		return FromRaw(uint64(GrantAllSafe))
	}

	// 3. Non member
	if !input.IsMember {
		return FromRaw(0)
	}

	// 4. Start with server default permissions
	value := FromRaw(uint64(input.DefaultPermissions))

	// 5. Apply assigned role overrides
	roles, err := SortRoles(input.Roles)

	if err != nil {
		return FromRaw(0)
	}

	for _, role := range roles {
		value = role.Override.Apply(value)
	}

	// 6. Apply publish restriction

	if !input.CanPublish {
		value.Revoke(Speak)
		value.Revoke(Video)
	}

	// 7. Apply receive restriction
	if !input.CanReceive {
		value.Revoke(Listen)
	}

	// 8. Apply timeout restriction
	if input.TimedOut {
		value.Restrict(ALLOW_IN_TIMEOUT)
	}

	// 9. Return final permissions
	return value

}
