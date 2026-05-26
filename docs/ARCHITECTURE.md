# Architecture

## Architecture Overview

ClinicFlow AI is a monorepo SaaS application with an Angular frontend, Go backend, PostgreSQL database, and Docker Compose local environment.

The current MVP architecture supports commercial clinic workflows: clinic configuration, service catalog management, lead CRM, AI-assisted WhatsApp replies, objection handling, manual follow-ups, and basic dashboard metrics.

The backend is the authority for authentication, authorization, tenant isolation, prompt safety, provider selection, and API contracts. The frontend is responsible for fast, simple operator workflows, but it must not own security or AI safety rules.

## Monorepo Structure

Status: Implemented.

Current top-level structure:

```txt
clinic-flow-ai/
├── AGENTS.md
├── README.md
├── .env.example
├── docker-compose.yml
├── e2e_test.sh
├── apps/
│   ├── frontend-angular/
│   └── backend-go/
├── database/
│   ├── migrations/
│   ├── seeds/
│   └── docs/
└── docs/
    ├── AGENT_RULES/
    ├── PROJECT_PLAN.md
    ├── ARCHITECTURE.md
    ├── API_CONTRACTS.md
    ├── DEVELOPMENT_STATUS.md
    └── DECISIONS_LOG.md
```

Generated local directories such as `node_modules/`, `dist/`, `.angular/`, and `.git/` are not architectural source directories.

## Frontend Architecture

Status: Implemented for the current MVP screens.

The frontend uses Angular with standalone components, Reactive Forms, Signals, Angular Router, route guards, HTTP interceptors, and Tailwind CSS.

Current source structure:

```txt
apps/frontend-angular/src/app/
├── core/
│   ├── auth/
│   ├── guards/
│   ├── interceptors/
│   ├── layout/
│   └── services/
├── features/
│   ├── ai-assistant/
│   ├── auth/
│   ├── clinics/
│   ├── dashboard/
│   ├── followups/
│   ├── leads/
│   └── services/
├── shared/
│   ├── components/
│   └── currency-format.ts
├── app.config.ts
├── app.routes.ts
└── app.ts
```

Implemented frontend responsibilities:

- Present authenticated application views.
- Use guards for private and guest routes.
- Use an interceptor to attach JWT authorization headers.
- Use `ApiService` for backend communication.
- Use Reactive Forms for login, clinic profile, service catalog, lead, follow-up, and AI input workflows.
- Use Signals for simple local UI state.
- Keep workflows simple for receptionists and clinic administrators.
- Support fast copy-to-clipboard actions for WhatsApp-ready AI responses.

Frontend non-responsibilities:

- It must not enforce tenant isolation as the only security layer.
- It must not own critical AI prompt safety rules.
- It must not store API keys or secrets.
- It must not hardcode backend URLs outside environment configuration.

Notes:

- The planned `content` feature is not implemented yet.
- `shared/` is intentionally light until repeated UI components justify more structure.

## Backend Architecture

Status: Implemented for current MVP modules.

The backend uses Go, REST APIs, JWT authentication, PostgreSQL persistence, and provider-agnostic AI integrations.

Current source structure:

```txt
apps/backend-go/
├── cmd/
│   ├── api/
│   └── migrate/
├── internal/
│   ├── ai/
│   ├── auth/
│   ├── clinics/
│   ├── config/
│   ├── dashboard/
│   ├── health/
│   ├── leads/
│   ├── migrations/
│   ├── services/
│   └── shared/
└── pkg/
    └── database/
```

Most business modules follow:

```txt
handler.go
service.go
repository.go
model.go
dto.go
routes.go
```

Current intentional deviations:

- Follow-up endpoints are implemented in `internal/leads` because follow-ups are modeled as leads with `next_action_at`.
- SQL migration execution support lives in `internal/migrations`; migration files live in `database/migrations`.
- `internal/content` is not implemented yet.
- `pkg/logger` is not present; logging currently uses standard Go logging.

Layer responsibilities:

- Handlers parse HTTP requests, validate request shape, call services, and map errors to API responses.
- Services implement business rules, authorization-aware use cases, and AI orchestration.
- Repositories handle SQL persistence, query scoping, transactions, and database mapping.
- Models represent internal business structures.
- DTOs represent request and response payloads.

Backend constraints:

- No SQL in HTTP handlers.
- No AI provider calls from handlers.
- No stack traces, SQL errors, or provider secrets in API responses.
- No diagnosis, prescription, or medical advice logic.
- All protected endpoints validate JWT and enforce clinic-level isolation in backend queries.

## Database Architecture

Status: Implemented for current MVP modules.

The database uses PostgreSQL, UUID primary keys, SQL migrations, seeds, foreign keys, indexes, and timestamp fields.

Implemented tables:

- `clinics`
- `users`
- `clinic_services`
- `leads`
- `lead_notes`
- `schema_migrations`

Current migration files live in `database/migrations/`. Demo data lives in `database/seeds/`.

Core tenant model:

