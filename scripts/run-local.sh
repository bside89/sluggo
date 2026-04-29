#!/usr/bin/env bash
# run-local.sh
# Starts only the PostgreSQL container, then runs the API locally.
# Usage: ./scripts/run-local.sh

set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "==> Copying .env.local to .env..."
cp .env.local .env

echo "==> Starting PostgreSQL container..."
docker compose up -d postgres

echo "==> Waiting for PostgreSQL to become healthy..."
until docker compose exec -T postgres pg_isready -U "${DB_USER:-sluggo}" -d "${DB_NAME:-sluggo}" > /dev/null 2>&1; do
  sleep 1
done

echo "==> PostgreSQL is ready."
echo "==> Starting API locally..."
go run ./cmd/api
