#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   ./scripts/restore.sh <backup-file>
#
# Uses .env.docker when present. Supports:
# - postgres mode via psql
# - sqlite mode via file copy

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <backup-file>"
  exit 1
fi

BACKUP_FILE="$1"
if [[ ! -f "$BACKUP_FILE" ]]; then
  echo "backup file not found: $BACKUP_FILE"
  exit 1
fi

if [[ -f ".env.docker" ]]; then
  # shellcheck disable=SC1091
  source ".env.docker"
fi

DB_DRIVER="${DATABASE_DRIVER:-postgres}"

if [[ "$DB_DRIVER" == "sqlite" ]]; then
  DB_PATH="${DATABASE_PATH:-/data/portlyn.db}"
  cp "$BACKUP_FILE" "$DB_PATH"
  echo "Restored sqlite backup to: $DB_PATH"
  exit 0
fi

DB_NAME="${POSTGRES_DB:-portlyn}"
DB_USER="${POSTGRES_USER:-portlyn}"
DB_HOST="${POSTGRES_HOST:-postgres}"
DB_PORT="${POSTGRES_PORT:-5432}"
DB_PASSWORD="${POSTGRES_PASSWORD:-}"

if [[ -n "${DATABASE_URL:-}" ]]; then
  PGPASSWORD="$DB_PASSWORD" psql "$DATABASE_URL" < "$BACKUP_FILE"
else
  PGPASSWORD="$DB_PASSWORD" psql \
    --host "$DB_HOST" \
    --port "$DB_PORT" \
    --username "$DB_USER" \
    --dbname "$DB_NAME" \
    < "$BACKUP_FILE"
fi

echo "Restore completed from: $BACKUP_FILE"
