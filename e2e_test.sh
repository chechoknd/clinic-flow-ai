#!/usr/bin/env bash
set -euo pipefail

API_BASE_URL="${API_BASE_URL:-http://127.0.0.1:18080}"
DEMO_EMAIL="${DEMO_EMAIL:-admin@sonrisaviva.demo}"
DEMO_PASSWORD="${DEMO_PASSWORD:-clinicflow123}"

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

api_get() {
  local path="$1"
  curl -fsS "${API_BASE_URL}${path}" \
    -H "Authorization: Bearer ${TOKEN}"
}

api_post() {
  local path="$1"
  local payload="$2"
  curl -fsS -X POST "${API_BASE_URL}${path}" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d "${payload}"
}

api_put() {
  local path="$1"
  local payload="$2"
  curl -fsS -X PUT "${API_BASE_URL}${path}" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d "${payload}"
}

api_delete() {
  local path="$1"
  curl -fsS -X DELETE "${API_BASE_URL}${path}" \
    -H "Authorization: Bearer ${TOKEN}"
}

assert_jq() {
  local payload="$1"
  local filter="$2"
  local message="$3"
  if ! jq -e "${filter}" <<<"${payload}" >/dev/null; then
    echo "Assertion failed: ${message}" >&2
    echo "Payload:" >&2
    jq . <<<"${payload}" >&2
    exit 1
  fi
}

require_command curl
require_command jq

printf 'Using API: %s\n' "${API_BASE_URL}"

LOGIN_RESPONSE=$(curl -fsS -X POST "${API_BASE_URL}/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${DEMO_EMAIL}\",\"password\":\"${DEMO_PASSWORD}\"}")

TOKEN=$(jq -r '.access_token // empty' <<<"${LOGIN_RESPONSE}")
if [[ -z "${TOKEN}" ]]; then
  echo "Login did not return an access token." >&2
  exit 1
fi

echo "Login OK."

ORIGINAL_CLINIC=""
SMOKE_SERVICE_ID=""
SMOKE_PROFESSIONAL_ID=""
SMOKE_PROFESSIONAL_PAYLOAD=""