- Clinic-owned entities include `clinic_id` where applicable.
- Backend queries scope data by authenticated user and clinic.
- Frontend filtering is not sufficient for tenant isolation.

Allowed lead data:

- Full name.
- Phone.
- Service of interest.
- Commercial status.
- Sales conversation notes.
- Next follow-up date.
- Basic source or channel.

Forbidden data:

- Diagnoses.
- Medical evolution notes.
- Clinical images.
- Prescriptions.
- Lab results.
- Medical records.
- Detailed symptoms stored as medical history.

Notes:

- Follow-ups are implemented through `leads.next_action_at`; there is no separate `followups` table.
- AI generation metadata is not persisted yet; there is no `ai_generations` table.
- `pgvector` may be prepared for future phases, but it is not required for MVP functionality.

## Currency Architecture

Status: Implemented for MVP baseline.

ClinicFlow AI needs multi-country commercial price display without becoming a payment, invoicing, or exchange-rate system. The approach is documented in `docs/CURRENCY_STRATEGY.md`.

Implemented architecture:

- Store default `country_code` and `currency_code` on each clinic.
- Validate values against a small backend-owned static currency catalog for the MVP.
- Keep service prices scoped to the clinic default currency.
- Do not perform automatic currency conversion or exchange-rate lookup.
- Expose currency metadata through clinic API responses so the frontend can format prices consistently.
- Expose service-level `currency_code` next to `price_from` for screens that render service prices.
- Use shared frontend formatting utilities instead of hardcoding symbols such as `COP`, `$`, or `S/` in feature templates.
- Keep monetary storage in PostgreSQL `NUMERIC(12,2)` for MVP compatibility, while avoiding new `float64`-based money handling in future backend code.

Initial planned currencies: `COP`, `PEN`, `ARS`, and `CLP`.

## Smart Lead Inbox Architecture

Status: Partially Implemented.

Smart Lead Inbox extends the current AI, leads, follow-up, and dashboard modules. The first implemented slice analyzes pasted conversations and turns them into reviewed commercial actions without storing raw conversations or sending messages automatically.

Implemented backend slice:

- `POST /api/ai/analyze-conversation` under `internal/ai`.
- Backend-owned prompt for commercial conversation analysis.
- Structured response with detected lead, service, intent, objections, summary, suggested reply, next action, and safety status.
- No raw conversation persistence.

Planned backend modules or extensions:

- Lead creation/update orchestration that only persists reviewed data.
- Optional inbound conversation/message storage, only after explicit retention and privacy rules are approved.
- Follow-up suggestion use case that proposes timing and next action without scheduling automatically.
- Dashboard action aggregation for daily priorities.

Implemented frontend slice:

- `features/inbox-ai/` route for manual conversation paste.
- Analysis result panel.
- Suggested reply copy-to-clipboard.
- Human-reviewed lead creation through the existing leads API.

Planned frontend modules or extensions:

- Dashboard action cards for daily assistant priorities.
- Review panels for extracted lead data, detected objections, suggested reply, and suggested follow-up.
- Copy-to-clipboard and explicit save/confirm buttons.

Architectural guardrails:

- No automatic message delivery.
- No native WhatsApp Cloud API dependency for the immediate MVP extension.
- No clinical records, diagnosis, prescriptions, or clinical decision support.
- Backend remains the source of truth for prompt safety, AI provider calls, auth, tenant isolation, and persistence.
- Frontend only presents and confirms reviewed actions.

## Planned Intelligent Dashboard Architecture

Status: Planned.

The dashboard is planned to evolve from a KPI summary into an action layer. It should help the receptionist answer `what should I handle now?` rather than only displaying totals.

Potential dashboard data sources:

- Lead status and creation time.
- `next_action_at` for overdue and due-today follow-ups.
- Service demand counts.
- AI conversation analysis outputs such as intent, objection type, and suggested next action when implemented.
- Human-reviewed attention state from the future inbox workflow.

Potential dashboard actions:

- Analyze conversation.
- Open Inbox AI.
- Create lead with AI.
- View overdue follow-ups.
- Review high-intent leads.
- Review detected objections.

This remains an assistant workflow. Dashboard recommendations must not trigger automatic patient messages.

## AI Integration Strategy

Status: Implemented for reply suggestions, objection handling, and manual follow-up messages.

The backend exposes provider-agnostic AI interfaces and selects providers through environment variables. Supported provider targets:

- `openai`
- `gemini`
- `deepseek`
- `mock` for deterministic local smoke testing only

Current MVP strategy:

- Use direct structured context injection.
- Avoid RAG/vector search as required MVP infrastructure.
- Keep prompts and safety rules backend-owned.
- Inject concise clinic, service, lead, and objection context into prompt templates.
- Validate AI output before returning it.
- Return structured JSON DTOs to the frontend.
- Planned Smart Lead Inbox conversation analysis will continue using backend-owned prompts, direct context injection, and structured JSON output. It must not require RAG/vector search for the immediate MVP extension.

Prompt context may include:

- Clinic name.
- Clinic type.
- City.
- Communication tone.
- Service details.
- Service benefits.
- Service FAQ.
- Common objections.
- Lead notes.

