# Development Status

## Current Phase

Phase 1 — Technical Base (Local Infrastructure)

## Completed Items

- Repository foundation exists.
- `AGENTS.md` defines project identity, scope, stack, safety rules, and agent governance.
- `docs/AGENT_RULES/` contains specialized rules for agent behavior, documentation, frontend, backend, database, AI safety, security, testing, and Git workflow.
- Core documentation files exist.
- MVP scope and exclusions are defined.
- PostgreSQL local infrastructure is configured through Docker Compose.
- Database documentation and structure initialized in `database/`.
- Backend Go module initialized in `apps/backend-go`.
- Backend API skeleton created with `GET /healthz` and `GET /readyz`.
- Backend database connection support added through `DATABASE_URL`.
- `/readyz` now validates PostgreSQL connectivity with a database ping.
- Initial core schema migration added for clinics, users, and clinic services.
- Base JWT login endpoint implemented for `POST /api/auth/login`.
- JWT authentication middleware and role authorization helper prepared for protected routes.
- Minimal Go migration runner added for ordered SQL migrations.
- Local demo seed added for one clinic admin login and sample dental services.
- Protected `GET /api/clinics/current` endpoint implemented using JWT claims.
- Protected `PUT /api/clinics/current` endpoint implemented for clinic administrators.
- Service Catalog module (CRUD) implemented with strict tenant isolation.
- Commercial Lead CRM (Leads Module) implemented with pagination and notes.
- AI Module implemented with provider abstraction, context injection, and safety validation.

## In-Progress Items

- Phase 1 technical base setup.
- Backend DB connection/readiness validation completed.
- Initializing Angular application shell.

## Pending Items

- Confirm final repository structure against the architecture docs.
- Initialize or verify Angular frontend base.
- Define `.env.example` placeholders.
- Add automated seed runner if seed usage grows beyond local/demo data.
- Apply JWT middleware to protected endpoints as business modules are implemented.
- Implement manual follow-up workflows.
- Implement dashboard summary.
- Add backend and frontend tests when code exists.

## Known Risks

- Scope creep into medical records, diagnosis, prescriptions, telemedicine, or hospital-system behavior.
- Premature WhatsApp Business Cloud API integration before validating manual assisted workflows.
- AI responses producing medical advice or unsafe claims if backend safety controls are weak.
- Multi-tenant data leakage if clinic scoping is not enforced in backend queries.
- Receptionist adoption risk if the UI requires too many steps or feels like a complex CRM.
- AI provider cost and reliability risk without token limits, timeouts, usage tracking, and fallback handling.

## Technical Debt

- Backend skeleton exists, but no business modules are implemented yet.
- No implemented frontend functionality is documented as complete yet.
- API contracts are initial planning contracts and must be updated during implementation.
- Migration execution tooling exists, but rollback/down migration support is not implemented yet.
- AI safety validators are planned but not implemented.
- Test suite status is pending until application code exists.

## Next Recommended Step

Implement `PUT /api/clinics/current` for clinic profile updates with `clinic_admin` authorization.

## Change Log

### 2026-05-20

- Implemented protected `GET /api/clinics/current` using JWT `clinic_id` claims.
- Added local demo seed data for clinic admin login and dental services.
- Added minimal Go migration runner for SQL files and `schema_migrations` tracking.
- Added reusable JWT authentication middleware and role authorization helper.
- Implemented base JWT login endpoint backed by the `users` table.
- Added initial PostgreSQL schema migration for `clinics`, `users`, and `clinic_services`.
- Connected backend readiness validation to PostgreSQL using `DATABASE_URL`.
- Prepared migration and seed README files for the upcoming auth, clinics, and users schema.
- Added backend Go API skeleton with health and readiness endpoints.
- Documented PostgreSQL local infrastructure as configured and backend skeleton as created.
- Populated initial project plan, architecture, API contracts, development status, and decisions log from the project specification and agent governance rules.
