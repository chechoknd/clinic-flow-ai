# Smart Schedule Database Plan

Status: Approved for first migration.

This document plans the database layer for the Smart Schedule / Agenda Inteligente direction. The first migration has been created in `database/migrations/20260601000100_create_smart_schedule_schema.sql`. Backend modules, frontend screens, and API handlers are still pending.

## Purpose

The database must support a visual commercial schedule for dental clinics while preserving the current product boundaries:

- Operational appointments, not clinical records.
- Dentists/professionals as schedule resources, not clinical profile records.
- Lead-to-appointment conversion with traceability.
- Manual confirmation and follow-up workflows.
- Strict clinic-level tenant isolation.
- No autonomous booking or WhatsApp sending.

## Current Database Baseline

Implemented tables:

- `clinics`
- `users`
- `clinic_services`
- `leads`
- `lead_notes`
- `schema_migrations`

Relevant current conventions:

- PostgreSQL.
- UUID primary keys through `gen_random_uuid()`.
- `TIMESTAMPTZ` for dates/times.
- `created_at` and `updated_at` where useful.
- `set_updated_at()` trigger already exists.
- Clinic-scoped entities include `clinic_id`.
- Lead statuses currently are Spanish labels:

```txt
Nuevo
Contactado
Interesado
Agendado
No Respondio
Perdido
Convertido
```

## Planned Tables

First Smart Schedule migration tables:

- `clinic_professionals`
- `professional_services`
- `appointments`

Planned later, not first migration unless approved:

- `appointment_events`
- `blocked_schedule_slots`
- `professional_time_off`

## Table: clinic_professionals

Status: Migration created.

Purpose: store operational dentists/professionals that can be assigned to appointments.

This table must not store clinical biography, medical credentials files, clinical notes, or sensitive health data.

Proposed columns:

| Column | Type | Required | Notes |
| --- | --- | --- | --- |
| `id` | `UUID` | yes | Primary key, default `gen_random_uuid()` |
| `clinic_id` | `UUID` | yes | FK to `clinics(id)` with `ON DELETE CASCADE` |
| `full_name` | `VARCHAR(150)` | yes | Display name |
| `role_or_specialty` | `VARCHAR(120)` | no | Operational specialty or role, e.g. `Ortodoncia` |
| `calendar_color` | `VARCHAR(20)` | no | UI identifier, e.g. `#2563EB` |
| `working_hours` | `JSONB` | yes | Weekly availability, default `{}` |
| `is_active` | `BOOLEAN` | yes | Default `TRUE` |
| `created_at` | `TIMESTAMPTZ` | yes | Default `NOW()` |
| `updated_at` | `TIMESTAMPTZ` | yes | Default `NOW()` |

Suggested constraints:

```sql
CONSTRAINT clinic_professionals_clinic_id_id_unique UNIQUE (clinic_id, id)
CONSTRAINT clinic_professionals_name_per_clinic_unique UNIQUE (clinic_id, full_name)
CONSTRAINT clinic_professionals_calendar_color_check CHECK (
    calendar_color IS NULL OR calendar_color ~ '^#[0-9A-Fa-f]{6}$'
)
```

Suggested indexes:

```sql
CREATE INDEX clinic_professionals_clinic_id_idx ON clinic_professionals(clinic_id);
CREATE INDEX clinic_professionals_active_idx ON clinic_professionals(clinic_id, is_active);
```

Working-hours shape should stay simple for MVP:

```json
{
  "monday": [{"start": "08:00", "end": "12:00"}, {"start": "14:00", "end": "18:00"}],
  "tuesday": [{"start": "08:00", "end": "12:00"}],
  "wednesday": [],
  "thursday": [],
  "friday": [],
  "saturday": [],
  "sunday": []
}
```

Validation note: JSONB shape validation can start in backend service logic. A database check for the full JSON structure is not required for the first migration.

## Table: professional_services

Status: Migration created.

Purpose: map which services a professional can perform.

This is a many-to-many relationship between `clinic_professionals` and `clinic_services`.

Proposed columns:

| Column | Type | Required | Notes |
| --- | --- | --- | --- |
| `clinic_id` | `UUID` | yes | Direct tenant scope |
| `professional_id` | `UUID` | yes | FK to `clinic_professionals(id)` |
| `service_id` | `UUID` | yes | FK to `clinic_services(id)` |
| `created_at` | `TIMESTAMPTZ` | yes | Default `NOW()` |

