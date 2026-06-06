#!/bin/bash
# =============================================================================
# VisiMon — Development Environment Starter
# Starts the full dev stack: agent, server, and dashboard
#
# Usage: ./scripts/dev.sh [--build] [--clean]
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(dirname "${SCRIPT_DIR}")"

cd "${REPO_DIR}"

# ---- Parse flags ----
BUILD=false
CLEAN=false
for arg in "$@"; do
    case "$arg" in
        --build) BUILD=true ;;
        --clean) CLEAN=true ;;
        *) echo "Unknown option: $arg"; exit 1 ;;
    esac
done

# ---- Clean up ----
if [ "$CLEAN" = true ]; then
    echo "[dev] Cleaning up..."
    docker compose -f docker/docker-compose.yml down -v 2>/dev/null || true
    echo "[dev] Clean complete."
    exit 0
fi

# ---- Check dependencies ----
command -v docker >/dev/null 2>&1 || { echo "[dev] Docker is required but not installed."; exit 1; }
command -v go >/dev/null 2>&1    || echo "[dev] WARNING: Go not found — building Docker images instead."

echo ""
echo "=========================================="
echo "  VisiMon Development Environment"
echo "=========================================="
echo ""

# ---- Build binaries locally (optional) ----
if [ "$BUILD" = true ]; then
    echo "[dev] Building agent binary..."
    cd "${REPO_DIR}/agent"
    go build -ldflags="-s -w" -o "${REPO_DIR}/bin/visimon-agent" ./cmd/visimon-agent/
    echo "[dev] Agent built at bin/visimon-agent"

    if [ -d "${REPO_DIR}/server/cmd" ]; then
        echo "[dev] Building server binary..."
        cd "${REPO_DIR}/server"
        go build -ldflags="-s -w" -o "${REPO_DIR}/bin/visimon-server" ./cmd/visimon-server/ 2>/dev/null || \
            echo "[dev] Server build skipped (code not ready)"
    fi

    cd "${REPO_DIR}"
fi

# ---- Start Docker Compose ----
echo "[dev] Starting Docker Compose environment..."
echo "[dev]  - Server:  http://localhost:8080"
echo "[dev]  - Dashboard: http://localhost:5173"
echo ""

DOCKER_BUILDKIT=1 docker compose -f docker/docker-compose.yml up --build -d

echo ""
echo "[dev] Environment started! Run 'docker compose -f docker/docker-compose.yml logs -f' to follow logs."
echo "[dev] To stop: docker compose -f docker/docker-compose.yml down"
echo "[dev] To clean: ./scripts/dev.sh --clean"