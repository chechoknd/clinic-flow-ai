# Currency Strategy

Status: Implemented for the MVP baseline. Currency conversion, exchange rates, payments, invoicing, and exact decimal refactors remain out of scope.

## Problem Summary

ClinicFlow AI currently stores service commercial prices through `clinic_services.price_from` and exposes them as `price_from` in service API responses. The frontend previously rendered service prices with a hardcoded `COP` label. The current MVP uses clinic currency metadata and a shared formatter in the Services screen.

The product is intended for clinics in multiple Latin American countries. The system must stop assuming a single currency and must support clinic-specific currency configuration for display, storage, validation, AI context, and future commercial amounts.

Initial target countries and currencies:

| Country | Country code | Currency code | Symbol | Decimals | Example locale |
| --- | --- | --- | --- | --- | --- |
| Colombia | CO | COP | $ | 0 | es-CO |
| Peru | PE | PEN | S/ | 2 | es-PE |
| Argentina | AR | ARS | $ | 2 | es-AR |
| Chile | CL | CLP | $ | 0 | es-CL |

No automatic conversion, exchange rates, multi-currency quoting, payment gateway, or invoicing is required for the MVP.

## Recommended Decision

Use clinic-level currency configuration as the source of truth.

For the MVP, store the selected currency directly on `clinics` and validate it against a small static currency catalog owned by the backend. The initial catalog should include `COP`, `PEN`, `ARS`, and `CLP`, with display metadata needed by the frontend.

Implemented clinic fields:

```txt
country_code     CHAR(2)      NOT NULL DEFAULT 'CO'
currency_code    CHAR(3)      NOT NULL DEFAULT 'COP'
```

Implemented currency metadata contract:

```json
{
  "code": "COP",
  "symbol": "$",
  "locale": "es-CO",
  "decimal_digits": 0,
  "thousand_separator": ".",
  "decimal_separator": ",",
  "symbol_position": "before",
  "default_country_code": "CO"
}
```

For storage, keep monetary values as PostgreSQL `NUMERIC(12,2)` in the MVP. Avoid Go `float64` in future DTO/domain work and move toward exact decimal handling or string-based request parsing before adding more monetary fields. Do not convert existing `price_from` to integer minor units in the MVP because COP and CLP commonly use zero decimal display while PEN and ARS use two decimals, and current data is already modeled as decimal-like SQL.

## Questions Answered

### 1. Where should currency configuration live?

Currency should live on the clinic profile because each tenant operates commercially in one default country/currency. This matches the existing tenant model and avoids adding unnecessary country or currency administration screens.

Implemented approach:

- Add `country_code` and `currency_code` to `clinics`.
- Keep a static backend currency catalog for validation and response metadata.
- Do not create `countries` or `currencies` tables for the MVP unless runtime-administered catalogs become a product requirement.

Why not only a country table:

- The app needs currency behavior, not country administration.
- A country can be a useful default, but the actual formatting decision should come from `currency_code` plus locale metadata.

Why not only frontend constants:

- Backend must validate accepted currency codes and return consistent API contracts.
- AI prompt context and future exports must not depend on frontend-only formatting rules.

### 2. What database changes are needed?

Implemented migration:

```sql
ALTER TABLE clinics
  ADD COLUMN country_code CHAR(2) NOT NULL DEFAULT 'CO',
  ADD COLUMN currency_code CHAR(3) NOT NULL DEFAULT 'COP';

ALTER TABLE clinics
  ADD CONSTRAINT clinics_country_code_check CHECK (country_code IN ('CO', 'PE', 'AR', 'CL')),
  ADD CONSTRAINT clinics_currency_code_check CHECK (currency_code IN ('COP', 'PEN', 'ARS', 'CLP'));
```

Existing `clinic_services.price_from NUMERIC(12,2)` can remain for MVP compatibility. Future migrations should use the same numeric approach or introduce a clearer money type convention before adding more amount columns.

If future data already exists, the migration should default existing clinics to `CO`/`COP`, then allow clinic admins to update the profile. For production, the default should be explicitly reviewed per tenant before rollout.

### 3. What backend Go changes are needed?

Implemented backend changes:

