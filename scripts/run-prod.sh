#!/usr/bin/env bash
# run-docker.sh
# Builds and starts all services (app + PostgreSQL) via Docker Compose.
# Usage: ./scripts/run-docker.sh

set -e

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "==> Copying .env.prod to .env..."
cp .env.prod .env

echo "==> Building and starting all services in Docker..."
docker compose up --build
