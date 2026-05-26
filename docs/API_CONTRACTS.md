# API Contracts

## Status

Status: Implemented for the current MVP backend/frontend surface.

This document describes the REST API currently implemented in the Go backend and consumed by the Angular frontend. Keep it updated when endpoints, payloads, authentication, pagination, roles, or error response shapes change.

## API Conventions

- Base path: `/api`
- Format: JSON request and response bodies.
- Authentication: JWT Bearer token for protected endpoints.
- IDs: UUID strings.
- Dates and times: ISO 8601 strings. Backend responses may include timezone offsets.
- Naming: JSON fields use `snake_case`.
- Tenant isolation: protected endpoints scope data to the authenticated clinic through JWT `clinic_id` claims.
- Medical safety: endpoints must not accept or return clinical diagnosis, prescription, medical-record, or treatment-decision data.

Example headers:

```http
Content-Type: application/json
Authorization: Bearer <jwt>
```

## Authentication Strategy

`POST /api/auth/login` authenticates an active operator with email and password, then returns a JWT. Protected endpoints require `Authorization: Bearer <jwt>`.

JWT-backed protected endpoints use these claims:

- User ID.
- Clinic ID.
- Role.
- Expiration.

Roles currently used by the API:

- `superadmin`
- `clinic_admin`
- `assistant`

Frontend guards improve UX, but backend validation is mandatory on every protected endpoint.

## Error Response Format

Implemented API errors use this shape:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request payload is invalid."
  }
}
```

Error `details` are not currently emitted. Responses must not expose stack traces, SQL errors, provider secrets, internal file paths, raw panic messages, or unsafe AI internals.

Common implemented error codes:

- `VALIDATION_ERROR`
- `INVALID_REQUEST`
- `UNAUTHORIZED`
- `FORBIDDEN`
- `NOT_FOUND`
- `AI_PROVIDER_ERROR`
- `AI_SAFETY_BLOCKED`
- `INTERNAL_ERROR`

## Pagination Convention

Lead and follow-up list endpoints support:

```txt
page=1
page_size=20
```

Paginated responses use:

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 0,
    "total_pages": 0
  }
}
```

## Endpoint Summary and Authentication

| Method | Endpoint | Auth Required | Roles |
| --- | --- | --- | --- |
| GET | `/healthz` | No | Public |
| GET | `/readyz` | No | Public |
| POST | `/api/auth/login` | No | Public |
| GET | `/api/clinics/current` | Yes | any authenticated user |
| PUT | `/api/clinics/current` | Yes | `clinic_admin` |
| GET | `/api/services` | Yes | any authenticated user |
| GET | `/api/services/:id` | Yes | any authenticated user |
| POST | `/api/services` | Yes | `clinic_admin` |
| PUT | `/api/services/:id` | Yes | `clinic_admin` |
| DELETE | `/api/services/:id` | Yes | `clinic_admin` |
| GET | `/api/leads` | Yes | any authenticated user |
| POST | `/api/leads` | Yes | any authenticated user |
| GET | `/api/leads/:id` | Yes | any authenticated user |
| PUT | `/api/leads/:id` | Yes | any authenticated user |
| GET | `/api/followups` | Yes | any authenticated user |
| POST | `/api/followups/:id/complete` | Yes | any authenticated user |
| POST | `/api/followups/:id/reschedule` | Yes | any authenticated user |
| POST | `/api/ai/analyze-conversation` | Yes | any authenticated user |
| POST | `/api/ai/reply-suggestion` | Yes | any authenticated user |
| POST | `/api/ai/objection-handler` | Yes | any authenticated user |
| POST | `/api/ai/follow-up-message` | Yes | any authenticated user |
| GET | `/api/dashboard/summary` | Yes | any authenticated user |

`POST /api/content/generate-post` is planned in MVP documentation but is not implemented in the current backend.

## Health Endpoints

### GET /healthz

Returns basic process health.

Response:

```json
{
  "status": "ok"
}
```

### GET /readyz

Returns database readiness. If the database connection is available and ping succeeds:

```json
{
  "status": "ready"
}
```

## Authentication Endpoints

### POST /api/auth/login

Authenticates a user and returns a JWT.

Request:

```json
{
  "email": "admin@sonrisaviva.demo",
  "password": "clinicflow123"
}
```

Response:

```json
{
  "access_token": "jwt-token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "22222222-2222-4222-8222-222222222222",
    "full_name": "Admin Demo",
    "email": "admin@sonrisaviva.demo",
    "role": "clinic_admin",
    "clinic_id": "11111111-1111-4111-8111-111111111111"
  }
}
```

