package authorization

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type fakeStore struct {
    allowed   bool
    policies  []string
    lastInput Request
}

func (f *fakeStore) Evaluate(subjectID string, roles []string, action, resource string) (bool, []string) {
    f.lastInput = Request{
        SubjectID: subjectID,
        Roles:     roles,
        Action:    action,
        Resource:  resource,
    }
    return f.allowed, f.policies
}

func (f *fakeStore) RegisterPolicy(policy Policy) error {
    return errors.New("not implemented")
}

func (f *fakeStore) Policies() []Policy {
    return nil
}

type fakeAuditor struct {
    decisions []audit.Decision
}

func (f *fakeAuditor) Record(decision audit.Decision) {
    f.decisions = append(f.decisions, decision)
}

func TestAuthorizeHandler_ReturnsBadRequestOnBindError(t *testing.T) {
    store := &fakeStore{}
    auditor := &fakeAuditor{}
    svc := NewService(store, auditor)

    req := httptest.NewRequest(http.MethodPost, "/v1/authorize", bytes.NewReader([]byte("{invalid json")))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()

    e := echo.New()
    c := e.NewContext(req, rec)

    err := svc.AuthorizeHandler(c)
    require.NoError(t, err)
    assert.Equal(t, http.StatusBadRequest, rec.Code)

    var resp ErrorResponse
    assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    assert.Equal(t, "Invalid request payload", resp.Error)
    assert.Empty(t, auditor.decisions)
}

func TestAuthorizeHandler_RecordsAuditForDeniedRequests(t *testing.T) {
    store := &fakeStore{
        allowed:  false,
        policies: []string{"viewer-cannot-edit"},
    }
    auditor := &fakeAuditor{}
    svc := NewService(store, auditor)

    payload := Request{
        SubjectID: "user-456",
        Roles:     []string{"viewer"},
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
    c.Set("actor", "token-from-header")

    err = svc.AuthorizeHandler(c)
    require.NoError(t, err)
    assert.Equal(t, http.StatusForbidden, rec.Code)

    var resp Response
    assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    assert.False(t, resp.Allowed)
    assert.Equal(t, "Access denied", resp.Message)
    assert.Equal(t, "subject lacks permission", resp.Reason)
    assert.Equal(t, payload.SubjectID, payload.SubjectID)
    assert.Equal(t, payload.Resource, payload.Resource)

    require.Len(t, auditor.decisions, 1)
    decision := auditor.decisions[0]
    assert.Equal(t, payload.SubjectID, decision.SubjectID)
    assert.Equal(t, payload.Action, decision.Action)
    assert.Equal(t, payload.Resource, decision.Resource)
    assert.False(t, decision.Allowed)
    assert.Equal(t, "token-from-header", decision.Actor)
}

func TestAuthorizeHandler_ReturnsSuccessAndRecordsAudit(t *testing.T) {
    store := &fakeStore{
        allowed:  true,
        policies: []string{"editor-can-edit"},
    }
    auditor := &fakeAuditor{}
    svc := NewService(store, auditor)

    payload := Request{
        SubjectID: "user-123",
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
    c.Set("actor", "token-from-header")

    err = svc.AuthorizeHandler(c)
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)

    var resp Response
    assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    assert.True(t, resp.Allowed)
    assert.Equal(t, "Access granted", resp.Message)
    assert.Empty(t, resp.Reason)
    assert.Equal(t, []string{"editor-can-edit"}, resp.PoliciesEvaluated)

    require.Len(t, auditor.decisions, 1)
    decision := auditor.decisions[0]
    assert.True(t, decision.Allowed)
    assert.Equal(t, "token-from-header", decision.Actor)
    assert.Equal(t, payload.SubjectID, decision.SubjectID)
    assert.Equal(t, payload.Action, decision.Action)
    assert.Equal(t, payload.Resource, decision.Resource)
}

func TestAuthorizeHandler_DefaultsActorToSystemWhenMissing(t *testing.T) {
    store := &fakeStore{
        allowed:  false,
        policies: []string{"viewer-cannot-edit"},
    }
    auditor := &fakeAuditor{}
    svc := NewService(store, auditor)

    payload := Request{
        SubjectID: "user-789",
        Roles:     []string{"guest"},
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
    assert.Equal(t, http.StatusForbidden, rec.Code)

    require.Len(t, auditor.decisions, 1)
    decision := auditor.decisions[0]
    assert.Equal(t, "system", decision.Actor)
}
