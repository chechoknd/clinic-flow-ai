# API Contracts

## Status

Status: Planned.

This document defines the initial REST API contract for the ClinicFlow AI MVP. It is a planning contract and must be updated when endpoints, payloads, authentication, pagination, or error formats change.

## API Conventions

- Base path: `/api`
- Format: JSON request and response bodies.
- Authentication: JWT Bearer token for protected endpoints.
- IDs: UUID strings.
- Dates and times: ISO 8601 strings.
- Naming: JSON fields use `snake_case`.
- Tenant isolation: protected endpoints must scope data to the authenticated clinic unless explicitly platform-level.
- Medical safety: endpoints must not accept or return clinical diagnosis, prescription, medical-record, or treatment-decision data.

Example headers:

```http
Content-Type: application/json
Authorization: Bearer <jwt>
```

## Authentication Strategy

`POST /api/auth/login` authenticates an active operator with email and password, then returns a JWT. Protected endpoints require `Authorization: Bearer <jwt>`.

JWT claims should support:

- User ID.
- Clinic ID.
- Role.
- Expiration.

Roles:

- `superadmin`
- `clinic_admin`
- `assistant`

Frontend guards improve UX, but backend validation is mandatory on every protected endpoint.

## Error Response Format

All API errors should use:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request payload is invalid.",
    "details": {
      "field": "email"
    }
  }
}
```

`details` is optional. Responses must not expose stack traces, SQL errors, provider secrets, internal file paths, raw panic messages, or unsafe AI internals.

Common error codes:

- `VALIDATION_ERROR`
- `UNAUTHORIZED`
- `FORBIDDEN`
- `NOT_FOUND`
- `CONFLICT`
- `AI_PROVIDER_ERROR`
- `AI_SAFETY_BLOCKED`
- `INTERNAL_ERROR`

## Pagination Convention

List endpoints should support:

```txt
page=1
page_size=20
```

Paginated responses should use:

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
| POST | `/api/auth/login` | No | Public |
| GET | `/api/clinics/current` | Yes | `clinic_admin`, `assistant` |
| PUT | `/api/clinics/current` | Yes | `clinic_admin` |
| GET | `/api/services` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/services` | Yes | `clinic_admin` |
| PUT | `/api/services/:id` | Yes | `clinic_admin` |
| DELETE | `/api/services/:id` | Yes | `clinic_admin` |
| GET | `/api/leads` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/leads` | Yes | `clinic_admin`, `assistant` |
| GET | `/api/leads/:id` | Yes | `clinic_admin`, `assistant` |
| PUT | `/api/leads/:id` | Yes | `clinic_admin`, `assistant` |
| GET | `/api/followups` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/followups/:id/complete` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/followups/:id/reschedule` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/ai/reply-suggestion` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/ai/objection-handler` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/ai/follow-up-message` | Yes | `clinic_admin`, `assistant` |
| POST | `/api/content/generate-post` | Yes | `clinic_admin`, `assistant` |
| GET | `/api/dashboard/summary` | Yes | `clinic_admin`, `assistant` |

`superadmin` permissions for clinic-scoped endpoints should be defined explicitly when platform administration is implemented.

## Core Endpoints

### POST /api/auth/login

Authenticates a user and returns a JWT.

Request:

```json
{
  "email": "assistant@clinic.example",
  "password": "password"
}
```

Response:

```json
{
  "access_token": "jwt-token",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "9b86b15a-43fe-4d71-90ef-12e8a3b3b5aa",
    "full_name": "Ana Gomez",
    "email": "assistant@clinic.example",
    "role": "assistant",
    "clinic_id": "5c673404-7b0f-4f39-8df3-839a7fdb63ef"
  }
}
```

### GET /api/clinics/current

Protected endpoint. Returns the authenticated user's clinic configuration using `clinic_id` from the JWT claims.

Response:

