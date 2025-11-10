# Message Service Makefile

.PHONY: all build run migrate test test-e2e test-coverage docker-up docker-down docker-prod docker-prod-down clean

# Default target
all: build test

# Build the application
build:
	@echo "Building..."
	@go build -o bin/message-service.exe cmd/api/main.go

# Run the application
run:
	@echo "Starting application..."
	@go run cmd/api/main.go

# Run migrations and seed data
migrate:
	@echo "Running migrations and seeding data..."
	@go run cmd/migrate/main.go -seed

# Test the application (unit tests only)
test:
	@echo "Running unit tests..."
	@go test ./... -v -short

# Run end-to-end tests
test-e2e:
	@echo "Running E2E tests..."
	@go test ./test/e2e/... -v -tags=e2e

# Test with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -cover -short

# Start infrastructure (PostgreSQL + Redis only)
docker-up:
	@echo "Starting infrastructure..."
	@docker-compose up -d
	@echo "Infrastructure started! (PostgreSQL + Redis)"

# Stop infrastructure
docker-down:
	@echo "Stopping infrastructure..."
	@docker-compose down

# Start production (PostgreSQL + Redis + App)
docker-prod:
	@echo "Starting production environment..."
	@docker-compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d --build
	@echo "Production started! API: http://localhost:8080"

# Stop production
docker-prod-down:
	@echo "Stopping production environment..."
	@docker-compose -f docker-compose.yaml -f docker-compose.prod.yaml down

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
