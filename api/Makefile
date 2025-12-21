# Multi-Tier App Makefile
# ========================

# Build variables
BINARY_NAME=api
BUILD_DIR=./bin
MAIN_PATH=./cmd/main.go

# Database variables
DB_HOST?=localhost
DB_PORT?=5432
DB_USER?=dummyuser
DB_PASSWORD?=dummypass
DB_NAME?=multi_tier_db
DB_SSLMODE?=disable
DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Migration variables
MIGRATIONS_PATH=./cmd/migrations

.PHONY: all build run test clean docker-up docker-down migrate-create migrate-up migrate-down migrate-force tidy lint help

# Default target
all: build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/ /'

## build: Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## run: Run the application
run:
	@echo "Starting the server..."
	go run $(MAIN_PATH)

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	go clean

## tidy: Run go mod tidy
tidy:
	@echo "Tidying modules..."
	go mod tidy

## lint: Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Docker targets
## docker-up: Start database container
docker-up:
	@echo "Starting database container..."
	docker compose up -d

## docker-down: Stop database container
docker-down:
	@echo "Stopping database container..."
	docker compose down

## docker-logs: View database logs
docker-logs:
	docker compose logs -f database

# Migration targets
## migrate-create: Create a new migration (usage: make migrate-create name=migration_name)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@mkdir -p $(MIGRATIONS_PATH)
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

## migrate-up: Apply all pending migrations
migrate-up:
	@echo "Applying migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

## migrate-down: Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

## migrate-down-all: Rollback all migrations
migrate-down-all:
	@echo "Rolling back all migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down -all

## migrate-force: Force migration version (usage: make migrate-force version=1)
migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: version is required. Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	@echo "Forcing migration version to: $(version)"
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)

## migrate-version: Show current migration version
migrate-version:
	@echo "Current migration version:"
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

## deps: Install required dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go get github.com/go-playground/validator/v10
	go get github.com/google/uuid
	@echo "Dependencies installed!"

## install-migrate: Install golang-migrate CLI
install-migrate:
	@echo "Installing golang-migrate..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "golang-migrate installed!"
