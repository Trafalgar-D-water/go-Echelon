package permissions

import "errors"

// Override represents a permission override.
//
// Allow contains permissions that should be granted.
// Deny contains permissions that should be revoked.
type Override struct {
	Allow Mask
	Deny  Mask
}

// Validate checks whether the override is valid.
//
// A permission cannot exist in both Allow and Deny.
func (o Override) validate() error {
	if o.Allow&o.Deny != 0 {
		return errors.New("permission cannot exist in both allow and deny")
	}

	return nil
}

// Apply applies the override to the current permission value.
//
// The order is important:
//
// 1. Allow permissions
// 2. Deny permissions
//
// Deny is applied second, so it wins if both operations
// were ever applied to the same permission.
func (o Override) apply(current PermissionValue) PermissionValue {
	current.mask |= o.Allow
	current.mask &^= o.Deny

	return current
}
