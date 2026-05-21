# Decisions Log

## 2026-05-21 — Use monorepo

### Context

ClinicFlow AI includes a frontend application, backend API, database migrations, seeds, Docker Compose configuration, and project documentation. These parts must evolve together during MVP development.

### Decision

Use a single monorepo for the Angular frontend, Go backend, database assets, infrastructure files, and documentation.

### Consequences

The monorepo keeps contracts, migrations, and application changes visible in one place. It also requires clear directory boundaries so frontend, backend, database, and documentation work do not become mixed.

## 2026-05-21 — Use Angular for frontend

### Context

The MVP needs a responsive web application for clinic administrators and reception staff. The product requires forms, protected routes, API integration, simple local state, and fast copy-to-clipboard workflows.

### Decision

Use Angular with standalone components, Reactive Forms, Signals, Angular Router, Guards, Interceptors, and Tailwind CSS.

### Consequences

Angular provides a structured frontend foundation for dashboard, clinic settings, service catalog, leads, AI assistant, content, and follow-up features. The frontend must remain simple and avoid unnecessary complexity such as NgRx unless explicitly needed later.

## 2026-05-21 — Use Go for backend

### Context

The backend must expose a REST API, enforce authentication and authorization, isolate tenant data, manage PostgreSQL access, and own AI provider integration and prompt safety.

### Decision

Use Go for the backend with Clean Architecture principles and module boundaries under `apps/backend-go/`.

### Consequences

Go supports a lightweight, performant API with explicit error handling and clear package boundaries. The team must keep handlers, services, repositories, models, DTOs, and routes separated to avoid mixing HTTP, business logic, SQL, and AI provider logic.

## 2026-05-21 — Use PostgreSQL for database

### Context

ClinicFlow AI needs reliable relational storage for clinics, users, services, leads, notes, follow-ups, and AI generation metadata. Multi-tenant isolation and query correctness are critical.

### Decision

Use PostgreSQL with UUID primary keys, SQL migrations, seeds, foreign keys, indexes, and timestamps where useful.

### Consequences

PostgreSQL gives a strong relational base for the MVP. Every clinic-owned entity must include tenant scoping where applicable. `pgvector` may be prepared for future phases, but it is not required for MVP logic.

## 2026-05-21 — Use Docker Compose for local development

### Context

Developers need a repeatable local environment for PostgreSQL and application services without manually configuring each dependency.

### Decision

Use Docker Compose for local development orchestration.

### Consequences

Docker Compose simplifies consistent setup across machines. Secrets must still be managed through local environment variables and must not be committed.

## 2026-05-21 — Use direct context injection instead of RAG for MVP

### Context

The MVP needs fast, safe, low-complexity AI assistance using clinic profile, service details, lead context, and objection context. A vector-search system would add infrastructure and operational complexity before product validation.

### Decision

Use direct structured context injection in backend-owned prompts for the MVP. Do not require RAG/vector search for MVP functionality.

### Consequences

This reduces cost, latency, and implementation complexity. Prompt context must remain concise and controlled. RAG can be revisited after MVP validation if clinics need richer knowledge bases.

## 2026-05-21 — Exclude WhatsApp Cloud API integration from MVP

### Context

Native bidirectional WhatsApp Business Cloud API integration introduces setup, compliance, webhook, automation, and operational complexity. The MVP can validate the core value through assisted copy-and-paste workflows.

### Decision

Exclude native WhatsApp Business Cloud API integration and autonomous WhatsApp bots from the MVP.

### Consequences

The MVP focuses on human-operated assistance: generate, copy, paste, and manually follow up through existing WhatsApp workflows. Automation can be evaluated after pilot validation.

## 2026-05-21 — Exclude medical records and diagnosis from MVP

### Context

The product is designed as a commercial communication assistant, not a hospital system or clinical system. Medical records, diagnosis, prescriptions, and clinical decision support create safety, regulatory, and scope risks.

### Decision

Exclude medical records, clinical histories, diagnoses, prescriptions, medical image interpretation, and clinical decision support from the MVP.

### Consequences

The product remains focused on sales, follow-up, marketing, and appointment-oriented communication. AI safety rules must prevent medical advice and direct patients toward professional evaluation when needed.

## 2026-05-21 — Use AGENTS.md plus docs/AGENT_RULES/ for agent governance

### Context

The project will likely be maintained with AI agent assistance. Agents need explicit product boundaries, technical rules, safety constraints, and documentation expectations.

### Decision

Use `AGENTS.md` as the central agent constitution and `docs/AGENT_RULES/` for detailed rules by topic.

### Consequences

Agents must read and follow the relevant rule files before making changes. Documentation, source code, database, AI, security, testing, and Git workflow decisions should remain aligned with these governance files.
