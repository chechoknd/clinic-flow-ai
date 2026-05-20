# Database Rules

## Conventions

- PostgreSQL.
- UUID primary keys.
- Explicit SQL migrations.
- Seeds for local/demo data.
- Timestamps on tables where useful.
- Foreign keys.
- Indexes for frequently filtered columns.
- Schema aligned with API and domain models.
- Avoid destructive migrations unless explicitly requested.

## Recommended Base Tables

- `clinics`
- `clinic_services`
- `users`
- `leads`
- `lead_notes`
- `ai_generations`
- `followups`

## Allowed Lead Data

- Full name.
- Phone.
- Service of interest.
- Commercial status.
- Notes related to sales conversation.
- Next follow-up date.
- Basic source/channel.

## Forbidden Medical Data

- Diagnoses.
- Medical evolution notes.
- Clinical images.
- Prescriptions.
- Lab results.
- Medical records.
- Detailed symptoms as medical history.

## Multi-Tenant Rules

Every entity must belong to a clinic when applicable. Backend queries must enforce data isolation between clinics. Never rely solely on frontend filters.

## pgvector

Optional `pgvector` extension may be prepared for future phases, but is not required for MVP logic.
