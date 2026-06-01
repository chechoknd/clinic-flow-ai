# Frontend Rules

## Architecture

Frontend lives in `apps/frontend-angular/`.

Expected Angular structure:

```
apps/frontend-angular/src/app/
├── core/
│   ├── auth/
│   ├── guards/
│   ├── interceptors/
│   ├── layout/
│   └── services/
├── shared/
│   ├── components/
│   ├── pipes/
│   ├── directives/
│   └── utils/
├── features/
│   ├── dashboard/
│   ├── clinics/
│   ├── services/
│   ├── professionals/
│   ├── schedule/
│   ├── leads/
│   ├── ai-assistant/
│   ├── content/
│   └── followups/
├── app.routes.ts
└── app.config.ts
```

## Coding Rules

- Use standalone components.
- Use lazy-loaded feature routes.
- Use Reactive Forms for forms.
- Use Angular Signals for simple local state.
- Use services for HTTP calls.
- Use interceptors for auth headers.
- Use route guards for private routes.
- Keep templates readable and components focused.
- Use Tailwind CSS utility classes consistently.
- Build responsive layouts from the beginning.
- Avoid NgRx unless explicitly requested.

Must not:
- Put all features in one component.
- Hardcode backend URLs or API keys in components.
- Duplicate UI logic across components.
- Mix unrelated feature logic in shared components.

## UI/UX Rules

ClinicFlow AI must feel simple, clean, commercial, and fast.

Primary user: receptionist or clinic administrator (non-technical).

Priorities:
1. Fewer clicks.
2. Clear action buttons.
3. Simple forms.
4. Fast copy-to-clipboard workflows.
5. Clear lead status and next action.
6. No dense dashboards.
7. No complex medical terminology unless needed.
8. No unnecessary visual noise.
9. Mobile-friendly layouts.

### Smart Schedule UX

Status: Planned.

The Smart Schedule should be the first operational surface once implemented. It should feel like a clear commercial agenda, not a hospital system.

Priorities:

1. Daily view first.
2. Weekly view for planning.
3. Appointment cards with time, contact, service, professional, and status.
4. Clear pending confirmation indicators.
5. Professional color identifiers.
6. Available slots that are easy to scan.
7. Fast create, confirm, cancel, complete, and reschedule actions.
8. Lead-to-appointment conversion from reviewed data.
9. No clinical notes or medical-history UI.

### AI Assistant UX

- Input: patient question or objection.
- Context: selected service.
- Output: ready-to-copy response.
- Variants: short, persuasive, technical, closing question.
- Button: copy to clipboard.
- Optional: save note to lead.
- Must feel like a useful tool, not a chatbot.

### Lead Pipeline UX

Statuses: `Nuevo → Contactado → Interesado → Agendado → No Respondió → Perdido → Convertido`

Make it easy to: see status, change status, add notes, set next action date, generate follow-up message, copy WhatsApp response.

## Accessibility

- Buttons must have readable text.
- Inputs must have labels.
- Forms must show validation messages.
- Color cannot be the only indicator.
- Text contrast must be readable.
- Keyboard navigation should not be broken.
- Loading states must be visible.
- Empty states must explain what to do next.

## Internationalization

First MVP can be Spanish-first. Avoid hardcoding text that makes future translation impossible. May later support Spanish, English, Portuguese.

## Testes

Frontend should include tests for: auth service, guards, form validation, lead status transitions, copy-to-clipboard behavior, API service methods.

Run: `npm install`, `npm run build`, `npm test`, `npm run lint`.
