# Backend Rules

## Architecture

Backend lives in `apps/backend-go/`.

Expected structure:

```
apps/backend-go/
├── cmd/api/main.go
├── internal/
│   ├── auth/
│   ├── clinics/
│   ├── services/
│   ├── leads/
│   ├── ai/
│   ├── content/
│   ├── followups/
│   ├── dashboard/
│   └── shared/
├── pkg/
│   ├── database/
│   └── logger/
└── migrations/
```

Each module inside `internal/` should follow: `handler.go`, `service.go`, `repository.go`, `model.go`, `dto.go`, `routes.go`.

## Layer Responsibilities

### Handler Layer

Responsible for: HTTP request parsing, request validation, calling services, returning HTTP responses, mapping errors to status codes.

Not allowed: SQL queries, business rules, direct AI provider calls, password hashing, JWT generation logic.

### Service Layer

Responsible for: business rules, use case orchestration, validation beyond request shape, calling repositories, calling AI interfaces.

Not allowed: raw HTTP response writing, direct frontend dependency.

### Repository Layer

Responsible for: database persistence, SQL queries, transaction handling, mapping rows to models.

Not allowed: business decisions, HTTP concerns, AI prompt generation.

### Model Layer

Domain entities and internal business structures.

### DTO Layer

Request/response payloads and API contracts.

## Backend Coding Rules

- Use clear package boundaries.
- Keep functions small and focused.
- Return errors explicitly.
- Avoid global mutable state.
- Use context-aware DB calls.
- Validate inputs.
- Hash passwords securely.
- Never log secrets, JWT tokens, or API keys.
- Never expose stack traces in API responses.
- Keep AI provider logic behind interfaces.
- Use environment variables for secrets and provider keys.
- Prefer simple, readable Go.

Must not:
- Put all logic in `main.go`.
- Mix SQL into HTTP handlers.
- Hardcode credentials, API keys, clinic IDs, or user IDs.
- Generate diagnosis or medical advice logic.
- Add complex frameworks without approval.

## REST API Rules

All routes prefixed with `/api`. Initial routes:

```
POST /api/auth/login
GET /api/clinics/current
PUT /api/clinics/current
GET /api/services
POST /api/services
PUT /api/services/:id
DELETE /api/services/:id
GET /api/leads
POST /api/leads
GET /api/leads/:id
PUT /api/leads/:id
POST /api/ai/reply-suggestion
POST /api/ai/objection-handler
POST /api/content/generate-post
GET /api/dashboard/summary
```

## Error Handling

Consistent error response format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request payload is invalid."
  }
}
```

Do not expose: SQL errors, stack traces, provider secrets, internal file paths, raw panic messages.

## Roles

- `superadmin`: manage platform-level concerns.
- `clinic_admin`: manage clinic profile, users, services, leads.
- `assistant`: view/update leads, generate AI responses, manage follow-ups.

## Testing

Backend should include tests for: auth service, password hashing, JWT generation/validation, clinic service logic, lead service logic, AI prompt builder, AI safety validator, repository behavior.

Run: `go test ./...`, `go vet ./...`, `go fmt ./...`.