AI output must never diagnose, interpret symptoms, recommend medication, guarantee outcomes, replace professional evaluation, or decide that a patient does not need a clinic visit.

## Security Architecture

Status: Implemented for current protected endpoints.

Security controls:

- JWT authentication for protected endpoints.
- Role-aware authorization for clinic profile and service catalog writes.
- Backend-enforced clinic-level data isolation.
- Secure password hashing.
- No secrets committed.
- Environment variables for secrets and provider keys.
- Sanitized API errors.
- Basic security headers and CORS middleware.
- Request body size limits and simple rate limiting middleware.

Forbidden logging:

- Passwords.
- JWT tokens.
- API keys.
- Full patient conversations.
- Sensitive notes.
- Raw AI prompts containing private data outside controlled local debug scenarios.

## Multi-Tenancy Strategy

Status: Implemented for current clinic-scoped modules.

ClinicFlow AI uses logical multi-tenancy by clinic.

- Each user belongs to a clinic, except future platform-level superadmin users.
- Each business entity belongs to a clinic where applicable.
- JWT claims identify user identity, role, and clinic context.
- Repository queries include clinic scoping.
- Service methods reject cross-clinic access through scoped repository operations.
- API responses must never leak objects from another clinic.

## Local Development Infrastructure

Status: Implemented.

Local development uses Docker Compose for PostgreSQL plus locally run Go and Angular processes.

Current local services and commands:

- PostgreSQL: `docker compose up -d postgres`
- Backend migrations: `go run ./cmd/migrate` from `apps/backend-go`
- Backend API: `go run ./cmd/api` from `apps/backend-go`
- Angular dev server: `npm start` from `apps/frontend-angular`
- Local API smoke: `./e2e_test.sh` with backend `AI_PROVIDER=mock`

## Environment Configuration

Environment variables cover:

- Backend bind address: `HTTP_ADDR`.
- Database connection: `DATABASE_URL`.
- JWT secret and expiration: `JWT_SECRET`, `JWT_EXPIRES_IN_SECONDS`.
- CORS: `ALLOWED_ORIGINS`.
- AI provider name and model: `AI_PROVIDER`, `AI_MODEL`.
- AI provider API keys: `OPENAI_API_KEY`, `GEMINI_API_KEY`, `DEEPSEEK_API_KEY`.
- Frontend API base URL through Angular environment configuration.

Rules:

- Commit placeholders only in `.env.example`.
- Do not commit `.env`, local secrets, provider keys, JWT secrets, or credentials.
- Keep frontend secrets out of browser-delivered code.
- Backend controls AI provider credentials and prompt safety.

## Phase 1 Completion Notes

Phase 1 local technical base is functionally complete for the implemented MVP surface:

- PostgreSQL local infrastructure, migrations, and seed data are in place.
- Backend modules are implemented for auth, clinic profile, services, leads, follow-ups, AI, dashboard, health, and readiness.
- Frontend authenticated shell and MVP workflows are connected to the API.
- Local smoke coverage validates the primary API workflows against PostgreSQL using deterministic AI output.
- API contracts and development status documentation reflect the implemented behavior.

Remaining Phase 1 follow-ups:

- Resolve GitHub remote authentication so commits can be pushed.
- Decide whether to add an automated seed runner beyond the current SQL seed flow.
- Keep expanding integration/e2e coverage as new modules are added.
- Keep content generation documented as planned, not implemented, until the module is built.

## Out-of-Scope Architecture for MVP

The MVP must not include architecture for:

- Clinical histories or medical records.
- Diagnosis or clinical decision support.
- Prescription generation.
- Medical image interpretation.
- Payment gateways.
- Electronic invoicing.
- Native mobile apps.
- Telemedicine or video calls.
- Full appointment scheduling.
- Native bidirectional WhatsApp Business Cloud API integration.
- Autonomous WhatsApp bots.
- Required RAG/vector-search infrastructure.

## Planned Integration Boundaries

Status: Future.

External sources may be considered after manual Smart Lead Inbox validation:

- n8n for optional workflow experiments.
- Gmail for email lead capture.
- Web forms.
- Meta Lead Ads.
- WhatsApp Business Cloud API.
- Browser extension for capturing selected text from WhatsApp Web or Gmail.

These integrations must not become immediate dependencies for the MVP extension. The first implementation path should support manual or semi-manual capture.

## Future Architecture Evolution

Potential post-MVP architecture additions:

- Smart Lead Inbox and manual conversation analysis.
- Intelligent dashboard action layer.
- WhatsApp Business Cloud API integration.
- Event-driven follow-up automation.
- RAG/vector search for clinic-specific knowledge bases.
- Advanced analytics and funnel attribution.
- Billing and subscription services.
- Multi-location clinic hierarchy.
- Chrome extension or WhatsApp Web companion.
- Queue-based AI workloads.
- Usage metering and plan limits.
- Additional clinical-adjacent vertical configurations while preserving non-diagnostic safety boundaries.
