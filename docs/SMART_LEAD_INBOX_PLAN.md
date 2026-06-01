# Smart Lead Inbox Plan

Status: In Progress.

Smart Lead Inbox, also called Inbox AI in some product notes, is the planned evolution of ClinicFlow AI from a manual CRM plus AI assistant into an intelligent commercial attention center for dental clinics.

After dentist validation, Inbox AI should align with the broader Smart Schedule direction. It should help turn conversations into reviewed leads, follow-ups, and appointment suggestions that feed the visual agenda.

This plan tracks the broader direction. The first slice, manual conversation analysis with human-reviewed lead creation, is implemented; the rest remains planned.

## Purpose

Smart Lead Inbox should help receptionists and clinic administrators handle commercial conversations faster and with better context, especially when a person is ready to schedule or needs confirmation/follow-up before becoming an appointment.

The product direction is:

```txt
AI suggests -> human reviews -> human replies
```

ClinicFlow AI must remain a commercial assistant. It must not become a medical diagnosis system, clinical decision support tool, autonomous WhatsApp bot, or medical-record platform.

## Current Problem

The current MVP works, but the daily workflow is still manual:

- The assistant receives a conversation in WhatsApp, email, Instagram, or another channel.
- The assistant copies part of the conversation.
- The assistant opens ClinicFlow AI.
- The assistant pastes the message into the AI Assistant.
- The assistant copies the generated response back to the original channel.
- The assistant manually creates or updates a lead and follow-up.

This proves the value of assisted responses, but it still forces the user to understand where to paste information, which lead to update, and which next action to schedule.

## Product Direction

Smart Lead Inbox should become a commercial inbox assisted by AI and connected to the planned Smart Schedule.

Future versions should allow clinic staff to:

- Paste a complete conversation.
- Ask AI to analyze the full commercial context.
- Detect the likely service of interest.
- Detect the lead's commercial intent.
- Detect objections such as price, fear, time, trust, or uncertainty.
- Generate a suggested WhatsApp-ready reply.
- Suggest the next commercial action.
- Suggest a follow-up date or task.
- Detect whether the person wants to schedule.
- Identify missing scheduling data such as preferred day, time, service, or contact phone.
- Suggest appointment creation for human review.
- Create a lead from the analyzed conversation.
- Update an existing lead from the analyzed conversation.
- Convert a reviewed lead into an appointment through the planned appointment flow.
- Keep a history of human-reviewed AI analyses.

The inbox should optimize for receptionist workflows: fewer steps, clear actions, fast copy-to-clipboard behavior, and safe commercial language.

## MVP Scope

The immediate planned MVP extension should stay manual or semi-manual.

Allowed for the Smart Lead Inbox MVP:

- Manual conversation paste.
- Manual source selection such as WhatsApp, Gmail, Instagram, web form, or other.
- AI conversation analysis for commercial attributes.
- AI extraction of lead name, phone, service interest, commercial status, objections, urgency, and suggested next action when present in the text.
- AI-generated response suggestions.
- Human review before saving or sending anything.
- Create lead from analysis after human confirmation.
- Update lead from analysis after human confirmation.
- Create suggested follow-up after human confirmation.
- Suggest appointment creation after human confirmation.
- Surface missing data needed before scheduling.
- Dashboard cards that prioritize what the user should handle today.

The MVP must keep message delivery outside the system. The user copies the approved response and sends it manually through the original channel.

## Out of Scope

The following remain out of scope for this phase:

- Autonomous WhatsApp bots.
- Automatic message sending.
- Native bidirectional WhatsApp Business Cloud API integration as a required MVP dependency.
- Diagnosis, clinical decision support, prescriptions, treatment recommendations, or medical advice.
- Clinical histories, medical records, lab results, clinical images, or detailed symptom histories.
- Telemedicine, video calls, clinical-grade scheduling, treatment planning, or autonomous appointment booking.
- Autonomous appointment booking.
- Clinical-grade scheduling, treatment planning, or clinical notes.
- RAG/vector search as required MVP infrastructure unless explicitly approved later.
- n8n as a central product dependency.
- Complex automation rules that bypass human approval.

## User Experience

### Dashboard Intelligent Actions

The current dashboard shows useful metrics. The planned dashboard should become more actionable.

Future dashboard actions should answer:

- What should the assistant handle today?
- Which appointments happen today?
- Which appointments are pending confirmation?
- Which schedule slots are still available?
- Which new leads have not received a response?
- Which follow-ups are overdue?
- Which leads show high intent?
- Which hot leads do not yet have an appointment?
- Which leads have price, fear, or timing objections?
- Which services are receiving the most attention?
- What action does AI recommend next?

