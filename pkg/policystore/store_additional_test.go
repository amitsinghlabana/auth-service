package policystore

import "testing"

import "github.com/stretchr/testify/require"

func TestPoliciesReturnsDefensiveCopy(t *testing.T) {
    store := NewInMemoryStore()

    policies := store.Policies()
    require.NotEmpty(t, policies)

    // Mutate the returned slice to ensure the backing store is not impacted.
    originalName := policies[0].Name
    policies[0].Name = "mutated-name"

    fresh := store.Policies()
    require.NotEmpty(t, fresh)
    require.Equal(t, originalName, fresh[0].Name, "expected underlying store to remain unchanged")
}

func TestEvaluateFailsWhenNoRoleMatchesPolicy(t *testing.T) {
    store := NewInMemoryStore()

    // The store has policies for "editor" and "viewer" only. Using a role that does not
    // intersect with any policy should result in a denial, even if the action and resource match.
    allowed, policies := store.Evaluate("user-000", []string{"guest"}, "article:view", "article:456")

    require.False(t, allowed, "expected store to deny access when roles do not match any policy")
    require.Empty(t, policies, "expected no policies to be returned when authorization is denied")
}
