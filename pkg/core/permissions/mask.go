package permissions

type Mask uint64

// Has checks whether the mask contains the given permission.
func (m Mask) Has(permission Permission) bool {
	return m&Mask(permission) != 0
}

// HasAll checks whether the mask contains all permissions
// represented by the required mask.
func (m Mask) HasAll(required Mask) bool {
	return m&required == required
}

// Add enables the given permission.
func (m *Mask) add(permission Permission) {
	*m |= Mask(permission)
}

// Remove disables the given permission.
func (m *Mask) remove(permission Permission) {
	*m &^= Mask(permission)
}

// GrantAllSafe enables all permissions in the safe permission range.
func (m *Mask) grantAllSafe() {
	*m = Mask(GrantAllSafe)
}

// Clear removes all permissions.
func (m *Mask) clear() {
	*m = 0
}
