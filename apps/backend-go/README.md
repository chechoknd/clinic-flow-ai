# ClinicFlow AI Backend

Go REST API for ClinicFlow AI.

## Requirements

- Go 1.22 or newer.
- PostgreSQL available locally through the repository Docker Compose setup.

## Configuration

- `HTTP_ADDR`: optional listen address. Defaults to `:8080`.
- `DATABASE_URL`: required for readiness checks, database connectivity, and auth repository access.
- `JWT_SECRET`: required to enable `POST /api/auth/login`. Use a strong secret outside version control.
- `JWT_EXPIRES_IN_SECONDS`: optional JWT lifetime in seconds. Defaults to `3600`.
- `ALLOWED_ORIGINS`: optional comma-separated CORS origins. Defaults to `*`.
- `AI_PROVIDER`: AI provider name. Supported values: `openai`, `gemini`, `deepseek`. Defaults to `openai`.
- `AI_MODEL`: optional model override for the selected AI provider.
- `OPENAI_API_KEY`, `GEMINI_API_KEY`, `DEEPSEEK_API_KEY`: provider keys. Use placeholders only for local validation when AI endpoints are not being called.

Local development database URL:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable"
```

The local credentials above match `docker-compose.yml` and are for local development only. Do not commit production secrets or create `.env` files.

## Local Development

Start PostgreSQL from the repository root:

```bash
docker compose up -d postgres
```

Run database migrations from `apps/backend-go`:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" go run ./cmd/migrate
```

Apply demo seed data from the repository root:

```bash
docker compose cp database/seeds/20260521000100_demo_core_data.sql postgres:/tmp/demo_core_data.sql
docker compose exec -T postgres psql -U clinicflow -d clinicflow_db -f /tmp/demo_core_data.sql
```

Demo login for local development:

```txt
email: admin@sonrisaviva.demo
password: clinicflow123
```

Run the API locally on port `18080` from `apps/backend-go`:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" \
JWT_SECRET="at-least-32-characters-very-secret-key" \
HTTP_ADDR=":18080" \
ALLOWED_ORIGINS="*" \
AI_PROVIDER="deepseek" \
AI_MODEL="deepseek-chat" \
DEEPSEEK_API_KEY="local-placeholder-key" \
go run ./cmd/api
```

## Commands

Run the API locally with the default HTTP address:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" JWT_SECRET="at-least-32-characters-very-secret-key" AI_PROVIDER="deepseek" DEEPSEEK_API_KEY="local-placeholder-key" go run ./cmd/api
```

Run the API locally on a custom port:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" JWT_SECRET="at-least-32-characters-very-secret-key" HTTP_ADDR=":18080" AI_PROVIDER="deepseek" DEEPSEEK_API_KEY="local-placeholder-key" go run ./cmd/api
```

Run database migrations from `apps/backend-go`:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" go run ./cmd/migrate
```

Use a custom migrations directory when needed:

```bash
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" go run ./cmd/migrate -dir ../../database/migrations
```

Run tests:

```bash
go test ./...
```

Run static validation:

```bash
go vet ./...
```

Format code:

```bash
gofmt -w $(find . -name "*.go")
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

Implemented backend scope currently includes health/readiness endpoints, JWT login, protected clinic profile endpoints, service catalog CRUD, commercial leads, and AI provider abstraction with safety validation.
