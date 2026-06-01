# Smart Schedule Plan

Status: Planned.

Smart Schedule, also called Agenda Inteligente, is the planned product center for ClinicFlow AI. It evolves the product from an AI-assisted lead CRM into a smart commercial schedule for dental clinics that connects appointments, dentists or professionals, leads, services, follow-ups, and AI assistance.

This plan is documentation and product alignment only. No frontend components, backend handlers, migrations, or tests have been implemented for this module yet.

## Purpose

Dentists and clinic teams understand value quickly when the product helps them manage the day around appointments. The schedule should become the operational surface where the clinic sees who is coming, who needs confirmation, which dentist is assigned, which leads can still become appointments, and which commercial actions should happen today.

The product direction becomes:

```txt
Smart commercial schedule -> lead conversion -> appointment confirmation -> manual follow-up -> AI-assisted communication
```

The AI rule remains unchanged:

```txt
AI suggests -> human reviews -> human responds
```

ClinicFlow AI must remain a commercial and operational assistant. It must not become a clinical history system, diagnosis system, prescription tool, telemedicine product, or autonomous WhatsApp bot.

## Main Users

- Receptionists who manage WhatsApp conversations, confirmations, rescheduling, and daily follow-ups.
- Clinic administrators who configure services, professionals, availability, and commercial priorities.
- Dentists or professionals who need clear visibility into their assigned appointments without clinical-record workflows.
- Clinic owners who want to understand conversion from leads into appointments.

## MVP Scope

The Smart Schedule MVP should include planning for:

- Daily schedule view.
- Weekly schedule view.
- Appointment cards with clear status and color cues.
- Dentist/professional assignment.
- Service assignment.
- Lead or contact reference.
- Start and end date/time or duration.
- Available slots based on professional working hours.
- Pending confirmation state.
- Rescheduling flow.
- Manual follow-up reminders connected to leads and appointments.
- Dashboard priorities based on the schedule.
- Human-reviewed AI suggestions for confirmations, follow-ups, and lead-to-appointment replies.

The MVP should keep scheduling operational and commercial. It should not add clinical notes, treatment plans, prescriptions, diagnosis, or medical history.

## Out of Scope

The Smart Schedule MVP must not include:

- Clinical histories.
- Diagnoses.
- Prescriptions.
- Medical decisions.
- Clinical notes.
- Clinical images.
- Telemedicine or video calls.
- Autonomous WhatsApp sending.
- AI medical advice.
- Automatic appointment booking without human review.
- Complex treatment planning.
- Insurance, billing, invoicing, or payment gateway workflows.

## Core Concepts

### Clinic

A clinic owns all operational data. Every schedule-related entity must be scoped by `clinic_id` and protected by backend authorization.

A clinic has many:

- Users.
- Services.
- Dentists/professionals.
- Leads.
- Appointments.

### Dentist / Professional

A professional is an operational schedule resource. The module should avoid clinical-profile detail.

Conceptual fields:

- `id`.
- `clinic_id`.
- `full_name`.
- `role_or_specialty`.
- `is_active`.
- `calendar_color`.
- `working_hours`.
- `created_at`.
- `updated_at`.

Relationships:

- Belongs to one clinic.
- Can provide many services.
- Has many appointments.

### Service

The current service catalog remains central. Services should connect to professionals and appointments.

Relationships:

- A service belongs to a clinic.
- A service may be provided by many professionals.
- A service may be linked to many leads.
- A service may be linked to many appointments.

### Lead

The lead CRM remains valid. A lead is a commercial opportunity and may become an appointment.

Relationships:

- A lead belongs to one clinic.
- A lead may reference one service of interest.
- A lead may have many notes.
- A lead may have follow-up dates.
- A lead may be converted into one or more appointments, depending on later business rules.

Lead statuses should remain commercial. `Agendado` can mean the lead has an appointment scheduled. `Convertido` can mean the clinic considers the commercial opportunity successfully converted according to agreed rules.

### Appointment

An appointment is an operational schedule item, not a clinical record.

Conceptual fields:

- `id`.
- `clinic_id`.
- `professional_id`.
- `lead_id` optional.
- `service_id`.
- `contact_name`.
- `contact_phone`.
- `starts_at`.
- `ends_at` or `duration_minutes`.
- `status`.
- `source`.
- `confirmation_status`.
- `admin_notes`.
- `created_by_user_id`.
- `created_at`.
- `updated_at`.

