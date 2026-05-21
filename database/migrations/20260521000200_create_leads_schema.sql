CREATE TABLE leads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    full_name VARCHAR(150) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    service_id UUID REFERENCES clinic_services(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Nuevo',
    source VARCHAR(50) NOT NULL DEFAULT 'whatsapp',
    next_action_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT leads_status_check CHECK (
        status IN ('Nuevo', 'Contactado', 'Interesado', 'Agendado', 'No Respondio', 'Perdido', 'Convertido')
    )
);

CREATE INDEX leads_clinic_id_idx ON leads(clinic_id);
CREATE INDEX leads_status_idx ON leads(status);
CREATE INDEX leads_service_id_idx ON leads(service_id);

CREATE TRIGGER leads_set_updated_at
BEFORE UPDATE ON leads
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE lead_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lead_id UUID NOT NULL REFERENCES leads(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX lead_notes_lead_id_idx ON lead_notes(lead_id);
