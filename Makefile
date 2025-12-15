.PHONY: help dev build test lint docker-up docker-down sqlc-generate

help:
	@echo "Available commands:"
	@echo "  make dev             - Run development server"
	@echo "  make build           - Build the application"
	@echo "  make test            - Run tests"
	@echo "  make lint            - Run linter"
	@echo "  make docker-up       - Start docker services"
	@echo "  make docker-down     - Stop docker services"
	@echo "  make sqlc-generate   - Generate sqlc code"

dev:
	@echo "Starting development server..."
	@go run cmd/api/main.go

build:
	@echo "Building application..."
	@go build -o build/api cmd/api/main.go

test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

docker-up:
	@echo "Starting docker services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping docker services..."
	@docker-compose down

sqlc-generate:
	@echo "Generating sqlc code..."
	@sqlc generate