- Add `CountryCode` and `CurrencyCode` to `internal/clinics` model, DTOs, repository queries, and update use case.
- Add backend validation for supported country/currency combinations.
- Add a small static currency catalog, likely under `internal/shared/currency` or `internal/clinics/currency.go`.
- Expose currency metadata in `GET /api/clinics/current`.
- Include currency metadata, or at minimum `currency_code`, in service responses where `price_from` appears.
- Update AI context so service price context includes currency code/symbol/locale, not only a raw number.
- `float64` monetary DTOs still exist for `price_from`; replacing them with exact decimal parsing remains future technical debt before adding more money fields.

Current risk to address:

- `internal/services` currently maps `NUMERIC(12,2)` to `*float64`. That is acceptable only as a short-lived MVP simplification. A currency implementation should avoid expanding this pattern.

### 4. What frontend Angular changes are needed?

Implemented frontend changes:

- Extend `ClinicProfile` with `country_code`, `currency_code`, and `currency` metadata.
- Add country/currency controls to the clinic profile form.
- Load clinic currency configuration once for authenticated flows, or retrieve it in screens that render money.
- Add a reusable currency formatter helper or pipe in `shared/`.
- Replace hardcoded `Desde COP {{ service.price_from || 0 }}` with the shared formatter.
- Use currency decimal digits to configure service price input hints and validation.
- Frontend tests cover updated clinic profile, service rendering dependencies, and API models.

### 5. How should currency be exposed in API contracts?

Recommended clinic response shape:

```json
{
  "id": "clinic-id",
  "name": "Sonrisa Viva Demo",
  "country_code": "CO",
  "currency_code": "COP",
  "currency": {
    "code": "COP",
    "symbol": "$",
    "locale": "es-CO",
    "decimal_digits": 0,
    "thousand_separator": ".",
    "decimal_separator": ",",
    "symbol_position": "before"
  }
}
```

Implemented service response addition:

```json
{
  "price_from": "250000.00",
  "currency_code": "COP"
}
```

For MVP compatibility, `price_from` may remain numeric in the first migration if changing it to a string is too disruptive. The preferred future-safe API shape is a string decimal to avoid JavaScript and Go float precision issues.

### 6. How to avoid hardcoded symbols like `$`?

Do not render currency symbols or codes directly in feature templates.

Use one shared formatter:

```txt
formatCurrencyAmount(amount, currencyMetadata)
```

or an Angular pipe:

```html
{{ service.price_from | clinicCurrency: clinic.currency }}
```

Templates should not contain `COP`, `$`, `S/`, decimal separator rules, or country-specific formatting logic.

### 7. How to handle existing prices in the future?

Migration strategy:

1. Add clinic currency fields with `CO`/`COP` defaults for existing local/demo data.
2. Keep existing `price_from` values unchanged.
3. Treat all existing service prices as belonging to the clinic currency.
4. For production tenants, run a one-time review list of clinics and their configured country/currency before enabling multi-country usage.
5. Do not attempt currency conversion.

### 8. Integer minor units or decimal?

Recommendation for MVP: keep PostgreSQL `NUMERIC(12,2)` and move backend/frontend contracts toward string decimal values when practical.

Rationale:

- Existing schema already uses `NUMERIC(12,2)`.
- PostgreSQL numeric is exact for money-like decimal storage.
- Currency conversion is out of scope.
- Some target currencies display zero decimals, but storing two decimal places is acceptable if validation prevents unwanted fractional input for zero-decimal currencies.

Integer minor units are a strong option for payment systems, invoicing, or high-volume accounting. ClinicFlow AI is not a payment or invoicing system in the MVP, so minor units would add migration and parsing complexity without enough benefit right now.

### 9. Recommended Angular formatting strategy

Use `Intl.NumberFormat` through a shared helper/pipe and backend-provided metadata.

Recommended behavior:

- `locale` from currency metadata, e.g. `es-CO`.
- `currency` from `currency_code`.
- `minimumFractionDigits` and `maximumFractionDigits` from `decimal_digits`.
- Prefer `currencyDisplay: 'symbol'` for normal UI.
- For ambiguous `$` currencies, show code in dense commercial contexts when needed, e.g. `COP $250.000` or `$250.000 COP` depending UX decision.

Example formatting expectation:

```txt
COP, es-CO, 0 decimals -> $ 250.000 or COP $250.000
PEN, es-PE, 2 decimals -> S/ 250.00
ARS, es-AR, 2 decimals -> $ 250.000,00 or ARS $250.000,00
CLP, es-CL, 0 decimals -> $250.000 or CLP $250.000
```

The exact visual style should be centralized so it can be adjusted once.

### 10. Required validations

Backend validations:

