package authorization

import (
    "github.com/example/rbac-service/internal/audit"
    "github.com/example/rbac-service/internal/middleware"
    "github.com/example/rbac-service/pkg/policystore"
    "github.com/labstack/echo/v4"
    "net/http"
)

type Service struct {
    store    policystore.PolicyStore
    auditor  audit.Auditor
}

func NewService(store policystore.PolicyStore, auditor audit.Auditor) *Service {
    return &Service{store: store, auditor: auditor}
}

func (s *Service) AuthorizeHandler(c echo.Context) error {
    var req Request
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request payload",
        })
    }

    actor := middleware.ActorFromContext(c)
    allowed, policies := s.store.Evaluate(req.SubjectID, req.Roles, req.Action, req.Resource)

    resp := Response{
        Allowed:           allowed,
        PoliciesEvaluated: policies,
    }

    if allowed {
        resp.Message = "Access granted"
    } else {
        resp.Message = "Access denied"
        resp.Reason = "subject lacks permission"
    }

    s.auditor.Record(audit.Decision{SubjectID: req.SubjectID, Action: req.Action, Resource: req.Resource, Allowed: allowed, Actor: actor})

    if allowed {
        return c.JSON(http.StatusOK, resp)
    }

    return c.JSON(http.StatusForbidden, resp)
}
