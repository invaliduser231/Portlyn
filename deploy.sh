#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
ENV_FILE="$ROOT_DIR/.env.docker"

if [ ! -f "$ENV_FILE" ]; then
  echo "Missing .env.docker. Copy .env.docker.example to .env.docker and set real values first." >&2
  exit 1
fi

set -a
. "$ENV_FILE"
set +a

if [ -z "${JWT_SECRET:-}" ] || [ -z "${ADMIN_EMAIL:-}" ] || [ -z "${ADMIN_PASSWORD:-}" ]; then
  echo "JWT_SECRET, ADMIN_EMAIL, and ADMIN_PASSWORD must be set in .env.docker" >&2
  exit 1
fi

if [ -z "${GRAFANA_ADMIN_PASSWORD:-}" ]; then
  echo "GRAFANA_ADMIN_PASSWORD must be set in .env.docker" >&2
  exit 1
fi

if [ -z "${DATABASE_URL:-}" ] && [ "${DATABASE_DRIVER:-postgres}" = "postgres" ] && [ -z "${POSTGRES_PASSWORD:-}" ]; then
  echo "POSTGRES_PASSWORD must be set when using the bundled PostgreSQL service" >&2
  exit 1
fi

docker compose --env-file "$ENV_FILE" up -d --build
