# Database Seeds

This directory contains SQL seeds for local development and testing.

## Current Status

Initial demo seed data exists for one clinic, one clinic admin user, and sample dental services.

## Intended Usage

Apply seeds after running migrations:

```bash
psql "postgres://clinicflow:clinicflow@localhost:5432/clinicflow_db?sslmode=disable" -v ON_ERROR_STOP=1 -f database/seeds/20260521000100_demo_core_data.sql
```

Demo login for local development only:

```txt
email: admin@sonrisaviva.demo
password: clinicflow123
```

Seeds are intended to populate the local database with initial demo data:

- One demo clinic admin user.
- One sample clinic.
- Sample dental services catalog.
- Sample commercial leads (future seed, not included yet).

## Principles

- Never include real patient data.
- Never include production credentials or secrets.
- Ensure consistency between clinic and related records.
- Keep demo data clearly fictional.
