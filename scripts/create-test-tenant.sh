#!/bin/bash

# Script to create a test tenant and initialize its schema

set -e

API_URL="${API_URL:-http://localhost:8080}"

echo "Creating test tenant..."

# Create tenant
RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/public/tenants" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Demo Company",
    "slug": "demo",
    "plan_id": null
  }')

echo "Response: $RESPONSE"

# Extract tenant ID
TENANT_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | sed 's/"id":"//')

if [ -z "$TENANT_ID" ]; then
  echo "Error: Failed to create tenant"
  exit 1
fi

echo "Tenant created with ID: $TENANT_ID"

# Initialize tenant schema
echo "Initializing tenant schema..."
curl -X POST "${API_URL}/api/v1/public/tenants/${TENANT_ID}/init"

echo ""
echo "Tenant initialized successfully!"
echo "Tenant ID: $TENANT_ID"
echo "Tenant Slug: demo"
echo ""
echo "You can now register a user:"
echo "curl -X POST ${API_URL}/api/v1/tenant/auth/register \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -H 'X-Tenant-ID: demo' \\"
echo "  -d '{\"email\":\"admin@demo.com\",\"password\":\"password123\",\"role\":\"admin\"}'"

