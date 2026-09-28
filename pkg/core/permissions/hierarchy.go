package permissions

// CanManageRole determines whether an actor can manage a target role.
//
// Smaller rank means higher authority.
//
// Owner:
//
//	handled by the caller because ownership is not represented by rank.
//
// Non-owner:
//
//	actor can manage only targets with a higher numeric rank.
func CanManageRole(actorRank, targetRank int) bool {
	return actorRank < targetRank
}

// CanManageMember determines whether an actor can manage a target member.
//
// Smaller rank means higher authority.
//
// The actor must have a strictly higher authority level than the target.
func CanManageMember(actorRank, targetRank int) bool {
	return actorRank < targetRank
}
