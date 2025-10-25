.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down create-tenant dev-up dev-down dev-logs dev-rebuild

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building API..."
	@go build -o bin/api ./cmd/api
	@echo "Building Worker..."
	@go build -o bin/worker ./cmd/worker

run-api: ## Run the API server
	@echo "Running API server..."
	@go run ./cmd/api/main.go

run-worker: ## Run the worker
	@echo "Running worker..."
	@go run ./cmd/worker/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -rf tmp/

# Comandos para PRODUCCIÓN
docker-up: ## Start Docker containers (PRODUCTION)
	@echo "Starting Docker containers (PRODUCTION)..."
	@docker-compose up -d

docker-down: ## Stop Docker containers (PRODUCTION)
	@echo "Stopping Docker containers (PRODUCTION)..."
	@docker-compose down

docker-logs: ## Show Docker logs (PRODUCTION)
	@docker-compose logs -f

docker-rebuild: ## Rebuild and restart Docker containers (PRODUCTION)
	@echo "Rebuilding Docker containers (PRODUCTION)..."
	@docker-compose up -d --build

# Comandos para DESARROLLO LOCAL con hot-reload
dev-up: ## Start development environment with hot-reload
	@echo "Starting development environment with hot-reload..."
	@docker-compose -f docker-compose.dev.yml up -d

dev-down: ## Stop development environment
	@echo "Stopping development environment..."
	@docker-compose -f docker-compose.dev.yml down

dev-logs: ## Show development logs
	@docker-compose -f docker-compose.dev.yml logs -f

dev-rebuild: ## Rebuild development containers
	@echo "Rebuilding development containers..."
	@docker-compose -f docker-compose.dev.yml up -d --build

dev-restart-api: ## Restart only the API service
	@echo "Restarting API service..."
	@docker-compose -f docker-compose.dev.yml restart api

dev-restart-frontend: ## Restart only the frontend service
	@echo "Restarting frontend service..."
	@docker-compose -f docker-compose.dev.yml restart frontend

migrate-up: ## Run database migrations (public schema)
	@echo "Running migrations..."
	@go run ./cmd/api/main.go # Migrations run on startup

create-tenant: ## Create a test tenant (usage: make create-tenant SLUG=test1)
	@echo "Creating tenant..."
	@curl -X POST http://localhost:8080/api/v1/public/tenants \
		-H "Content-Type: application/json" \
		-d '{"name":"Test Tenant","slug":"$(SLUG)","plan_id":null}'

init-tenant: ## Initialize tenant schema (usage: make init-tenant ID=uuid)
	@echo "Initializing tenant schema..."
	@curl -X POST http://localhost:8080/api/v1/public/tenants/$(ID)/init

migrate-tenants: ## Run migrations on all existing tenants
	@echo "Running migrations on all tenants..."
	@./scripts/migrate-tenants.sh

insert-accounts: ## Insert chart of accounts in tenants with empty tables
	@echo "Inserting chart of accounts..."
	@./scripts/insert-accounts-all-tenants.sh

check-accounts: ## Check accounting tables status in all tenants
	@echo "Checking accounting tables..."
	@./scripts/check-accounting-tables.sh

recreate-demo: ## Recreate demo tenant from scratch
	@echo "Recreating demo tenant..."
	@./scripts/recreate-demo-tenant.sh

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

mod-tidy: ## Tidy go modules
	@go mod tidy

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.DEFAULT_GOAL := help

