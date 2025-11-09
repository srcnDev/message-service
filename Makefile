# Simple Makefile for Go project

# Generate Swagger documentation
swag:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# Build the application
build:
	@echo "Building..."
	@go build -o bin/message-service.exe cmd/api/main.go

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

.PHONY: build run test clean install swag
