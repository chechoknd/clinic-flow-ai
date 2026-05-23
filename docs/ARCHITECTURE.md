# Architecture

## Architecture Overview

ClinicFlow AI is planned as a monorepo SaaS application with an Angular frontend, Go backend, PostgreSQL database, and Docker Compose local environment.

The MVP architecture supports commercial clinic workflows: clinic configuration, service catalog management, lead CRM, AI-assisted WhatsApp replies, objection handling, manual follow-ups, marketing content generation, and basic dashboard metrics.

The backend is the authority for authentication, authorization, tenant isolation, prompt safety, provider selection, and API contracts. The frontend is responsible for fast, simple operator workflows, but it must not own security or AI safety rules.

## Monorepo Structure

Planned repository structure:

```txt
clinic-flow-ai/
├── AGENTS.md
├── README.md
├── .env.example
├── docker-compose.yml
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

## Frontend Architecture

Status: In Progress.

The frontend uses Angular with standalone components, Reactive Forms, Signals, Angular Router, route guards, HTTP interceptors, and Tailwind CSS.

Expected structure:

```txt
apps/frontend-angular/src/app/
├── core/
│   ├── auth/
│   ├── guards/
│   ├── interceptors/
│   ├── layout/
│   └── services/
├── shared/
│   ├── components/
│   ├── pipes/
│   ├── directives/
│   └── utils/
├── features/
│   ├── dashboard/
│   ├── clinics/
│   ├── services/
│   ├── leads/
│   ├── ai-assistant/
│   ├── content/
│   └── followups/
├── app.routes.ts
└── app.config.ts
```

Frontend responsibilities:

- Present authenticated application views.
- Use guards for private routes.
- Use interceptors to attach JWT authorization headers.
- Use services for API communication.
- Use Reactive Forms for clinic, service, lead, and AI input forms.
- Use Signals for simple local UI state.
- Keep workflows simple for receptionists and clinic administrators.
- Support fast copy-to-clipboard actions for WhatsApp-ready messages.

Frontend non-responsibilities:

- It must not enforce tenant isolation as the only security layer.
- It must not own critical AI prompt safety rules.
- It must not store API keys or secrets.
- It must not hardcode backend URLs outside environment configuration.

## Backend Architecture

Status: Planned.

The backend will use Go, REST APIs, Clean Architecture, JWT authentication, PostgreSQL persistence, and provider-agnostic AI integrations.

Expected structure:

```txt
apps/backend-go/
├── cmd/api/main.go
├── internal/
│   ├── auth/
│   ├── clinics/
│   ├── services/
│   ├── leads/
│   ├── ai/
│   ├── content/
│   ├── followups/
│   ├── dashboard/
│   └── shared/
├── pkg/
│   ├── database/
│   └── logger/
└── migrations/
```

Each backend module should follow:

```txt
handler.go
service.go
repository.go
model.go
dto.go
routes.go
```

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
- All protected endpoints validate JWT and enforce clinic-level isolation.

## Database Architecture

Status: Planned.

The database will use PostgreSQL, UUID primary keys, SQL migrations, seeds, foreign keys, indexes, and timestamp fields where useful.

Recommended base tables:

- `clinics`
- `users`
- `clinic_services`
- `leads`
- `lead_notes`
- `followups`
- `ai_generations`

Core tenant model:

- Clinic-owned entities must include `clinic_id` where applicable.
- Backend queries must scope data by authenticated user and clinic.
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

`pgvector` may be prepared for future phases, but it is not required for MVP functionality.

## AI Integration Strategy

Status: Planned.

The backend will expose provider-agnostic AI interfaces and select providers through environment variables. Supported provider targets include OpenAI, Gemini, DeepSeek, or compatible APIs.

MVP strategy:

- Use direct structured context injection.
- Avoid RAG/vector search as required MVP infrastructure.
- Keep prompts and safety rules backend-owned.
- Inject concise clinic, service, lead, and objection context into prompt templates.
- Validate AI output where possible before returning it.
- Return safe fallback messages when output appears risky.

Prompt context may include:

- Clinic name.
- Clinic type.
- City.
- Communication tone.
- Operating WhatsApp.
- Service details.
- Service benefits.
- Service FAQ.
- Common objections.
- Lead status.
- Optional commercial notes.

AI output must never diagnose, interpret symptoms, recommend medication, guarantee outcomes, replace professional evaluation, or decide that a patient does not need a clinic visit.

## Security Architecture

Status: Planned.

Security controls:

- JWT authentication for protected endpoints.
- Role-aware authorization for `superadmin`, `clinic_admin`, and `assistant`.
- Backend-enforced clinic-level data isolation.
- Secure password hashing.
- No secrets committed.
- Environment variables for secrets and provider keys.
- Sanitized API errors.
- Restricted logging.

Forbidden logging:

- Passwords.
- JWT tokens.
- API keys.
- Full patient conversations.
- Sensitive notes.
- Raw AI prompts containing private data outside controlled local debug scenarios.

## Multi-Tenancy Strategy

Status: Planned.

ClinicFlow AI uses logical multi-tenancy by clinic.

- Each user belongs to a clinic, except platform-level superadmin users.
- Each business entity belongs to a clinic where applicable.
- JWT claims should identify user identity, role, and clinic context.
- Repository queries must include clinic scoping.
- Service methods must reject cross-clinic access.
- API responses must never leak objects from another clinic.

## Local Development Infrastructure

Status: Planned.

Local development will use Docker Compose to provide reproducible services, especially PostgreSQL. The repository may include `.env.example` with placeholders only. Real `.env` files and secrets must not be committed.

Expected local services:

- PostgreSQL.
- Backend API.
- Angular development server.
- Optional database tooling if approved later.

## Environment Configuration

Environment variables should cover:

- Backend port and environment.
- Database connection settings.
- JWT secret and expiration.
- AI provider name.
- AI provider API key.
- AI model name.
- Provider timeout and token limits.
- Frontend API base URL.

Rules:

- Commit placeholders only in `.env.example`.
- Do not commit `.env`, local secrets, provider keys, JWT secrets, or credentials.
- Keep frontend secrets out of browser-delivered code.
- Backend controls AI provider credentials and prompt safety.

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

## Future Architecture Evolution

Potential post-MVP architecture additions:

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
