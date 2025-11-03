#!/usr/bin/env bash

set -euo pipefail

API_URL=${API_URL:-"http://localhost:8080"}
TOKEN=${TOKEN:-""}
TENANT_ID=${TENANT_ID:-""}
TENANT_SCHEMA=${TENANT_SCHEMA:-"tenant_demo"}
SUPPLIER_COUNT=${SUPPLIER_COUNT:-3}

print_usage() {
  cat <<EOF
Usage: $(basename "$0") [--api-url URL] [--token TOKEN] [--tenant-id UUID] [--tenant-schema SCHEMA] [--count N]

Env vars (defaults):
  API_URL        (default: http://localhost:8080)
  TOKEN          (required for protected endpoints)
  TENANT_ID      (optional; if provided, runs /tenants/{id}/init)
  TENANT_SCHEMA  (default: tenant_demo; used as X-Tenant-ID and psql checks)
  SUPPLIER_COUNT (default: 3)

Examples:
  TENANT_ID=... TOKEN=... ./scripts/init-suppliers.sh
  API_URL=http://localhost:8080 TOKEN=... TENANT_SCHEMA=tenant_acme ./scripts/init-suppliers.sh --count 5
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --api-url)
      API_URL="$2"; shift 2 ;;
    --token)
      TOKEN="$2"; shift 2 ;;
    --tenant-id)
      TENANT_ID="$2"; shift 2 ;;
    --tenant-schema)
      TENANT_SCHEMA="$2"; shift 2 ;;
    --count)
      SUPPLIER_COUNT="$2"; shift 2 ;;
    -h|--help)
      print_usage; exit 0 ;;
    *) echo "Unknown arg: $1"; print_usage; exit 1 ;;
  esac
done

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required" >&2; exit 1
fi

echo "API_URL        = ${API_URL}"
echo "TENANT_SCHEMA  = ${TENANT_SCHEMA}"
echo "TENANT_ID      = ${TENANT_ID:-<not-set>}"
echo "SUPPLIER_COUNT = ${SUPPLIER_COUNT}"

# 1) Aplicar migraciones del tenant si TENANT_ID está disponible
if [[ -n "${TENANT_ID}" ]]; then
  echo "Initializing tenant schema via API for TENANT_ID=${TENANT_ID} ..."
  curl -sS -X POST "${API_URL}/api/v1/public/tenants/${TENANT_ID}/init" \
    -H "Content-Type: application/json" || {
      echo "Failed to init tenant ${TENANT_ID}" >&2; exit 1; 
    }
  echo "Tenant initialized."
else
  echo "TENANT_ID not set; skipping /tenants/{id}/init call. Ensure the tenant schema has migrations applied."
fi

# 2) Crear proveedores de prueba (requiere TOKEN)
if [[ -z "${TOKEN}" ]]; then
  echo "TOKEN is required to create suppliers (protected endpoint). Set TOKEN env or --token." >&2
  exit 1
fi

echo "Creating ${SUPPLIER_COUNT} sample suppliers in schema ${TENANT_SCHEMA} ..."

for i in $(seq 1 "${SUPPLIER_COUNT}"); do
  NAME="Proveedor Prueba ${i}"
  TAX="1310000${i}"
  EMAIL="proveedor${i}@demo.com"
  PHONE="809-555-00${i}"

  curl -sS -X POST "${API_URL}/api/v1/tenant/suppliers" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "X-Tenant-ID: ${TENANT_SCHEMA}" \
    -H "Content-Type: application/json" \
    -d "{\
      \"name\": \"${NAME}\",\
      \"tax_id\": \"${TAX}\",\
      \"email\": \"${EMAIL}\",\
      \"phone\": \"${PHONE}\",\
      \"city\": \"Santo Domingo\",\
      \"country\": \"DO\"\
    }" >/dev/null || {
      echo "Failed to create supplier ${i}" >&2; exit 1;
    }
  echo " - Created: ${NAME}"
done

echo "Done. You can list suppliers with:"
echo "curl -s ${API_URL}/api/v1/tenant/suppliers -H 'Authorization: Bearer ${TOKEN}' -H 'X-Tenant-ID: ${TENANT_SCHEMA}' | jq ."