Planned dashboard sections:

- `Que atender hoy`: priority queue for the assistant.
- `Citas de hoy`: appointments grouped by time, professional, and status.
- `Pendientes de confirmacion`: appointments requiring manual confirmation.
- `Espacios disponibles`: available slots that can be offered to hot leads.
- `Leads nuevos sin respuesta`: new opportunities requiring first contact.
- `Leads calientes sin cita`: high-intent leads not linked to an appointment.
- `Seguimientos vencidos`: leads whose next action date has passed.
- `Alta intencion`: leads likely close to booking a valuation.
- `Objeciones detectadas`: leads where AI found price, fear, time, or trust objections.
- `Servicios mas consultados`: demand signal by service.
- `Recomendaciones AI`: safe next actions, not automatic actions.

Planned quick actions:

- Analyze conversation.
- Open Inbox AI.
- Create lead with AI.
- View pending follow-ups.
- Generate suggested reply.
- Create appointment from reviewed lead.
- Open Smart Schedule.

### Inbox AI Flow

The planned manual Inbox AI flow is:

1. User opens Smart Lead Inbox.
2. User pastes a conversation from WhatsApp, Gmail, Instagram, a web form, or another source.
3. User optionally selects source and existing lead.
4. AI analyzes the conversation.
5. AI returns structured commercial analysis.
6. User reviews extracted data.
7. User edits anything that is incomplete or wrong.
8. User chooses one or more actions:
   - Copy suggested reply.
   - Create new lead.
   - Update existing lead.
   - Add note to lead.
   - Schedule follow-up.
   - Create appointment suggestion for review when Smart Schedule is implemented.
9. User manually sends the approved message outside ClinicFlow AI.

The UI should make it clear that the AI output is a draft and must be reviewed.

## AI Responsibilities

AI may assist with commercial interpretation only.

Allowed AI responsibilities:

- Summarize the commercial conversation.
- Extract contact information if explicitly present.
- Identify likely service interest.
- Classify commercial intent such as low, medium, or high.
- Classify commercial objections such as price, fear, time, trust, or availability.
- Suggest lead status.
- Suggest next action.
- Suggest follow-up timing.
- Detect scheduling intent.
- Suggest missing scheduling data before appointment creation.
- Suggest appointment creation with service, duration, and preferred time when explicitly supported by the conversation.
- Generate a safe response draft for WhatsApp.
- Flag when a question should be answered by inviting the patient to a professional valuation.

Forbidden AI responsibilities:

- Diagnose.
- Interpret symptoms.
- Recommend medication or dosage.
- Prescribe treatments.
- Decide urgency based on symptoms.
- Decide medical urgency.
- Guarantee outcomes.
- Tell a patient they do not need a clinic visit.
- Store or produce clinical-history content.
- Create, confirm, cancel, or reschedule appointments automatically.

## Human Review Requirements

Human approval is mandatory.

ClinicFlow AI must not send messages automatically in this phase. The AI response is a draft. A clinic staff member must review, edit if needed, copy, and send the message manually.

This rule exists for:

- Safety.
- Quality control.
- Brand tone control.
- Scope discipline.
- Avoiding premature WhatsApp automation complexity.

## Backend Implications

Status: Planned.

Potential backend areas:

- New AI use case for conversation analysis.
- Structured DTOs for extracted commercial fields.
- Optional inbound message or conversation storage after human confirmation.
- Lead creation/update orchestration from reviewed analysis.
- Appointment suggestion and lead-to-appointment orchestration from reviewed analysis.
- Follow-up suggestion generation.
- Dashboard action summary endpoint or extension of dashboard summary.
- AI generation metadata for audit, cost, and quality tracking.

Backend constraints:

- All endpoints must require JWT unless public by design.
- All persisted data must be scoped by clinic.
- Do not persist raw full conversations by default unless a product decision approves retention rules.
- Do not create appointments automatically from AI output.
- Do not expose raw provider prompts or unsafe AI internals.
- Keep prompt ownership and safety validation in backend.

## Frontend Implications

Status: Planned.

Potential frontend areas:

- New `inbox-ai` or `smart-inbox` feature route.
- Conversation paste form.
- Source selector.
- Optional existing lead selector.
- Analysis result panel.
- Suggested response panel with copy-to-clipboard.
- Editable extracted lead fields.
- Create/update lead actions.
- Suggested follow-up confirmation.
- Suggested appointment review panel connected to the planned Smart Schedule.
- Missing scheduling data prompts.
- Dashboard action cards and shortcuts.