```json
{
  "id": "5c673404-7b0f-4f39-8df3-839a7fdb63ef",
  "name": "Sonrisa Viva",
  "clinic_type": "odontologia",
  "city": "Bogota",
  "phone": "+573001112233",
  "whatsapp": "+573001112233",
  "address": "Calle 123 #45-67",
  "opening_hours": {
    "monday_friday": "08:00-18:00",
    "saturday": "08:00-13:00"
  },
  "communication_tone": "amable",
  "general_faq": [
    {
      "question": "Atienden urgencias?",
      "answer": "Si, se recomienda escribir por WhatsApp para validar disponibilidad."
    }
  ]
}
```

### PUT /api/clinics/current

Updates clinic configuration.

Request:

```json
{
  "name": "Sonrisa Viva",
  "city": "Bogota",
  "phone": "+573001112233",
  "whatsapp": "+573001112233",
  "address": "Calle 123 #45-67",
  "opening_hours": {
    "monday_friday": "08:00-18:00",
    "saturday": "08:00-13:00"
  },
  "communication_tone": "profesional",
  "general_faq": []
}
```

Response:

```json
{
  "id": "5c673404-7b0f-4f39-8df3-839a7fdb63ef",
  "name": "Sonrisa Viva",
  "clinic_type": "odontologia",
  "city": "Bogota",
  "communication_tone": "profesional",
  "updated_at": "2026-05-20T15:30:00Z"
}
```

### GET /api/services

Lists clinic services.

Response:

```json
{
  "data": [
    {
      "id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
      "name": "Blanqueamiento dental",
      "description": "Tratamiento estetico para mejorar el tono de la sonrisa.",
      "duration_minutes": 60,
      "price_from": 250000,
      "benefits": ["Mejora estetica visible", "Valoracion personalizada"],
      "faq": [],
      "common_objections": ["Esta muy caro", "Me da miedo"]
    }
  ]
}
```

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

Response:

