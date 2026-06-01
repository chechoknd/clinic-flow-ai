# API Contracts

## Status

Status: Implemented for the current MVP backend/frontend surface. Smart Schedule endpoints are planned/proposed only and are not implemented.

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
| GET | `/api/professionals` | Yes | any authenticated user |
| GET | `/api/professionals/:id` | Yes | any authenticated user |
| POST | `/api/professionals` | Yes | `clinic_admin` |
| PUT | `/api/professionals/:id` | Yes | `clinic_admin` |
| GET | `/api/appointments` | Yes | any authenticated user |
| GET | `/api/appointments/:id` | Yes | any authenticated user |
| POST | `/api/appointments` | Yes | any authenticated user |
| PUT | `/api/appointments/:id` | Yes | any authenticated user |
| POST | `/api/appointments/:id/status` | Yes | any authenticated user |
| POST | `/api/appointments/:id/reschedule` | Yes | any authenticated user |
| POST | `/api/leads/:id/convert-to-appointment` | Yes | any authenticated user |
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

Planned Smart Schedule endpoints are documented later in this file. They are not part of the implemented API surface yet.

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
  "price_from": 250000,
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
    "price_from": 250000,
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
  "price_from": 250000,
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
  "price_from": 280000,
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

## Professional Endpoints

Status: Implemented.

Professionals are operational schedule resources for the Smart Schedule. They are scoped to the authenticated clinic and must not store clinical histories, diagnoses, prescriptions, clinical notes, or medical-record content.

Professional response shape:

```json
{
  "id": "professional-id",
  "clinic_id": "clinic-id",
  "full_name": "Dra. Ana Gomez",
  "role_or_specialty": "Ortodoncia",
  "calendar_color": "#2563EB",
  "working_hours": {
    "monday": [{"start": "08:00", "end": "12:00"}]
  },
  "is_active": true,
  "service_ids": ["service-id"],
  "created_at": "2026-06-01T14:00:00Z",
  "updated_at": "2026-06-01T14:00:00Z"
}
```

### GET /api/professionals

Lists professionals for the authenticated clinic.

Query parameters:

```txt
is_active=true
service_id=33333333-3333-4333-8333-333333333331
```

Response:

```json
[
  {
    "id": "professional-id",
    "clinic_id": "clinic-id",
    "full_name": "Dra. Ana Gomez",
    "role_or_specialty": "Ortodoncia",
    "calendar_color": "#2563EB",
    "working_hours": {},
    "is_active": true,
    "service_ids": ["service-id"],
    "created_at": "2026-06-01T14:00:00Z",
    "updated_at": "2026-06-01T14:00:00Z"
  }
]
```

### GET /api/professionals/:id

Returns one professional scoped to the authenticated clinic.

Response uses the professional response shape.

### POST /api/professionals

Creates an operational professional.

Request:

```json
{
  "full_name": "Dra. Ana Gomez",
  "role_or_specialty": "Ortodoncia",
  "calendar_color": "#2563EB",
  "working_hours": {
    "monday": [{"start": "08:00", "end": "12:00"}]
  },
  "service_ids": ["33333333-3333-4333-8333-333333333331"]
}
```

Response: `201 Created`, professional response shape.

### PUT /api/professionals/:id

Updates an operational professional.

Request:

```json
{
  "full_name": "Dra. Ana Gomez",
  "role_or_specialty": "Ortodoncia",
  "calendar_color": "#2563EB",
  "working_hours": {},
  "is_active": true,
  "service_ids": ["33333333-3333-4333-8333-333333333331"]
}
```

Response: professional response shape.

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

Future Smart Lead Inbox schedule integration should allow AI to suggest appointment creation, missing scheduling data, and confirmation/follow-up messaging. Appointment creation must still happen through a separate human-reviewed appointment endpoint.

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

## Appointment Endpoints

Status: Implemented.

Appointments are operational schedule records. They connect clinic, professional, service, optional lead, contact, date/time, status, confirmation state, source, and commercial/admin notes only.

Appointment response shape:

```json
{
  "id": "appointment-id",
  "clinic_id": "clinic-id",
  "professional": {
    "id": "professional-id",
    "full_name": "Dra. Ana Gomez",
    "calendar_color": "#2563EB"
  },
  "lead": {
    "id": "lead-id",
    "full_name": "Maria Perez"
  },
  "service": {
    "id": "service-id",
    "name": "Ortodoncia"
  },
  "contact_name": "Maria Perez",
  "contact_phone": "+573009998877",
  "starts_at": "2026-06-02T14:00:00-05:00",
  "ends_at": "2026-06-02T15:00:00-05:00",
  "status": "pending_confirmation",
  "source": "whatsapp",
  "confirmation_status": "pending",
  "admin_notes": "Confirmar asistencia en la manana.",
  "created_at": "2026-06-01T14:00:00Z",
  "updated_at": "2026-06-01T14:00:00Z"
}
```

### GET /api/appointments

Lists appointments for the authenticated clinic.

Query parameters:

```txt
date=2026-06-02
date_from=2026-06-02T00:00:00-05:00
date_to=2026-06-08T23:59:59-05:00
professional_id=professional-id
service_id=service-id
status=pending_confirmation
lead_id=lead-id
```

Response:

```json
{
  "data": []
}
```

### GET /api/appointments/:id

Returns one appointment scoped to the authenticated clinic.

Response uses the appointment response shape.

### POST /api/appointments

Creates an appointment after human review.

Request:

```json
{
  "professional_id": "professional-id",
  "lead_id": "optional-lead-id",
  "service_id": "service-id",
  "contact_name": "Maria Perez",
  "contact_phone": "+573009998877",
  "starts_at": "2026-06-02T14:00:00-05:00",
  "duration_minutes": 60,
  "status": "pending_confirmation",
  "source": "whatsapp",
  "admin_notes": "Solicito horario en la tarde."
}
```

