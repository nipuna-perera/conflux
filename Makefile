# Development automation and build commands
# Provides convenient commands for common development tasks
# Manages Docker containers, database migrations, and testing

.PHONY: help dev build test clean migrate docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Start development environment with hot reload
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

build: ## Build production images
	docker-compose -f docker-compose.prod.yml build

test: ## Run all tests
	cd backend && go test ./...
	cd frontend && npm test

lint: ## Run linters for both backend and frontend
	@echo "Running backend linter..."
	cd backend && golangci-lint run --config .golangci.yml
	@echo "Running frontend linter..."
	cd frontend && npm run lint

clean: ## Clean up containers and volumes
	docker-compose down -v
	docker system prune -f

migrate: ## Run database migrations
	cd backend && go run cmd/migrate/main.go

docker-up: ## Start Docker services
	docker-compose up -d

docker-down: ## Stop Docker services
	docker-compose down

backend-dev: ## Run backend in development mode
	cd backend && go run cmd/server/main.go

frontend-dev: ## Run frontend in development mode
	cd frontend && npm run dev
