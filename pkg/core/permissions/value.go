package permissions

type PermissionValue struct {
	Mask Mask
}

func FromRaw(raw uint64) PermissionValue {
	return PermissionValue{
		Mask: Mask(raw),
	}
}

func (v PermissionValue) IntoRaw() uint64 {
	return uint64(v.Mask)
}

func (v *PermissionValue) Allow(permission Permission) {
	v.Mask.Add(permission)
}

func (v *PermissionValue) Revoke(permission Permission) {
	v.Mask.Remove(permission)
}

func (v *PermissionValue) RevokeAll() {
	v.Mask.Clear()
}

func (v *PermissionValue) Restrict(mask Mask) {
	v.Mask &= mask
}

func (v PermissionValue) Has(permission Permission) bool {
	return v.Mask.Has(permission)
}
