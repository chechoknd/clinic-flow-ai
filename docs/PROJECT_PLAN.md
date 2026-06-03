# Project Plan

## Product Vision

ClinicFlow AI is a SaaS platform for small and medium private clinics that improves commercial attention, appointment-oriented operations, lead follow-up, and patient conversion through an AI-assisted Smart Schedule.

The first MVP focuses on dental clinics and avoids heavy clinical, hospital, or medical-record workflows. The product promise is:

> Help clinics organize their daily agenda, respond better, recover interested patients, and convert more conversations into booked appointments.

ClinicFlow AI is evolving from an AI-assisted commercial CRM into a smart commercial schedule for dental clinics. The schedule connects appointments, dentists/professionals, leads, services, follow-ups, and AI-assisted communication.

ClinicFlow AI remains a commercial and operational assistant for clinic staff. It is not a medical diagnosis system, prescription tool, telemedicine platform, hospital system, medical records system, or autonomous WhatsApp bot in the MVP.

## Problem Statement

Private clinics receive many interested patients through WhatsApp, especially around pricing, appointment availability, treatment fears, and service details. Reception teams often lack time, structure, schedule visibility, and commercial support to answer consistently, follow up at the right time, and convert interest into booked evaluations.

This creates a commercial "black hole": leads ask a question, receive a weak or delayed answer, and disappear without structured follow-up.

ClinicFlow AI addresses this by giving staff a smart operational agenda plus fast, safe, context-aware message suggestions, objection handling, follow-up reminders, and simple CRM visibility.

## Target Users

- Small and medium private clinics.
- Dental clinics and independent dental offices for the first MVP.
- Clinic owners who need more predictable lead conversion.
- Clinic administrators who manage service catalogs, staff, and daily commercial performance.
- Assistants or receptionists who answer WhatsApp messages, confirm appointments, reschedule, and follow up with leads.
- Dentists/professionals who need visibility into their assigned schedule without clinical-record workflows.

## Initial Niche: Dental Clinics

The MVP is optimized for dental clinics because they commonly have:

- High-value private treatments such as orthodontics, implants, smile design, whitening, and oral rehabilitation.
- Repetitive WhatsApp questions about prices, fears, availability, and procedure expectations.
- Reception staff overloaded with operational tasks.
- Strong dependence on trust-building messages before a patient books an evaluation.
- Ongoing need for educational and promotional social media content.

Initial dental services may include dental cleaning, whitening, orthodontics, smile design, implants, oral rehabilitation, dental emergencies, initial evaluation, and preventive controls.

## Product Direction: Smart Schedule, Smart Lead Inbox, and Intelligent Dashboard

Status: Partially Implemented.

ClinicFlow AI is now evolving around the Smart Schedule / Agenda Inteligente as the product center. The first schedule slice has been implemented across database, backend, and Angular screens.

The new product direction is:

> A smart commercial schedule for dental clinics that connects appointments, dentists, leads, services, follow-ups, and AI assistance.

The schedule should help clinic staff manage:

- Appointments.
- Dentists/professionals.
- Leads that can become appointments.
- Follow-ups.
- Appointment status.
- Available time slots.
- Pending confirmations.
- Daily commercial priorities.

ClinicFlow AI is planned to evolve toward a Smart Lead Inbox experience: a commercial inbox assisted by AI where staff can paste conversations, analyze commercial intent, detect service interest and objections, generate safe response drafts, and decide the next follow-up action.

This direction keeps the existing product rule:

```txt
AI suggests -> human reviews -> human replies
```

The system must not send autonomous WhatsApp messages in the MVP. The assistant or clinic staff member remains responsible for reviewing, editing, copying, and sending any response through the original channel.

The dashboard is also planned to evolve from mostly metrics into an intelligent action layer that helps staff decide what to handle today: today's appointments, pending confirmations, available schedule opportunities, hot leads without appointment, overdue follow-ups, detected objections, service demand, and AI-recommended next actions.

See `docs/SMART_SCHEDULE_PLAN.md` for the detailed Smart Schedule plan.
See `docs/SMART_LEAD_INBOX_PLAN.md` for the detailed product plan.

## MVP Scope

The MVP includes the following planned capabilities:

- Authentication and roles.
- Clinic profile and configuration.
- Clinic service catalog.
- Smart Schedule / Intelligent Agenda first slice.
- Dentists/professionals operational module.
- Appointment management connected to leads and services.
- Simple commercial CRM for leads.
- AI reply assistant for WhatsApp.
- AI objection handler.
- Manual assisted follow-ups.
- Basic dashboard.
- Smart Lead Inbox for manual conversation analysis.
- Intelligent dashboard action layer first slice.
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
- Clinical-grade scheduling, treatment planning, or autonomous appointment booking.
- Native bidirectional WhatsApp Business Cloud API integration.
- Autonomous WhatsApp bots.
- RAG/vector search as required MVP infrastructure.

If a task introduces excluded scope, it must be treated as a scope conflict.

## Main Modules

