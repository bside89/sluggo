#!/usr/bin/env bash
# update-swagger.sh
# Regenerates Swagger docs without requiring a globally installed swag binary.
# Usage: ./scripts/update-swagger.sh

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if command -v swag >/dev/null 2>&1; then
  SWAG_CMD=(swag)
elif [[ -x "$(go env GOPATH)/bin/swag" ]]; then
  SWAG_CMD=("$(go env GOPATH)/bin/swag")
else
  SWAG_CMD=(go run github.com/swaggo/swag/cmd/swag@latest)
fi

echo "==> Generating Swagger docs..."
"${SWAG_CMD[@]}" init -g cmd/api/main.go -o docs --parseInternal --parseDependency

echo "==> Swagger docs updated in ./docs"
