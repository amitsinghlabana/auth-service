package authorization_test

import (
    "bytes"
    "encoding/json"
    "github.com/example/rbac-service/cmd/rbac-service"
    "github.com/example/rbac-service/internal/authorization"
    "github.com/example/rbac-service/internal/audit"
    "github.com/example/rbac-service/pkg/policystore"
    "github.com/labstack/echo/v4"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAuthorizationAllow(t *testing.T) {
    store := policystore.NewInMemoryStore()
    auditor := audit.NewLogger()
    svc := authorization.NewService(store, auditor)

    payload := authorization.Request{
        SubjectID: "user-123",
        Roles:     []string{"editor"},
        Action:    "article:update",
        Resource:  "article:456",
    }

    body, _ := json.Marshal(payload)
    req := httptest.NewRequest(http.MethodPost, "/v1/authorize", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    req.Header.Set(echo.HeaderAuthorization, "token-abc")

    e := echo.New()
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := svc.AuthorizeHandler(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)

    var resp authorization.Response
    assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    assert.True(t, resp.Allowed)
}

func TestAuthorizationDeny(t *testing.T) {
    store := policystore.NewInMemoryStore()
    auditor := audit.NewLogger()
    svc := authorization.NewService(store, auditor)

    payload := authorization.Request{
        SubjectID: "user-456",
        Roles:     []string{"viewer"},
        Action:    "article:update",
        Resource:  "article:456",
    }

    body, _ := json.Marshal(payload)
    req := httptest.NewRequest(http.MethodPost, "/v1/authorize", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    req.Header.Set(echo.HeaderAuthorization, "token-abc")

    e := echo.New()
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := svc.AuthorizeHandler(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusForbidden, rec.Code)

    var resp authorization.Response
    assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    assert.False(t, resp.Allowed)
}
