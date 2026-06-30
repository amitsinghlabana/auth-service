package middleware

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

const actorKey = "actor"

func Authentication() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            token := c.Request().Header.Get("Authorization")
            if token == "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error":   "missing_authorization_header",
                    "message": "Authorization header is required",
                })
            }

            c.Set(actorKey, token)
            return next(c)
        }
    }
}

func ActorFromContext(c echo.Context) string {
    actor, ok := c.Get(actorKey).(string)
    if !ok {
        return "system"
    }
    return actor
}
