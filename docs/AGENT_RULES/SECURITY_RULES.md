# Security Rules

## No Secrets Committed

Never commit:
- `.env`, `.env.local`, `.env.production`, `.env.development`
- `*.pem`, `*.key`, `*.p12`, `*.pfx`
- `credentials.json`, `service-account.json`
- `openai_api_key.txt`, `gemini_api_key.txt`, `deepseek_api_key.txt`

Allowed: `.env.example` with placeholder values only.

Always inspect `.gitignore` before adding files. If it does not protect secrets, update it immediately.

## JWT Safety

- Never log raw JWT tokens.
- Never expose JWT secrets.
- Validate JWT on every protected backend endpoint.

## Password Handling

- Hash passwords securely.
- Never log passwords.
- Never expose password hashes in API responses.

## Logging Restrictions

Allowed:
- Request path, status code, execution time, error code, module name, provider name (without key), user ID if needed, clinic ID if needed.

Forbidden:
- Passwords, JWT tokens, API keys, full patient conversations, sensitive notes, raw AI prompts with private data (unless debug mode locally only).

## Sensitive Data Restrictions

The MVP must not store clinical history or sensitive health records.

Allowed lead data:
- Full name, phone, service of interest, commercial status, sales conversation notes, next follow-up date, basic source/channel.

Allowed appointment data:
- Contact name and phone, linked lead, linked service, linked professional, appointment date/time, appointment status, confirmation status, commercial/admin notes, basic source/channel.

Not allowed:
- Diagnoses, medical evolution notes, clinical images, prescriptions, lab results, medical records, detailed symptoms as medical history.
- Clinical notes, treatment plans, or medical decisions inside appointment records.

## Multi-Tenant Data Isolation

Every business entity belongs to a clinic. Backend queries must enforce data isolation. A user from Clinic A must never access data from Clinic B.

This must be enforced in backend queries, not only in frontend filters.

Smart Schedule entities must validate clinic ownership for every linked object: professional, service, lead, appointment, and user.

## Backend Authorization Enforcement

Frontend guard is not enough. Every protected endpoint must validate:
- JWT exists.
- JWT is valid.
- User belongs to the clinic.
- User has required role or permission.
