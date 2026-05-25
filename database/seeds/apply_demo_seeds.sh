#!/usr/bin/env bash
set -euo pipefail

DB_CONTAINER="${DB_CONTAINER:-postgres}"
DB_NAME="${POSTGRES_DB:-clinicflow_db}"
DB_USER="${POSTGRES_USER:-clinicflow}"
SEED_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SEED_DIR}/../.." && pwd)"

apply_seed() {
  local source_file="$1"
  local target_file="/tmp/$(basename "${source_file}")"

  docker compose -f "${ROOT_DIR}/docker-compose.yml" cp "${source_file}" "${DB_CONTAINER}:${target_file}"
  docker compose -f "${ROOT_DIR}/docker-compose.yml" exec -T "${DB_CONTAINER}" \
    psql -U "${DB_USER}" -d "${DB_NAME}" -v ON_ERROR_STOP=1 -f "${target_file}"
}

apply_seed "${SEED_DIR}/20260521000100_demo_core_data.sql"
apply_seed "${SEED_DIR}/20260525000100_demo_ux_data.sql"

printf 'Demo seeds applied: core data + UX data.\n'