## Currency Contract

Status: Implemented for MVP baseline.

The API exposes clinic-level country/currency configuration and service-level `currency_code` next to `price_from`. The MVP does not perform exchange-rate conversion.

Clinic response additions:

```json
{
  "country_code": "CO",
  "currency_code": "COP",
  "currency": {
    "code": "COP",
    "symbol": "$",
    "locale": "es-CO",
    "decimal_digits": 0,
    "thousand_separator": ".",
    "decimal_separator": ",",
    "symbol_position": "before"
  }
}
```

Clinic update additions:

```json
{
  "country_code": "CO",
  "currency_code": "COP"
}
```

Supported MVP country/currency pairs:

| country_code | currency_code | Locale | Decimal digits |
| --- | --- | --- | --- |
| CO | COP | es-CO | 0 |
| PE | PEN | es-PE | 2 |
| AR | ARS | es-AR | 2 |
| CL | CLP | es-CL | 0 |

Service create/update requests and responses include `price_from` as a string decimal value plus `currency_code` in responses. The backend accepts either a JSON string or number for compatibility, but responses use strings to avoid floating-point money precision issues.

No endpoint should perform exchange-rate conversion in the MVP.

## Clinic Endpoints

### GET /api/clinics/current

Returns the authenticated user's clinic configuration using `clinic_id` from JWT claims.

Response:

```json
{
  "id": "11111111-1111-4111-8111-111111111111",
  "name": "Sonrisa Viva Demo",
  "clinic_type": "odontologia",
  "city": "Bogota",
  "country_code": "CO",
  "currency_code": "COP",
  "currency": {
    "code": "COP",
    "symbol": "$",
    "locale": "es-CO",
    "decimal_digits": 0,
    "thousand_separator": ".",
    "decimal_separator": ",",
    "symbol_position": "before"
  },
  "phone": "+573001112233",
  "whatsapp": "+573001112233",
  "address": "Calle 123 #45-67",
  "opening_hours": {
    "monday_friday": "08:00-18:00",
    "saturday": "08:00-13:00"
  },
  "general_faq": [
    {
      "question": "Atienden urgencias?",
      "answer": "Si, podemos validar disponibilidad por WhatsApp."
    }
  ],
  "communication_tone": "amable"
}
```

### PUT /api/clinics/current

Updates commercial clinic configuration. `clinic_type` is not updatable through this endpoint.

Allowed `communication_tone` values:

```txt
amable
profesional
cercano
juvenil
elegante
```

Request:

```json
{
  "name": "Sonrisa Viva Demo",
  "city": "Bogota",
  "country_code": "CO",
  "currency_code": "COP",
  "phone": "+573001112233",
  "whatsapp": "+573001112233",
  "address": "Calle 123 #45-67",
  "opening_hours": {
    "monday_friday": "08:00-18:00",
    "saturday": "08:00-13:00"
  },
  "general_faq": [],
  "communication_tone": "profesional"
}
```

Response uses the same shape as `GET /api/clinics/current`.

## Service Catalog Endpoints

Service catalog endpoints are scoped to the authenticated clinic.

Service response shape:

```json
{
  "id": "33333333-3333-4333-8333-333333333331",
  "clinic_id": "11111111-1111-4111-8111-111111111111",
  "name": "Blanqueamiento dental",
  "description": "Tratamiento estetico para mejorar el tono de la sonrisa con valoracion previa.",
  "duration_minutes": 60,
  "price_from": "250000.00",
  "currency_code": "COP",
  "benefits": ["Mejora estetica visible", "Valoracion personalizada"],
  "faq": [
    {
      "question": "Debilita los dientes?",
      "answer": "La valoracion permite confirmar si el paciente es apto."
    }
  ],
  "common_objections": ["Esta muy caro", "Me da miedo"],
  "is_active": true
}
```

### GET /api/services

Lists clinic services.

Response:

```json
[
  {
    "id": "33333333-3333-4333-8333-333333333331",
    "clinic_id": "11111111-1111-4111-8111-111111111111",
    "name": "Blanqueamiento dental",
    "description": "Tratamiento estetico para mejorar el tono de la sonrisa con valoracion previa.",
    "duration_minutes": 60,
    "price_from": "250000.00",
    "currency_code": "COP",
    "benefits": ["Mejora estetica visible", "Valoracion personalizada"],
    "faq": [],
    "common_objections": ["Esta muy caro", "Me da miedo"],
    "is_active": true
  }
]
```

### GET /api/services/:id

Returns one service scoped to the authenticated clinic.

