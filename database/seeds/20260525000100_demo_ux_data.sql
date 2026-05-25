-- Extended demo UX seed data for local development only.
-- Requires 20260521000100_demo_core_data.sql to be applied first.
-- Data is fictional and commercial-only. Do not use real patient data.

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
    '33333333-3333-4333-8333-333333333333',
    '11111111-1111-4111-8111-111111111111',
    'Limpieza dental profesional',
    'Limpieza preventiva con valoracion inicial y orientacion comercial sobre cuidados posteriores.',
    45,
    90000,
    '["Cita rapida","Ayuda a mantener una sonrisa saludable","Recomendaciones de cuidado general"]'::jsonb,
    '[{"question":"Cada cuanto se recomienda?","answer":"La frecuencia se confirma en la valoracion segun cada persona."}]'::jsonb,
    '["No tengo tiempo","Lo puedo hacer despues"]'::jsonb,
    TRUE
),
(
    '33333333-3333-4333-8333-333333333334',
    '11111111-1111-4111-8111-111111111111',
    'Implantes dentales',
    'Valoracion comercial para conocer opciones de reposicion dental con plan profesional.',
    60,
    180000,
    '["Plan personalizado","Opciones de financiacion por confirmar","Acompanamiento profesional"]'::jsonb,
    '[{"question":"Puedo saber el precio exacto por WhatsApp?","answer":"Podemos dar un rango inicial y confirmar el plan despues de la valoracion."}]'::jsonb,
    '["Es costoso","Quiero pensarlo","Necesito financiacion"]'::jsonb,
    TRUE
),
(
    '33333333-3333-4333-8333-333333333335',
    '11111111-1111-4111-8111-111111111111',
    'Diseño de sonrisa',
    'Valoracion estetica para definir alternativas de mejora de sonrisa segun expectativas comerciales.',
    60,
    150000,
    '["Enfoque estetico personalizado","Opciones segun presupuesto","Explicacion clara del proceso"]'::jsonb,
    '[{"question":"Me pueden mostrar opciones?","answer":"Si, en la valoracion pueden explicarte alternativas segun tu objetivo."}]'::jsonb,
    '["Me preocupa que se vea artificial","Quiero comparar precios"]'::jsonb,
    TRUE
),
(
    '33333333-3333-4333-8333-333333333336',
    '11111111-1111-4111-8111-111111111111',
    'Carillas dentales',
    'Consulta de valoracion para revisar alternativas esteticas y resolver dudas comerciales.',
    60,
    220000,
    '["Alternativas esteticas","Valoracion personalizada","Acompanamiento en la decision"]'::jsonb,
    '[{"question":"Danan los dientes?","answer":"El profesional explica las opciones y condiciones durante la valoracion."}]'::jsonb,
    '["Me da miedo el procedimiento","No se si soy candidato"]'::jsonb,
    TRUE
),
(
    '33333333-3333-4333-8333-333333333337',
    '11111111-1111-4111-8111-111111111111',
    'Valoracion odontológica',
    'Cita inicial para escuchar el objetivo del paciente y orientar el siguiente paso comercial.',
    30,
    70000,
    '["Primer paso claro","Resolucion de dudas","Orientacion sobre servicios disponibles"]'::jsonb,
    '[{"question":"La valoracion se descuenta del tratamiento?","answer":"El equipo comercial puede confirmar promociones vigentes al agendar."}]'::jsonb,
    '["Solo quiero informacion","No se que servicio necesito"]'::jsonb,
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

INSERT INTO leads (
    id,
    clinic_id,
    full_name,
    phone,
    service_id,
    status,
    source,
    next_action_at,
    created_at
) VALUES
(
    '44444444-4444-4444-8444-444444444401',
    '11111111-1111-4111-8111-111111111111',
    'Laura Martinez Demo',
    '+573100000001',
    '33333333-3333-4333-8333-333333333331',
    'Nuevo',
    'whatsapp',
    NOW() - interval '2 hours',
    NOW() - interval '1 day'
),
(
    '44444444-4444-4444-8444-444444444402',
    '11111111-1111-4111-8111-111111111111',
    'Andres Rojas Demo',
    '+573100000002',
    '33333333-3333-4333-8333-333333333333',
    'Nuevo',
    'instagram',
    NOW() + interval '3 hours',
    NOW() - interval '8 hours'
),
(
    '44444444-4444-4444-8444-444444444403',
    '11111111-1111-4111-8111-111111111111',
    'Camila Torres Demo',
    '+573100000003',
    '33333333-3333-4333-8333-333333333332',
    'Contactado',
    'whatsapp',
    NOW() + interval '1 day',
    NOW() - interval '3 days'
),
(
    '44444444-4444-4444-8444-444444444404',
    '11111111-1111-4111-8111-111111111111',
    'Diego Salazar Demo',
    '+573100000004',
    '33333333-3333-4333-8333-333333333334',
    'Contactado',
    'facebook',
    NOW() - interval '1 day',
    NOW() - interval '4 days'
),
(
    '44444444-4444-4444-8444-444444444405',
    '11111111-1111-4111-8111-111111111111',
    'Natalia Perez Demo',
    '+573100000005',
    '33333333-3333-4333-8333-333333333335',
    'Interesado',
    'whatsapp',
    NOW() + interval '2 hours',
    NOW() - interval '2 days'
),
(
    '44444444-4444-4444-8444-444444444406',
    '11111111-1111-4111-8111-111111111111',
    'Felipe Gomez Demo',
    '+573100000006',
    '33333333-3333-4333-8333-333333333336',
    'Interesado',
    'web',
    NOW() + interval '2 days',
    NOW() - interval '5 days'
),
(
    '44444444-4444-4444-8444-444444444407',
    '11111111-1111-4111-8111-111111111111',
    'Valentina Castro Demo',
    '+573100000007',
    '33333333-3333-4333-8333-333333333337',
    'Agendado',
    'whatsapp',
    NOW() + interval '5 hours',
    NOW() - interval '6 days'
),
(
    '44444444-4444-4444-8444-444444444408',
    '11111111-1111-4111-8111-111111111111',
    'Sebastian Moreno Demo',
    '+573100000008',
    '33333333-3333-4333-8333-333333333331',
    'Agendado',
    'llamada',
    NOW() + interval '3 days',
    NOW() - interval '7 days'
),
(
    '44444444-4444-4444-8444-444444444409',
    '11111111-1111-4111-8111-111111111111',
    'Isabella Ramirez Demo',
    '+573100000009',
    '33333333-3333-4333-8333-333333333333',
    'No Respondio',
    'whatsapp',
    NOW() - interval '3 hours',
    NOW() - interval '8 days'
),
(
    '44444444-4444-4444-8444-444444444410',
    '11111111-1111-4111-8111-111111111111',
    'Mateo Herrera Demo',
    '+573100000010',
    '33333333-3333-4333-8333-333333333334',
    'No Respondio',
    'instagram',
    NOW() + interval '1 day 4 hours',
    NOW() - interval '9 days'
),
(
    '44444444-4444-4444-8444-444444444411',
    '11111111-1111-4111-8111-111111111111',
    'Sofia Vargas Demo',
    '+573100000011',
    '33333333-3333-4333-8333-333333333335',
    'Perdido',
    'facebook',
    NULL,
    NOW() - interval '10 days'
),
(
    '44444444-4444-4444-8444-444444444412',
    '11111111-1111-4111-8111-111111111111',
    'Juan Pablo Medina Demo',
    '+573100000012',
    '33333333-3333-4333-8333-333333333336',
    'Perdido',
    'whatsapp',
    NULL,
    NOW() - interval '11 days'
),
(
    '44444444-4444-4444-8444-444444444413',
    '11111111-1111-4111-8111-111111111111',
    'Daniela Ortiz Demo',
    '+573100000013',
    '33333333-3333-4333-8333-333333333332',
    'Convertido',
    'whatsapp',
    NULL,
    NOW() - interval '12 days'
),
(
    '44444444-4444-4444-8444-444444444414',
    '11111111-1111-4111-8111-111111111111',
    'Cristian Lopez Demo',
    '+573100000014',
    '33333333-3333-4333-8333-333333333337',
    'Convertido',
    'web',
    NULL,
    NOW() - interval '13 days'
)
ON CONFLICT (id) DO UPDATE SET
    clinic_id = EXCLUDED.clinic_id,
    full_name = EXCLUDED.full_name,
    phone = EXCLUDED.phone,
    service_id = EXCLUDED.service_id,
    status = EXCLUDED.status,
    source = EXCLUDED.source,
    next_action_at = EXCLUDED.next_action_at,
    created_at = EXCLUDED.created_at,
    updated_at = NOW();

INSERT INTO lead_notes (id, lead_id, body, created_at) VALUES
('55555555-5555-4555-8555-555555555401', '44444444-4444-4444-8444-444444444401', 'Consulta precio aproximado y disponibilidad esta semana.', NOW() - interval '1 day'),
('55555555-5555-4555-8555-555555555402', '44444444-4444-4444-8444-444444444402', 'Pidio informacion general y prefiere contacto por WhatsApp.', NOW() - interval '8 hours'),
('55555555-5555-4555-8555-555555555403', '44444444-4444-4444-8444-444444444403', 'Ya recibio informacion inicial y quiere comparar horarios.', NOW() - interval '3 days'),
('55555555-5555-4555-8555-555555555404', '44444444-4444-4444-8444-444444444404', 'Solicito rango de inversion y opciones de pago disponibles.', NOW() - interval '4 days'),
('55555555-5555-4555-8555-555555555405', '44444444-4444-4444-8444-444444444405', 'Alta intencion comercial; quiere resolver dudas antes de agendar.', NOW() - interval '2 days'),
('55555555-5555-4555-8555-555555555406', '44444444-4444-4444-8444-444444444406', 'Interesado en opciones esteticas y en conocer tiempos aproximados.', NOW() - interval '5 days'),
('55555555-5555-4555-8555-555555555407', '44444444-4444-4444-8444-444444444407', 'Cita acordada pendiente de recordatorio manual.', NOW() - interval '6 days'),
('55555555-5555-4555-8555-555555555408', '44444444-4444-4444-8444-444444444408', 'Agendo valoracion y pidio confirmacion un dia antes.', NOW() - interval '7 days'),
('55555555-5555-4555-8555-555555555409', '44444444-4444-4444-8444-444444444409', 'No respondio al ultimo mensaje; intentar recuperacion comercial.', NOW() - interval '8 days'),
('55555555-5555-4555-8555-555555555410', '44444444-4444-4444-8444-444444444410', 'No respondio; reintentar con mensaje corto y claro.', NOW() - interval '9 days'),
('55555555-5555-4555-8555-555555555411', '44444444-4444-4444-8444-444444444411', 'Oportunidad cerrada por presupuesto; mantener tono amable.', NOW() - interval '10 days'),
('55555555-5555-4555-8555-555555555412', '44444444-4444-4444-8444-444444444412', 'Decidio no continuar por ahora.', NOW() - interval '11 days'),
('55555555-5555-4555-8555-555555555413', '44444444-4444-4444-8444-444444444413', 'Lead convertido despues de seguimiento comercial.', NOW() - interval '12 days'),
('55555555-5555-4555-8555-555555555414', '44444444-4444-4444-8444-444444444414', 'Lead convertido desde formulario web.', NOW() - interval '13 days')
ON CONFLICT (id) DO UPDATE SET
    lead_id = EXCLUDED.lead_id,
    body = EXCLUDED.body,
    created_at = EXCLUDED.created_at;
