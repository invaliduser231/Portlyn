#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   ./scripts/backup.sh [output-dir]
#
# Uses .env.docker when present. Supports:
# - postgres mode via pg_dump
# - sqlite mode via file copy

OUT_DIR="${1:-./backups}"
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
mkdir -p "$OUT_DIR"

if [[ -f ".env.docker" ]]; then
  # shellcheck disable=SC1091
  source ".env.docker"
fi

DB_DRIVER="${DATABASE_DRIVER:-postgres}"

if [[ "$DB_DRIVER" == "sqlite" ]]; then
  DB_PATH="${DATABASE_PATH:-/data/portlyn.db}"
  cp "$DB_PATH" "$OUT_DIR/portlyn-sqlite-${TIMESTAMP}.db"
  echo "Created sqlite backup: $OUT_DIR/portlyn-sqlite-${TIMESTAMP}.db"
  exit 0
fi

DB_NAME="${POSTGRES_DB:-portlyn}"
DB_USER="${POSTGRES_USER:-portlyn}"
DB_HOST="${POSTGRES_HOST:-postgres}"
DB_PORT="${POSTGRES_PORT:-5432}"
DB_PASSWORD="${POSTGRES_PASSWORD:-}"

if [[ -n "${DATABASE_URL:-}" ]]; then
  PGPASSWORD="$DB_PASSWORD" pg_dump "$DATABASE_URL" > "$OUT_DIR/portlyn-postgres-${TIMESTAMP}.sql"
else
  PGPASSWORD="$DB_PASSWORD" pg_dump \
    --host "$DB_HOST" \
    --port "$DB_PORT" \
    --username "$DB_USER" \
    --dbname "$DB_NAME" \
    > "$OUT_DIR/portlyn-postgres-${TIMESTAMP}.sql"
fi

echo "Created postgres backup: $OUT_DIR/portlyn-postgres-${TIMESTAMP}.sql"
