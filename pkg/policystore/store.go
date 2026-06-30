package policystore

import (
    "fmt"
    "sync"
)

type Policy struct {
    Name       string
    Actions    []string
    Resources  []string
    Roles      []string
}

type DataSource interface {
    Load() []Policy
}

type PolicyStore interface {
    Evaluate(subjectID string, roles []string, action, resource string) (bool, []string)
    RegisterPolicy(policy Policy) error
    Policies() []Policy
}

type InMemoryStore struct {
    mu       sync.RWMutex
    policies []Policy
    watchers []chan struct{}
}

func NewInMemoryStore() *InMemoryStore {
    policies := []Policy{
        {Name: "editor-can-edit", Actions: []string{"article:update"}, Resources: []string{"article:456"}, Roles: []string{"editor"}},
        {Name: "viewer-can-read", Actions: []string{"article:view"}, Resources: []string{"article:456"}, Roles: []string{"viewer", "editor"}},
    }

    return &InMemoryStore{policies: policies}
}

func (s *InMemoryStore) Evaluate(subjectID string, roles []string, action, resource string) (bool, []string) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    matched := []string{}
    for _, policy := range s.policies {
        if !contains(policy.Actions, action) {
            continue
        }
        if !contains(policy.Resources, resource) {
            continue
        }
        if intersects(policy.Roles, roles) {
            matched = append(matched, policy.Name)
        }
    }

    return len(matched) > 0, matched
}

func (s *InMemoryStore) RegisterPolicy(policy Policy) error {
    // Human-in-the-loop gate: require review flag
    if !policyReviewApproved() {
        return ErrReviewRequired
    }

    s.mu.Lock()
    defer s.mu.Unlock()
    s.policies = append(s.policies, policy)
    s.emitChange()
    return nil
}

func (s *InMemoryStore) emitChange() {
    for _, ch := range s.watchers {
        select {
        case ch <- struct{}{}:
        default:
        }
    }
}

func (s *InMemoryStore) Policies() []Policy {
    s.mu.RLock()
    defer s.mu.RUnlock()
    copy := make([]Policy, len(s.policies))
    copy = append(copy[:0], s.policies...)
    return copy
}

func contains(list []string, item string) bool {
    for _, v := range list {
        if v == item {
            return true
        }
    }
    return false
}

func intersects(a, b []string) bool {
    set := make(map[string]struct{})
    for _, v := range a {
        set[v] = struct{}{}
    }
    for _, v := range b {
        if _, ok := set[v]; ok {
            return true
        }
    }
    return false
}

var (
    ErrReviewRequired = fmt.Errorf("policy write requires human review completion")
)

func policyReviewApproved() bool {
    // Stubbed. In practice, this would verify a gate (e.g. manual flag, ticket ID, etc.)
    return false
}
