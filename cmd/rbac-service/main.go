package main

import (
    "log"

    "github.com/example/rbac-service/internal/audit"
    "github.com/example/rbac-service/internal/authorization"
    "github.com/example/rbac-service/internal/middleware"
    "github.com/example/rbac-service/pkg/policystore"
    echomiddleware "github.com/labstack/echo/v4/middleware"
    "github.com/labstack/echo/v4"
)

func main() {
    store := policystore.NewInMemoryStore()
    auditSvc := audit.NewLogger()
    authSvc := authorization.NewService(store, auditSvc)

    e := echo.New()
    e.Use(middleware.RequestID)
    e.Use(middleware.Recover())
    e.Use(middleware.AccessLogger)

    api := e.Group("/v1")
    api.Use(middleware.Authentication())
    api.POST("/authorize", authSvc.AuthorizeHandler)

    log.Println("starting RBAC authorization service on :8080")
    if err := e.Start(":8080"); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}