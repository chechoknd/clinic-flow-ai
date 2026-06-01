ALTER TABLE users
ADD CONSTRAINT users_clinic_id_id_unique UNIQUE (clinic_id, id);

ALTER TABLE clinic_services
ADD CONSTRAINT clinic_services_clinic_id_id_unique UNIQUE (clinic_id, id);

ALTER TABLE leads
ADD CONSTRAINT leads_clinic_id_id_unique UNIQUE (clinic_id, id);

CREATE TABLE clinic_professionals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    full_name VARCHAR(150) NOT NULL,
    role_or_specialty VARCHAR(120),
    calendar_color VARCHAR(20),
    working_hours JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT clinic_professionals_clinic_id_id_unique UNIQUE (clinic_id, id),
    CONSTRAINT clinic_professionals_name_per_clinic_unique UNIQUE (clinic_id, full_name),
    CONSTRAINT clinic_professionals_calendar_color_check CHECK (
        calendar_color IS NULL OR calendar_color ~ '^#[0-9A-Fa-f]{6}$'
    )
);

CREATE INDEX clinic_professionals_clinic_id_idx ON clinic_professionals(clinic_id);
CREATE INDEX clinic_professionals_active_idx ON clinic_professionals(clinic_id, is_active);

CREATE TRIGGER clinic_professionals_set_updated_at
BEFORE UPDATE ON clinic_professionals
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE professional_services (
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL,
    service_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (professional_id, service_id),
    CONSTRAINT professional_services_professional_fk FOREIGN KEY (clinic_id, professional_id)
        REFERENCES clinic_professionals(clinic_id, id) ON DELETE CASCADE,
    CONSTRAINT professional_services_service_fk FOREIGN KEY (clinic_id, service_id)
        REFERENCES clinic_services(clinic_id, id) ON DELETE CASCADE
);

CREATE INDEX professional_services_clinic_id_idx ON professional_services(clinic_id);
CREATE INDEX professional_services_service_id_idx ON professional_services(clinic_id, service_id);
CREATE INDEX professional_services_professional_id_idx ON professional_services(clinic_id, professional_id);

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    professional_id UUID NOT NULL,
    lead_id UUID REFERENCES leads(id) ON DELETE SET NULL,
    service_id UUID NOT NULL,
    contact_name VARCHAR(150) NOT NULL,
    contact_phone VARCHAR(30),
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_confirmation',
    confirmation_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    source VARCHAR(50) NOT NULL DEFAULT 'manual',
    admin_notes TEXT,
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT appointments_professional_fk FOREIGN KEY (clinic_id, professional_id)
        REFERENCES clinic_professionals(clinic_id, id) ON DELETE RESTRICT,
    CONSTRAINT appointments_service_fk FOREIGN KEY (clinic_id, service_id)
        REFERENCES clinic_services(clinic_id, id) ON DELETE RESTRICT,
    CONSTRAINT appointments_status_check CHECK (
        status IN (
            'scheduled',
            'confirmed',
            'pending_confirmation',
            'rescheduled',
            'no_show',
            'cancelled',
            'completed',
            'converted_from_lead'
        )
    ),
    CONSTRAINT appointments_confirmation_status_check CHECK (
        confirmation_status IN ('pending', 'confirmed', 'not_required', 'failed')
    ),
    CONSTRAINT appointments_time_check CHECK (ends_at > starts_at)
);

CREATE INDEX appointments_clinic_date_idx ON appointments(clinic_id, starts_at);
CREATE INDEX appointments_professional_date_idx ON appointments(clinic_id, professional_id, starts_at);
CREATE INDEX appointments_status_date_idx ON appointments(clinic_id, status, starts_at);
CREATE INDEX appointments_confirmation_date_idx ON appointments(clinic_id, confirmation_status, starts_at);
CREATE INDEX appointments_lead_id_idx ON appointments(clinic_id, lead_id);
CREATE INDEX appointments_service_date_idx ON appointments(clinic_id, service_id, starts_at);
CREATE INDEX appointments_active_professional_time_idx
ON appointments(clinic_id, professional_id, starts_at, ends_at)
WHERE status IN (
    'scheduled',
    'confirmed',
    'pending_confirmation',
    'rescheduled',
    'converted_from_lead'
);

CREATE TRIGGER appointments_set_updated_at
BEFORE UPDATE ON appointments
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
