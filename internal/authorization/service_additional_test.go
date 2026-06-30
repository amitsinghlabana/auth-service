package authorization

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/require"
)

type stubStore struct {
    allowed   bool
    policies  []string
    lastInput Request
}

func (s *stubStore) Evaluate(subjectID string, roles []string, action, resource string) (bool, []string) {
    s.lastInput = Request{
        SubjectID: subjectID,
        Roles:     roles,
        Action:    action,
        Resource:  resource,
    }
    return s.allowed, s.policies
}

func (s *stubStore) RegisterPolicy(policy Policy) error {
    return ErrReviewRequired
}

func (s *stubStore) Policies() []Policy {
    return nil
}

type stubAuditor struct {
    decisions []audit.Decision
}

func (a *stubAuditor) Record(decision audit.Decision) {
    a.decisions = append(a.decisions, decision)
}

func TestAuthorizeHandler_AllowsAccessWithSystemActor(t *testing.T) {
    store := &stubStore{
        allowed:  true,
        policies: []string{"policy-1"},
    }
    auditor := &stubAuditor{}
    svc := NewService(store, auditor)

    payload := Request{
        SubjectID: "user-abc",
        Roles:     []string{"editor"},
        Action:    "article:update",
        Resource:  "article:456",
    }
    body, err := json.Marshal(payload)
    require.NoError(t, err)

    req := httptest.NewRequest(http.MethodPost, "/v1/authorize", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)

    err = svc.AuthorizeHandler(c)
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, rec.Code)

    var resp Response
    require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    require.True(t, resp.Allowed)
    require.Equal(t, "Access granted", resp.Message)
    require.Equal(t, store.policies, resp.PoliciesEvaluated)

    require.Len(t, auditor.decisions, 1)
    require.Equal(t, "system", auditor.decisions[0].Actor)
}

func TestAuthorizeHandler_DeniesAccessAndReturnsPolicies(t *testing.T) {
    store := &stubStore{
        allowed:  false,
        policies: []string{"viewer-only"},
    }
    auditor := &stubAuditor{}
    svc := NewService(store, auditor)

    payload := Request{
        SubjectID: "user-xyz",
        Roles:     []string{"viewer"},
        Action:    "article:update",
        Resource:  "article:456",
    }
    body, err := json.Marshal(payload)
    require.NoError(t, err)

    req := httptest.NewRequest(http.MethodPost, "/v1/authorize", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    req.Header.Set(echo.HeaderAuthorization, "token-123")
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)
    c.Set("actor", "token-123")

    err = svc.AuthorizeHandler(c)
    require.NoError(t, err)
    require.Equal(t, http.StatusForbidden, rec.Code)

    var resp Response
    require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    require.False(t, resp.Allowed)
    require.Equal(t, "Access denied", resp.Message)
    require.Equal(t, "subject lacks permission", resp.Reason)
    require.Equal(t, store.policies, resp.PoliciesEvaluated)

    require.Len(t, auditor.decisions, 1)
    require.Equal(t, "token-123", auditor.decisions[0].Actor)
}
