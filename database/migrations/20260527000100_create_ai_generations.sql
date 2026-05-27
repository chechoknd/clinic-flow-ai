CREATE TABLE ai_generations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES clinics(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    feature VARCHAR(80) NOT NULL,
    provider VARCHAR(80) NOT NULL,
    model VARCHAR(120) NOT NULL,
    status VARCHAR(40) NOT NULL,
    safety_status VARCHAR(80),
    error_code VARCHAR(80),
    input_tokens INTEGER,
    output_tokens INTEGER,
    input_char_count INTEGER NOT NULL DEFAULT 0,
    output_char_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT ai_generations_status_check CHECK (
        status IN ('success', 'provider_error', 'response_error', 'safety_blocked')
    )
);

CREATE INDEX ai_generations_clinic_created_idx ON ai_generations(clinic_id, created_at DESC);
CREATE INDEX ai_generations_feature_idx ON ai_generations(feature);
CREATE INDEX ai_generations_status_idx ON ai_generations(status);