- `country_code` must be one of `CO`, `PE`, `AR`, `CL`.
- `currency_code` must be one of `COP`, `PEN`, `ARS`, `CLP`.
- Country/currency combination must be valid for the MVP catalog.
- Monetary amounts must be `>= 0`.
- For zero-decimal currencies (`COP`, `CLP`), reject fractional values at API boundary.
- For two-decimal currencies (`PEN`, `ARS`), allow up to two decimals.
- All protected updates remain scoped by authenticated clinic.

Frontend validations:

- Clinic profile must require country/currency selection.
- Service price input must reject negative values.
- Service price input should respect decimal digits for the clinic currency.
- UI should show clear help text, e.g. `Precio en COP` or `Precio en PEN`.

Database validations:

- Check constraints for initial supported country and currency codes.
- Existing `price_from >= 0` check remains valid.
- If strict zero-decimal DB enforcement is desired later, it should be implemented carefully because `price_from` lives in `clinic_services` while decimal rules live on `clinics`.

## Implementation Plan by Phases

### Phase 1 — Documentation and product contract

Status: Completed.

- Approved clinic-level currency configuration.
- Approved initial currency catalog values.
- Updated architecture and API docs.
- Implemented the MVP baseline in database, backend, frontend, and seeds.

### Phase 2 — Database and seed migration

Status: Completed.

- Added `country_code` and `currency_code` columns to `clinics`.
- Defaulted existing demo clinic to `CO`/`COP`.
- Updated demo seeds to include country/currency.
- Validated migration and seeds against local PostgreSQL.

### Phase 3 — Backend contract support

Status: Completed for MVP.

- Updated clinic model, DTO, repository, service validation, and tests.
- Added static currency catalog and validation helper.
- Updated service responses to include `currency_code`.
- Updated AI context to include formatted commercial price with currency code.
- Updated `docs/API_CONTRACTS.md`.

### Phase 4 — Frontend display and forms

Status: Completed for MVP.

- Extended API models.
- Added clinic country/currency controls.
- Added shared currency formatter helper.
- Replaced hardcoded `COP` in Services screen.
- Updated relevant frontend tests.

### Phase 5 — Future amount fields

Only if product needs it:

- Add commercial budget fields to leads or quotes.
- Use the same clinic currency unless a future explicit multi-currency quote requirement is approved.
- Keep exchange rates and conversion out of scope until product validation requires them.

## Future Files Likely to Change

Database:

- `database/migrations/<new>_add_clinic_currency.sql`
- `database/seeds/20260521000100_demo_core_data.sql`
- `database/seeds/20260525000100_demo_ux_data.sql`

Backend:

- `apps/backend-go/internal/clinics/model.go`
- `apps/backend-go/internal/clinics/dto.go`
- `apps/backend-go/internal/clinics/repository.go`
- `apps/backend-go/internal/clinics/service.go`
- `apps/backend-go/internal/services/model.go`
- `apps/backend-go/internal/services/dto.go`
- `apps/backend-go/internal/services/service.go`
- `apps/backend-go/internal/ai/service.go`
- New helper such as `apps/backend-go/internal/shared/currency.go`

Frontend:

- `apps/frontend-angular/src/app/core/services/api.models.ts`
- `apps/frontend-angular/src/app/features/clinics/clinic.page.ts`
- `apps/frontend-angular/src/app/features/clinics/clinic.page.html`
- `apps/frontend-angular/src/app/features/services/services.page.ts`
- `apps/frontend-angular/src/app/features/services/services.page.html`
- New shared formatter/pipe under `apps/frontend-angular/src/app/shared/`

Docs:

- `docs/ARCHITECTURE.md`
- `docs/API_CONTRACTS.md`
- `docs/DEVELOPMENT_STATUS.md`
- `docs/DECISIONS_LOG.md`
- `docs/CURRENCY_STRATEGY.md`

## Risks and Care Points

- Do not introduce exchange rates, payment gateway behavior, invoicing, or accounting logic as part of this feature.
- Avoid expanding `float64` monetary handling. It is already present, but new money work should move toward exact decimal handling.
- Avoid frontend-only currency logic; backend must validate and expose the clinic configuration.
- Avoid hardcoded symbols and currency codes in templates, AI prompts, seeds, and tests.
- Be careful with ambiguous `$` symbols across COP, ARS, and CLP; commercial UI may need to show code plus symbol.
- AI responses that mention prices should receive currency-aware context from backend, not infer currency from city text.
- Existing service prices should be interpreted as the clinic default currency, not converted.
