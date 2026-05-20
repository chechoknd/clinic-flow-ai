# ClinicFlow AI — Agent Constitution

## 1. Project Identity

ClinicFlow AI is a SaaS platform focused on helping small and medium private clinics improve commercial attention, lead follow-up, and patient conversion through AI-assisted WhatsApp communication.

The first MVP focuses on dental clinics.

The product is **not** a hospital system, a clinical diagnosis system, a medical records system, or an autonomous WhatsApp bot in the MVP.

**Core business promise:**

> Help clinics respond better, recover interested patients, and sell more treatments through WhatsApp.

---

## 2. MVP Scope — Included

- Authentication and roles (superadmin, clinic_admin, assistant)
- Clinic profile and configuration
- Clinic service catalog
- Simple commercial CRM for leads
- AI reply assistant for WhatsApp
- AI objection handler
- Manual assisted follow-ups
- Basic dashboard
- AI content generator for marketing
- PostgreSQL database
- Docker Compose local environment
- Angular frontend
- Go backend
- REST API
- JWT authentication
- Environment-based AI provider configuration

---

## 3. MVP Scope — Strictly Excluded

If a task introduces any of these, stop and report the scope conflict:

- Clinical histories, medical records, diagnoses, prescriptions
- Clinical decision support, medical image interpretation
- Payment gateways, electronic invoicing
- Native mobile apps, telemedicine, video calls
- Full appointment scheduling system
- Native bidirectional WhatsApp Business Cloud API integration
- Autonomous WhatsApp bots
- RAG/vector search as required MVP infrastructure

---

## 4. Mandatory Reading Strategy

**Always read:**
- `AGENTS.md` (this file)
- `README.md`
- `docs/DEVELOPMENT_STATUS.md`

**Read depending on task:**
- Backend tasks → `docs/AGENT_RULES/BACKEND_RULES.md`
- Frontend tasks → `docs/AGENT_RULES/FRONTEND_RULES.md`
- Database tasks → `docs/AGENT_RULES/DATABASE_RULES.md`
- AI tasks → `docs/AGENT_RULES/AI_SAFETY_RULES.md`
- Security/auth/secrets tasks → `docs/AGENT_RULES/SECURITY_RULES.md`
- Git/branch/commit tasks → `docs/AGENT_RULES/GITFLOW_RULES.md`
- Documentation tasks → `docs/AGENT_RULES/DOCUMENTATION_RULES.md`
- Testing tasks → `docs/AGENT_RULES/TESTING_RULES.md`
- General agent behavior → `docs/AGENT_RULES/AGENT_OPERATING_RULES.md`

---

## 5. Technology Stack

**Frontend:** Angular, standalone components, Reactive Forms, Signals, Router, Guards, Interceptors, Tailwind CSS. No NgRx unless requested.

**Backend:** Go, REST API, Clean Architecture, JWT, PostgreSQL, provider abstraction for AI APIs (OpenAI, Gemini, DeepSeek).

**Database:** PostgreSQL, UUID primary keys, SQL migrations, seeds. pgvector prepared for future phases but not required for MVP.

**Infrastructure:** Docker Compose, environment variables, no secrets committed.

**AI:** Provider-agnostic backend interfaces. Prompt safety controlled by backend. No diagnosis or medical advice.

---

## 6. Repository Structure

```
clinic-flow-ai/
├── AGENTS.md
├── README.md
├── .gitignore
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
    ├── AGENT_RULES/           # Detailed rule documents
    ├── PROJECT_PLAN.md
    ├── ARCHITECTURE.md
    ├── DEVELOPMENT_STATUS.md
    ├── API_CONTRACTS.md
    └── DECISIONS_LOG.md
```

---

## 7. Non-Negotiable Safety Rules

- **No secrets committed.** Never commit `.env`, keys, or credentials.
- **No medical data.** Never store diagnoses, prescriptions, or clinical records.
- **Multi-tenant isolation.** Backend must enforce clinic-level data isolation.
- **Backend authorization.** JWT validation on every protected endpoint. Frontend guard alone is insufficient.
- **AI cannot diagnose.** The AI assistant is a commercial communication tool, not a doctor.
- **No WhatsApp Cloud API in the MVP.** No autonomous WhatsApp bots.
- **Never expose stack traces or SQL errors in API responses.**
- **Never log passwords, JWT tokens, or API keys.**

---

## 8. Git Workflow Summary

- Never work directly on `main`.
- Before starting: `git status`, `git branch --show-current`, `git pull`.
- Branch format: `feature/<desc>`, `fix/<desc>`, `chore/<desc>`, `docs/<desc>`, `refactor/<desc>`.
- Commit format: `type(scope): short description`. Types: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`, `style`, `build`, `ci`.
- Push: `git push -u origin HEAD`.
- PRs when available: include summary, files changed, testing, risks, follow-ups.

---

## 9. Agent Response Format

After every task, respond with:

```md
## Summary
Short explanation.

## Files Changed
- `path`: explanation.

## Validation
Commands executed / not executed with reasons.

## Risks or Notes
Any concern, limitation, or pending issue.

## Next Suggested Step
One clear next step.
```

---

## 10. References

Detailed rules are organized in `docs/AGENT_RULES/`:

| File | Contains |
|------|----------|
| `AGENT_OPERATING_RULES.md` | Golden rule, DoD, forbidden/preferred behaviors |
| `SECURITY_RULES.md` | Secrets, JWT, logging, data isolation |
| `AI_SAFETY_RULES.md` | AI constraints, prompt ownership, output validation |
| `BACKEND_RULES.md` | Go architecture, layers, coding rules, endpoints |
| `FRONTEND_RULES.md` | Angular structure, UX, accessibility, i18n |
| `DATABASE_RULES.md` | PostgreSQL conventions, migrations, table design |
| `GITFLOW_RULES.md` | Branch naming, commits, PRs |
| `DOCUMENTATION_RULES.md` | Docs maintenance, status tracking |
| `TESTING_RULES.md` | Backend/frontend test expectations |

---

## 11. Final Reminder

This is not a hospital system. This is not clinical diagnosis. This is not a WhatsApp bot in the MVP.

This is a commercial assistant for private clinics that helps human staff respond better, follow up faster, and convert more interested patients into booked appointments.

Every technical decision must support that goal.
