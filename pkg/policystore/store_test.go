package policystore

import "testing"

func TestInMemoryStoreEvaluateAllow(t *testing.T) {
    store := NewInMemoryStore()
    allowed, policies := store.Evaluate("user-123", []string{"editor"}, "article:update", "article:456")

    if !allowed {
        t.Fatalf("expected authorization to be allowed")
    }

    if len(policies) != 1 || policies[0] != "editor-can-edit" {
        t.Fatalf("expected matching policy 'editor-can-edit', got %v", policies)
    }
}

func TestInMemoryStoreEvaluateDeny(t *testing.T) {
    store := NewInMemoryStore()
    allowed, policies := store.Evaluate("user-789", []string{"guest"}, "article:update", "article:456")

    if allowed {
        t.Fatalf("expected authorization to be denied")
    }

    if len(policies) != 0 {
        t.Fatalf("expected no policies to match, got %v", policies)
    }
}

func TestInMemoryStoreRegisterPolicyRequiresReview(t *testing.T) {
    store := NewInMemoryStore()
    err := store.RegisterPolicy(Policy{Name: "test", Actions: []string{"test"}, Resources: []string{"test:1"}, Roles: []string{"test-role"}})

    if err != ErrReviewRequired {
        t.Fatalf("expected ErrReviewRequired, got %v", err)
    }

    if len(store.Policies()) != 2 {
        t.Fatalf("expected policy count to remain unchanged, got %d", len(store.Policies()))
    }
}