cleanup() {
  if [[ -n "${SMOKE_PROFESSIONAL_ID}" && -n "${SMOKE_PROFESSIONAL_PAYLOAD}" ]]; then
    PROFESSIONAL_DEACTIVATE_PAYLOAD=$(jq '. + {is_active: false}' <<<"${SMOKE_PROFESSIONAL_PAYLOAD}")
    api_put "/api/professionals/${SMOKE_PROFESSIONAL_ID}" "${PROFESSIONAL_DEACTIVATE_PAYLOAD}" >/dev/null 2>&1 || true
  fi
  if [[ -n "${SMOKE_SERVICE_ID}" ]]; then
    api_delete "/api/services/${SMOKE_SERVICE_ID}" >/dev/null 2>&1 || true
  fi
  if [[ -n "${ORIGINAL_CLINIC}" ]]; then
    RESTORE_PAYLOAD=$(jq '{
      name,
      city,
      phone,
      whatsapp,
      address,
      opening_hours: (.opening_hours // {}),
      general_faq: (.general_faq // []),
      communication_tone
    }' <<<"${ORIGINAL_CLINIC}")
    api_put "/api/clinics/current" "${RESTORE_PAYLOAD}" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

echo -e "\n1. GET /api/clinics/current"
ORIGINAL_CLINIC=$(api_get "/api/clinics/current")
echo "${ORIGINAL_CLINIC}" | jq .
assert_jq "${ORIGINAL_CLINIC}" '.id and .name and .whatsapp' "clinic profile should include id, name, and whatsapp"

echo -e "\n2. PUT /api/clinics/current"
CLINIC_UPDATE_PAYLOAD=$(jq '{
  name: (.name + " Smoke"),
  city,
  phone,
  whatsapp,
  address,
  opening_hours: (.opening_hours // {}),
  general_faq: (.general_faq // []),
  communication_tone: "profesional"
}' <<<"${ORIGINAL_CLINIC}")
CLINIC_UPDATE_RESPONSE=$(api_put "/api/clinics/current" "${CLINIC_UPDATE_PAYLOAD}")
echo "${CLINIC_UPDATE_RESPONSE}" | jq .
assert_jq "${CLINIC_UPDATE_RESPONSE}" '.name | endswith(" Smoke")' "clinic update should persist changed name"

echo -e "\n3. POST /api/services"
SERVICE_CREATE_PAYLOAD=$(jq -n '{
  name: "Servicio Smoke E2E",
  description: "Servicio creado por e2e_test.sh para validar catalogo.",
  duration_minutes: 30,
  price_from: 99000,
  benefits: ["Validacion local", "Contexto para AI"],
  faq: [],
  common_objections: ["Precio"]
}')
SERVICE_CREATE_RESPONSE=$(api_post "/api/services" "${SERVICE_CREATE_PAYLOAD}")
echo "${SERVICE_CREATE_RESPONSE}" | jq .
SMOKE_SERVICE_ID=$(jq -r '.id // empty' <<<"${SERVICE_CREATE_RESPONSE}")
if [[ -z "${SMOKE_SERVICE_ID}" ]]; then
  echo "Service creation did not return an id." >&2
  exit 1
fi

echo -e "\n4. PUT /api/services/${SMOKE_SERVICE_ID}"
SERVICE_UPDATE_PAYLOAD=$(jq -n '{
  name: "Servicio Smoke E2E Actualizado",
  description: "Servicio actualizado por e2e_test.sh.",
  duration_minutes: 45,
  price_from: 125000,
  benefits: ["Validacion local actualizada"],
  faq: [],
  common_objections: ["Tiempo"],
  is_active: true
}')
SERVICE_UPDATE_RESPONSE=$(api_put "/api/services/${SMOKE_SERVICE_ID}" "${SERVICE_UPDATE_PAYLOAD}")
echo "${SERVICE_UPDATE_RESPONSE}" | jq .
assert_jq "${SERVICE_UPDATE_RESPONSE}" '.name == "Servicio Smoke E2E Actualizado"' "service update should persist changed name"

echo -e "\n5. GET /api/services"
SERVICES_RESPONSE=$(api_get "/api/services")
echo "${SERVICES_RESPONSE}" | jq .
if ! jq -e --arg id "${SMOKE_SERVICE_ID}" 'if type == "array" then any(.[]; .id == $id) else any(.data[]; .id == $id) end' <<<"${SERVICES_RESPONSE}" >/dev/null; then
  echo "Assertion failed: service list should include smoke service" >&2
  exit 1
fi

echo -e "\n6. POST /api/leads"
NEXT_ACTION_AT=$(date -u -d '+1 day' '+%Y-%m-%dT%H:%M:%SZ')
LEAD_PAYLOAD=$(jq -n --arg service_id "${SMOKE_SERVICE_ID}" --arg next_action_at "${NEXT_ACTION_AT}" '{
  full_name: "Lead Smoke Test",
  phone: "+573009998888",
  service_id: $service_id,
  status: "Nuevo",
  source: "whatsapp",
  notes: "Lead creado por e2e_test.sh para validar el flujo local.",
  next_action_at: $next_action_at
}')
LEAD_RESPONSE=$(api_post "/api/leads" "${LEAD_PAYLOAD}")
echo "${LEAD_RESPONSE}" | jq .
LEAD_ID=$(jq -r '.id // empty' <<<"${LEAD_RESPONSE}")
if [[ -z "${LEAD_ID}" ]]; then
  echo "Lead creation did not return an id." >&2
  exit 1
fi

echo -e "\n7. PUT /api/leads/${LEAD_ID}"
UPDATED_NEXT_ACTION_AT=$(date -u -d '+2 days' '+%Y-%m-%dT%H:%M:%SZ')
LEAD_UPDATE_PAYLOAD=$(jq -n --arg next_action_at "${UPDATED_NEXT_ACTION_AT}" '{
  status: "Contactado",
  note: "Nota agregada por smoke e2e.",
  next_action_at: $next_action_at
}')
LEAD_UPDATE_RESPONSE=$(api_put "/api/leads/${LEAD_ID}" "${LEAD_UPDATE_PAYLOAD}")
echo "${LEAD_UPDATE_RESPONSE}" | jq .
assert_jq "${LEAD_UPDATE_RESPONSE}" '.status == "Contactado"' "lead update should return Contactado status"

echo -e "\n8. GET /api/leads/${LEAD_ID}"
LEAD_DETAIL_RESPONSE=$(api_get "/api/leads/${LEAD_ID}")
echo "${LEAD_DETAIL_RESPONSE}" | jq .
assert_jq "${LEAD_DETAIL_RESPONSE}" '.notes | length >= 1' "lead detail should include notes"

echo -e "\n9. GET /api/followups"
FOLLOWUPS_RESPONSE=$(api_get "/api/followups")
echo "${FOLLOWUPS_RESPONSE}" | jq .
if ! jq -e --arg id "${LEAD_ID}" 'any(.data[]; .id == $id)' <<<"${FOLLOWUPS_RESPONSE}" >/dev/null; then
  echo "Assertion failed: followups should include lead with next action" >&2
  exit 1
fi

echo -e "\n10. POST /api/followups/${LEAD_ID}/reschedule"
RESCHEDULE_AT=$(date -u -d '+3 days' '+%Y-%m-%dT%H:%M:%SZ')
RESCHEDULE_PAYLOAD=$(jq -n --arg next_action_at "${RESCHEDULE_AT}" '{
  next_action_at: $next_action_at,
  note: "Reprogramado por smoke e2e."
}')
RESCHEDULE_RESPONSE=$(api_post "/api/followups/${LEAD_ID}/reschedule" "${RESCHEDULE_PAYLOAD}")
echo "${RESCHEDULE_RESPONSE}" | jq .
assert_jq "${RESCHEDULE_RESPONSE}" '.status == "rescheduled"' "follow-up reschedule should return rescheduled status"

echo -e "\n11. POST /api/ai/reply-suggestion"
AI_REPLY_PAYLOAD=$(jq -n --arg lead_id "${LEAD_ID}" --arg service_id "${SMOKE_SERVICE_ID}" '{
  lead_id: $lead_id,
  service_id: $service_id,
  patient_message: "Quiero saber precio y disponibilidad.",
  desired_tone: "cercano"
}')
AI_REPLY_RESPONSE=$(api_post "/api/ai/reply-suggestion" "${AI_REPLY_PAYLOAD}")
echo "${AI_REPLY_RESPONSE}" | jq .
assert_jq "${AI_REPLY_RESPONSE}" '.variants.short and .safety_status' "AI reply should include variants and safety status"

echo -e "\n12. POST /api/ai/objection-handler"
AI_OBJECTION_PAYLOAD=$(jq -n --arg lead_id "${LEAD_ID}" --arg service_id "${SMOKE_SERVICE_ID}" '{
  lead_id: $lead_id,
  service_id: $service_id,
  objection: "Me parece costoso."
}')
AI_OBJECTION_RESPONSE=$(api_post "/api/ai/objection-handler" "${AI_OBJECTION_PAYLOAD}")
echo "${AI_OBJECTION_RESPONSE}" | jq .
assert_jq "${AI_OBJECTION_RESPONSE}" '.suggested_message and .safety_status' "AI objection handler should include suggested message"

echo -e "\n13. POST /api/ai/follow-up-message"
AI_FOLLOWUP_PAYLOAD=$(jq -n --arg lead_id "${LEAD_ID}" --arg service_id "${SMOKE_SERVICE_ID}" '{
  lead_id: $lead_id,
  service_id: $service_id,
  last_contact_note: "Pidio informacion y quedo en revisar."
}')
AI_FOLLOWUP_RESPONSE=$(api_post "/api/ai/follow-up-message" "${AI_FOLLOWUP_PAYLOAD}")
echo "${AI_FOLLOWUP_RESPONSE}" | jq .
assert_jq "${AI_FOLLOWUP_RESPONSE}" '.suggested_message and .next_step and .safety_status' "AI follow-up should include suggested message and next step"

echo -e "\n14. POST /api/followups/${LEAD_ID}/complete"
COMPLETE_PAYLOAD=$(jq -n '{
  status: "Interesado",
  note: "Seguimiento completado por smoke e2e."
}')
COMPLETE_RESPONSE=$(api_post "/api/followups/${LEAD_ID}/complete" "${COMPLETE_PAYLOAD}")
echo "${COMPLETE_RESPONSE}" | jq .
assert_jq "${COMPLETE_RESPONSE}" '.status == "completed"' "follow-up complete should return completed status"

echo -e "\n15. GET /api/dashboard/summary"
DASHBOARD_RESPONSE=$(api_get "/api/dashboard/summary")
echo "${DASHBOARD_RESPONSE}" | jq .
assert_jq "${DASHBOARD_RESPONSE}" '.leads_total >= 1' "dashboard should report at least one lead"

SCHEDULE_SERVICE_ID=$(jq -r --arg smoke "${SMOKE_SERVICE_ID}" '
  def items: if type == "array" then . else (.data // []) end;
  ([items[] | select(.id != $smoke and (.is_active // true)) | .id][0] // $smoke)
' <<<"${SERVICES_RESPONSE}")
if [[ -z "${SCHEDULE_SERVICE_ID}" || "${SCHEDULE_SERVICE_ID}" == "null" ]]; then
  SCHEDULE_SERVICE_ID="${SMOKE_SERVICE_ID}"
fi

SMOKE_RUN_ID=$(date -u '+%Y%m%d%H%M%S')
SCHEDULE_DATE=$(date -u -d '+7 days' '+%Y-%m-%d')
APPOINTMENT_START="${SCHEDULE_DATE}T14:00:00Z"
APPOINTMENT_RESCHEDULE_START="${SCHEDULE_DATE}T15:00:00Z"
CONVERTED_APPOINTMENT_START="${SCHEDULE_DATE}T16:00:00Z"

WORKING_HOURS=$(jq -n '{
  monday: [{start: "08:00", end: "18:00"}],
  tuesday: [{start: "08:00", end: "18:00"}],
  wednesday: [{start: "08:00", end: "18:00"}],
  thursday: [{start: "08:00", end: "18:00"}],
  friday: [{start: "08:00", end: "18:00"}],
  saturday: [{start: "08:00", end: "14:00"}],
  sunday: [{start: "08:00", end: "14:00"}]
}')

SMOKE_PROFESSIONAL_PAYLOAD=$(jq -n \
  --arg full_name "Dra. Smoke E2E ${SMOKE_RUN_ID}" \
  --arg role "Agenda Operativa" \
  --arg color "#0F766E" \
  --arg service_id "${SCHEDULE_SERVICE_ID}" \
  --argjson working_hours "${WORKING_HOURS}" '{
    full_name: $full_name,
    role_or_specialty: $role,
    calendar_color: $color,
    working_hours: $working_hours,
    service_ids: [$service_id]
  }')

echo -e "\n16. POST /api/professionals"
PROFESSIONAL_RESPONSE=$(api_post "/api/professionals" "${SMOKE_PROFESSIONAL_PAYLOAD}")
echo "${PROFESSIONAL_RESPONSE}" | jq .
SMOKE_PROFESSIONAL_ID=$(jq -r '.id // empty' <<<"${PROFESSIONAL_RESPONSE}")
if [[ -z "${SMOKE_PROFESSIONAL_ID}" ]]; then
  echo "Professional creation did not return an id." >&2
  exit 1
fi
assert_jq "${PROFESSIONAL_RESPONSE}" '.is_active == true and (.service_ids | length >= 1)' "professional should be active and linked to a service"

SMOKE_PROFESSIONAL_PAYLOAD=$(jq '. + {is_active: true}' <<<"${SMOKE_PROFESSIONAL_PAYLOAD}")

echo -e "\n17. GET /api/professionals/${SMOKE_PROFESSIONAL_ID}"
PROFESSIONAL_DETAIL_RESPONSE=$(api_get "/api/professionals/${SMOKE_PROFESSIONAL_ID}")
echo "${PROFESSIONAL_DETAIL_RESPONSE}" | jq .
if ! jq -e --arg id "${SMOKE_PROFESSIONAL_ID}" '.id == $id' <<<"${PROFESSIONAL_DETAIL_RESPONSE}" >/dev/null; then
  echo "Assertion failed: professional detail should return created professional" >&2
  exit 1
fi

echo -e "\n18. GET /api/schedule/availability"
AVAILABILITY_RESPONSE=$(api_get "/api/schedule/availability?professional_id=${SMOKE_PROFESSIONAL_ID}&service_id=${SCHEDULE_SERVICE_ID}&date_from=${SCHEDULE_DATE}&date_to=${SCHEDULE_DATE}&duration_minutes=30")
echo "${AVAILABILITY_RESPONSE}" | jq .
if ! jq -e --arg id "${SMOKE_PROFESSIONAL_ID}" '.professional_id == $id and (.slots | length >= 1)' <<<"${AVAILABILITY_RESPONSE}" >/dev/null; then
  echo "Assertion failed: availability should return at least one slot for the smoke professional" >&2
  exit 1
fi

echo -e "\n19. POST /api/appointments"
APPOINTMENT_PAYLOAD=$(jq -n \
  --arg professional_id "${SMOKE_PROFESSIONAL_ID}" \
  --arg service_id "${SCHEDULE_SERVICE_ID}" \
  --arg starts_at "${APPOINTMENT_START}" '{
    professional_id: $professional_id,
    service_id: $service_id,
    contact_name: "Contacto Agenda Smoke",
    contact_phone: "+573009997777",
    starts_at: $starts_at,
    duration_minutes: 30,
    status: "pending_confirmation",
    source: "whatsapp",
    admin_notes: "Cita operativa creada por smoke e2e."
  }')
APPOINTMENT_RESPONSE=$(api_post "/api/appointments" "${APPOINTMENT_PAYLOAD}")
echo "${APPOINTMENT_RESPONSE}" | jq .
APPOINTMENT_ID=$(jq -r '.id // empty' <<<"${APPOINTMENT_RESPONSE}")
if [[ -z "${APPOINTMENT_ID}" ]]; then
  echo "Appointment creation did not return an id." >&2
  exit 1
fi
assert_jq "${APPOINTMENT_RESPONSE}" '.status == "pending_confirmation" and .confirmation_status == "pending"' "appointment should start pending confirmation"

echo -e "\n20. GET /api/appointments?date=${SCHEDULE_DATE}&professional_id=${SMOKE_PROFESSIONAL_ID}"
APPOINTMENTS_RESPONSE=$(api_get "/api/appointments?date=${SCHEDULE_DATE}&professional_id=${SMOKE_PROFESSIONAL_ID}")
echo "${APPOINTMENTS_RESPONSE}" | jq .
if ! jq -e --arg id "${APPOINTMENT_ID}" 'any(.data[]; .id == $id)' <<<"${APPOINTMENTS_RESPONSE}" >/dev/null; then
  echo "Assertion failed: appointment list should include smoke appointment" >&2
  exit 1
fi

echo -e "\n21. POST /api/appointments/${APPOINTMENT_ID}/reschedule"
APPOINTMENT_RESCHEDULE_PAYLOAD=$(jq -n --arg starts_at "${APPOINTMENT_RESCHEDULE_START}" '{
  starts_at: $starts_at,
  duration_minutes: 30,
  admin_note: "Reprogramada por smoke e2e."
}')
APPOINTMENT_RESCHEDULE_RESPONSE=$(api_post "/api/appointments/${APPOINTMENT_ID}/reschedule" "${APPOINTMENT_RESCHEDULE_PAYLOAD}")
echo "${APPOINTMENT_RESCHEDULE_RESPONSE}" | jq .
assert_jq "${APPOINTMENT_RESCHEDULE_RESPONSE}" '.status == "rescheduled"' "appointment reschedule should return rescheduled status"

echo -e "\n22. POST /api/appointments/${APPOINTMENT_ID}/status"
APPOINTMENT_STATUS_PAYLOAD=$(jq -n '{
  status: "confirmed",
  admin_note: "Confirmada manualmente por smoke e2e."
}')
APPOINTMENT_STATUS_RESPONSE=$(api_post "/api/appointments/${APPOINTMENT_ID}/status" "${APPOINTMENT_STATUS_PAYLOAD}")
echo "${APPOINTMENT_STATUS_RESPONSE}" | jq .
assert_jq "${APPOINTMENT_STATUS_RESPONSE}" '.status == "confirmed"' "appointment status update should return confirmed status"

echo -e "\n23. POST /api/leads for schedule conversion"
SCHEDULE_LEAD_PAYLOAD=$(jq -n --arg service_id "${SCHEDULE_SERVICE_ID}" '{
  full_name: "Lead Agenda Smoke",
  phone: "+573009996666",
  service_id: $service_id,
  status: "Interesado",
  source: "whatsapp",
  notes: "Lead operativo creado por smoke e2e para convertir a cita."
}')
SCHEDULE_LEAD_RESPONSE=$(api_post "/api/leads" "${SCHEDULE_LEAD_PAYLOAD}")
echo "${SCHEDULE_LEAD_RESPONSE}" | jq .
SCHEDULE_LEAD_ID=$(jq -r '.id // empty' <<<"${SCHEDULE_LEAD_RESPONSE}")
if [[ -z "${SCHEDULE_LEAD_ID}" ]]; then
  echo "Schedule lead creation did not return an id." >&2
  exit 1
fi

echo -e "\n24. POST /api/leads/${SCHEDULE_LEAD_ID}/convert-to-appointment"
CONVERT_PAYLOAD=$(jq -n \
  --arg professional_id "${SMOKE_PROFESSIONAL_ID}" \
  --arg service_id "${SCHEDULE_SERVICE_ID}" \
  --arg starts_at "${CONVERTED_APPOINTMENT_START}" '{
    professional_id: $professional_id,
    service_id: $service_id,
    starts_at: $starts_at,
    duration_minutes: 30,
    status: "pending_confirmation",
    admin_notes: "Conversion humana validada por smoke e2e.",
    update_lead_status: true
  }')
CONVERT_RESPONSE=$(api_post "/api/leads/${SCHEDULE_LEAD_ID}/convert-to-appointment" "${CONVERT_PAYLOAD}")
echo "${CONVERT_RESPONSE}" | jq .
if ! jq -e --arg lead_id "${SCHEDULE_LEAD_ID}" '.lead_id == $lead_id and .lead_status == "Agendado" and .appointment_status == "pending_confirmation"' <<<"${CONVERT_RESPONSE}" >/dev/null; then
  echo "Assertion failed: lead conversion should create appointment and update lead status" >&2
  exit 1
fi

echo -e "\n25. GET /api/dashboard/schedule-summary?date=${SCHEDULE_DATE}"
SCHEDULE_SUMMARY_RESPONSE=$(api_get "/api/dashboard/schedule-summary?date=${SCHEDULE_DATE}")
echo "${SCHEDULE_SUMMARY_RESPONSE}" | jq .
assert_jq "${SCHEDULE_SUMMARY_RESPONSE}" '.todays_appointments >= 2 and .appointments_pending_confirmation >= 1 and (.appointments_by_professional | length >= 1)' "schedule summary should include smoke appointments"

echo -e "\nSmoke test completed. Created lead: ${LEAD_ID}. Created schedule professional: ${SMOKE_PROFESSIONAL_ID}"
