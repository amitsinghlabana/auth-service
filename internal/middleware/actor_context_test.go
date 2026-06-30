package middleware

import (
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/require"
)

func TestActorFromContext_ReturnsActorWhenPresent(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    c.Set(actorKey, "token-abc")

    actor := ActorFromContext(c)
    require.Equal(t, "token-abc", actor)
}

func TestActorFromContext_DefaultsToSystemWhenMissing(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    actor := ActorFromContext(c)
    require.Equal(t, "system", actor)
}