Response uses the service response shape.

### POST /api/services

Creates a service in the authenticated clinic.

Request:

```json
{
  "name": "Blanqueamiento dental",
  "description": "Tratamiento estetico para mejorar el tono de la sonrisa.",
  "duration_minutes": 60,
  "price_from": "250000.00",
  "benefits": ["Mejora estetica visible", "Valoracion personalizada"],
  "faq": [
    {
      "question": "Debilita los dientes?",
      "answer": "La valoracion permite confirmar si el paciente es apto."
    }
  ],
  "common_objections": ["Esta muy caro", "Me da miedo"]
}
```

Response: `201 Created`, service response shape.

### PUT /api/services/:id

Updates a clinic service.

Request:

```json
{
  "name": "Blanqueamiento dental",
  "description": "Tratamiento estetico supervisado por profesionales.",
  "duration_minutes": 60,
  "price_from": "280000.00",
  "benefits": ["Mejora estetica visible"],
  "faq": [],
  "common_objections": ["Esta muy caro"],
  "is_active": true
}
```

Response: service response shape.

### DELETE /api/services/:id

Deletes a clinic service. Existing leads referencing that service keep tenant isolation and the database sets `service_id` to `null`.

Response: `204 No Content` with an empty body.

## Lead Endpoints

Allowed lead statuses:

```txt
Nuevo
Contactado
Interesado
Agendado
No Respondio
Perdido
Convertido
```

### GET /api/leads

Lists leads with pagination and optional filters.

Query parameters:

```txt
page=1
page_size=20
status=Nuevo
service_id=33333333-3333-4333-8333-333333333331
```

Response:

