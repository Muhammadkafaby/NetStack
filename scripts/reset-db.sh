#!/bin/bash
# =============================================================================
# VisiMon — Reset Local Dev Database
# Deletes the SQLite database used in local development.
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(dirname "${SCRIPT_DIR}")"

DB_FILE="${REPO_DIR}/data/visimon.db"

if [ -f "${DB_FILE}" ]; then
    echo "[reset-db] Removing database at ${DB_FILE}..."
    rm -f "${DB_FILE}"
    echo "[reset-db] Database removed."
else
    echo "[reset-db] No database found at ${DB_FILE} — nothing to reset."
fi

# Also clean Docker volumes if they exist
if docker ps -a --format '{{.Names}}' 2>/dev/null | grep -q 'visimon'; then
    echo "[reset-db] Docker containers found. To also reset Docker data:"
    echo "  docker compose -f ${REPO_DIR}/docker/docker-compose.yml down -v"
fi