.PHONY: run build test clean docs swagger

# Variables
APP_NAME=fiber-api
BUILD_DIR=bin

# Run the application
run:
	go run cmd/api/main.go

# Build the application
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Generate swagger docs
swagger:
	swag init -g cmd/api/main.go -o docs

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Development mode with hot reload (requires air)
dev:
	air

# Database operations
db-init:
	rm -f users.db

# Help
help:
	@echo "Available commands:"
	@echo "  run     - Run the application"
	@echo "  build   - Build the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean build artifacts"
	@echo "  swagger - Generate swagger docs"
	@echo "  deps    - Install dependencies"
	@echo "  fmt     - Format code"
	@echo "  lint    - Lint code"
	@echo "  dev     - Run in development mode with hot reload"
	@echo "  db-init - Initialize database (removes existing db)"
	@echo "  help    - Show this help message"
