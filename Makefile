# Makefile for Donation System

.PHONY: help build run test clean docker-up docker-down docker-logs docker-build docker-rebuild

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-up     - Start all services with Docker"
	@echo "  docker-db     - Start only database services"
	@echo "  docker-down   - Stop Docker services"
	@echo "  docker-clean  - Stop and remove volumes"
	@echo "  docker-rebuild- Rebuild and restart all services"
	@echo "  docker-logs   - View all service logs"
	@echo "  docker-shell  - Execute shell in app container"
	@echo "  docker-psql   - Execute psql in postgres container"
	@echo "  swagger-up    - Start Swagger UI"
	@echo "  swagger-restart- Restart Swagger UI"
	@echo "  swagger-logs  - View Swagger UI logs"
	@echo "  swagger-validate- Validate OpenAPI spec"
	@echo "  swagger-open  - Open Swagger UI in browser"
	@echo "  deps          - Download dependencies"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/mediashar cmd/api/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t mediashar:latest .

# Start all services with Docker
docker-up:
	@echo "Starting all services..."
	docker-compose up -d

# Start only database services
docker-db:
	@echo "Starting database services..."
	docker-compose up -d postgres pgadmin adminer

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

# Stop and remove volumes
docker-clean:
	@echo "Cleaning Docker services and volumes..."
	docker-compose down -v

# Rebuild and restart all services
docker-rebuild:
	@echo "Rebuilding and restarting services..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# View all service logs
docker-logs:
	@echo "Viewing all service logs..."
	docker-compose logs -f

# View app logs only
docker-logs-app:
	@echo "Viewing app logs..."
	docker-compose logs -f app

# View database logs only
docker-logs-db:
	@echo "Viewing database logs..."
	docker-compose logs -f postgres

# Execute shell in app container
docker-shell:
	@echo "Opening shell in app container..."
	docker-compose exec app sh

# Execute psql in postgres container
docker-psql:
	@echo "Opening psql in postgres container..."
	docker-compose exec postgres psql -U postgres -d donation_system

# Development setup (database only)
dev-setup: docker-db deps
	@echo "Development environment ready!"
	@echo "Database: PostgreSQL running on localhost:5432"
	@echo "pgAdmin: http://localhost:8082 (admin@mediashar.com / admin123)"
	@echo "Adminer: http://localhost:8081"
	@echo "Run 'make run' to start the application locally"

# Full Docker setup
docker-setup: docker-up
	@echo "Full Docker environment ready!"
	@echo "Application: http://localhost:8080"
	@echo "Swagger UI: http://localhost:8081"
	@echo "Database: PostgreSQL running on localhost:5432"
	@echo "pgAdmin: http://localhost:8082 (admin@mediashar.com / admin123)"

# Swagger-specific commands
swagger-up:
	@echo "Starting Swagger UI..."
	docker-compose up -d swagger-ui

swagger-down:
	@echo "Stopping Swagger UI..."
	docker-compose stop swagger-ui

swagger-restart:
	@echo "Restarting Swagger UI..."
	docker-compose restart swagger-ui

swagger-logs:
	@echo "Viewing Swagger UI logs..."
	docker-compose logs -f swagger-ui

swagger-validate:
	@echo "Validating OpenAPI specification..."
	docker run --rm -v $(PWD)/docs:/docs mikefarah/yq eval docs/swagger.yaml > /dev/null && echo "✅ swagger.yaml is valid" || echo "❌ swagger.yaml has syntax errors"

swagger-open:
	@echo "Opening Swagger UI in browser..."
	@command -v xdg-open >/dev/null 2>&1 && xdg-open http://localhost:8081 || \
	command -v open >/dev/null 2>&1 && open http://localhost:8081 || \
	echo "Please open http://localhost:8081 in your browser" 