Response: `201 Created`, appointment response shape.

### PUT /api/appointments/:id

Updates appointment operational data.

Request:

```json
{
  "professional_id": "professional-id",
  "service_id": "service-id",
  "contact_name": "Maria Perez",
  "contact_phone": "+573009998877",
  "starts_at": "2026-06-02T15:00:00-05:00",
  "duration_minutes": 60,
  "status": "rescheduled",
  "confirmation_status": "pending",
  "admin_notes": "Reprogramada por solicitud del contacto."
}
```

Response: appointment response shape.

### POST /api/appointments/:id/status

Updates appointment status and optionally replaces/sets the admin note.

Request:

```json
{
  "status": "confirmed",
  "admin_note": "Confirmo asistencia por WhatsApp."
}
```

Response: appointment response shape.

### POST /api/appointments/:id/reschedule

Reschedules an appointment after human confirmation.

Request:

```json
{
  "starts_at": "2026-06-03T10:00:00-05:00",
  "duration_minutes": 60,
  "admin_note": "Reprogramada para manana en la manana."
}
```

Response: appointment response shape.

### POST /api/leads/:id/convert-to-appointment

Creates an appointment from a lead after human review.

Request:

```json
{
  "professional_id": "professional-id",
  "service_id": "service-id",
  "starts_at": "2026-06-02T14:00:00-05:00",
  "duration_minutes": 60,
  "status": "pending_confirmation",
  "admin_notes": "Lead pidio valoracion inicial.",
  "update_lead_status": true
}
```

Response:

```json
{
  "lead_id": "lead-id",
  "lead_status": "Agendado",
  "appointment_id": "appointment-id",
  "appointment_status": "pending_confirmation"
}
```

Notes:

- Appointment creation validates clinic ownership for professional, service, lead, and user.
- The backend rejects overlapping active appointments for the same professional.
- Appointment notes must remain commercial/admin only.
- No message is sent automatically.

## Proposed Smart Schedule API Contracts

Status: Planned / Proposed for remaining schedule availability and dashboard endpoints. Professional and appointment endpoints are implemented and documented above.

These contracts document the remaining planned Smart Schedule / Agenda Inteligente surface. They are intentionally not implemented yet. All endpoints require JWT auth, clinic-level tenant isolation, sanitized errors, and backend authorization.

The Smart Schedule is operational and commercial. It must not accept or return clinical histories, diagnoses, prescriptions, treatment plans, clinical images, clinical notes, or AI medical advice.

### Planned Appointment Status Values

Stable API values:

```txt
scheduled
confirmed
pending_confirmation
rescheduled
no_show
cancelled
completed
converted_from_lead
```

Suggested Spanish UI labels:

```txt
Programada
Confirmada
Pendiente de confirmacion
Reprogramada
No asistio
Cancelada
Completada
Convertida desde lead
```

### GET /api/schedule/availability

Returns available slots for a professional, date range, and optional service duration.

Proposed query parameters:

```txt
professional_id=professional-id
service_id=service-id
date_from=2026-06-02
date_to=2026-06-08
duration_minutes=60
```

Proposed response:

```json
{
  "professional_id": "professional-id",
  "slots": [
    {
      "starts_at": "2026-06-02T14:00:00-05:00",
      "ends_at": "2026-06-02T15:00:00-05:00"
    }
  ]
}
```

Notes:

- Availability is operational and based on configured working hours plus existing appointments.
- It must not imply clinical urgency or medical prioritization.

### GET /api/dashboard/schedule-summary

Returns schedule-centered daily priorities for the authenticated clinic.

Proposed query parameters:

```txt
date=2026-06-02
```

Proposed response:

```json
{
  "date": "2026-06-02",
  "todays_appointments": 18,
  "appointments_pending_confirmation": 5,
  "available_slots": 7,
  "hot_leads_without_appointment": 4,
  "overdue_followups": 2,
  "appointments_by_professional": [
    {
      "professional_id": "professional-id",
      "professional_name": "Dra. Ana Gomez",
      "appointment_count": 6
    }
  ],
  "leads_converted_to_appointments": 3,
  "top_services_by_schedule_demand": [
    {
      "service_id": "service-id",
      "service_name": "Ortodoncia",
      "appointment_count": 8
    }
  ]
}
```

## Proposed Smart Schedule AI Extensions

Status: Planned / Proposed. Not implemented.

Potential extension to `POST /api/ai/analyze-conversation` response:

```json
{
  "scheduling_intent": "high",
  "missing_scheduling_data": ["preferred_date", "preferred_time"],
  "suggested_appointment": {
    "service_id": "service-id",
    "professional_id": "",
    "duration_minutes": 60,
    "preferred_time_text": "manana en la tarde"
  },
  "suggested_next_action": "Pedir horario preferido y ofrecer opciones de valoracion"
}
```

Notes:

- AI may suggest appointment creation or missing scheduling data.
- AI must not create, confirm, cancel, or reschedule an appointment automatically.
- Human review is mandatory before any appointment endpoint is called.

## Planned But Not Implemented

- `POST /api/content/generate-post`.
- `POST /api/leads/from-conversation` for proposed human-reviewed lead creation/update from analysis.
- `POST /api/inbound/messages` for optional future manual inbound message storage after retention rules are approved.
- `POST /api/followups/suggest` for proposed follow-up suggestions requiring human confirmation.
- `GET /api/schedule/availability`.
- `GET /api/dashboard/schedule-summary`.

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
- Keep remaining Smart Schedule endpoints marked as planned until database schema, frontend flows, backend validation, and tenant isolation rules are approved.

- Do not add clinical-history, diagnosis, prescription, or medical-record endpoints.