UX constraints:

- Keep the first screen operational, not a marketing page.
- Avoid dense dashboards.
- Prioritize clear next actions.
- Make AI confidence and human review requirements visible.
- Keep all copy Spanish-first for the MVP.

## API Implications

Status: Proposed contracts only. Not implemented.

Potential endpoints:

- `POST /api/ai/analyze-conversation`
- `POST /api/leads/from-conversation`
- `POST /api/leads/:id/convert-to-appointment`
- `POST /api/inbound/messages`
- `POST /api/followups/suggest`

These endpoints require product and technical validation before implementation. They must not imply automatic message sending.

## Future Integrations

Future integrations may be considered after validating the manual Inbox AI workflow.

Possible future sources:

- n8n.
- Gmail.
- Web forms.
- Meta Lead Ads.
- WhatsApp Business Cloud API.
- Browser extension for capturing selected text from WhatsApp Web or Gmail.

These are not immediate MVP dependencies. The product should first prove value with manual or semi-manual capture.

## Suggested Implementation Phases

### Phase A — Documentation and Product Alignment

Status: Implemented.

- Document Smart Lead Inbox.
- Document Dashboard Inteligente.
- Register decision in Decisions Log.
- Update Development Status.
- Align API contracts as planned, not implemented.

### Phase B — Manual Conversation Analysis

Status: Partially Implemented.

- Implemented screen to paste a conversation.
- Implemented AI conversation analysis endpoint.
- Implemented AI extraction of commercial data.
- Implemented AI response suggestion display and copy.
- Implemented lead creation from reviewed analysis through the existing leads API.
- Pending: update existing lead from analysis, persisted analysis history, and richer follow-up confirmation.

### Phase C — Dashboard Action Layer

Status: Planned.

- Add `Que atender hoy` dashboard section.
- Add today's appointments and pending confirmations.
- Add hot leads without appointment.
- Add available schedule opportunities after Smart Schedule exists.
- Surface priority leads.
- Surface overdue follow-ups.
- Surface new leads without response.
- Surface AI recommendations.

### Phase D — Inbox AI Schedule Integration

Status: Planned.

- Detect scheduling intent.
- Show missing appointment data.
- Suggest appointment creation after human review.
- Preserve traceability between analysis, lead, appointment, and follow-up.

### Phase E — Inbox AI MVP

Status: Planned.

- Build a simple inbox-style view.
- Track attention states.
- Show suggested replies.
- Show suggested follow-ups.
- Keep history of analysis after human confirmation.

### Phase F — Optional Integrations

Status: Future.

- Evaluate n8n.
- Evaluate Gmail import.
- Evaluate web form capture.
- Evaluate Meta Lead Ads capture.
- Evaluate WhatsApp Business Cloud API.
- Evaluate browser extension.

## Acceptance Criteria

A future implementation should be considered acceptable only when:

- Staff can paste a conversation and receive structured commercial analysis.
- AI output avoids clinical diagnosis, prescriptions, and medical advice.
- Staff can review and edit extracted fields before saving.
- Staff can copy the suggested response manually.
- The system does not send messages automatically.
- Lead creation/update remains clinic-scoped and authenticated.
- Appointment suggestions require human confirmation and a separate schedule action.
- Follow-up suggestions require human confirmation.
- Dashboard actions help the user decide what to handle next.
- Inbox AI can support the Smart Schedule without becoming an autonomous booking tool.
- API contracts and docs are updated before or with implementation.

## Risks and Guardrails

Risks:

- Users may paste sensitive or clinical data into conversation analysis.
- AI may infer too much from incomplete conversation context.
- Dashboard recommendations may feel like automation if wording is not clear.
- Appointment suggestions may be mistaken for automatic booking if the UI is not explicit.
- Premature integrations may distract from validating the manual workflow.
- Full conversation retention may create privacy and compliance concerns.

Guardrails:

- Keep analysis commercial.
- Warn users not to paste clinical histories or sensitive medical records.
- Require human confirmation before creating appointments.
- Keep scheduling suggestions operational and avoid clinical urgency decisions.
- Use backend safety prompts and output validation.
- Prefer storing summarized commercial notes over full raw conversations unless retention is explicitly approved.
- Require human confirmation for every save, follow-up, or message send action.
- Keep WhatsApp API and autonomous bots out of the immediate MVP.
