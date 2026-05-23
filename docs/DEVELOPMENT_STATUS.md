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
- Manual follow-up workflows implemented using lead `next_action_at`, including pending list, complete, reschedule, and AI follow-up message generation.
- Dashboard summary endpoint implemented with lead totals, status counts, top services, follow-up counts, and conversion rate.
- Local backend quick start documented in root and backend README files.
- Angular frontend base initialized in `apps/frontend-angular` with standalone routing and Tailwind CSS.
- Frontend authentication shell implemented with login page, JWT storage, auth guard, guest guard, and HTTP auth interceptor.
- Initial private frontend layout and lazy-loaded MVP screens added for dashboard, leads, follow-ups, services, clinic profile, and AI assistant.
- Frontend lead workflow implemented for create, status update, notes, AI reply copy, and follow-up scheduling.
- Frontend follow-up screen connected to complete, reschedule, and AI follow-up message workflows.
- Frontend service catalog workflow implemented for create, update, active state, and delete actions.
- Frontend clinic profile workflow implemented for editing commercial profile and communication tone.
- Local end-to-end smoke script covers login, clinic profile update, service catalog CRUD, leads, follow-ups, dashboard, and AI-assisted endpoints with a deterministic mock AI provider.
- API contracts documentation updated to match the implemented backend/frontend response shapes and smoke coverage.

## In-Progress Items

- Phase 1 technical base setup.
- Backend DB connection/readiness validation completed.
- Frontend MVP screens connected to primary API workflows and API error states.
- Local API smoke validated against PostgreSQL, seed data, JWT auth, and mock AI provider.

## Pending Items

- Confirm final repository structure against the architecture docs and prepare Phase 1 completion notes.
- Add automated seed runner if seed usage grows beyond local/demo data.
- Apply JWT middleware to protected endpoints as business modules are implemented.
- Add broader integration coverage for backend workflows and frontend tests when frontend code exists.

## Known Risks

- Scope creep into medical records, diagnosis, prescriptions, telemedicine, or hospital-system behavior.
- Premature WhatsApp Business Cloud API integration before validating manual assisted workflows.
- AI responses producing medical advice or unsafe claims if backend safety controls are weak.
- Multi-tenant data leakage if clinic scoping is not enforced in backend queries.
- Receptionist adoption risk if the UI requires too many steps or feels like a complex CRM.
- AI provider cost and reliability risk without token limits, timeouts, usage tracking, and fallback handling.

## Technical Debt

- Backend has implemented MVP modules, but content generation, AI usage metadata, and richer AI error handling remain pending.
- Frontend screens currently cover first navigation, lead workflows, AI reply generation, follow-up actions, service catalog CRUD, and clinic profile editing.
- API contracts now reflect implemented MVP endpoints, but should continue to be updated when response shapes change.
- Migration execution tooling exists, but rollback/down migration support is not implemented yet.
- AI safety validation exists, but provider error taxonomy and usage tracking need improvement.
- Backend unit tests exist for implemented modules; e2e coverage should be expanded.

## Next Recommended Step

Confirm final repository structure against `docs/ARCHITECTURE.md` and prepare Phase 1 completion notes.

## Change Log

### 2026-05-23

- Added lead management frontend workflow with create lead, status updates, notes, AI reply generation, and follow-up scheduling.
- Added follow-up frontend actions for complete, reschedule, and AI follow-up message generation.
- Added service catalog frontend CRUD workflow.
- Added clinic profile frontend update workflow.
- Added deterministic mock AI provider for local smoke testing.
- Expanded `e2e_test.sh` to validate login, clinic profile update, services, leads, follow-ups, dashboard, and AI-assisted endpoints.
- Updated `docs/API_CONTRACTS.md` to describe implemented endpoint contracts and planned-but-not-implemented APIs.
- Validated frontend with `npm test -- --watch=false`, `npm run lint`, and `npm run build`.
- Validated backend smoke with `go test ./...`, `go vet ./...`, migrations, demo seed, readiness check, and `./e2e_test.sh`.

### 2026-05-22

- Initialized Angular frontend application with Tailwind CSS.
- Added frontend auth shell, protected layout, lazy routes, and first MVP screens.
- Validated frontend with `npm run build` and `npm test -- --watch=false`.

### 2026-05-21

- Documented local backend quick start commands for PostgreSQL, migrations, seeds, API execution, and health checks.

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
