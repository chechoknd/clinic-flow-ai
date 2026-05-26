ALTER TABLE clinics
    ADD COLUMN country_code CHAR(2) NOT NULL DEFAULT 'CO',
    ADD COLUMN currency_code CHAR(3) NOT NULL DEFAULT 'COP';

ALTER TABLE clinics
    ADD CONSTRAINT clinics_country_code_check CHECK (country_code IN ('CO', 'PE', 'AR', 'CL')),
    ADD CONSTRAINT clinics_currency_code_check CHECK (currency_code IN ('COP', 'PEN', 'ARS', 'CLP')),
    ADD CONSTRAINT clinics_country_currency_check CHECK (
        (country_code = 'CO' AND currency_code = 'COP') OR
        (country_code = 'PE' AND currency_code = 'PEN') OR
        (country_code = 'AR' AND currency_code = 'ARS') OR
        (country_code = 'CL' AND currency_code = 'CLP')
    );
