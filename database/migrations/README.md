# Database Migrations

This directory contains the SQL migration files for ClinicFlow AI.

## Structure

Migrations should be named following the pattern: `YYYYMMDDHHMMSS_description.sql`.

## Principles

- Use explicit SQL.
- Favor non-destructive changes.
- Ensure `clinic_id` is present in all tenant-owned tables.
- Use UUID for primary keys.

## Tooling

Migration execution tool to be defined during Phase 1 (Go-based migration tool recommended).
