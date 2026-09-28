package permissions

type PermissionValue struct {
	mask Mask
}

func FromRaw(raw uint64) PermissionValue {
	return PermissionValue{
		mask: Mask(raw),
	}
}

func (v PermissionValue) IntoRaw() uint64 {
	return uint64(v.mask)
}

func (v *PermissionValue) allow(permission Permission) {
	v.mask.add(permission)
}

func (v *PermissionValue) revoke(permission Permission) {
	v.mask.remove(permission)
}

func (v *PermissionValue) revokeAll() {
	v.mask.clear()
}

func (v *PermissionValue) restrict(mask Mask) {
	v.mask &= mask
}

func (v PermissionValue) Has(permission Permission) bool {
	return v.mask.Has(permission)
}
