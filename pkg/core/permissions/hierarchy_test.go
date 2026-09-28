package permissions

import "testing"

func TestCanManageRole(t *testing.T) {
	tests := []struct {
		name       string
		actorRank  int
		targetRank int
		want       bool
	}{
		{
			name:       "actor can manage lower authority role",
			actorRank:  5,
			targetRank: 10,
			want:       true,
		},
		{
			name:       "actor cannot manage same rank",
			actorRank:  5,
			targetRank: 5,
			want:       false,
		},
		{
			name:       "actor cannot manage higher authority role",
			actorRank:  5,
			targetRank: 2,
			want:       false,
		},
		{
			name:       "highest authority can manage lower role",
			actorRank:  2,
			targetRank: 10,
			want:       true,
		},
		{
			name:       "lower authority cannot manage higher role",
			actorRank:  10,
			targetRank: 2,
			want:       false,
		},
		{
			name:       "adjacent lower authority",
			actorRank:  5,
			targetRank: 6,
			want:       true,
		},
		{
			name:       "adjacent higher authority",
			actorRank:  5,
			targetRank: 4,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanManageRole(tt.actorRank, tt.targetRank)

			if got != tt.want {
				t.Fatalf(
					"expected %v, got %v",
					tt.want,
					got,
				)
			}
		})
	}
}

func TestCanManageMember(t *testing.T) {
	tests := []struct {
		name       string
		actorRank  int
		targetRank int
		want       bool
	}{
		{
			name:       "actor can manage lower authority member",
			actorRank:  5,
			targetRank: 10,
			want:       true,
		},
		{
			name:       "actor cannot manage same rank member",
			actorRank:  5,
			targetRank: 5,
			want:       false,
		},
		{
			name:       "actor cannot manage higher authority member",
			actorRank:  5,
			targetRank: 2,
			want:       false,
		},
		{
			name:       "highest authority can manage lower member",
			actorRank:  2,
			targetRank: 10,
			want:       true,
		},
		{
			name:       "lower authority cannot manage higher member",
			actorRank:  10,
			targetRank: 2,
			want:       false,
		},
		{
			name:       "adjacent lower authority",
			actorRank:  5,
			targetRank: 6,
			want:       true,
		},
		{
			name:       "adjacent higher authority",
			actorRank:  5,
			targetRank: 4,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanManageMember(tt.actorRank, tt.targetRank)

			if got != tt.want {
				t.Fatalf(
					"expected %v, got %v",
					tt.want,
					got,
				)
			}
		})
	}
}
