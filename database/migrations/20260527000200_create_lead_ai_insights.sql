CREATE TABLE lead_ai_insights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    lead_id UUID NOT NULL REFERENCES leads(id) ON DELETE CASCADE,
    ai_generation_id UUID REFERENCES ai_generations(id) ON DELETE SET NULL,
    intent VARCHAR(20) NOT NULL,
    detected_objections JSONB NOT NULL DEFAULT '[]'::jsonb,
    commercial_summary TEXT,
    suggested_next_action TEXT,
    source VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT lead_ai_insights_intent_check CHECK (intent IN ('low', 'medium', 'high'))
);

CREATE INDEX lead_ai_insights_clinic_created_idx ON lead_ai_insights(clinic_id, created_at DESC);
CREATE INDEX lead_ai_insights_lead_created_idx ON lead_ai_insights(lead_id, created_at DESC);
CREATE INDEX lead_ai_insights_intent_idx ON lead_ai_insights(intent);
