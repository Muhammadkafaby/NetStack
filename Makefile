# =============================================================================
# VisiMon — Makefile
# Build, test, and deploy automation
# =============================================================================

SHELL := /bin/bash
.PHONY: all build build-agent build-server test lint clean run run-agent run-server \
        docker-build docker-build-agent docker-build-server dev dev-clean \
        fmt vet check

# ---- Variables ----
BIN_DIR          := bin
AGENT_SRC        := ./agent/cmd/visimon-agent
SERVER_SRC       := ./server/cmd/visimon-server
AGENT_BIN        := $(BIN_DIR)/visimon-agent
SERVER_BIN       := $(BIN_DIR)/visimon-server

# ---- Default ----
all: fmt vet build test

# ---- Build ----
build: build-agent build-server

build-agent:
	@echo "Building agent..."
	@mkdir -p $(BIN_DIR)
	cd agent && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../$(AGENT_BIN) ./cmd/visimon-agent/
	@echo "  -> $(AGENT_BIN)"

build-server:
	@echo "Building server..."
	@mkdir -p $(BIN_DIR)
	@if [ -d server/cmd ]; then \
		cd server && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../$(SERVER_BIN) ./cmd/visimon-server/ 2>/dev/null && \
		echo "  -> $(SERVER_BIN)" || echo "  -> Server source not found, skipping"; \
	else \
		echo "  -> No server/cmd directory, skipping"; \
	fi

# ---- Test ----
test:
	@echo "Running tests..."
	@cd agent && go test ./... -v -count=1 2>&1 | tail -20
	@if [ -d server ]; then cd server && go test ./... -v -count=1 2>&1 | tail -20; fi

# ---- Lint ----
fmt:
	@echo "Formatting code..."
	@cd agent && go fmt ./...
	@if [ -d server ]; then cd server && go fmt ./...; fi

vet:
	@echo "Running go vet..."
	@cd agent && go vet ./...
	@if [ -d server ]; then cd server && go vet ./...; fi

lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd agent && golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, running go vet instead..."; \
		cd agent && go vet ./...; \
	fi

check: fmt vet lint

# ---- Run ----
run: run-agent

run-agent: build-agent
	@echo "Starting VisiMon agent..."
	$(AGENT_BIN) --config ./agent.yaml --server http://localhost:8080 --interval 1

run-server: build-server
	@echo "Starting VisiMon server..."
	$(SERVER_BIN)

# ---- Docker ----
docker-build: docker-build-agent docker-build-server

docker-build-agent:
	@echo "Building Docker image: visimon-agent..."
	docker build -f docker/Dockerfile.agent -t visimon/agent:latest .

docker-build-server:
	@echo "Building Docker image: visimon-server..."
	docker build -f docker/Dockerfile.server -t visimon/server:latest .

docker-build-all: docker-build
	@echo "All Docker images built."

# ---- Dev ----
dev:
	@echo "Starting development environment..."
	@./scripts/dev.sh --build

dev-clean:
	@echo "Cleaning development environment..."
	@./scripts/dev.sh --clean

# ---- Clean ----
clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f visimon.yaml agent.yaml
	docker compose -f docker/docker-compose.yml down -v 2>/dev/null || true

# ---- Release (goreleaser) ----
release:
	@echo "Running goreleaser..."
	@goreleaser release --clean

release-snapshot:
	@echo "Running goreleaser (snapshot)..."
	@goreleaser release --snapshot --clean

# ---- Help ----
help:
	@echo "VisiMon Makefile"
	@echo "================"
	@echo "build          - Build agent and server binaries"
	@echo "build-agent    - Build agent binary only"
	@echo "build-server   - Build server binary only (if server code exists)"
	@echo "test           - Run all tests"
	@echo "lint           - Lint Go code"
	@echo "fmt            - Format Go code"
	@echo "vet            - Run go vet"
	@echo "check          - Run fmt, vet, lint"
	@echo "run            - Run agent (requires running server)"
	@echo "docker-build   - Build all Docker images"
	@echo "dev            - Start full dev environment"
	@echo "dev-clean      - Stop and clean dev environment"
	@echo "clean          - Remove binaries and Docker volumes"
	@echo "release        - Run goreleaser (requires GITHUB_TOKEN)"
	@echo ""