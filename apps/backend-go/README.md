# ClinicFlow AI Backend

Minimal Go API skeleton for ClinicFlow AI.

## Requirements

- Go 1.22 or newer.

## Commands

Run the API locally:

```bash
go run ./cmd/api
```

Run tests:

```bash
go test ./...
```

Format code:

```bash
gofmt -w ./cmd ./internal
```

## Health Endpoints

- `GET /healthz` returns process health.
- `GET /readyz` returns readiness for the current skeleton.

The server listens on `:8080` by default. Set `HTTP_ADDR` to override it, for example `HTTP_ADDR=:8081`.

## Current Scope

This skeleton does not implement authentication, clinics, services, leads, AI, database access, or business logic yet.
