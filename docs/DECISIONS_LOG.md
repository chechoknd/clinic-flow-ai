# Decisions Log

## 2026-05-20 — Use monorepo

### Context

ClinicFlow AI includes a frontend application, backend API, database migrations, seeds, Docker Compose configuration, and project documentation. These parts must evolve together during MVP development.

### Decision

Use a single monorepo for the Angular frontend, Go backend, database assets, infrastructure files, and documentation.

### Consequences

The monorepo keeps contracts, migrations, and application changes visible in one place. It also requires clear directory boundaries so frontend, backend, database, and documentation work do not become mixed.

## 2026-05-20 — Use Angular for frontend

### Context

The MVP needs a responsive web application for clinic administrators and reception staff. The product requires forms, protected routes, API integration, simple local state, and fast copy-to-clipboard workflows.

### Decision

Use Angular with standalone components, Reactive Forms, Signals, Angular Router, Guards, Interceptors, and Tailwind CSS.

### Consequences

Angular provides a structured frontend foundation for dashboard, clinic settings, service catalog, leads, AI assistant, content, and follow-up features. The frontend must remain simple and avoid unnecessary complexity such as NgRx unless explicitly needed later.

## 2026-05-20 — Use Go for backend

### Context

The backend must expose a REST API, enforce authentication and authorization, isolate tenant data, manage PostgreSQL access, and own AI provider integration and prompt safety.

### Decision

Use Go for the backend with Clean Architecture principles and module boundaries under `apps/backend-go/`.

### Consequences

Go supports a lightweight, performant API with explicit error handling and clear package boundaries. The team must keep handlers, services, repositories, models, DTOs, and routes separated to avoid mixing HTTP, business logic, SQL, and AI provider logic.

## 2026-05-20 — Use PostgreSQL for database

### Context

ClinicFlow AI needs reliable relational storage for clinics, users, services, leads, notes, follow-ups, and AI generation metadata. Multi-tenant isolation and query correctness are critical.

### Decision

Use PostgreSQL with UUID primary keys, SQL migrations, seeds, foreign keys, indexes, and timestamps where useful.

### Consequences

PostgreSQL gives a strong relational base for the MVP. Every clinic-owned entity must include tenant scoping where applicable. `pgvector` may be prepared for future phases, but it is not required for MVP logic.

## 2026-05-20 — Use Docker Compose for local development

### Context

Developers need a repeatable local environment for PostgreSQL and application services without manually configuring each dependency.

### Decision

Use Docker Compose for local development orchestration.

### Consequences

Docker Compose simplifies consistent setup across machines. Secrets must still be managed through local environment variables and must not be committed.

## 2026-05-20 — Use direct context injection instead of RAG for MVP

### Context

The MVP needs fast, safe, low-complexity AI assistance using clinic profile, service details, lead context, and objection context. A vector-search system would add infrastructure and operational complexity before product validation.

### Decision

Use direct structured context injection in backend-owned prompts for the MVP. Do not require RAG/vector search for MVP functionality.

### Consequences

This reduces cost, latency, and implementation complexity. Prompt context must remain concise and controlled. RAG can be revisited after MVP validation if clinics need richer knowledge bases.

## 2026-05-20 — Exclude WhatsApp Cloud API integration from MVP

### Context

Native bidirectional WhatsApp Business Cloud API integration introduces setup, compliance, webhook, automation, and operational complexity. The MVP can validate the core value through assisted copy-and-paste workflows.

### Decision

Exclude native WhatsApp Business Cloud API integration and autonomous WhatsApp bots from the MVP.

### Consequences

The MVP focuses on human-operated assistance: generate, copy, paste, and manually follow up through existing WhatsApp workflows. Automation can be evaluated after pilot validation.

## 2026-05-20 — Exclude medical records and diagnosis from MVP

### Context

The product is designed as a commercial communication assistant, not a hospital system or clinical system. Medical records, diagnosis, prescriptions, and clinical decision support create safety, regulatory, and scope risks.

### Decision

Exclude medical records, clinical histories, diagnoses, prescriptions, medical image interpretation, and clinical decision support from the MVP.

