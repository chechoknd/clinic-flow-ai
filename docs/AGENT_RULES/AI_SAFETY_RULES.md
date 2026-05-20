# AI Safety Rules

## AI Assistant Role

The AI assistant must behave as a commercial communication assistant, not a doctor.

## Forbidden AI Behavior

The AI must never:
- Diagnose.
- Interpret symptoms.
- Recommend medication.
- Suggest drug dosage.
- Guarantee medical results.
- Claim a treatment is risk-free.
- Replace professional evaluation.
- Interpret medical images.
- Decide urgency based on symptoms.
- Tell a patient they do not need to visit the clinic.

## Required Safe AI Behavior

The AI must:
- Encourage professional evaluation.
- Use empathetic language.
- Be transparent with "prices from" when available.
- Avoid absolute promises.
- Suggest booking a consultation.
- Keep WhatsApp messages short.
- Use emojis moderately.
- Adapt to clinic communication tone.
- Respect clinic service context.
- Respect objection-handling strategy.
- Avoid fear-based manipulation.

## Prompt Ownership

Prompt templates must live in backend-controlled files or backend services. The frontend must not own critical prompt safety rules.

The backend is the source of truth for:
- System prompts.
- Safety constraints.
- Provider selection.
- AI response validation.
- Token limits.
- Prompt context injection.

## AI Output Validation

When possible, the backend should validate AI output before returning it. Plan validators to detect and block:
- Medication names, dosage instructions, diagnostic claims, guaranteed outcomes, emergency misdirection, overly aggressive sales language.

If risky output is detected, return a safe fallback message.

## Cost Control

Design for:
- Token limits.
- Provider timeout.
- Provider error fallback.
- Optional daily/monthly usage counters.
- Storing AI generation metadata.
- Avoiding unnecessarily long prompts.
- Avoiding repeated calls when not needed.

Suggested AI generation metadata: `id`, `clinic_id`, `user_id`, `feature`, `provider`, `model`, `input_tokens`, `output_tokens`, `created_at`.

## Provider Abstraction

Support OpenAI, Gemini, DeepSeek, or compatible providers through environment variables. Provider-specific code must be isolated behind interfaces.

## No RAG Required for MVP

Use direct structured context injection. Do not implement RAG unless explicitly requested.

## Context Injection Strategy

Prompt context may include: clinic name, clinic type, city, communication tone, service details, service benefits, service FAQ, common objections, lead status, optional lead notes.

Keep context concise to control cost and latency.
