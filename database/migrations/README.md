# Database Migrations

This directory contains SQL migration files for ClinicFlow AI.

## Current Status

Migration tooling and the initial schema are not implemented yet. The next expected step is to add initial migrations for authentication, clinics, and users.

## Structure

Migrations should be named following the pattern: `YYYYMMDDHHMMSS_description.sql`.

## Principles

- Use explicit SQL.
- Favor non-destructive changes.
- Ensure `clinic_id` is present in all tenant-owned tables.
- Use UUID for primary keys.
- Add timestamps where useful.
- Add indexes for frequently filtered columns.

## Tooling

Migration execution tool to be defined during Phase 1. A Go-based migration tool is recommended.
