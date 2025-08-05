# Makefile para o projeto GoLedger Challenge

.PHONY: help test test-unit test-integration test-coverage build run clean deps fmt lint

# Default target
help:
	@echo "Available commands:"
	@echo "  make test           - Run all tests"
	@echo "  make test-unit      - Run only unit tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make deps           - Download dependencies"
	@echo "  make fmt            - Format Go code"
	@echo "  make lint           - Run linting"
	@echo "  make clean          - Clean build artifacts"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format Go code
fmt:
	go fmt ./...

# Run linting (requires golangci-lint to be installed)
lint:
	golangci-lint run

# Run all tests
test:
	go test -v ./...

# Run unit tests only (excluding integration tests)
test-unit:
	go test -v -short ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in a specific package
test-models:
	go test -v ./internal/models/...

test-handlers:
	go test -v ./internal/handlers/...

test-repositories:
	go test -v ./internal/repositories/...

test-usecases:
	go test -v ./internal/usecases/...

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Clean build artifacts
clean:
	rm -f bin/api
	rm -f coverage.out
	rm -f coverage.html

# Database related commands
db-up:
	docker compose up -d postgres

db-down:
	docker compose down

db-reset: db-down db-up
	@echo "Database reset complete"

# Development helpers
dev-setup: deps fmt
	@echo "Development environment ready"

# CI pipeline
ci: deps fmt test-coverage
	@echo "CI pipeline completed"