```json
{
  "data": [
    {
      "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
      "full_name": "Maria Perez",
      "phone": "+573009998877",
      "service_id": "33333333-3333-4333-8333-333333333331",
      "service_name": "Blanqueamiento dental",
      "status": "Nuevo",
      "source": "whatsapp",
      "next_action_at": "2026-05-22T14:00:00Z",
      "created_at": "2026-05-20T15:30:00Z",
      "updated_at": "2026-05-20T15:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

### POST /api/leads

Creates a commercial lead and optionally creates the first note.

Request:

```json
{
  "full_name": "Maria Perez",
  "phone": "+573009998877",
  "service_id": "33333333-3333-4333-8333-333333333331",
  "status": "Nuevo",
  "source": "whatsapp",
  "notes": "Pregunta por precio desde Instagram.",
  "next_action_at": "2026-05-22T14:00:00Z"
}
```

Response: `201 Created`, lead list item shape.

### GET /api/leads/:id

Returns one lead scoped to the authenticated clinic.

Response:

```json
{
  "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "full_name": "Maria Perez",
  "phone": "+573009998877",
  "service": {
    "id": "33333333-3333-4333-8333-333333333331",
    "name": "Blanqueamiento dental"
  },
  "status": "Interesado",
  "source": "whatsapp",
  "notes": [
    {
      "id": "0d9319a3-9b65-4ca6-9f4d-b38f5e4bce78",
      "body": "Solicito informacion de precio.",
      "created_at": "2026-05-20T15:30:00Z"
    }
  ],
  "next_action_at": "2026-05-22T14:00:00Z",
  "created_at": "2026-05-20T15:30:00Z"
}
```

### PUT /api/leads/:id

Updates lead status, optionally appends a note, and optionally sets the next action date.

Request:

```json
{
  "status": "Interesado",
  "note": "Quiere agendar valoracion esta semana.",
  "next_action_at": "2026-05-22T14:00:00Z"
}
```

Response:

```json
{
  "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "status": "Interesado"
}
```

## Follow-Up Endpoints

Follow-ups are modeled as leads with a non-null `next_action_at`.

### GET /api/followups

Lists leads with a pending manual follow-up.

Query parameters:

```txt
due=all|today|overdue|upcoming
page=1
page_size=20
```

Response uses the standard paginated lead list shape from `GET /api/leads`.

### POST /api/followups/:id/complete

Marks a manual follow-up as completed, clears `next_action_at`, optionally updates status, and optionally appends a lead note.

Request:

```json
{
  "status": "Contactado",
  "note": "Se envio mensaje de seguimiento por WhatsApp."
}
```

Response:

```json
{
  "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "status": "completed"
}
```

### POST /api/followups/:id/reschedule

Schedules or reschedules the next manual follow-up for a lead.

Request:

```json
{
  "next_action_at": "2026-05-23T15:00:00Z",
  "note": "Reprogramado para manana."
}
```

Response:

```json
{
  "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "status": "rescheduled"
}
```

## AI Endpoints

AI endpoints are commercial-assistance endpoints. They must not diagnose, prescribe, produce clinical decision support, or operate as autonomous WhatsApp bots. The backend builds prompts from clinic, service, and lead context and validates AI output safety.

The backend supports configured providers `openai`, `gemini`, `deepseek`, plus `mock` for local smoke testing only.

### POST /api/ai/analyze-conversation

Status: Implemented.

Analyzes a manually pasted commercial conversation and returns structured commercial insight plus a safe response draft. This endpoint does not store the raw conversation and does not send messages. Any lead creation still requires human review and a separate `POST /api/leads` request.

Request:

```json
{
  "conversation_text": "Paciente: Hola, cuanto cuesta el blanqueamiento?\nClinica: Hola, con gusto...",
  "source": "whatsapp",
  "service_id": "33333333-3333-4333-8333-333333333331"
}
```

Response:

```json
{
  "analysis_id": "analysis-11111111-1111-4111-8111-111111111111",
  "detected_lead": {
    "full_name": "Maria Perez",
    "phone": "+573009998877"
  },
  "detected_service": {
    "service_id": "",
    "service_name": "Blanqueamiento dental",
    "confidence": "medium"
  },
  "intent": "high",
  "detected_objections": ["precio"],
  "suggested_status": "Interesado",
  "commercial_summary": "Pregunta por precio y muestra interes en conocer el tratamiento.",
  "suggested_reply": "Claro, podemos orientarte y agendar una valoracion para revisar tu caso.",
  "suggested_next_action": "Responder y proponer valoracion",
  "suggested_follow_up_at": "",
  "safety_status": "passed"
}
```

Notes:

- `conversation_text` is analyzed as commercial context, not stored as a medical record.
- Extracted fields are suggestions until a human reviews them.
- The frontend uses reviewed fields with the existing `POST /api/leads` endpoint when creating a lead.

### POST /api/ai/reply-suggestion

Generates safe WhatsApp-ready reply suggestions.

Request:

```json
{
  "lead_id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "service_id": "33333333-3333-4333-8333-333333333331",
  "patient_message": "El blanqueamiento debilita los dientes?",
  "desired_tone": "amable"
}
```

Response:

```json
{
  "generation_id": "gen-11111111-1111-4111-8111-111111111111",
  "variants": {
    "short": "Claro, con gusto te ayudamos con la informacion.",
    "persuasive": "Podemos orientarte y revisar la mejor opcion segun tu necesidad.",
    "technical": "Primero realizamos una valoracion para confirmar el plan adecuado.",
    "closing_question": "Quieres que te ayudemos a coordinar una valoracion?"
  },
  "safety_status": "passed"
}
```

### POST /api/ai/objection-handler

Generates a safe commercial response to an objection.

Request:

```json
{
  "lead_id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "service_id": "33333333-3333-4333-8333-333333333331",
  "objection": "En otro lado es mas barato"
}
```

Response:

```json
{
  "generation_id": "gen-11111111-1111-4111-8111-111111111111",
  "objection_type": "precio",
  "recommended_strategy": "Validar interes y reforzar valor sin presionar.",
  "suggested_message": "Entiendo tu inquietud. Podemos revisar opciones y resolver tus dudas antes de tomar una decision.",
  "closing_question": "Quieres que te compartamos los pasos para una valoracion?",
  "safety_status": "passed"
}
```

### POST /api/ai/follow-up-message

Generates a safe WhatsApp-ready re-engagement message for a lead.

Request:

```json
{
  "lead_id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "service_id": "33333333-3333-4333-8333-333333333331",
  "last_contact_note": "Pregunto por precio y no respondio despues de la primera respuesta."
}
```

Response:

```json
{
  "generation_id": "gen-11111111-1111-4111-8111-111111111111",
  "suggested_message": "Hola, queriamos saber si pudiste revisar la informacion. Estamos atentos para ayudarte con el siguiente paso.",
  "recommended_timing": "Hoy en horario laboral",
  "next_step": "Enviar mensaje manual por WhatsApp",
  "safety_status": "passed"
}
```

## Dashboard Endpoints

### GET /api/dashboard/summary

Returns basic commercial metrics for the authenticated clinic.

Response:

```json
{
  "leads_total": 128,
  "leads_by_status": {
    "Nuevo": 18,
    "Contactado": 22,
    "Interesado": 31,
    "Agendado": 20,
    "No Respondio": 15,
    "Perdido": 12,
    "Convertido": 10
  },
  "top_services": [
    {
      "service_id": "33333333-3333-4333-8333-333333333331",
      "service_name": "Blanqueamiento dental",
      "lead_count": 35
    }
  ],
  "pending_followups_today": 9,
  "overdue_followups": 2,
  "upcoming_followups": 14,
  "conversion_rate": 0.078
}
```

## Local Smoke Coverage

`./e2e_test.sh` currently validates these flows against the local Go API and PostgreSQL demo data:

- Login.
- Clinic profile read/update and cleanup restore.
- Service create/update/list/delete cleanup.
- Lead create/update/detail.
- Follow-up list/reschedule/complete.
- AI reply, objection, and follow-up message generation with `AI_PROVIDER=mock`.
- Dashboard summary.

## Proposed Smart Lead Inbox API Contracts

Status: Planned / Proposed. Not implemented.

These contracts document remaining proposed Smart Lead Inbox capabilities. They are intentionally not part of the implemented API surface yet and must be validated before development. All proposed endpoints require JWT auth, clinic-level tenant isolation, sanitized errors, and backend-owned AI safety controls.

These endpoints must not send WhatsApp messages or any other patient message automatically. They only support analysis, draft generation, reviewed lead updates, and suggested follow-ups.

### POST /api/leads/from-conversation

Creates or updates a lead from human-reviewed conversation analysis.

Proposed request:

```json
{
  "analysis_id": "analysis-id",
  "action": "create",
  "lead_id": "optional-existing-lead-id",
  "reviewed_lead": {
    "full_name": "Maria Perez",
    "phone": "+573009998877",
    "service_id": "33333333-3333-4333-8333-333333333331",
    "status": "Interesado",
    "source": "whatsapp",
    "note": "Pregunta por precio y quiere conocer disponibilidad.",
    "next_action_at": "2026-05-25T14:00:00Z"
  }
}
```

Proposed response:

```json
{
  "id": "lead-id",
  "status": "Interesado",
  "created_from_analysis": true
}
```

Notes:

- A human must review and submit this request.
- Backend must enforce clinic isolation and validate allowed lead fields.
- Do not store diagnosis, prescription, or clinical-history content in the note.

### POST /api/inbound/messages

Optionally stores a manually captured inbound commercial message or conversation fragment after retention rules are approved.

Proposed request:

```json
{
  "source": "whatsapp",
  "lead_id": "optional-lead-id",
  "message_text": "Hola, quiero saber el precio del blanqueamiento",
  "received_at": "2026-05-24T15:00:00Z"
}
```

Proposed response:

```json
{
  "id": "message-id",
  "status": "stored"
}
```

Notes:

- This endpoint is optional and should not be implemented until privacy and retention rules are clear.
- Full conversation logging can create sensitive-data risk. Prefer summarized commercial notes when possible.

### POST /api/followups/suggest

Suggests a manual follow-up action from lead context or conversation analysis.

Proposed request:

```json
{
  "lead_id": "lead-id",
  "analysis_id": "optional-analysis-id",
  "last_contact_note": "Respondio que lo va a pensar por precio."
}
```

Proposed response:

```json
{
  "suggested_next_action": "Enviar mensaje amable de seguimiento",
  "suggested_follow_up_at": "2026-05-25T14:00:00Z",
  "suggested_message": "Hola, queriamos saber si pudiste revisar la informacion. Estamos atentos para ayudarte.",
  "safety_status": "passed"
}
```

Notes:

- The response is a suggestion only.
- Scheduling a follow-up requires human confirmation through an implemented lead/follow-up update flow.
- Sending a message remains manual.

## Planned But Not Implemented

- `POST /api/content/generate-post`.
- `POST /api/leads/from-conversation` for proposed human-reviewed lead creation/update from analysis.
- `POST /api/inbound/messages` for optional future manual inbound message storage after retention rules are approved.
- `POST /api/followups/suggest` for proposed follow-up suggestions requiring human confirmation.

- User management endpoints.
- Superadmin platform administration endpoints.
- Lead deletion endpoint.
- AI usage tracking and billing metadata.
- Native WhatsApp Business Cloud API integration, which remains excluded from MVP.

## Notes for Future Changes

- Add explicit superadmin endpoints when platform administration is implemented.
- Keep AI provider selection provider-agnostic and environment-based.
- Keep WhatsApp Business Cloud API integration out of the MVP contract unless scope is explicitly changed.
- Keep Smart Lead Inbox endpoints marked as planned until backend, frontend, persistence, safety, and privacy decisions are approved.

- Do not add clinical-history, diagnosis, prescription, or medical-record endpoints.