Allowed appointment notes are commercial/admin only, for example:

- "Confirmar asistencia en la manana."
- "Solicito horario despues de las 4 pm."
- "Viene por valoracion inicial de ortodoncia."

Forbidden appointment notes:

- Diagnoses.
- Prescriptions.
- Clinical evolution notes.
- Medical history.
- Lab results.
- Detailed symptoms stored as clinical information.

### Inbox AI Analysis

Inbox AI may suggest appointment-related actions from a pasted conversation, but it must not schedule automatically.

Potential analysis output:

- Detected scheduling intent.
- Detected service interest.
- Missing data before scheduling.
- Suggested appointment action.
- Suggested reply draft.
- Suggested follow-up.
- Suggested appointment duration from service context, if available.

Any lead update or appointment creation requires human review.

## Appointment Statuses

Use Spanish-first labels in the frontend while keeping API values stable and documented.

Planned appointment statuses:

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

Suggested Spanish labels:

```txt
Programada
Confirmada
Pendiente de confirmacion
Reprogramada
No asistio
Cancelada
Completada
Convertida desde lead
```

Status rules to validate later:

- A new appointment usually starts as `scheduled` or `pending_confirmation`.
- Confirmation changes status or confirmation state only after human action.
- Rescheduling should preserve traceability through notes or an appointment event history in a later phase.
- `completed` means the operational appointment occurred, not that a clinical outcome was recorded.
- `converted_from_lead` should be used carefully. It may be better as a source/event flag instead of a long-term status if it conflicts with normal appointment lifecycle.

## Main Workflows

### Create Appointment Manually

1. User opens Smart Schedule.
2. User clicks an available slot or `Nueva cita`.
3. User selects professional.
4. User selects service.
5. User enters contact or selects an existing lead.
6. User chooses date and time.
7. User adds admin notes only if needed.
8. User saves appointment.

### Convert Lead Into Appointment

1. User opens a lead or Inbox AI analysis.
2. User confirms that the person wants to schedule.
3. User selects professional, service, date, and time.
4. System creates appointment after human confirmation.
5. Lead status updates to `Agendado`.
6. Lead and appointment remain linked for traceability.

### Confirm Appointment

1. Dashboard or schedule shows `Pendiente de confirmacion`.
2. User opens appointment.
3. AI may suggest a short confirmation message.
4. User reviews, copies, and sends manually.
5. User marks appointment as confirmed after response.

### Reschedule Appointment

1. User opens appointment.
2. User selects new slot and optional reason as admin note.
3. System updates appointment to the new time.
4. User may create a manual follow-up if confirmation is still needed.

### Handle No-Show

1. User marks appointment as `no_show`.
2. System can suggest a respectful reactivation message.
3. User reviews and sends manually.
4. Lead may return to a follow-up state if appropriate.

## Visual Requirements

The schedule should be easy for non-technical reception staff.

Frontend planning principles:

- Daily view is the default.
- Weekly view is available for planning.
- Appointment cards show time, patient/contact name, service, professional, and status.
- Professionals use consistent colors.
- Statuses use clear labels and icons, not color alone.
- Available slots should be visually distinct from booked appointments.
- Pending confirmations should be easy to spot.
- Rescheduling should be done through direct, low-friction actions.
- The interface should avoid dense medical-system complexity.
- Mobile layout should prioritize today, pending confirmations, and next available slots.

## Data Relationships

Suggested product model:

```txt
Clinic
  has many Users
  has many Services
  has many Professionals
  has many Leads
  has many Appointments

Professional
  belongs to Clinic
  can provide many Services
  has many Appointments

Lead
  belongs to Clinic
  may be interested in one Service
  may have many FollowUps through next_action_at and notes
  may be converted into Appointment

Appointment
  belongs to Clinic
  belongs to Professional
  may belong to Lead
  belongs to Service
  has status
  has scheduled date/time
  may have confirmation/follow-up state

Inbox AI Analysis
  belongs to Clinic
  may create/update Lead after human review
  may suggest Appointment after human review
```

## Database Planning

Status: Planned. Do not create migrations until implementation is approved.

Potential tables:

- `clinic_professionals`.
- `professional_services`.
- `appointments`.
- Optional later: `appointment_events` for audit/history.

Detailed database planning is documented in `docs/SMART_SCHEDULE_DATABASE_PLAN.md`.

Recommended indexes:

