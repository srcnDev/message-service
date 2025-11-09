# Simple Makefile for Go project

# Generate Swagger documentation
swag:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# Build the application
build:
	@echo "Building..."
	@go build -o bin/message-service.exe cmd/api/main.go

# Build migration tool
build-migrate:
	@echo "Building migration tool..."
	@go build -o bin/migrate.exe cmd/migrate/main.go

# Run migrations only
migrate: build-migrate
	@echo "Running migrations..."
	@./bin/migrate.exe

# Run migrations and seed data
migrate-seed: build-migrate
	@echo "Running migrations and seeding..."
	@./bin/migrate.exe -seed

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Reset database and seed
db-reset:
	@echo "Resetting database..."
	@docker-compose down -v
	@docker-compose up -d psql redis
	@sleep 5
	@echo "Running migrations and seeding..."
	@$(MAKE) migrate-seed
	@echo "Database ready with seed data!"

.PHONY: build build-migrate migrate migrate-seed run test clean install swag db-reset
