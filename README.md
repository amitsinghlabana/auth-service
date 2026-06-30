# RBAC Authorization Service

## Overview

The RBAC Authorization Service provides role- and permission-based access control for downstream services. It exposes a secure REST API for permission checks and policy administration, with accessibility support (screen-reader announcements and labeled fields for any UI components).

## Prerequisites

- Go 1.18 or later
- Git

## Setup & Run

```bash
# Clone the repo
git clone https://github.com/your-org/rbac-service.git
cd rbac-service

# Fetch dependencies
go mod download

# Build and run the service
go run cmd/rbac-service/main.go
```

By default, the service listens on port **8080**. You can override with the `PORT` environment variable:

```bash
PORT=9090 go run cmd/rbac-service/main.go
```

## API

### Health Check

```http
GET /health
```

- **200 OK** – Service is healthy.

### Authorization Check

```http
POST /api/v1/authorize
Content-Type: application/json

{
  "user_id": "alice",
  "action": "read",
  "resource": "documents:123"
}
```

**Responses**:

- **200 OK** – `{ "allowed": true }` or `{ "allowed": false }`
- **401 Unauthorized** – missing or invalid auth token
- **400 Bad Request** – malformed payload

### Policy Management

#### List Policies
```http
GET /api/v1/policies
Authorization: Bearer <admin-token>
```

- **200 OK** – JSON array of policies

#### Create/Update Policy
```http
POST /api/v1/policies
Content-Type: application/json
Authorization: Bearer <admin-token>

{
  "role": "editor",
  "permissions": [
    { "action": "update", "resource": "documents:*" }
  ]
}
```

- **201 Created** or **200 OK** – on success
- **403 Forbidden** – if current user lacks admin privileges

#### Delete Policy
```http
DELETE /api/v1/policies/{policy_id}
Authorization: Bearer <admin-token>
```

- **204 No Content** – on success

## Testing

Run all unit tests and linters with:

```bash
# Execute every test in the repo
go test ./...
```

Ensure code coverage and policy-watch behaviors by adjusting mocks in `pkg/policystore` and running specific test packages:

```bash
go test ./internal/authorization
go test ./internal/middleware
go test ./pkg/policystore
go test ./tests/authorization
```