### Consequences

The product remains focused on sales, follow-up, marketing, and appointment-oriented communication. AI safety rules must prevent medical advice and direct patients toward professional evaluation when needed.

## 2026-05-20 — Use AGENTS.md plus docs/AGENT_RULES/ for agent governance

### Context

The project will likely be maintained with AI agent assistance. Agents need explicit product boundaries, technical rules, safety constraints, and documentation expectations.

### Decision

Use `AGENTS.md` as the central agent constitution and `docs/AGENT_RULES/` for detailed rules by topic.

### Consequences

Agents must read and follow the relevant rule files before making changes. Documentation, source code, database, AI, security, testing, and Git workflow decisions should remain aligned with these governance files.

## 2026-05-22 — Initialize Angular frontend shell

### Context

The backend already exposes core MVP endpoints for authentication, clinic profile, services, leads, follow-ups, dashboard, and AI assistance. The product needs a navigable frontend foundation before detailed workflows can be implemented.

### Decision

Initialize `apps/frontend-angular` as an Angular application with standalone components, lazy-loaded routes, Tailwind CSS, JWT auth shell, protected layout, and first MVP screens connected through API service wrappers.

### Consequences

The frontend can now evolve screen by screen against the Go API. The current shell prioritizes navigation, read flows, AI response generation, and copy-to-clipboard behavior; create/update forms and deeper error handling remain pending.

## 2026-05-24 — Plan Smart Lead Inbox before implementation

### Context

The current MVP validates AI-assisted commercial replies, objection handling, lead management, and manual follow-ups. The next product direction is to reduce manual copy/paste friction by letting staff paste complete conversations and receive structured commercial analysis, response drafts, and suggested next actions.

This direction introduces potential scope risks around message automation, full conversation storage, WhatsApp integration, clinical data, and AI overreach.

### Decision

Document Smart Lead Inbox / Inbox AI and Dashboard Inteligente as planned product evolution before implementing code.

The planned product rule is:

```txt
AI suggests -> human reviews -> human replies
```

ClinicFlow AI will continue to exclude autonomous WhatsApp bots, automatic message sending, diagnosis, prescriptions, clinical decision support, medical records, and required RAG/vector search from the immediate MVP extension.

Initial implementation, when approved, should start with manual or semi-manual conversation capture: paste conversation, analyze commercial context, review extracted data, copy suggested response, and optionally create/update lead or follow-up after human confirmation.

### Consequences

The project now has a documented direction toward an intelligent commercial inbox and a more actionable dashboard without changing the current implemented scope.

Future implementation must validate privacy and retention rules before storing raw conversations. Integrations such as n8n, Gmail, Meta Lead Ads, WhatsApp Business Cloud API, and browser extensions remain future options, not immediate dependencies.

## 2026-05-24 — Implement first Smart Lead Inbox slice

### Context

The Smart Lead Inbox direction was approved to start implementation after documentation alignment. The immediate product need is a manual conversation-analysis workflow that feels more useful than copying isolated patient messages into separate AI tools.

The first implementation must remain inside MVP guardrails: no autonomous WhatsApp bot, no automatic message sending, no required external integrations, no clinical records, and no diagnosis or prescription behavior.

### Decision

Implement the first Smart Lead Inbox slice as a manual, human-reviewed workflow:

- Add `POST /api/ai/analyze-conversation` for authenticated users.
- Analyze pasted commercial conversations through the backend AI provider abstraction.
- Return structured commercial suggestions: detected lead data, detected service, intent, objections, suggested reply, next action, optional follow-up date, and safety status.
- Add an Angular `Inbox AI` screen where staff can paste a conversation, review the AI result, copy the suggested reply, and create a lead through the existing leads API.
- Do not store raw conversations in this slice.
- Do not send messages automatically.

### Consequences

ClinicFlow AI now has a usable manual Inbox AI workflow without changing the product into an autonomous bot or clinical system.

Next implementation decisions should focus on whether to prioritize updating existing leads from analysis, persisting reviewed analysis history, or adding dashboard action cards.
