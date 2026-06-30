# OAuth‑Enabled Auth Service

## Overview

Our authentication service now supports processing OAuth tokens from external providers. It validates bearer tokens via the provider's introspection or userinfo endpoint, synchronizes user data with our internal model, and issues internal session tokens or JWTs.

## Running the Service

1. **Set required environment variables**:

   ```bash
   export OAUTH_PROVIDER_URL=https://provider.example.com
   export OAUTH_CLIENT_ID=your-client-id
   export OAUTH_CLIENT_SECRET=your-client-secret
   export SESSION_SIGNING_KEY=$(openssl rand -hex 32)
   ```

2. **Build and run**:

   ```bash
   go build -o bin/auth-service ./cmd/rbac-service
   ./bin/auth-service --config path/to/config.yaml
   ```

   Or, to run directly:

   ```bash
   go run ./cmd/rbac-service --config path/to/config.yaml
   ```

## API

### POST /auth/oauth

Validate an external OAuth token and issue an internal session token.

- **Headers**: `Authorization: Bearer <external-oauth-token>`
- **Request Body**: *none*

#### Responses

| Status | Body                                               | Description                         |
|:------:|:---------------------------------------------------|:------------------------------------|
| 200    | `{ "token": "<internal_jwt>", "expires_in": 3600 }` | Successfully issued session token   |
| 401    | `{ "error": "unauthorized" }`                   | Invalid or expired OAuth token      |
| 500    | `{ "error": "internal_server_error" }`          | Unexpected error validating token   |

## Testing

Run unit and integration tests across all modules:

```bash
go test ./internal/middleware ./internal/authorization ./pkg/policystore ./tests/authorization -v
```

Ensure that all tests pass before merging changes.