CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE clinics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    clinic_type VARCHAR(50) NOT NULL DEFAULT 'odontologia',
    city VARCHAR(100) NOT NULL,
    phone VARCHAR(30),
    whatsapp VARCHAR(30) NOT NULL,
    address TEXT,
    opening_hours JSONB NOT NULL DEFAULT '{}'::jsonb,
    general_faq JSONB NOT NULL DEFAULT '[]'::jsonb,
    communication_tone VARCHAR(50) NOT NULL DEFAULT 'amable',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT clinics_communication_tone_check CHECK (
        communication_tone IN ('amable', 'profesional', 'cercano', 'juvenil', 'elegante')
    )
);

CREATE TRIGGER clinics_set_updated_at
BEFORE UPDATE ON clinics
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE RESTRICT,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_role_check CHECK (role IN ('superadmin', 'clinic_admin', 'assistant')),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE INDEX users_clinic_id_idx ON users(clinic_id);
CREATE INDEX users_role_idx ON users(role);

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE clinic_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    duration_minutes INTEGER,
    price_from NUMERIC(12,2),
    benefits JSONB NOT NULL DEFAULT '[]'::jsonb,
    faq JSONB NOT NULL DEFAULT '[]'::jsonb,
    common_objections JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT clinic_services_duration_minutes_check CHECK (
        duration_minutes IS NULL OR duration_minutes > 0
    ),
    CONSTRAINT clinic_services_price_from_check CHECK (
        price_from IS NULL OR price_from >= 0
    ),
    CONSTRAINT clinic_services_name_per_clinic_unique UNIQUE (clinic_id, name)
);

CREATE INDEX clinic_services_clinic_id_idx ON clinic_services(clinic_id);
CREATE INDEX clinic_services_active_idx ON clinic_services(is_active);

CREATE TRIGGER clinic_services_set_updated_at
BEFORE UPDATE ON clinic_services
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
