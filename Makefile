# =============================================================================
# AI-Powered Research Paper Analyzer — Makefile
# =============================================================================
# Usage:
#   make help        Show this help message
#   make dev         Start all services via Docker Compose
#   make build       Build Docker images for all services
#   make test        Run tests for backend and frontend
#   make clean       Stop containers and remove volumes
# =============================================================================

.PHONY: help dev build stop clean test test-backend test-frontend lint seed migrate logs backend-dev frontend-dev

# Default target
help: ## Show this help message
	@echo ""
	@echo "  AI-Powered Research Paper Analyzer"
	@echo "  ====================================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ── Docker Compose ──────────────────────────────────────────────────────────

dev: ## Start all services in development mode (docker-compose up)
	docker-compose up --build

dev-detached: ## Start all services in background
	docker-compose up --build -d

build: ## Build Docker images without starting
	docker-compose build

stop: ## Stop all running containers
	docker-compose down

clean: ## Stop containers, remove volumes and orphaned containers
	docker-compose down -v --remove-orphans
	@echo "✓ All containers stopped and volumes removed."

logs: ## Tail logs from all services
	docker-compose logs -f

logs-backend: ## Tail logs from backend only
	docker-compose logs -f backend

logs-frontend: ## Tail logs from frontend only
	docker-compose logs -f frontend

# ── Local Development (without Docker) ──────────────────────────────────────

backend-dev: ## Run backend locally with hot-reload (requires Air)
	cd backend && air

frontend-dev: ## Run frontend locally with Vite dev server
	cd frontend && npm run dev

# ── Testing ─────────────────────────────────────────────────────────────────

test: test-backend test-frontend ## Run all tests

test-backend: ## Run backend Go tests with coverage
	cd backend && go test -v -cover ./...

test-frontend: ## Run frontend tests
	cd frontend && npm test -- --watchAll=false

# ── Linting ─────────────────────────────────────────────────────────────────

lint: lint-backend lint-frontend ## Run all linters

lint-backend: ## Lint backend Go code
	cd backend && golangci-lint run ./...

lint-frontend: ## Lint frontend TypeScript code
	cd frontend && npm run lint

# ── Database ────────────────────────────────────────────────────────────────

migrate: ## Run database migrations (via backend auto-migrate)
	cd backend && go run cmd/migrate/main.go

seed: ## Seed the database with sample data
	@echo "Seeding database..."
	docker-compose exec postgres psql -U postgres -d research_paper_analyzer -f /dev/stdin < database/schema.sql
	@echo "✓ Database seeded successfully."

db-shell: ## Open a psql shell to the database
	docker-compose exec postgres psql -U postgres -d research_paper_analyzer

db-reset: ## Reset the database (drop and recreate)
	docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS research_paper_analyzer;"
	docker-compose exec postgres psql -U postgres -c "CREATE DATABASE research_paper_analyzer;"
	@echo "✓ Database reset. Run 'make seed' to re-seed."

# ── Utilities ───────────────────────────────────────────────────────────────

env: ## Copy .env.example to .env (will not overwrite existing)
	@if [ ! -f .env ]; then cp .env.example .env && echo "✓ .env created from .env.example"; else echo "⚠ .env already exists, skipping."; fi

check-deps: ## Verify required tools are installed
	@echo "Checking dependencies..."
	@command -v docker >/dev/null 2>&1        && echo "  ✓ docker"        || echo "  ✗ docker (required)"
	@command -v docker-compose >/dev/null 2>&1 && echo "  ✓ docker-compose" || echo "  ✗ docker-compose (required)"
	@command -v go >/dev/null 2>&1             && echo "  ✓ go"             || echo "  ✗ go (optional for local dev)"
	@command -v node >/dev/null 2>&1           && echo "  ✓ node"           || echo "  ✗ node (optional for local dev)"
	@echo ""
