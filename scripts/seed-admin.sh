#!/usr/bin/env bash
set -euo pipefail

GATEWAY_URL="${GATEWAY_URL:-http://localhost:8080}"

curl -sS -X POST "${GATEWAY_URL}/query" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation($input: RegisterInput!) { register(input: $input) { token user { id email role } } }",
    "variables": {
      "input": {
        "email": "admin@eventhub.io",
        "name": "Platform Admin",
        "password": "AdminPass123!"
      }
    }
  }' | jq .

echo "Admin is auto-seeded when SEED_ADMIN=true (default in docker-compose)."
