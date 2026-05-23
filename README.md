# ClinicFlow AI

## Local Backend Quick Start

Prerequisites:

- Go 1.22 or newer.
- Docker and Docker Compose.

Start PostgreSQL:

```bash
docker compose up -d postgres
```

Run backend migrations:

```bash
cd apps/backend-go
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" go run ./cmd/migrate
cd ../..
```

Apply demo seed data from the repository root:

```bash
docker compose cp database/seeds/20260521000100_demo_core_data.sql postgres:/tmp/demo_core_data.sql
docker compose exec -T postgres psql -U clinicflow -d clinicflow_db -f /tmp/demo_core_data.sql
```

Run the backend API:

```bash
cd apps/backend-go
DATABASE_URL="postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" \
JWT_SECRET="at-least-32-characters-very-secret-key" \
HTTP_ADDR=":18080" \
ALLOWED_ORIGINS="*" \
AI_PROVIDER="deepseek" \
AI_MODEL="deepseek-chat" \
DEEPSEEK_API_KEY="local-placeholder-key" \
go run ./cmd/api
```

Validate the API:

```bash
curl -sS http://127.0.0.1:18080/healthz
curl -sS http://127.0.0.1:18080/readyz
```

For a fully local smoke test without external AI calls, run the API with `AI_PROVIDER=mock` and execute:

```bash
./e2e_test.sh
```

Demo login:

```txt
email: admin@sonrisaviva.demo
password: clinicflow123
```


## Local Frontend Quick Start

Prerequisites:

- Node.js 22 or newer.
- Backend API running on `http://127.0.0.1:18080` for authenticated API flows.

Run the Angular frontend:

```bash
cd apps/frontend-angular
npm install
npm start
```

Open `http://localhost:4200`.

## Documentation Map

- `AGENTS.md`: compact master instructions for AI agents.
- `docs/AGENT_RULES/`: detailed operational rules for agents.
- `docs/PROJECT_PLAN.md`: product plan and MVP roadmap.
- `docs/ARCHITECTURE.md`: technical architecture.
- `docs/API_CONTRACTS.md`: REST API contracts.
- `docs/DEVELOPMENT_STATUS.md`: current progress and pending work.
- `docs/DECISIONS_LOG.md`: architecture and product decisions.
