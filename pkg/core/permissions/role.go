package permissions

import (
	"fmt"
	"sort"
)

// SortRoles returns roles ordered deterministically.
//
// Ordering rules:
//  1. Higher numeric rank comes first.
//  2. If two roles have the same rank, role ID is used
//     as a deterministic tie-breaker.
//  3. Duplicate role IDs are rejected.
//
// The input slice is not modified.
func sortRoles(roles []RolePermissionInput) ([]RolePermissionInput, error) {
	sorted := append([]RolePermissionInput(nil), roles...)

	seen := make(map[string]struct{}, len(sorted))

	for _, role := range sorted {
		if _, exists := seen[role.ID]; exists {
			return nil, fmt.Errorf("duplicate role ID: %s", role.ID)
		}

		seen[role.ID] = struct{}{}
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Rank != sorted[j].Rank {
			return sorted[i].Rank > sorted[j].Rank
		}

		return sorted[i].ID < sorted[j].ID
	})

	return sorted, nil
}