Suggested primary key:

```sql
PRIMARY KEY (professional_id, service_id)
```

Suggested foreign keys:

```sql
clinic_id REFERENCES clinics(id) ON DELETE CASCADE
professional_id REFERENCES clinic_professionals(id) ON DELETE CASCADE
service_id REFERENCES clinic_services(id) ON DELETE CASCADE
```

Suggested indexes:

```sql
CREATE INDEX professional_services_clinic_id_idx ON professional_services(clinic_id);
CREATE INDEX professional_services_service_id_idx ON professional_services(clinic_id, service_id);
CREATE INDEX professional_services_professional_id_idx ON professional_services(clinic_id, professional_id);
```

Tenant validation rule:

- The migration adds composite tenant constraints for professional and service assignment.
- Backend must still verify that `clinic_id`, `professional_id`, and `service_id` all belong to the authenticated clinic before insert/update.

## Table: appointments

Status: Migration created.

Purpose: store operational appointments for the visual agenda.

Appointments are commercial/admin records. They must not store diagnosis, prescriptions, treatment plans, clinical images, or clinical notes.

Proposed columns:

| Column | Type | Required | Notes |
| --- | --- | --- | --- |
| `id` | `UUID` | yes | Primary key, default `gen_random_uuid()` |
| `clinic_id` | `UUID` | yes | FK to `clinics(id)` with `ON DELETE CASCADE` |
| `professional_id` | `UUID` | yes | FK to `clinic_professionals(id)` with `ON DELETE RESTRICT` |
| `lead_id` | `UUID` | no | FK to `leads(id)` with `ON DELETE SET NULL` |
| `service_id` | `UUID` | yes | FK to `clinic_services(id)` with `ON DELETE RESTRICT` |
| `contact_name` | `VARCHAR(150)` | yes | Snapshot for agenda display |
| `contact_phone` | `VARCHAR(30)` | no | Snapshot for confirmation/follow-up |
| `starts_at` | `TIMESTAMPTZ` | yes | Appointment start |
| `ends_at` | `TIMESTAMPTZ` | yes | Appointment end |
| `status` | `VARCHAR(50)` | yes | Stable API value |
| `confirmation_status` | `VARCHAR(50)` | yes | Default `pending` |
| `source` | `VARCHAR(50)` | yes | Default `manual` |
| `admin_notes` | `TEXT` | no | Commercial/admin only |
| `created_by_user_id` | `UUID` | no | FK to `users(id)` with `ON DELETE SET NULL` |
| `created_at` | `TIMESTAMPTZ` | yes | Default `NOW()` |
| `updated_at` | `TIMESTAMPTZ` | yes | Default `NOW()` |

Recommended appointment statuses:

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

Recommended confirmation statuses:

```txt
pending
confirmed
not_required
failed
```

Suggested constraints:

```sql
CONSTRAINT appointments_status_check CHECK (
    status IN (
        'scheduled',
        'confirmed',
        'pending_confirmation',
        'rescheduled',
        'no_show',
        'cancelled',
        'completed',
        'converted_from_lead'
    )
)

CONSTRAINT appointments_confirmation_status_check CHECK (
    confirmation_status IN ('pending', 'confirmed', 'not_required', 'failed')
)

CONSTRAINT appointments_time_check CHECK (ends_at > starts_at)
```

Migration note:

- `professional_id` and `service_id` use composite FKs with `clinic_id` to enforce tenant consistency at the database level.
- `lead_id` and `created_by_user_id` use nullable single-column FKs with `ON DELETE SET NULL`; backend service validation remains mandatory to ensure they belong to the authenticated clinic before insert/update.

Suggested indexes:

```sql
CREATE INDEX appointments_clinic_date_idx ON appointments(clinic_id, starts_at);
CREATE INDEX appointments_professional_date_idx ON appointments(clinic_id, professional_id, starts_at);
CREATE INDEX appointments_status_date_idx ON appointments(clinic_id, status, starts_at);
CREATE INDEX appointments_confirmation_date_idx ON appointments(clinic_id, confirmation_status, starts_at);
CREATE INDEX appointments_lead_id_idx ON appointments(clinic_id, lead_id);
CREATE INDEX appointments_service_date_idx ON appointments(clinic_id, service_id, starts_at);
```

