# ClinicFlow AI Backend

Minimal Go API skeleton for ClinicFlow AI.

## Requirements

- Go 1.22 or newer.
- PostgreSQL available locally through the repository Docker Compose setup.

## Configuration

- `HTTP_ADDR`: optional listen address. Defaults to `:8080`.
- `DATABASE_URL`: required for readiness checks, database connectivity, and auth repository access.
- `JWT_SECRET`: required to enable `POST /api/auth/login`. Use a strong secret outside version control.
- `JWT_EXPIRES_IN_SECONDS`: optional JWT lifetime in seconds. Defaults to `3600`.

Local development example:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable"
```

The local credentials above match `docker-compose.yml` and are for local development only. Do not commit production secrets or create `.env` files.

## Commands

Run the API locally with the default HTTP address:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" JWT_SECRET="local-dev-secret" go run ./cmd/api
```

Run the API locally on a custom port:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" JWT_SECRET="local-dev-secret" HTTP_ADDR=":18080" go run ./cmd/api
```

Run tests:

```bash
go test ./...
```

Format code:

```bash
gofmt -w ./...
```

## Health Endpoints

`GET /healthz` checks only that the API process is responding:

```bash
curl -sS http://127.0.0.1:18080/healthz
```

Expected response:

```json
{"status":"ok"}
```

`GET /readyz` checks PostgreSQL connectivity with `PingContext`:

```bash
curl -sS http://127.0.0.1:18080/readyz
```

Expected response when PostgreSQL is reachable:

```json
{"status":"ready"}
```

If `DATABASE_URL` is missing or PostgreSQL is unreachable, `/readyz` returns HTTP `503` with:

```json
{"status":"not_ready"}
```

## Current Scope

This skeleton implements the base `POST /api/auth/login` endpoint plus reusable JWT authentication and role-authorization middleware. It does not implement clinics, services, leads, AI provider logic, or business modules yet.
