# Database Migrations

This directory contains SQL migration files for ClinicFlow AI.

## Current Status

Initial core schema migration exists for `clinics`, `users`, and `clinic_services`. Migration execution tooling is still pending.

## Structure

Migrations should be named following the pattern: `YYYYMMDDHHMMSS_description.sql`.

## Principles

- Use explicit SQL.
- Favor non-destructive changes.
- Ensure `clinic_id` is present in all tenant-owned tables.
- Use UUID for primary keys.
- Add timestamps where useful.
- Add indexes for frequently filtered columns.
- Do not store clinical records, diagnoses, prescriptions, or medical history.

## Tooling

Migration execution tool to be defined during Phase 1. A Go-based migration tool is recommended.
