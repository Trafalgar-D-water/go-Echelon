package permissions

import "testing"

func TestResolveChannel(t *testing.T) {
	tests := []struct {
		name  string
		input ChannelPermissionInput
		want  PermissionValue
	}{
		{
			name: "server permissions inherited",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(SendMessage),
				),
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory) |
					uint64(SendMessage),
			),
		},

		{
			name: "everyone channel allow",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory),
				),
				DefaultOverride: Override{
					Allow: Mask(SendMessage),
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory) |
					uint64(SendMessage),
			),
		},

		{
			name: "everyone channel deny",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(SendMessage),
				),
				DefaultOverride: Override{
					Deny: Mask(SendMessage),
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},

		{
			name: "channel role allow",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory),
				),
				RoleOverrides: []RolePermissionInput{
					{
						ID:   "moderator",
						Rank: 10,
						Override: Override{
							Allow: Mask(SendMessage),
						},
					},
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory) |
					uint64(SendMessage),
			),
		},

		{
			name: "channel role deny",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(SendMessage),
				),
				RoleOverrides: []RolePermissionInput{
					{
						ID:   "moderator",
						Rank: 10,
						Override: Override{
							Deny: Mask(SendMessage),
						},
					},
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},

		{
			name: "multiple channel roles",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory),
				),
				RoleOverrides: []RolePermissionInput{
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
							Allow: Mask(ManageMessages),
						},
					},
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory) |
					uint64(SendMessage) |
					uint64(ManageMessages),
			),
		},

		{
			name: "member-specific override",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(SendMessage),
				),
				MemberOverride: Override{
					Deny: Mask(SendMessage),
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},

		{
			name: "member override wins over role override",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory),
				),
				RoleOverrides: []RolePermissionInput{
					{
						ID:   "moderator",
						Rank: 10,
						Override: Override{
							Allow: Mask(SendMessage),
						},
					},
				},
				MemberOverride: Override{
					Deny: Mask(SendMessage),
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},

		{
			name: "timeout",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(SendMessage) |
						uint64(Connect) |
						uint64(Speak),
				),
				MemberOverride: Override{
					Allow: Mask(SendMessage),
				},
				CanPublish: true,
				CanReceive: true,
				TimedOut:   true,
			},
			want: FromRaw(uint64(ALLOW_IN_TIMEOUT)),
		},

		{
			name: "publish restriction",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(Speak) |
						uint64(Video),
				),
				CanPublish: false,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},

		{
			name: "receive restriction",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(Listen),
				),
				CanPublish: true,
				CanReceive: false,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},

		{
			name: "no ViewChannel",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ReadMessageHistory) |
						uint64(SendMessage) |
						uint64(ManageMessages),
				),
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(0),
		},

		{
			name: "SendMessage denied by channel",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(SendMessage),
				),
				DefaultOverride: Override{
					Deny: Mask(SendMessage),
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},

		{
			name: "server allows SendMessage but channel denies it",
			input: ChannelPermissionInput{
				ServerPermissions: FromRaw(
					uint64(ViewChannel) |
						uint64(ReadMessageHistory) |
						uint64(SendMessage),
				),
				DefaultOverride: Override{
					Deny: Mask(SendMessage),
				},
				CanPublish: true,
				CanReceive: true,
			},
			want: FromRaw(
				uint64(ViewChannel) |
					uint64(ReadMessageHistory),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveChannel(tt.input)

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
