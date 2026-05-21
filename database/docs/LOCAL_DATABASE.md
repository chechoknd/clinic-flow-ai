# Local Database Infrastructure

This document describes the local database setup for ClinicFlow AI.

## Technology Stack

- **Engine:** PostgreSQL 16 (Alpine-based Docker image).
- **Primary Keys:** UUID v4 (recommended).
- **Isolation:** logical multi-tenancy using `clinic_id`.

## Local Connection

- **Host:** `localhost`
- **Port:** `5432`
- **User:** `clinicflow`
- **Password:** `clinicflow` (see `.env.example`)
- **Database:** `clinicflow_db`
- **SSL:** Disabled for local development.

## Docker Compose Configuration

The PostgreSQL service is managed by Docker Compose:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: clinic-flow-ai-postgres
    environment:
      POSTGRES_DB: clinicflow_db
      POSTGRES_USER: clinicflow
      POSTGRES_PASSWORD: clinicflow
    ports:
      - "5432:5432"
    volumes:
      - clinic_flow_ai_postgres_data:/var/lib/postgresql/data
```

## Data Persistence

A Docker volume named `clinic_flow_ai_postgres_data` is used to persist data between container restarts.

## Security Note

Credentials in `docker-compose.yml` are for local development only. Production credentials must never be committed and should be managed via environment variables.
