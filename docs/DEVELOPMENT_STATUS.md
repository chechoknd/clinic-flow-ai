# Development Status

## Current Phase

Phase B — Manual Conversation Analysis

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
- Local end-to-end smoke script covers login, clinic profile update, service catalog CRUD, currency/price validation, leads, follow-ups, dashboard, and AI-assisted endpoints with a deterministic mock AI provider.
- API contracts documentation updated to match the implemented backend/frontend response shapes and smoke coverage.
- Architecture documentation reconciled with the actual repository structure and Phase 1 completion state.
- Smart Lead Inbox product direction documented as planned in `docs/SMART_LEAD_INBOX_PLAN.md`.
- Dashboard Inteligente evolution documented as planned, focused on daily assistant actions instead of metrics only.
- Proposed Smart Lead Inbox API contracts documented as planned and not implemented.
- Implemented `POST /api/ai/analyze-conversation` for manual pasted conversation analysis.
- Added Angular Inbox AI screen for pasted conversation analysis, suggested reply copy, human-reviewed lead creation, and existing lead follow-up updates through the existing leads API.
- Improved frontend UX around the receptionist daily workflow: actionable dashboard queue, clearer navigation, lead status counts, follow-up urgency groups, and a primary reviewed action in Inbox AI.
- Added optional UX demo seed data with five extra services and two fictional leads per commercial status, plus a Docker Compose seed helper script.
- Implemented MVP multi-currency baseline with clinic-level country/currency configuration, backend validation, service currency metadata, frontend formatting, and updated docs.

## In-Progress Items

- Phase B manual conversation analysis and first daily-workflow UX enhancements are implemented and ready for local review.

## Pending Items

- Review the actionable dashboard, Inbox AI create/update lead workflow, and urgency grouping with real local usage.
- Define retention rules before storing full inbound conversations.
- Define exact dashboard action scoring before coding priority queues.
- Keep seed data fictional and update the demo seed helper if new demo datasets are added.
- Apply JWT middleware to any future protected endpoints as new modules are added.
- Add broader integration coverage for backend workflows and frontend tests when frontend code exists.

## Known Risks

- Scope creep into medical records, diagnosis, prescriptions, telemedicine, or hospital-system behavior.
- Premature WhatsApp Business Cloud API integration before validating manual assisted workflows.
- Smart Lead Inbox scope creep into autonomous bots or automatic message sending.
- Full conversation storage creating privacy risk if retention rules are not defined.
- AI responses producing medical advice or unsafe claims if backend safety controls are weak.
- Multi-tenant data leakage if clinic scoping is not enforced in backend queries.
- Receptionist adoption risk if the UI requires too many steps or feels like a complex CRM.
- AI provider cost and reliability risk without token limits, timeouts, usage tracking, and fallback handling.

## Technical Debt

- Backend has implemented MVP modules, but content generation, AI usage metadata, and richer AI error handling remain pending.
- Frontend screens currently cover first navigation, lead workflows, AI reply generation, follow-up actions, service catalog CRUD, and clinic profile editing.
- Smart Lead Inbox now has a first frontend route, conversation analysis endpoint, and frontend actions to create or update leads from reviewed analysis. Dashboard now has a frontend-computed action queue from current leads/follow-ups; conversation analysis persistence and remaining proposed API endpoints are not implemented.
- API contracts now reflect implemented MVP endpoints, but should continue to be updated when response shapes change.
- Migration execution tooling exists, but rollback/down migration support is not implemented yet.
- AI safety validation exists, but provider error taxonomy and usage tracking need improvement.
- Backend unit tests exist for implemented modules; e2e coverage should be expanded.
- Currency support is implemented for the MVP baseline; service prices now use string decimal API contracts and exact backend validation.

## Next Recommended Step

Run a local product review of the updated receptionist workflow, then decide whether the dashboard action queue needs a dedicated backend endpoint for richer scoring.

## Change Log

### 2026-05-25

- Improved the daily receptionist workflow in Angular without adding new backend endpoints.
- Renamed the main dashboard navigation to `Atender hoy` and added a dashboard action queue for overdue follow-ups, today's follow-ups, and new leads.
- Added lead status counts, lead urgency highlights, and quick response actions in the Leads screen.
- Grouped Follow-ups into `Vencidos`, `Hoy`, and `Proximos` for clearer prioritization.
- Simplified Inbox AI save actions around one primary human-reviewed action based on whether an existing lead is selected.
- Added optional UX demo seed data and `database/seeds/apply_demo_seeds.sh` to populate more services and leads for local review.
- Implemented clinic-level currency support for Colombia, Peru, Argentina, and Chile without exchange rates, payments, or invoicing.
- Replaced service `price_from` float DTO handling with string decimal request/response contracts and exact backend validation.
- Expanded `e2e_test.sh` to validate COP decimal rejection, PEN decimal acceptance, exact string `price_from` responses, and clinic currency restoration.

### 2026-05-24

- Documented planned Smart Lead Inbox / Inbox AI direction.
- Implemented first Smart Lead Inbox slice: manual conversation analysis endpoint and Angular Inbox AI screen.
- Enhanced Inbox AI so reviewed analysis can create a new lead or update an existing lead follow-up from the same screen.
- Documented planned Dashboard Inteligente action layer.
- Added proposed API contracts for conversation analysis, lead creation from conversation, optional inbound messages, and follow-up suggestions.
- Recorded human-review requirement: AI suggests, human reviews, human replies.
- Confirmed autonomous WhatsApp bots, automatic message sending, clinical records, diagnosis, prescriptions, and required RAG remain out of scope.

### 2026-05-23

- Added lead management frontend workflow with create lead, status updates, notes, AI reply generation, and follow-up scheduling.
- Added follow-up frontend actions for complete, reschedule, and AI follow-up message generation.
- Added service catalog frontend CRUD workflow.
- Added clinic profile frontend update workflow.
- Added deterministic mock AI provider for local smoke testing.
- Expanded `e2e_test.sh` to validate login, clinic profile update, services, leads, follow-ups, dashboard, and AI-assisted endpoints.
- Updated `docs/API_CONTRACTS.md` to describe implemented endpoint contracts and planned-but-not-implemented APIs.
- Updated `docs/ARCHITECTURE.md` to match the actual repository structure and Phase 1 implementation notes.
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
