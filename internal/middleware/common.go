package middleware

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
    return middleware.RequestID()(next)
}

func Recover() echo.MiddlewareFunc {
    return middleware.Recover()
}

func AccessLogger(next echo.HandlerFunc) echo.HandlerFunc {
    return middleware.Logger()(next)
}
