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

echo -e "\n1. GET /api/clinics/current"
api_get "/api/clinics/current" | jq .

echo -e "\n2. GET /api/services"
SERVICES_RESPONSE=$(api_get "/api/services")
echo "${SERVICES_RESPONSE}" | jq .
SERVICE_ID=$(jq -r 'if type == "array" then .[0].id // empty else .data[0].id // empty end' <<<"${SERVICES_RESPONSE}")
if [[ -z "${SERVICE_ID}" ]]; then
  echo "No service_id available. Run demo seeds before this smoke test." >&2
  exit 1
fi

echo -e "\n3. POST /api/leads"
LEAD_PAYLOAD=$(jq -n --arg service_id "${SERVICE_ID}" '{
  full_name: "Lead Smoke Test",
  phone: "+573009998888",
  service_id: $service_id,
  status: "Nuevo",
  source: "whatsapp",
  notes: "Lead creado por e2e_test.sh para validar el flujo local."
}')
LEAD_RESPONSE=$(api_post "/api/leads" "${LEAD_PAYLOAD}")
echo "${LEAD_RESPONSE}" | jq .
LEAD_ID=$(jq -r '.id // empty' <<<"${LEAD_RESPONSE}")
if [[ -z "${LEAD_ID}" ]]; then
  echo "Lead creation did not return an id." >&2
  exit 1
fi

echo -e "\n4. GET /api/leads"
api_get "/api/leads" | jq .

echo -e "\n5. GET /api/dashboard/summary"
api_get "/api/dashboard/summary" | jq .

echo -e "\nSmoke test completed. Created lead: ${LEAD_ID}"
