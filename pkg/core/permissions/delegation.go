package permissions

// CanDelegatePermissions determines whether an actor can grant the
// requested permissions to a target role.
func CanDelegatePermissions(input PermissionDelegationInput) bool {
	// Owner can manage normally.
	if input.IsOwner {
		return true
	}

	// Actor cannot modify an equal or higher role.
	if !CanManageRole(input.ActorRank, input.TargetRank) {
		return false
	}

	// Actor cannot grant permissions they do not possess.
	if !input.ActorPermissions.HasAll(input.GrantedPermissions) {
		return false
	}

	return true
}
