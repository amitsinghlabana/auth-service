# OAuth-Enabled Auth Service Design

## Overview
Implement an OAuth token processing service using Go and the standard `net/http` stack plus minimal, vetted dependencies (e.g., `github.com/google/uuid` or `golang.org/x/oauth2` only if truly necessary) to stay aligned with [S3]. Go suits the service due to its efficient concurrency, static binary, and strong standard library support for networking and crypto. The service will accept OAuth tokens from trusted external providers, validate them, fetch user metadata, and issue our internal session tokens or user context.

## Architecture
- **API gateway**: Exposes REST endpoints, validates incoming request structure, enforces authorization, and handles error translation to keep responses compliant with [S2].
- **OAuth verifier**: Validates incoming bearer tokens by calling the OAuth provider’s introspection/userinfo endpoint, caching JWKS, and ensuring tokens are not expired or revoked.
- **User sync**: Reconciles the verified OAuth identity with our internal user model, creating or updating users as needed.
- **Session manager**: Issues internal JWTs or session tokens from secrets kept strictly in environment variables per [S1].

## Component Breakdown
1. **cmd/server/main.go** – Entry point that initializes config, logger, OAuth client, and HTTP server.
2. **internal/config/config.go** – Loads configuration from environment variables (OAuth provider URL, client ID/secret, JWT signing key, cache settings) with no embedded secrets [S1].
3. **internal/server/server.go** – Sets up router (using Go `http.ServeMux` or minimal router) with middleware for logging, metrics, and error handling.
4. **internal/handlers/oauth.go** – Handles `/oauth/verify` (POST) and `/oauth/refresh` endpoints, managing payload validation and response formatting.
5. **internal/oauth/verifier.go** – Encapsulates OAuth provider interactions: token introspection/userinfo calls, JWK fetching, signature verification, caching.
6. **internal/session/session.go** – Generates signed JWTs for our ecosystem with configurable algorithm and expiration pulled from env vars.
7. **internal/errors/errors.go** – Centralized error types ensuring authentication failures return the generic `401 "invalid email or password"` message no matter the root cause [S2].
8. **internal/log/log.go** – Provides structured logging wrapper.
9. **tests/** – Unit tests for handlers/verifier/session ensuring correct logic and security.

## API Contract
| Endpoint | Method | Request | Response | Notes |
| --- | --- | --- | --- | --- |
| `/oauth/verify` | POST | `{ "token": "<bearer token>" }` | `200 { "user_id": "<uuid>", "roles": [...], "session_token": "<jwt>" }` | Rejects malformed requests early; uses generic 401 for invalid credentials [S2]. |
| `/oauth/refresh` | POST | `{ "refresh_token": "..." }` | `200 { "session_token": "<jwt>" }` | Optional refresh; depends on session management policy. |

All secret configuration (OAuth client secret, JWT signing key) is provided via environment variables and never leaked back to clients [S1].

## Sequence Diagram
```mermaid
sequenceDiagram
    participant Client
    participant API
    participant OAuthProvider

    Client->>API: POST /oauth/verify { token }
    API->>OAuthProvider: Validate token/introspect
    OAuthProvider-->>API: Token details
    API->>API: Map to internal user; issue JWT
    API-->>Client: 200 + session_token
```

## Project structure
```
auth-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── errors/
│   │   └── errors.go
│   ├── handlers/
│   │   └── oauth.go
│   ├── log/
│   │   └── log.go
│   ├── oauth/
│   │   └── verifier.go
│   ├── server/
│   │   └── server.go
│   └── session/
│       └── session.go
├── tests/
│   ├── handlers/
│   │   └── oauth_test.go
│   ├── oauth/
│   │   └── verifier_test.go
│   └── session/
│       └── session_test.go
├── go.mod
├── go.sum
└── README.md
```