```json
{
  "id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
  "name": "Blanqueamiento dental",
  "created_at": "2026-05-20T15:30:00Z"
}
```

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
  "common_objections": ["Esta muy caro"]
}
```

Response:

```json
{
  "id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
  "updated_at": "2026-05-20T15:30:00Z"
}
```

### DELETE /api/services/:id

Deletes or deactivates a clinic service.

Response:

```json
{
  "deleted": true
}
```

### GET /api/leads

Lists leads with pagination and optional filters.

Query parameters:

```txt
page=1
page_size=20
status=Nuevo
service_id=c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d
```

Response:

```json
{
  "data": [
    {
      "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
      "full_name": "Maria Perez",
      "phone": "+573009998877",
      "service_id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
      "service_name": "Blanqueamiento dental",
      "status": "Nuevo",
      "source": "whatsapp",
      "next_action_at": "2026-05-22T14:00:00Z",
      "created_at": "2026-05-20T15:30:00Z"
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

Creates a commercial lead.

Request:

```json
{
  "full_name": "Maria Perez",
  "phone": "+573009998877",
  "service_id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
  "status": "Nuevo",
  "source": "whatsapp",
  "notes": "Pregunta por precio desde Instagram.",
  "next_action_at": "2026-05-22T14:00:00Z"
}
```

Response:

```json
{
  "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "status": "Nuevo",
  "created_at": "2026-05-20T15:30:00Z"
}
```

### GET /api/leads/:id

Returns one lead scoped to the authenticated clinic.

Response:

```json
{
  "id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "full_name": "Maria Perez",
  "phone": "+573009998877",
  "service": {
    "id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
    "name": "Blanqueamiento dental"
  },
  "status": "Interesado",
  "source": "whatsapp",
  "notes": [
    {
      "id": "0d9319a3-9b65-4ca6-9f4d-b38f5e4bce78",
      "body": "Solicito informacion de precio desde.",
      "created_at": "2026-05-20T15:30:00Z"
    }
  ],
  "next_action_at": "2026-05-22T14:00:00Z"
}
```

### PUT /api/leads/:id

Updates lead status, notes, or next action date.

Allowed statuses:

```txt
Nuevo
Contactado
Interesado
Agendado
No Respondio
Perdido
Convertido
```

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
  "status": "Interesado",
  "updated_at": "2026-05-20T15:30:00Z"
}
```

### POST /api/ai/reply-suggestion

Generates safe WhatsApp-ready reply suggestions. The backend builds the prompt using clinic and service context.

Request:

```json
{
  "lead_id": "bce2f64d-3582-4eb3-a5f2-12c8b46b2181",
  "service_id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
  "patient_message": "El blanqueamiento debilita los dientes?",
  "desired_tone": "amable"
}
```

Response:

```json
{
  "generation_id": "7ed7c244-5734-44ff-95e6-5c8ebc1550bc",
  "variants": {
    "short": "Hola Maria, el blanqueamiento debe valorarse de forma personalizada para confirmar que sea adecuado para ti. Podemos agendar una valoracion y explicarte el proceso con calma.",
    "persuasive": "Hola Maria, entiendo tu duda. En la valoracion revisamos tu caso y te explicamos que opcion es adecuada para cuidar tu sonrisa y lograr un resultado natural. Te gustaria que revisemos disponibilidad esta semana?",
    "technical": "El blanqueamiento se indica despues de una valoracion profesional para revisar el estado dental y definir si eres candidata. La recomendacion final debe darla el odontologo en consulta.",
    "closing_question": "Te queda mejor una valoracion en la manana o en la tarde?"
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
  "service_id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
  "objection": "En otro lado es mas barato"
}
```

Response:

```json
{
  "generation_id": "9c9b14ad-c73e-4eba-a7d8-25989e2b1a0c",
  "objection_type": "price_comparison",
  "recommended_strategy": "Reforzar valor, evaluacion personalizada y confianza profesional sin atacar a otros proveedores.",
  "suggested_message": "Te entiendo. En estos tratamientos el precio puede variar segun la valoracion, materiales y seguimiento. Lo ideal es que podamos revisarte y explicarte una opcion adecuada para tu caso antes de que tomes una decision.",
  "closing_question": "Quieres que te ayude a separar una valoracion?",
  "safety_status": "passed"
}
```

### POST /api/content/generate-post

Generates marketing content for clinic channels.

Request:

```json
{
  "service_id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
  "content_type": "instagram_post",
  "topic": "Beneficios de una valoracion antes del blanqueamiento",
  "tone": "cercano"
}
```

Response:

```json
{
  "generation_id": "94b50a52-7730-4024-af6c-b731ea0c53dd",
  "title": "Antes de blanquear tu sonrisa, revisemos tu caso",
  "body": "Un blanqueamiento dental debe comenzar con una valoracion. Asi podemos revisar tus necesidades, resolver dudas y orientarte con una opcion adecuada.",
  "call_to_action": "Escribenos por WhatsApp y agenda tu valoracion.",
  "safety_status": "passed"
}
```


### GET /api/followups

Lists leads with a pending manual follow-up based on `next_action_at`.

Query parameters:

```txt
due=all|today|overdue|upcoming
page=1
page_size=20
```

Response uses the standard paginated lead response shape.

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
  "id": "lead-id",
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
  "id": "lead-id",
  "status": "rescheduled"
}
```

### POST /api/ai/follow-up-message

Generates a safe WhatsApp-ready re-engagement message for a pending lead.

Request:

```json
{
  "lead_id": "lead-id",
  "service_id": "service-id",
  "last_contact_note": "Pregunto por precio y no respondio despues de la primera respuesta."
}
```

Response:

```json
{
  "generation_id": "gen-clinic-id",
  "suggested_message": "Hola Maria, queria saber si te quedo alguna duda sobre el blanqueamiento dental...",
  "recommended_timing": "Dia siguiente en horario laboral.",
  "next_step": "Si responde, agendar valoracion. Si no responde, reprogramar seguimiento.",
  "safety_status": "passed"
}
```

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
      "service_id": "c35bb8e5-73de-4f93-aa9a-42c7aa2cf67d",
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

## Notes for Future Changes

- Add explicit superadmin endpoints when platform administration is implemented.
- Add user management endpoints after the authentication base is stable.
- Add AI usage limits and billing metadata when subscription plans are introduced.
- Keep WhatsApp Business Cloud API integration out of the MVP contract unless scope is explicitly changed.
- Do not add clinical-history, diagnosis, prescription, or medical-record endpoints.
