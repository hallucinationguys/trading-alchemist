.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down seed sqlc-generate docker-build docker-dev docker-prod

# Variables
BINARY_NAME=trading-alchemist
DOCKER_COMPOSE_FILE=docker/compose/docker-compose.yml
DOCKER_COMPOSE_DEV_FILE=docker/compose/docker-compose.override.yml
DOCKERFILE=docker/app/Dockerfile
MIGRATION_PATH=internal/infrastructure/database/migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5433/trading_alchemist_db?sslmode=disable

# === SIMPLIFIED DEVELOPMENT WORKFLOW ===
setup: config-copy-dev docker-up migrate-up sqlc-generate swagger-generate ## Complete development environment setup (run once)
	@echo ""
	@echo "🎉 Development environment setup complete!"
	@echo ""
	@echo "To start developing:"
	@echo "  make dev    # Start backend application via Docker"
	@echo ""

dev: ## Start backend application via Docker (daily development)
	@echo "🚀 Starting backend application via Docker..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) up

stop: ## Stop development environment
	@echo "🛑 Stopping development environment..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) down

restart: stop dev ## Restart development environment

# === END SIMPLIFIED WORKFLOW ===

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Quick Start (for new developers):'
	@echo '  devbox shell  # Enter development environment (tools pre-installed)'
	@echo '  make setup    # One-time setup (configs, docker, migrations)'
	@echo '  make dev      # Start backend application via Docker'
	@echo ''
	@echo 'Daily Development:'
	@echo '  make dev      # Start backend + database via Docker'
	@echo '  make stop     # Stop all Docker services'
	@echo ''
	@echo 'All Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Local development commands
build: ## Build the application
	GOFLAGS='-mod=mod' go build -o bin/$(BINARY_NAME) cmd/api/main.go



test: ## Run tests
	GOFLAGS='-mod=mod' APP_ENV=test go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/ tmp/

clean-air: ## Clean Air temporary files
	rm -rf tmp/ build-errors.log

# Docker commands
docker-build: ## Build Docker image
	docker build -f $(DOCKERFILE) -t $(BINARY_NAME):latest .

# Docker Compose commands
docker-up: ## Start all services with Docker Compose
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down: ## Stop all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-logs: ## View Docker logs
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-dev: ## Start development environment with hot-reload
	docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) up 

docker-dev-build: ## Build and start development environment
	docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) up --build

docker-dev-logs: ## Follow logs for development environment
	docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) logs -f

docker-dev-shell: ## Shell into development container
	docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) exec app sh

docker-dev-stop: ## Stop development environment
	docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) down



docker-restart: ## Restart all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) restart

docker-clean: ## Clean Docker resources
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans
	docker system prune -f

# Database commands
migrate-up: ## Run database migrations up
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down

migrate-create: ## Create new migration file (usage: make migrate-create name=create_table_name)
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

# Database commands
seed: ## Seed the database with initial data
	go run cmd/seeder/main.go

# Database commands for Docker
docker-migrate-up: ## Run migrations in Docker environment
	docker-compose exec app migrate -path $(MIGRATION_PATH) -database "postgres://postgres:postgres@postgres:5432/trading_alchemist_db?sslmode=disable" up

docker-migrate-down: ## Run migrations down in Docker environment
	docker-compose exec app migrate -path $(MIGRATION_PATH) -database "postgres://postgres:postgres@postgres:5432/trading_alchemist_db?sslmode=disable" down

# Code generation
sqlc-generate: ## Generate SQLC code
	sqlc generate

swagger-generate: ## Generate Swagger documentation
	GOFLAGS='-mod=mod' swag init -g docs/swagger.go -o docs --parseInternal --parseDependency

# Configuration management
config-copy-dev: ## Copy development config from example
	@if [ ! -f configs/env.dev ]; then \
		cp configs/env.dev.example configs/env.dev; \
		echo "Created configs/env.dev from example"; \
	else \
		echo "configs/env.dev already exists"; \
	fi

config-copy-test: ## Copy test config from example
	@if [ ! -f configs/env.test ]; then \
		cp configs/env.test.example configs/env.test; \
		echo "Created configs/env.test from example"; \
	else \
		echo "configs/env.test already exists"; \
	fi

config-copy-all: config-copy-dev config-copy-test ## Copy all config files from examples

# Development utilities
deps: ## Install dependencies
	go mod tidy
	go mod download

docker-dev-test: ## Test Docker development environment
	@echo "Testing Docker development environment..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) -f $(DOCKER_COMPOSE_DEV_FILE) ps
	@echo "Checking if app is responding..."
	curl -f http://localhost:8080/health || echo "App not yet ready, check logs with 'make docker-dev-logs'"

# Health checks
health: ## Check application health
	curl -f http://localhost:8080/health || exit 1

docker-health: ## Check Docker services health
	docker-compose ps

# Testing in Docker
docker-test: ## Run tests in Docker container
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golang:1.24.3-alpine \
		go test -v ./...

# Monitoring
docker-stats: ## Show Docker container stats
	docker stats $(shell docker-compose ps -q)

.DEFAULT_GOAL := help 