- Authentication: login, JWT issuance, role-aware access.
- Clinics: clinic profile, city, operating WhatsApp, address, hours, FAQs, communication tone.
- Services: dental service catalog with descriptions, benefits, "price from", FAQs, and common objections.
- Professionals: implemented operational module for dentists/professionals, specialties or roles, active state, service assignment, availability, and calendar colors.
- Smart Schedule: first visual daily/weekly agenda slice implemented with appointments, statuses, professional assignment, services, lead/contact reference, confirmations, available slots, and rescheduling.
- Appointments: implemented operational records connected to clinic, professional, service, optional lead/contact, date/time, status, source, confirmation state, and commercial/admin notes only.
- Leads: simple commercial pipeline with status, service of interest, notes, source, and next action date.
- AI Assistant: safe WhatsApp-ready reply suggestions based on clinic and service context.
- Objection Handler: categorization and response strategy for commercial objections.
- Follow-ups: manual reminders and AI-assisted re-engagement messages for stale leads.
- Content Generator: marketing copy, post ideas, Reel/TikTok scripts, carousel outlines, and WhatsApp campaign text.
- Dashboard: basic commercial KPIs plus first schedule-centered action cards for today's appointments, pending confirmations, available slots, hot leads without appointment, and overdue follow-ups.
- Smart Lead Inbox: first manual conversation analysis slice implemented for pasted conversations, commercial extraction, response suggestion, and reviewed lead creation. Planned evolution: detect scheduling intent and suggest reviewed appointment creation.


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

### Phase 6: Smart Lead Inbox Product Alignment

Document and validate the Smart Lead Inbox direction, Dashboard Inteligente action layer, conceptual API contracts, safety rules, and implementation phases before writing code.

### Phase 7: Manual Conversation Analysis

Partially implemented: manual conversation paste screen, AI conversation analysis, extracted commercial fields, suggested response, and reviewed lead creation through the existing leads API. Pending: update existing lead from analysis and persisted analysis history.

### Phase 8: Smart Schedule Product Alignment

Implemented documentation alignment for the Smart Schedule direction, professionals model, appointment model, lead-to-appointment flow, dashboard changes, safety rules, and draft API contracts.

### Phase 9: Smart Schedule Database and Backend First Slice

Partially implemented: tables for professionals, professional services, and appointments; backend professionals endpoints; appointment CRUD/status/reschedule/lead-conversion endpoints; schedule availability endpoint; and schedule-centered dashboard summary.

### Phase 10: Smart Schedule Frontend First Slice

Partially implemented: Angular schedule view, professional management, appointment create/edit flows, lead links into schedule, pending confirmation states, and dashboard schedule priorities.

### Phase 11: Dashboard Action Layer

Partially implemented: dashboard has a schedule-first tab with daily appointment, pending confirmation, hot lead, overdue follow-up, professional, and service demand signals. Further action scoring remains pending.

### Phase 12: Inbox AI Schedule Integration

Align Inbox AI with the schedule-centered workflow: detect scheduling intent, identify missing scheduling data, suggest appointment creation, and keep human-reviewed lead/appointment traceability.

### Phase 13: MVP Closure and Pilot

Finalize dashboard summary, integration testing under Docker Compose, safety review, documentation updates, pilot readiness checklist, and pilot execution with 5 to 10 independent dental clinics.

## Success Criteria for MVP

- A clinic can configure its profile, communication tone, and service catalog.
- A clinic can manage professionals/dentists as operational schedule resources.
- Staff can view a visual daily/weekly agenda in the first Smart Schedule slice.
- Staff can create, edit, confirm, cancel, reschedule, and complete operational appointments without storing clinical data.
- Staff can create, view, filter, and update commercial leads.
- Staff can convert a lead into an appointment after human review.
- Staff can generate safe WhatsApp reply suggestions using clinic and service context.
- Staff can generate objection-handling responses that avoid diagnosis and medical advice.
- Staff can manage manual follow-ups without native WhatsApp API automation.
- Staff can paste a conversation, receive AI commercial analysis, review it, and manually create or update a lead when Smart Lead Inbox is implemented.
- Staff can use the first intelligent dashboard action layer to review appointments pending confirmation, available slots, hot leads without appointment, and overdue follow-ups.

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
- Inbox AI scope creep: keep conversation capture manual or semi-manual first, and require human approval before saving lead updates, follow-ups, or sending messages.
- Smart Schedule scope creep: keep appointment data operational and commercial; do not add clinical notes, treatment plans, diagnosis, prescriptions, or autonomous booking.
- Schedule complexity risk: start with daily/weekly views, simple professional assignment, service duration, and clear statuses before advanced calendar features.

- Multi-tenant data leakage: enforce clinic isolation in backend queries and authorization checks.
- AI provider cost or reliability issues: use provider abstraction, token limits, timeouts, fallbacks, and usage metadata.
- Product positioning confusion: communicate ClinicFlow AI as a commercial assistant, not "AI for medicine."

## Future Expansion Opportunities

- Smart Lead Inbox with manual conversation paste, AI commercial analysis, and human-reviewed lead creation/update.
- Smart Schedule with professionals, appointments, availability, confirmations, and lead-to-appointment conversion.
- Intelligent dashboard action layer for daily assistant priorities.
- WhatsApp Business Cloud API integration after MVP validation.

- Chrome extension or lightweight overlay for faster WhatsApp Web workflows.
- RAG/vector search for richer clinic knowledge bases when justified by usage.
- Advanced analytics for conversion rates, response quality, and service demand.
- Additional verticals such as ophthalmology, aesthetics and dermatology, physiotherapy, and veterinary clinics.
- Subscription management and payment gateway integration.
- Multi-location clinic support.
- Automated campaign segmentation.
- Native mobile apps only if validated by demand.
