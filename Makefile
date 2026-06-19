.PHONY: help run build clean test lint swagger db-up db-down db-reset migrate-up migrate-down dev

APP_NAME    := go-gin-api
APP_PORT    := 8080
GO          := go
SWAG        := $(HOME)/go/bin/swag
DOCKER      := docker compose

.DEFAULT_GOAL := help

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

run: swagger ## Run the API server
	$(GO) run main.go

build: ## Build the binary
	$(GO) build -o bin/$(APP_NAME) main.go

clean: ## Remove build artifacts
	rm -rf bin/

test: ## Run tests
	$(GO) test -v -race ./...

test-coverage: ## Run tests with coverage
	$(GO) test -v -race -coverprofile=coverage/coverage.out -covermode=atomic -coverpkg=./... ./...
	$(GO) tool cover -html=coverage/coverage.out -o coverage.html
	$(GO) tool cover -func=coverage/coverage.out

lint: ## Run go vet
	$(GO) vet ./...

swagger: ## Generate swagger docs
	PATH="$(HOME)/go/bin:$$PATH" $(SWAG) init -g cmd/api/main.go --parseDependency --parseInternal

db-up: ## Start PostgreSQL container
	$(DOCKER) up -d

db-down: ## Stop PostgreSQL container
	$(DOCKER) down

db-reset: ## Stop, remove volumes, and restart PostgreSQL
	$(DOCKER) down -v
	$(DOCKER) up -d

migrate-up: ## Run database migrations up
	$(GO) run main.go -migrate=up

migrate-down: ## Rollback all database migrations
	$(GO) run main.go -migrate=down
sonar-scan:
	$(DOCKER) -f docker-compose-scan.yml up -d 
dev: db-up swagger ## Start PostgreSQL, run migrations, and start the API
	@echo "waiting for PostgreSQL to be ready..."
	@sleep 2
	$(GO) run main.go -migrate=up
	$(GO) run main.go
