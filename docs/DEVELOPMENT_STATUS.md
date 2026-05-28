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
- Added an Inbox AI human-review preview showing the detected lead, reviewed action mode, service/status, and next step before saving.
- Added Inbox AI review-completion feedback before saving reviewed lead actions.
- Added visible Inbox AI missing-field chips for reviewed lead name, WhatsApp, and status.
- Added a post-save Inbox AI link to open the reviewed lead detail directly.
- Added contextual lead detail feedback when a reviewed lead is opened from Inbox AI.
- Improved frontend UX around the receptionist daily workflow: actionable dashboard queue, clearer navigation, lead status counts, follow-up urgency groups, and a primary reviewed action in Inbox AI.
- Added optional UX demo seed data with five extra services and two fictional leads per commercial status, plus a Docker Compose seed helper script.
- Implemented MVP multi-currency baseline with clinic-level country/currency configuration, backend validation, service currency metadata, frontend formatting, and updated docs.
- Implemented backend-prioritized dashboard action queue endpoint and connected the Angular dashboard to it.
- Implemented backend AI generation metadata persistence for provider/model/feature/status/safety audit without storing raw prompts or conversations.
- Implemented reviewed Inbox AI metadata persistence on lead create/update and connected high-intent/objection signals to dashboard priority actions.
- Improved the Leads screen with a commercial detail panel that loads notes, follow-up context, and reviewed AI insights from `GET /api/leads/:id`.
- Added manual contact shortcuts in lead detail for opening WhatsApp and copying the lead phone number.
- Added quick post-contact actions in lead detail for common outcomes and support for clearing `next_action_at` when a lead leaves the follow-up queue.

## In-Progress Items

- Review the actionable dashboard, Inbox AI workflow, and lead detail panel with real local usage.

## Pending Items

- Define retention rules before storing full inbound conversations.
- Expand dashboard action scoring later with aging, follow-up history, and conversion outcome signals.
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

- Backend has implemented MVP modules and AI usage metadata, but content generation and richer AI error handling remain pending.
- Frontend screens currently cover first navigation, lead workflows, AI reply generation, follow-up actions, service catalog CRUD, and clinic profile editing.
- Smart Lead Inbox now has a first frontend route, conversation analysis endpoint, frontend actions to create or update leads from reviewed analysis, and reviewed AI metadata persistence. Dashboard now consumes a backend-prioritized action queue for follow-up, new-lead, high-intent, and objection signals. The Leads screen can show notes and reviewed AI insights in a commercial detail panel; remaining proposed API endpoints are not implemented.
- API contracts now reflect implemented MVP endpoints, but should continue to be updated when response shapes change.
- Migration execution tooling exists, but rollback/down migration support is not implemented yet.
- AI safety validation and basic usage metadata exist, but provider error taxonomy still needs improvement.
- Backend unit tests exist for implemented modules; e2e coverage should be expanded.
- Currency support is implemented for the MVP baseline; service prices now use string decimal API contracts and exact backend validation.

## Next Recommended Step

Review the updated Inbox AI to dashboard to lead-detail workflow locally and tune priority wording/order from real receptionist usage.

## Change Log

### 2026-05-27

- Added `GET /api/dashboard/actions` with backend-owned daily prioritization for overdue follow-ups, today's follow-ups, and new leads.
- Updated Angular `Atender hoy` to consume the dashboard action queue endpoint instead of assembling actions from separate leads and follow-ups calls.
- Expanded API service tests and smoke coverage for dashboard priority actions.
- Added `ai_generations` metadata persistence for AI feature, provider, model, status, safety status, user, clinic, and character counts without storing raw prompt or conversation content.
- Added `lead_ai_insights` metadata persistence for reviewed Inbox AI analysis and dashboard prioritization of high-intent leads and detected objections.
- Added a Leads commercial detail panel for selected lead notes, follow-up context, and reviewed AI insights.
- Added manual WhatsApp and phone copy actions to the lead detail panel without introducing WhatsApp API automation.
- Added post-contact quick actions to update lead status, notes, follow-up dates, and closed outcomes from the commercial detail panel.
- Added a recent commercial notes section in the lead detail panel, ordered by newest contact activity.
- Added a manual copy action for a concise commercial lead summary from the lead detail panel.
- Added dashboard queue summary counters for overdue, today, high-intent, objection, and new-lead actions.
- Added dashboard queue filters so staff can focus the visible action list by urgency or AI signal.
- Added quick commercial notes from dashboard action cards using the existing lead update endpoint.
- Added in-session saved-note indicators on dashboard action cards after a quick commercial note is stored.
- Added direct Dashboard links from AI insight actions to Inbox AI with lead context.

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
- Added a human-review preview in Inbox AI before saving analyzed lead actions.
- Added a review-completion indicator in Inbox AI based on the reviewed lead form validity.
- Added field-level review hints in Inbox AI when required reviewed lead data is missing or invalid.
- Added direct navigation from Inbox AI to the reviewed lead detail after create or update.
- Added an Inbox AI context message in the lead detail panel for reviewed lead handoff.
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
