package middleware

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestAuthenticationMiddleware_AllowsAuthorizedRequest(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set(echo.HeaderAuthorization, "token-abc")
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    nextCalled := false
    handler := Authentication()(func(ctx echo.Context) error {
        nextCalled = true
        assert.Equal(t, "token-abc", ActorFromContext(ctx))
        return ctx.NoContent(http.StatusNoContent)
    })

    err := handler(c)
    assert.NoError(t, err)
    assert.True(t, nextCalled, "expected next handler to be invoked")
    assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAuthenticationMiddleware_MissingHeaderReturnsUnauthorized(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    handler := Authentication()(func(ctx echo.Context) error {
        t.Fatal("handler should not be invoked when Authorization header is missing")
        return nil
    })

    err := handler(c)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusUnauthorized, rec.Code)

    var payload map[string]string
    assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
    assert.Equal(t, "missing_authorization_header", payload["error"])
    assert.Equal(t, "Authorization header is required", payload["message"])
}
