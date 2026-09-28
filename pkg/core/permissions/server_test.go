package permissions

import "testing"

func TestResolveServer(t *testing.T) {
	tests := []struct {
		name  string
		input ServerPermissionInput
		want  PermissionValue
	}{
		{
			name: "owner",
			input: ServerPermissionInput{
				IsOwner:            true,
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
			},
			want: FromRaw(uint64(GrantAllSafe)),
		},
		{
			name: "privileged user",
			input: ServerPermissionInput{
				IsPrivileged: true,
				IsMember:     false,
			},
			want: FromRaw(uint64(GrantAllSafe)),
		},
		{
			name: "non-member",
			input: ServerPermissionInput{
				IsMember:           false,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
			},
			want: FromRaw(0),
		},
		{
			name: "normal member",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
			},
			want: FromRaw(uint64(DEFAULT_PERMISSION_SERVER)),
		},
		{
			name: "member with one role",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
				Roles: []RolePermissionInput{
					{
						ID:   "moderator",
						Rank: 10,
						Override: Override{
							Allow: Mask(ManageMessages),
						},
					},
				},
			},
			want: FromRaw(
				uint64(DEFAULT_PERMISSION_SERVER) |
					uint64(ManageMessages),
			),
		},
		{
			name: "member with multiple roles",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
				Roles: []RolePermissionInput{
					{
						ID:   "moderator",
						Rank: 10,
						Override: Override{
							Allow: Mask(ManageMessages),
						},
					},
					{
						ID:   "admin",
						Rank: 2,
						Override: Override{
							Allow: Mask(ManageChannel),
						},
					},
				},
			},
			want: FromRaw(
				uint64(DEFAULT_PERMISSION_SERVER) |
					uint64(ManageMessages) |
					uint64(ManageChannel),
			),
		},
		{
			name: "conflicting roles",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
				Roles: []RolePermissionInput{
					{
						ID:   "moderator",
						Rank: 10,
						Override: Override{
							Allow: Mask(SendMessage),
						},
					},
					{
						ID:   "admin",
						Rank: 2,
						Override: Override{
							Deny: Mask(SendMessage),
						},
					},
				},
			},
			want: FromRaw(
				uint64(DEFAULT_PERMISSION_SERVER) &
					^uint64(SendMessage),
			),
		},
		{
			name: "allow and deny",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
				Roles: []RolePermissionInput{
					{
						ID:   "moderator",
						Rank: 10,
						Override: Override{
							Allow: Mask(ManageMessages),
							Deny:  Mask(SendMessage),
						},
					},
				},
			},
			want: FromRaw(
				(uint64(DEFAULT_PERMISSION_SERVER) |
					uint64(ManageMessages)) &
					^uint64(SendMessage),
			),
		},
		{
			name: "timeout",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         true,
				TimedOut:           true,
			},
			want: FromRaw(uint64(ALLOW_IN_TIMEOUT)),
		},
		{
			name: "cannot publish",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         false,
				CanReceive:         true,
			},
			want: FromRaw(
				uint64(DEFAULT_PERMISSION_SERVER) &
					^uint64(Speak) &
					^uint64(Video),
			),
		},
		{
			name: "cannot receive",
			input: ServerPermissionInput{
				IsMember:           true,
				DefaultPermissions: DEFAULT_PERMISSION_SERVER,
				CanPublish:         true,
				CanReceive:         false,
			},
			want: FromRaw(
				uint64(DEFAULT_PERMISSION_SERVER) &
					^uint64(Listen),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveServer(tt.input)

			if got.IntoRaw() != tt.want.IntoRaw() {
				t.Fatalf(
					"expected %v, got %v",
					tt.want.IntoRaw(),
					got.IntoRaw(),
				)
			}
		})
	}
}

func TestResolveServer_RoleOrderDoesNotMatter(t *testing.T) {
	rolesA := []RolePermissionInput{
		{
			ID:   "moderator",
			Rank: 10,
			Override: Override{
				Allow: Mask(SendMessage),
			},
		},
		{
			ID:   "admin",
			Rank: 2,
			Override: Override{
				Deny: Mask(SendMessage),
			},
		},
	}

	rolesB := []RolePermissionInput{
		{
			ID:   "admin",
			Rank: 2,
			Override: Override{
				Deny: Mask(SendMessage),
			},
		},
		{
			ID:   "moderator",
			Rank: 10,
			Override: Override{
				Allow: Mask(SendMessage),
			},
		},
	}

	inputA := ServerPermissionInput{
		IsMember:           true,
		DefaultPermissions: DEFAULT_PERMISSION_SERVER,
		CanPublish:         true,
		CanReceive:         true,
		Roles:              rolesA,
	}

	inputB := ServerPermissionInput{
		IsMember:           true,
		DefaultPermissions: DEFAULT_PERMISSION_SERVER,
		CanPublish:         true,
		CanReceive:         true,
		Roles:              rolesB,
	}

	gotA := ResolveServer(inputA)
	gotB := ResolveServer(inputB)

	if gotA.IntoRaw() != gotB.IntoRaw() {
		t.Fatalf(
			"changing role order changed result: A=%v B=%v",
			gotA.IntoRaw(),
			gotB.IntoRaw(),
		)
	}
}