## Overlap Prevention

Recommended MVP rule: prevent overlapping active appointments for the same professional.

Active statuses for overlap checks:

```txt
scheduled
confirmed
pending_confirmation
rescheduled
converted_from_lead
```

Non-blocking statuses:

```txt
cancelled
no_show
completed
```

Preferred first implementation:

- Enforce overlap checks in the backend service layer inside a transaction.
- Query for overlapping active appointments by `clinic_id`, `professional_id`, and time range.
- Return a `VALIDATION_ERROR` or conflict-style domain error when an overlap exists.

Optional later database-level guard:

- Use PostgreSQL range types and an exclusion constraint with `btree_gist`.
- Do not add this in the first migration unless the team accepts the extra extension and operational complexity.

## Lead-To-Appointment Conversion

When a lead becomes an appointment:

1. Backend validates authenticated clinic.
2. Backend validates the lead belongs to the clinic.
3. Backend validates the selected service belongs to the clinic.
4. Backend validates the selected professional belongs to the clinic.
5. Backend optionally validates the professional can provide the service through `professional_services`.
6. Backend creates the appointment.
7. Backend may update lead status to `Agendado` when `update_lead_status=true`.
8. Backend should append a commercial lead note such as `Cita creada para 2026-06-02 14:00`.

No clinical data should be copied into `appointments.admin_notes` or `lead_notes.body`.

## Tenant Isolation Rules

Every planned table includes `clinic_id` directly. This is intentional even when the clinic can be reached through another relation.

Required backend checks:

- `clinic_professionals.clinic_id` must match authenticated clinic.
- `professional_services.clinic_id` must match authenticated clinic.
- `appointments.clinic_id` must match authenticated clinic.
- Appointment `professional_id`, `service_id`, `lead_id`, and `created_by_user_id` must belong to the same clinic when present.
- List queries must always filter by `clinic_id`.
- Detail/update/delete queries must always include `clinic_id` in the lookup.

Frontend filters are not tenant isolation.

## Appointment Notes Safety

Allowed `admin_notes` examples:

- "Confirmar asistencia en la manana."
- "Solicito horario despues de las 4 pm."
- "Pidio reprogramar por disponibilidad."
- "Viene por valoracion inicial de ortodoncia."

Forbidden `admin_notes` examples:

- Diagnosis.
- Prescription.
- Clinical evolution note.
- Treatment plan.
- Lab result.
- Clinical image interpretation.
- Detailed symptoms stored as medical history.

Backend should apply the same safety posture used for lead notes: validate obvious unsafe content where practical and keep API errors sanitized.

## Suggested First Migration Shape

The first migration now:

1. Create `clinic_professionals`.
2. Add trigger `clinic_professionals_set_updated_at`.
3. Create `professional_services`.
4. Create `appointments`.
5. Add trigger `appointments_set_updated_at`.
6. Add indexes listed above.

Do not create:

- `appointment_events`.
- blocked slots.
- time-off tables.
- external calendar sync tables.
- WhatsApp integration tables.

Those are later phases.

## Open Questions

- Should `converted_from_lead` be a long-term appointment status or a `source` value? Current recommendation: keep it documented but consider `source='lead_conversion'` during implementation to avoid lifecycle confusion.
- Should `confirmation_status='confirmed'` automatically imply `status='confirmed'`, or should they remain separate? Current recommendation: keep both for MVP but enforce a clear service rule.
- Should contact phone be required for appointments created from leads? Current recommendation: require it when confirmation workflows depend on WhatsApp.
- Should professional availability live only in `working_hours` JSONB for MVP? Current recommendation: yes, then add time-off/blocked-slot tables later.
- Should appointments allow no professional assignment? Current recommendation: no for MVP, because professional assignment is core to the schedule value.

## Next Implementation Step

The migration has been created:

```txt
database/migrations/20260601000100_create_smart_schedule_schema.sql
```

Next, update backend models and repositories in planned modules:

- `internal/professionals`
- `internal/appointments`
- `internal/leads` conversion use case
- `internal/dashboard` schedule summary
