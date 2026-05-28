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

api_request_status() {
  local method="$1"
  local path="$2"
  local payload="${3:-}"
  local response

  if [[ -n "${payload}" ]]; then
    response=$(curl -sS -w $'\n%{http_code}' -X "${method}" "${API_BASE_URL}${path}" \
      -H "Authorization: Bearer ${TOKEN}" \
      -H "Content-Type: application/json" \
      -d "${payload}")
  else
    response=$(curl -sS -w $'\n%{http_code}' -X "${method}" "${API_BASE_URL}${path}" \
      -H "Authorization: Bearer ${TOKEN}")
  fi

  HTTP_STATUS="${response##*$'\n'}"
  HTTP_BODY="${response%$'\n'*}"
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
CURRENCY_SERVICE_ID=""

cleanup() {
  if [[ -n "${CURRENCY_SERVICE_ID}" ]]; then
    api_delete "/api/services/${CURRENCY_SERVICE_ID}" >/dev/null 2>&1 || true
  fi
  if [[ -n "${SMOKE_SERVICE_ID}" ]]; then
    api_delete "/api/services/${SMOKE_SERVICE_ID}" >/dev/null 2>&1 || true
  fi
  if [[ -n "${ORIGINAL_CLINIC}" ]]; then
    RESTORE_PAYLOAD=$(jq '{
      name,
      city,
      country_code,
      currency_code,
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
ORIGINAL_COUNTRY_CODE=$(jq -r '.country_code' <<<"${ORIGINAL_CLINIC}")
ORIGINAL_CURRENCY_CODE=$(jq -r '.currency_code' <<<"${ORIGINAL_CLINIC}")

echo -e "\n2. PUT /api/clinics/current"
CLINIC_UPDATE_PAYLOAD=$(jq '{
  name: (.name + " Smoke"),
  city,
  country_code,
  currency_code,
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
  price_from: "99000",
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
  price_from: "125000",
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

echo -e "\n5b. Currency price validation"
COP_CLINIC_PAYLOAD=$(jq '{
  name,
  city,
  country_code: "CO",
  currency_code: "COP",
  phone,
  whatsapp,
  address,
  opening_hours: (.opening_hours // {}),
  general_faq: (.general_faq // []),
  communication_tone
}' <<<"${ORIGINAL_CLINIC}")
COP_CLINIC_RESPONSE=$(api_put "/api/clinics/current" "${COP_CLINIC_PAYLOAD}")
echo "${COP_CLINIC_RESPONSE}" | jq .
assert_jq "${COP_CLINIC_RESPONSE}" '.country_code == "CO" and .currency_code == "COP"' "clinic should switch to COP"

COP_DECIMAL_PAYLOAD=$(jq -n '{
  name: "Servicio Decimal COP E2E",
  description: "Debe ser rechazado porque COP no permite decimales.",
  duration_minutes: 30,
  price_from: "120000.50",
  benefits: [],
  faq: [],
  common_objections: []
}')
api_request_status "POST" "/api/services" "${COP_DECIMAL_PAYLOAD}"
echo "${HTTP_BODY}" | jq .
if [[ "${HTTP_STATUS}" != "400" ]]; then
  echo "Assertion failed: COP decimal price should return 400, got ${HTTP_STATUS}" >&2
  exit 1
fi
assert_jq "${HTTP_BODY}" '.error.message == "price_from has too many decimal places for clinic currency"' "COP decimal rejection should explain precision"

PEN_CLINIC_PAYLOAD=$(jq '{
  name,
  city,
  country_code: "PE",
  currency_code: "PEN",
  phone,
  whatsapp,
  address,
  opening_hours: (.opening_hours // {}),
  general_faq: (.general_faq // []),
  communication_tone
}' <<<"${ORIGINAL_CLINIC}")
PEN_CLINIC_RESPONSE=$(api_put "/api/clinics/current" "${PEN_CLINIC_PAYLOAD}")
echo "${PEN_CLINIC_RESPONSE}" | jq .
assert_jq "${PEN_CLINIC_RESPONSE}" '.country_code == "PE" and .currency_code == "PEN"' "clinic should switch to PEN"

PEN_DECIMAL_PAYLOAD=$(jq -n '{
  name: "Servicio Decimal PEN E2E",
  description: "Debe ser aceptado porque PEN permite dos decimales.",
  duration_minutes: 30,
  price_from: "120.50",
  benefits: ["Decimal exacto"],
  faq: [],
  common_objections: []
}')
PEN_DECIMAL_RESPONSE=$(api_post "/api/services" "${PEN_DECIMAL_PAYLOAD}")
echo "${PEN_DECIMAL_RESPONSE}" | jq .
CURRENCY_SERVICE_ID=$(jq -r '.id // empty' <<<"${PEN_DECIMAL_RESPONSE}")
if [[ -z "${CURRENCY_SERVICE_ID}" ]]; then
  echo "Currency service creation did not return an id." >&2
  exit 1
fi
assert_jq "${PEN_DECIMAL_RESPONSE}" '.price_from == "120.50" and .currency_code == "PEN"' "PEN decimal service should return exact string price and PEN currency"
api_delete "/api/services/${CURRENCY_SERVICE_ID}" >/dev/null
CURRENCY_SERVICE_ID=""

RESTORE_CLINIC_PAYLOAD=$(jq '{
  name,
  city,
  country_code,
  currency_code,
  phone,
  whatsapp,
  address,
  opening_hours: (.opening_hours // {}),
  general_faq: (.general_faq // []),
  communication_tone
}' <<<"${ORIGINAL_CLINIC}")
RESTORE_CLINIC_RESPONSE=$(api_put "/api/clinics/current" "${RESTORE_CLINIC_PAYLOAD}")
echo "${RESTORE_CLINIC_RESPONSE}" | jq .
if ! jq -e --arg country "${ORIGINAL_COUNTRY_CODE}" --arg currency "${ORIGINAL_CURRENCY_CODE}" '.country_code == $country and .currency_code == $currency' <<<"${RESTORE_CLINIC_RESPONSE}" >/dev/null; then
  echo "Assertion failed: clinic currency should be restored after currency validation" >&2
  echo "Payload:" >&2
  jq . <<<"${RESTORE_CLINIC_RESPONSE}" >&2
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

echo -e "\n7. GET /api/dashboard/actions"
DASHBOARD_ACTIONS_RESPONSE=$(api_get "/api/dashboard/actions?limit=20")
echo "${DASHBOARD_ACTIONS_RESPONSE}" | jq .
if ! jq -e --arg id "${LEAD_ID}" 'any(.data[]; .lead_id == $id and .type == "new_lead")' <<<"${DASHBOARD_ACTIONS_RESPONSE}" >/dev/null; then
  echo "Assertion failed: dashboard actions should include newly created lead" >&2
  exit 1
fi

echo -e "\n8. PUT /api/leads/${LEAD_ID}"
UPDATED_NEXT_ACTION_AT=$(date -u -d '+2 days' '+%Y-%m-%dT%H:%M:%SZ')
LEAD_UPDATE_PAYLOAD=$(jq -n --arg next_action_at "${UPDATED_NEXT_ACTION_AT}" '{
  status: "Contactado",
  note: "Nota agregada por smoke e2e.",
  contact_outcome: "asked_price",
  next_action_at: $next_action_at
}')
LEAD_UPDATE_RESPONSE=$(api_put "/api/leads/${LEAD_ID}" "${LEAD_UPDATE_PAYLOAD}")
echo "${LEAD_UPDATE_RESPONSE}" | jq .
assert_jq "${LEAD_UPDATE_RESPONSE}" '.status == "Contactado"' "lead update should return Contactado status"

echo -e "\n9. GET /api/leads/${LEAD_ID}"
LEAD_DETAIL_RESPONSE=$(api_get "/api/leads/${LEAD_ID}")
echo "${LEAD_DETAIL_RESPONSE}" | jq .
assert_jq "${LEAD_DETAIL_RESPONSE}" '.notes | length >= 1' "lead detail should include notes"
assert_jq "${LEAD_DETAIL_RESPONSE}" '.notes[0].contact_outcome == "asked_price"' "lead detail should include contact outcome metadata"

echo -e "\n10. GET /api/followups"
FOLLOWUPS_RESPONSE=$(api_get "/api/followups")
echo "${FOLLOWUPS_RESPONSE}" | jq .
if ! jq -e --arg id "${LEAD_ID}" 'any(.data[]; .id == $id)' <<<"${FOLLOWUPS_RESPONSE}" >/dev/null; then
  echo "Assertion failed: followups should include lead with next action" >&2
  exit 1
fi

echo -e "\n11. POST /api/followups/${LEAD_ID}/reschedule"
RESCHEDULE_AT=$(date -u -d '+3 days' '+%Y-%m-%dT%H:%M:%SZ')
RESCHEDULE_PAYLOAD=$(jq -n --arg next_action_at "${RESCHEDULE_AT}" '{
  next_action_at: $next_action_at,
  note: "Reprogramado por smoke e2e."
}')
RESCHEDULE_RESPONSE=$(api_post "/api/followups/${LEAD_ID}/reschedule" "${RESCHEDULE_PAYLOAD}")
echo "${RESCHEDULE_RESPONSE}" | jq .
assert_jq "${RESCHEDULE_RESPONSE}" '.status == "rescheduled"' "follow-up reschedule should return rescheduled status"

echo -e "\n12. POST /api/ai/reply-suggestion"
AI_REPLY_PAYLOAD=$(jq -n --arg lead_id "${LEAD_ID}" --arg service_id "${SMOKE_SERVICE_ID}" '{
  lead_id: $lead_id,
  service_id: $service_id,
  patient_message: "Quiero saber precio y disponibilidad.",
  desired_tone: "cercano"
}')
AI_REPLY_RESPONSE=$(api_post "/api/ai/reply-suggestion" "${AI_REPLY_PAYLOAD}")
echo "${AI_REPLY_RESPONSE}" | jq .
assert_jq "${AI_REPLY_RESPONSE}" '.variants.short and .safety_status' "AI reply should include variants and safety status"

echo -e "\n13. POST /api/ai/objection-handler"
AI_OBJECTION_PAYLOAD=$(jq -n --arg lead_id "${LEAD_ID}" --arg service_id "${SMOKE_SERVICE_ID}" '{
  lead_id: $lead_id,
  service_id: $service_id,
  objection: "Me parece costoso."
}')
AI_OBJECTION_RESPONSE=$(api_post "/api/ai/objection-handler" "${AI_OBJECTION_PAYLOAD}")
echo "${AI_OBJECTION_RESPONSE}" | jq .
assert_jq "${AI_OBJECTION_RESPONSE}" '.suggested_message and .safety_status' "AI objection handler should include suggested message"

echo -e "\n14. POST /api/ai/follow-up-message"
AI_FOLLOWUP_PAYLOAD=$(jq -n --arg lead_id "${LEAD_ID}" --arg service_id "${SMOKE_SERVICE_ID}" '{
  lead_id: $lead_id,
  service_id: $service_id,
  last_contact_note: "Pidio informacion y quedo en revisar."
}')
AI_FOLLOWUP_RESPONSE=$(api_post "/api/ai/follow-up-message" "${AI_FOLLOWUP_PAYLOAD}")
echo "${AI_FOLLOWUP_RESPONSE}" | jq .
assert_jq "${AI_FOLLOWUP_RESPONSE}" '.suggested_message and .next_step and .safety_status' "AI follow-up should include suggested message and next step"

echo -e "\n15. POST /api/ai/analyze-conversation"
AI_ANALYZE_PAYLOAD=$(jq -n --arg lead_id "${LEAD_ID}" --arg service_id "${SMOKE_SERVICE_ID}" '{
  lead_id: $lead_id,
  service_id: $service_id,
  source: "whatsapp",
  conversation_text: "Paciente: Hola, quiero saber precio y disponibilidad para una valoracion. Clinica: Con gusto te ayudamos con informacion general."
}')
AI_ANALYZE_RESPONSE=$(api_post "/api/ai/analyze-conversation" "${AI_ANALYZE_PAYLOAD}")
echo "${AI_ANALYZE_RESPONSE}" | jq .
assert_jq "${AI_ANALYZE_RESPONSE}" '.analysis_id and .suggested_reply and .safety_status' "AI conversation analysis should include analysis id, suggested reply, and safety status"

echo -e "\n16. PUT /api/leads/${LEAD_ID} with reviewed AI analysis"
REVIEWED_ANALYSIS_UPDATE_PAYLOAD=$(jq -n --argjson analysis "${AI_ANALYZE_RESPONSE}" '{
  status: "Interesado",
  note: "Analisis AI revisado por humano en smoke e2e.",
  contact_outcome: "interested",
  reviewed_ai_analysis: {
    analysis_id: $analysis.analysis_id,
    intent: $analysis.intent,
    detected_objections: $analysis.detected_objections,
    commercial_summary: $analysis.commercial_summary,
    suggested_next_action: $analysis.suggested_next_action,
    source: "whatsapp"
  }
}')
REVIEWED_ANALYSIS_UPDATE_RESPONSE=$(api_put "/api/leads/${LEAD_ID}" "${REVIEWED_ANALYSIS_UPDATE_PAYLOAD}")
echo "${REVIEWED_ANALYSIS_UPDATE_RESPONSE}" | jq .
assert_jq "${REVIEWED_ANALYSIS_UPDATE_RESPONSE}" '.status == "Interesado"' "reviewed AI analysis update should keep lead interested"

echo -e "\n17. GET /api/leads/${LEAD_ID} with AI insights"
LEAD_AI_DETAIL_RESPONSE=$(api_get "/api/leads/${LEAD_ID}")
echo "${LEAD_AI_DETAIL_RESPONSE}" | jq .
assert_jq "${LEAD_AI_DETAIL_RESPONSE}" '.ai_insights | length >= 1' "lead detail should include reviewed AI insights"

echo -e "\n18. GET /api/dashboard/actions with AI insight priority"
DASHBOARD_AI_ACTIONS_RESPONSE=$(api_get "/api/dashboard/actions?limit=20")
echo "${DASHBOARD_AI_ACTIONS_RESPONSE}" | jq .
if ! jq -e --arg id "${LEAD_ID}" 'any(.data[]; .lead_id == $id and (.type == "high_intent" or .type == "detected_objection"))' <<<"${DASHBOARD_AI_ACTIONS_RESPONSE}" >/dev/null; then
  echo "Assertion failed: dashboard actions should include reviewed AI insight priority" >&2
  exit 1
fi

echo -e "\n19. POST /api/followups/${LEAD_ID}/complete"
COMPLETE_PAYLOAD=$(jq -n '{
  status: "Interesado",
  note: "Seguimiento completado por smoke e2e.",
  contact_outcome: "follow_up_requested"
}')
COMPLETE_RESPONSE=$(api_post "/api/followups/${LEAD_ID}/complete" "${COMPLETE_PAYLOAD}")
echo "${COMPLETE_RESPONSE}" | jq .
assert_jq "${COMPLETE_RESPONSE}" '.status == "completed"' "follow-up complete should return completed status"

echo -e "\n20. GET /api/dashboard/summary"
DASHBOARD_RESPONSE=$(api_get "/api/dashboard/summary")
echo "${DASHBOARD_RESPONSE}" | jq .
assert_jq "${DASHBOARD_RESPONSE}" '.leads_total >= 1' "dashboard should report at least one lead"

echo -e "\nSmoke test completed. Created lead: ${LEAD_ID}"