- `clinic_professionals(clinic_id, is_active)`.
- `professional_services(clinic_id, professional_id)`.
- `appointments(clinic_id, starts_at)`.
- `appointments(clinic_id, professional_id, starts_at)`.
- `appointments(clinic_id, status, starts_at)`.
- `appointments(clinic_id, lead_id)`.

Tenant isolation:

- Every planned table should include `clinic_id` directly, even when reachable through another relation.
- Backend queries must always scope by authenticated `clinic_id`.
- Foreign keys should prevent cross-clinic relationships where practical, and service-layer validation must enforce clinic ownership.

## API Planning

Status: Planned. Draft contracts live in `docs/API_CONTRACTS.md`.

Planned API areas:

- Professionals CRUD.
- Professional service assignment.
- Appointment CRUD.
- Appointment rescheduling and status updates.
- Lead-to-appointment conversion.
- Schedule availability lookup.
- Dashboard schedule summary.

All endpoints require JWT unless explicitly public. All write operations require clinic scoping and role-aware authorization.

## Frontend Planning

Status: Planned.

Planned feature areas:

- `features/schedule/` for daily/weekly agenda.
- `features/professionals/` for professional management.
- Appointment create/edit modal or screen.
- Lead conversion action from lead detail and Inbox AI.
- Dashboard cards for schedule priorities.

The first usable screen should be the daily agenda, not a marketing page.

## Dashboard Evolution

The dashboard should evolve from KPI summary to daily operational priorities.

Planned sections:

- Today's appointments.
- Appointments pending confirmation.
- Empty slots or available schedule opportunities.
- Hot leads without appointment.
- Overdue follow-ups.
- Appointments by dentist/professional.
- Leads converted to appointments.
- Services with most demand.

Dashboard recommendations must not trigger automatic messages.

## AI Responsibilities

AI may help with commercial/admin scheduling communication.

Allowed:

- Suggest confirmation messages.
- Suggest rescheduling messages.
- Detect scheduling intent in pasted conversations.
- Identify missing scheduling data such as preferred time, service, or contact phone.
- Suggest that a lead may be ready for appointment creation.
- Suggest follow-up text after no-show or cancellation.

Forbidden:

- Diagnose.
- Interpret symptoms.
- Recommend treatment or medication.
- Decide urgency based on symptoms.
- Tell a person they do or do not need care.
- Create an appointment automatically without human confirmation.
- Send WhatsApp messages automatically.

## Suggested Implementation Phases

### Phase A — Documentation and Product Alignment

Status: In Progress.

- Update product vision.
- Define Smart Schedule scope.
- Define professionals.
- Define appointment model.
- Define lead-to-appointment flow.
- Update safety rules and docs.

### Phase B — Database Planning

Status: Planned.

- Plan `clinic_professionals`.
- Plan `professional_services`.
- Plan `appointments`.
- Define appointment statuses.
- Define indexes and tenant separation.
- Do not create migrations until approved.

### Phase C — Backend API Planning

Status: Planned.

- Plan endpoints for professionals.
- Plan endpoints for appointments.
- Plan lead-to-appointment conversion endpoint.
- Plan dashboard schedule summary endpoint.
- Define service-layer validation and authorization rules.

### Phase D — Frontend Planning

Status: Planned.

- Plan Smart Schedule UI.
- Plan professional management UI.
- Plan appointment creation and editing.
- Plan lead conversion to appointment.
- Plan dashboard changes.

### Phase E — Implementation Later

Status: Future.

- Implement only after documentation review and approval.
- Start with database migrations and backend contracts.
- Add focused frontend schedule workflows.
- Expand smoke/e2e coverage after endpoints exist.

## Future Improvements

- Appointment event history.
- Drag-and-drop rescheduling if usability validates it.
- Recurring working hours by professional.
- Holiday and blocked-slot management.
- Multi-location schedules.
- Optional WhatsApp Business Cloud API integration after MVP validation and explicit approval.
- Calendar export or external calendar sync.
- Deeper analytics for lead-to-appointment conversion.

## Acceptance Criteria For Future Implementation

A future implementation should be considered acceptable only when:

- Staff can view daily and weekly appointments.
- Staff can create and edit appointments linked to professionals, services, and optional leads.
- Staff can convert a reviewed lead into an appointment.
- Staff can see pending confirmations and available slots.
- Appointment data is clinic-scoped and authenticated.
- Admin notes remain commercial/admin only.
- AI suggestions require human review and do not send messages automatically.
- No clinical histories, diagnoses, prescriptions, or medical advice are added.
