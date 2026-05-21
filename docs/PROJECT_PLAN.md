# Project Plan

## Product Vision

ClinicFlow AI is a SaaS platform for small and medium private clinics that improves commercial attention, lead follow-up, and patient conversion through AI-assisted WhatsApp communication.

The first MVP focuses on dental clinics and avoids heavy clinical, hospital, or medical-record workflows. The product promise is:

> Help clinics respond better, recover interested patients, and sell more treatments through WhatsApp.

ClinicFlow AI is a commercial communication assistant for clinic staff. It is not a medical diagnosis system, prescription tool, telemedicine platform, hospital system, medical records system, or autonomous WhatsApp bot in the MVP.

## Problem Statement

Private clinics receive many interested patients through WhatsApp, especially around pricing, appointment availability, treatment fears, and service details. Reception teams often lack time, structure, and commercial support to answer consistently, follow up at the right time, and convert interest into booked evaluations.

This creates a commercial "black hole": leads ask a question, receive a weak or delayed answer, and disappear without structured follow-up.

ClinicFlow AI addresses this by giving staff fast, safe, context-aware message suggestions, objection handling, follow-up reminders, and simple CRM visibility.

## Target Users

- Small and medium private clinics.
- Dental clinics and independent dental offices for the first MVP.
- Clinic owners who need more predictable lead conversion.
- Clinic administrators who manage service catalogs, staff, and daily commercial performance.
- Assistants or receptionists who answer WhatsApp messages and follow up with leads.

## Initial Niche: Dental Clinics

The MVP is optimized for dental clinics because they commonly have:

- High-value private treatments such as orthodontics, implants, smile design, whitening, and oral rehabilitation.
- Repetitive WhatsApp questions about prices, fears, availability, and procedure expectations.
- Reception staff overloaded with operational tasks.
- Strong dependence on trust-building messages before a patient books an evaluation.
- Ongoing need for educational and promotional social media content.

Initial dental services may include dental cleaning, whitening, orthodontics, smile design, implants, oral rehabilitation, dental emergencies, initial evaluation, and preventive controls.

## MVP Scope

The MVP includes the following planned capabilities:

- Authentication and roles.
- Clinic profile and configuration.
- Clinic service catalog.
- Simple commercial CRM for leads.
- AI reply assistant for WhatsApp.
- AI objection handler.
- Manual assisted follow-ups.
- Basic dashboard.
- AI content generator for marketing.
- PostgreSQL database.
- Docker Compose local environment.
- Angular frontend.
- Go backend.
- REST API.
- JWT authentication.
- Environment-based AI provider configuration.

## MVP Exclusions

The MVP strictly excludes:

- Clinical histories.
- Medical records.
- Diagnoses.
- Prescriptions.
- Clinical decision support.
- Medical image interpretation.
- Payment gateways.
- Electronic invoicing.
- Native mobile apps.
- Telemedicine.
- Video calls.
- Full appointment scheduling system.
- Native bidirectional WhatsApp Business Cloud API integration.
- Autonomous WhatsApp bots.
- RAG/vector search as required MVP infrastructure.

If a task introduces excluded scope, it must be treated as a scope conflict.

## Main Modules

- Authentication: login, JWT issuance, role-aware access.
- Clinics: clinic profile, city, operating WhatsApp, address, hours, FAQs, communication tone.
- Services: dental service catalog with descriptions, benefits, "price from", FAQs, and common objections.
- Leads: simple commercial pipeline with status, service of interest, notes, source, and next action date.
- AI Assistant: safe WhatsApp-ready reply suggestions based on clinic and service context.
- Objection Handler: categorization and response strategy for commercial objections.
- Follow-ups: manual reminders and AI-assisted re-engagement messages for stale leads.
- Content Generator: marketing copy, post ideas, Reel/TikTok scripts, carousel outlines, and WhatsApp campaign text.
- Dashboard: basic commercial KPIs for leads, service demand, pending follow-ups, and AI usage.

## User Roles

- `superadmin`: manages platform-level concerns, aggregated business visibility, and subscription administration.
- `clinic_admin`: manages clinic configuration, users, services, leads, and clinic analytics.
- `assistant`: operates daily CRM workflows, updates lead status, generates AI responses, and manages manual follow-ups.

## Development Phases

### Phase 0: Documentation and Repository Foundation

Establish repository structure, agent governance, project documentation, product scope, architecture direction, API contracts, and development status tracking.

### Phase 1: Technical Base

Set up the Go REST API foundation, PostgreSQL connection pool, JWT authentication, Angular application shell, protected routing, HTTP interceptors, Docker Compose services, and environment variable conventions.

### Phase 2: Clinics and Services

Implement clinic profile configuration, service catalog storage, SQL migrations, seeds, backend CRUD endpoints, and Angular forms/views for clinic and service management.

### Phase 3: CRM Leads Core

Implement the commercial lead model, lead pipeline statuses, lead notes, next action dates, filtering, pagination, backend persistence, and Angular lead management workflows.

### Phase 4: AI Assistant

Implement AI provider abstraction, backend-owned prompts, direct structured context injection, safety validation, reply suggestion endpoint, objection handling endpoint, and frontend copy-to-clipboard workflows.

### Phase 5: Marketing and Follow-ups

Implement AI content generation, manual follow-up reminders, follow-up message suggestions, priority indicators, and operational views for pending re-engagement.

### Phase 6: MVP Closure and Pilot

Finalize dashboard summary, integration testing under Docker Compose, safety review, documentation updates, pilot readiness checklist, and pilot execution with 5 to 10 independent dental clinics.

## Success Criteria for MVP

- A clinic can configure its profile, communication tone, and service catalog.
- Staff can create, view, filter, and update commercial leads.
- Staff can generate safe WhatsApp reply suggestions using clinic and service context.
- Staff can generate objection-handling responses that avoid diagnosis and medical advice.
- Staff can manage manual follow-ups without native WhatsApp API automation.
- Clinic administrators can see basic commercial KPIs.
- AI prompts and provider configuration are controlled by the backend.
- The system enforces role-aware access and clinic-level data isolation.
- Local development works through Docker Compose without committed secrets.
- Pilot clinics can use the system for daily WhatsApp-assisted commercial workflows.

## Risks and Mitigations

- Low adoption by reception staff: keep workflows simple, fast, and focused on copy-to-clipboard actions.
- AI medical safety risk: enforce backend-owned prompt safety, output validation, and fallback responses.
- Scope creep into clinical systems: keep medical records, diagnosis, prescriptions, and clinical decision support out of MVP.
- WhatsApp automation complexity: use manual assisted follow-ups and exclude native bidirectional WhatsApp Cloud API integration from MVP.
- Multi-tenant data leakage: enforce clinic isolation in backend queries and authorization checks.
- AI provider cost or reliability issues: use provider abstraction, token limits, timeouts, fallbacks, and usage metadata.
- Product positioning confusion: communicate ClinicFlow AI as a commercial assistant, not "AI for medicine."

## Future Expansion Opportunities

- WhatsApp Business Cloud API integration after MVP validation.
- Chrome extension or lightweight overlay for faster WhatsApp Web workflows.
- RAG/vector search for richer clinic knowledge bases when justified by usage.
- Advanced analytics for conversion rates, response quality, and service demand.
- Additional verticals such as ophthalmology, aesthetics and dermatology, physiotherapy, and veterinary clinics.
- Subscription management and payment gateway integration.
- Multi-location clinic support.
- Automated campaign segmentation.
- Native mobile apps only if validated by demand.
