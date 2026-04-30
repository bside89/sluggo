#!/usr/bin/env bash
# run-local.sh
# Starts only the PostgreSQL container, then runs the API locally.
# Usage: ./scripts/run-local.sh

set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "==> Copying .env.local to .env..."
cp .env.local .env

echo "==> Starting PostgreSQL and Redis container..."
docker compose up -d postgres redis

echo "==> Waiting for PostgreSQL and Redis to become healthy..."
until docker compose exec -T postgres pg_isready -U "${DB_USER:-sluggo}" -d "${DB_NAME:-sluggo}" > /dev/null 2>&1; do
  sleep 1
done
until docker exec $(compose_cmd ps -q redis) redis-cli ping > /dev/null 2>&1; do
  sleep 1
done

echo "==> PostgreSQL and Redis are ready."
echo "==> Starting API locally..."

go run ./cmd/api
