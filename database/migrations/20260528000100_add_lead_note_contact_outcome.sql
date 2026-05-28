ALTER TABLE lead_notes
ADD COLUMN contact_outcome VARCHAR(50);

ALTER TABLE lead_notes
ADD CONSTRAINT lead_notes_contact_outcome_check CHECK (
    contact_outcome IS NULL OR contact_outcome IN (
        'attempted_no_answer',
        'asked_price',
        'interested',
        'scheduled',
        'lost_price',
        'lost_timing',
        'lost_trust',
        'converted',
        'follow_up_requested',
        'other'
    )
);

CREATE INDEX lead_notes_contact_outcome_idx ON lead_notes(contact_outcome);
