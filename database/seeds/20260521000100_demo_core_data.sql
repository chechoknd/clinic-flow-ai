-- Demo seed data for local development only.
-- Demo login: admin@sonrisaviva.demo / clinicflow123

INSERT INTO clinics (
    id,
    name,
    clinic_type,
    city,
    phone,
    whatsapp,
    address,
    opening_hours,
    general_faq,
    communication_tone
) VALUES (
    '11111111-1111-4111-8111-111111111111',
    'Sonrisa Viva Demo',
    'odontologia',
    'Bogota',
    '+573001112233',
    '+573001112233',
    'Calle 123 #45-67',
    '{"monday_friday":"08:00-18:00","saturday":"08:00-13:00"}'::jsonb,
    '[{"question":"Atienden urgencias?","answer":"Si, podemos validar disponibilidad por WhatsApp."}]'::jsonb,
    'amable'
) ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    clinic_type = EXCLUDED.clinic_type,
    city = EXCLUDED.city,
    phone = EXCLUDED.phone,
    whatsapp = EXCLUDED.whatsapp,
    address = EXCLUDED.address,
    opening_hours = EXCLUDED.opening_hours,
    general_faq = EXCLUDED.general_faq,
    communication_tone = EXCLUDED.communication_tone;

INSERT INTO users (
    id,
    clinic_id,
    full_name,
    email,
    password_hash,
    role,
    is_active
) VALUES (
    '22222222-2222-4222-8222-222222222222',
    '11111111-1111-4111-8111-111111111111',
    'Admin Demo',
    'admin@sonrisaviva.demo',
    '$2a$10$NMzKrHgbh9V1Dllt5dBCd.lz9KnykcQLtoLeb0C63LYpUEQwrNyt.',
    'clinic_admin',
    TRUE
) ON CONFLICT (email) DO UPDATE SET
    clinic_id = EXCLUDED.clinic_id,
    full_name = EXCLUDED.full_name,
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    is_active = EXCLUDED.is_active;

INSERT INTO clinic_services (
    id,
    clinic_id,
    name,
    description,
    duration_minutes,
    price_from,
    benefits,
    faq,
    common_objections,
    is_active
) VALUES
(
    '33333333-3333-4333-8333-333333333331',
    '11111111-1111-4111-8111-111111111111',
    'Blanqueamiento dental',
    'Tratamiento estetico para mejorar el tono de la sonrisa con valoracion previa.',
    60,
    250000,
    '["Mejora estetica visible","Valoracion personalizada"]'::jsonb,
    '[{"question":"Debilita los dientes?","answer":"La valoracion permite confirmar si el paciente es apto."}]'::jsonb,
    '["Esta muy caro","Me da miedo"]'::jsonb,
    TRUE
),
(
    '33333333-3333-4333-8333-333333333332',
    '11111111-1111-4111-8111-111111111111',
    'Ortodoncia',
    'Valoracion y plan de tratamiento para alinear la sonrisa.',
    45,
    120000,
    '["Plan personalizado","Seguimiento profesional"]'::jsonb,
    '[{"question":"Cuanto dura?","answer":"Depende del caso y se define despues de la valoracion."}]'::jsonb,
    '["Es muy largo","En otro lado es mas barato"]'::jsonb,
    TRUE
)
ON CONFLICT (clinic_id, name) DO UPDATE SET
    description = EXCLUDED.description,
    duration_minutes = EXCLUDED.duration_minutes,
    price_from = EXCLUDED.price_from,
    benefits = EXCLUDED.benefits,
    faq = EXCLUDED.faq,
    common_objections = EXCLUDED.common_objections,
    is_active = EXCLUDED.is_active;
