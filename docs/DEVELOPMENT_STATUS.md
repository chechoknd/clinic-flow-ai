# Development Status

## Current Phase

Phase B — Smart Schedule First Slice Validation and Documentation Alignment

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
- Local end-to-end smoke script covers login, clinic profile update, service catalog CRUD, leads, follow-ups, dashboard, AI-assisted endpoints, professionals, schedule availability, appointments, lead-to-appointment conversion, and schedule dashboard summary with a deterministic mock AI provider.
- API contracts documentation updated to match the implemented backend/frontend response shapes and smoke coverage.
- Architecture documentation reconciled with the actual repository structure and Phase 1 completion state.
- Smart Lead Inbox product direction documented as planned in `docs/SMART_LEAD_INBOX_PLAN.md`.
- Dashboard Inteligente evolution documented as planned, focused on daily assistant actions instead of metrics only.
- Proposed Smart Lead Inbox API contracts documented as planned and not implemented.
- Implemented `POST /api/ai/analyze-conversation` for manual pasted conversation analysis.
- Added Angular Inbox AI screen for pasted conversation analysis, suggested reply copy, human-reviewed lead creation, and existing lead follow-up updates through the existing leads API.
- Smart Schedule database plan documented in `docs/SMART_SCHEDULE_DATABASE_PLAN.md` for professionals, professional-service assignment, appointments, indexes, overlap rules, and tenant isolation.
- Smart Schedule schema migration created and locally applied in `database/migrations/20260601000100_create_smart_schedule_schema.sql`.
- Backend professionals module implemented with clinic-scoped list, detail, create, and update endpoints.
- Backend appointments module implemented with clinic-scoped list, detail, create, update, status update, reschedule, overlap validation, and lead-to-appointment conversion.
- Backend schedule availability endpoint implemented through `GET /api/schedule/availability`.
- Backend schedule-centered dashboard summary endpoint implemented through `GET /api/dashboard/schedule-summary`.
- Angular Smart Schedule, professionals, and schedule-first dashboard screens exist and are wired to API service methods.
- Clinic profile now includes country/currency configuration for CO/COP, PE/PEN, AR/ARS, and CL/CLP.
- Angular follow-ups screen groups pending actions by overdue, today, and upcoming urgency while keeping manual human action.
- Angular frontend now has an initial lightweight ES/EN internationalization service, persisted language selection, and translated login/layout shell.
- Angular Dashboard, Smart Schedule, Leads, Follow-ups, Clinic, Professionals, and Services screens now use the initial ES/EN dictionary for primary labels, actions, status names, empty states, and operational messages.

## In-Progress Items

- Smart Schedule first slice is being validated end-to-end and reconciled in documentation.
- Inbox AI schedule integration remains in product planning and must keep human review before appointment creation.

## Pending Items

- Polish Angular daily/weekly Smart Schedule workflows with real receptionist usage.
- Review schedule availability and dashboard summary behavior against local demo data.
- Extend Inbox AI so it can suggest scheduling intent, missing appointment data, and human-reviewed appointment creation.
- Review the Smart Lead Inbox create/update lead workflow with real local usage.
- Align Inbox AI so it can detect scheduling intent, missing appointment data, and suggested appointment creation while requiring human review.
- Define retention rules before storing full inbound conversations.
- Define exact dashboard action scoring before coding priority queues, including today's appointments, pending confirmations, available slots, hot leads without appointment, and overdue follow-ups.
- Add automated seed runner if seed usage grows beyond local/demo data.
- Apply JWT middleware to any future protected endpoints as new modules are added.
- Add broader integration coverage for backend workflows and frontend tests when frontend code exists.
- Continue migrating the remaining feature screens to the ES/EN translation dictionary in small batches.

## Known Risks

- Scope creep into medical records, diagnosis, prescriptions, telemedicine, or hospital-system behavior.
- Smart Schedule scope creep into clinical scheduling, treatment planning, clinical notes, or autonomous booking.
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
- Smart Lead Inbox now has a first frontend route, conversation analysis endpoint, and frontend actions to create or update leads from reviewed analysis. Dashboard action layer, conversation analysis persistence, and remaining proposed API endpoints are not implemented.
- Smart Schedule database migration, professionals backend module, appointments backend module, availability endpoint, schedule dashboard summary, Angular schedule/professionals screens, and smoke-script coverage exist; local smoke execution against a running API and UX hardening remain pending.
- API contracts now reflect implemented MVP endpoints, but should continue to be updated when response shapes change.
- Migration execution tooling exists, but rollback/down migration support is not implemented yet.
- AI safety validation exists, but provider error taxonomy and usage tracking need improvement.
- Backend unit tests exist for implemented modules; e2e coverage should be expanded.

## Next Recommended Step

Validate the Smart Schedule flow end-to-end locally, then harden the Angular schedule UX and Inbox AI schedule handoff.

## Change Log

### 2026-06-03

- Added clinic country/currency fields to backend clinic DTOs, repository, service validation, and SQL migration.
- Added Angular clinic country/currency controls and updated frontend API models/specs.
- Improved the Angular follow-ups screen with urgency groups, readable date formatting, lead detail links, and split date/time rescheduling inputs.
- Validated backend tests, frontend build, frontend tests, backend vet, and local `/api/clinics/current` currency response.
- Started frontend internationalization with a reusable Signals-based ES/EN service, language persistence, translated login page, translated navigation labels, and tests.
- Extended the initial frontend internationalization pass to Dashboard, Leads, Smart Schedule, Follow-ups, Clinic, Professionals, and Services primary workflows.

### 2026-06-02

- Reconciled Smart Schedule documentation to mark the first slice as partially implemented instead of planned-only.
- Updated API contracts for implemented professionals, appointments, schedule availability, and schedule dashboard summary endpoints.
- Expanded `e2e_test.sh` to cover professionals, schedule availability, appointment create/list/reschedule/status update, lead-to-appointment conversion, and schedule dashboard summary.
- Fixed the Inbox AI frontend spec router provider so the full Angular test suite can run with `RouterLink`.
- Polished the Angular Smart Schedule screen with daily operational summary cards, safer empty states, Spanish status/source labels, service-based duration autofill, weekly date stepping, and frontend appointment payload alignment with `duration_minutes`.
- Improved Inbox AI to Smart Schedule handoff with reviewed schedule query params, commercial notes transfer, lead linkage after reviewed creation, and Schedule hydration from handoff params.

### 2026-06-01

- Started product realignment from AI-assisted commercial CRM toward a Smart Schedule-centered product for dental clinics.
- Planned Smart Schedule / Agenda Inteligente as the new core workflow connecting appointments, professionals, leads, services, follow-ups, and AI assistance.
- Documented that schedule data must remain commercial/admin only and must not introduce clinical histories, diagnoses, prescriptions, clinical notes, telemedicine, autonomous WhatsApp sending, or AI medical advice.
- Marked professionals/dentists, appointments, availability, lead-to-appointment conversion, and schedule-centered dashboard APIs as planned, not implemented.
- Added Smart Schedule database planning for `clinic_professionals`, `professional_services`, `appointments`, tenant isolation, overlap prevention, and lead-to-appointment conversion.
- Created the first Smart Schedule SQL migration for professionals, professional service assignment, appointments, indexes, triggers, status constraints, and schedule tenant relationships.
- Validated the migration against local PostgreSQL and confirmed backend tests still pass.
- Implemented the first Smart Schedule backend module for professionals with authenticated clinic-scoped endpoints.
- Implemented backend appointment endpoints with operational appointment CRUD, status changes, rescheduling, overlap checks, and lead conversion.

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
