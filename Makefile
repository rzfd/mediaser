# MediaShar Microservices Project Makefile
# Comprehensive task automation for microservices architecture

.PHONY: help build test clean dev docker-build docker-up docker-down docker-test docker-logs swagger platform deps lint fmt vet security install-tools migrate-up migrate-down backup restore quick-test health-check frontend-serve frontend-serve-python frontend-serve-node frontend-dev frontend-test frontend-open proto-install proto-gen proto-clean monitoring-up monitoring-down monitoring-logs metrics-test dev-up dev-down health

# Default target
.DEFAULT_GOAL := help

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[1;33m
BLUE   := \033[0;34m
RED    := \033[0;31m
NC     := \033[0m # No Color

# Project variables
APP_NAME := mediashar
DOCKER_COMPOSE := docker-compose
FRONTEND_DIR := frontend
FRONTEND_PORT := 8000
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
DONATION_SERVICE=donation-service
PAYMENT_SERVICE=payment-service
NOTIFICATION_SERVICE=notification-service
API_GATEWAY=api-gateway

##@ Help

help: ## Display available commands
	@echo ""
	@echo "$(GREEN)🚀 MediaShar Microservices Commands$(NC)"
	@echo "====================================="
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(GREEN)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""

##@ Development

build: build-microservices ## Build all microservices

build-microservices: build-donation-service build-payment-service build-notification-service build-api-gateway ## Build all microservices

build-donation-service: ## Build donation microservice
	@echo "$(YELLOW)📦 Building donation service...$(NC)"
	@$(GOBUILD) $(LDFLAGS) -o bin/$(DONATION_SERVICE) ./cmd/donation-service
	@echo "$(GREEN)✅ Donation service built$(NC)"

build-payment-service: ## Build payment microservice
	@echo "$(YELLOW)📦 Building payment service...$(NC)"
	@$(GOBUILD) $(LDFLAGS) -o bin/$(PAYMENT_SERVICE) ./cmd/payment-service
	@echo "$(GREEN)✅ Payment service built$(NC)"

build-notification-service: ## Build notification microservice
	@echo "$(YELLOW)📦 Building notification service...$(NC)"
	@$(GOBUILD) $(LDFLAGS) -o bin/$(NOTIFICATION_SERVICE) ./cmd/notification-service
	@echo "$(GREEN)✅ Notification service built$(NC)"

build-api-gateway: ## Build API gateway
	@echo "$(YELLOW)📦 Building API gateway...$(NC)"
	@$(GOBUILD) $(LDFLAGS) -o bin/$(API_GATEWAY) ./cmd/api-gateway
	@echo "$(GREEN)✅ API gateway built$(NC)"

dev: up ## Start development environment

##@ Local Development

run-donation-service: ## Run donation service locally
	@echo "$(YELLOW)🚀 Starting donation service locally...$(NC)"
	@$(GOBUILD) -o $(DONATION_SERVICE) ./cmd/donation-service
	@./$(DONATION_SERVICE)

run-payment-service: ## Run payment service locally
	@echo "$(YELLOW)🚀 Starting payment service locally...$(NC)"
	@$(GOBUILD) -o $(PAYMENT_SERVICE) ./cmd/payment-service
	@./$(PAYMENT_SERVICE)

run-notification-service: ## Run notification service locally
	@echo "$(YELLOW)🚀 Starting notification service locally...$(NC)"
	@$(GOBUILD) -o $(NOTIFICATION_SERVICE) ./cmd/notification-service
	@./$(NOTIFICATION_SERVICE)

run-api-gateway: ## Run API gateway locally
	@echo "$(YELLOW)🚀 Starting API gateway locally...$(NC)"
	@$(GOBUILD) -o $(API_GATEWAY) ./cmd/api-gateway
	@./$(API_GATEWAY)

##@ Frontend

frontend-serve: frontend-serve-python ## Serve frontend testing interface (default: Python)

frontend-serve-python: ## Serve frontend with Python HTTP server
	@echo "$(YELLOW)🌐 Starting frontend with Python server on port $(FRONTEND_PORT)...$(NC)"
	@if [ -d "$(FRONTEND_DIR)" ]; then \
		echo "$(BLUE)ℹ️  Frontend available at: http://localhost:$(FRONTEND_PORT)$(NC)"; \
		echo "$(BLUE)ℹ️  Press Ctrl+C to stop$(NC)"; \
		cd $(FRONTEND_DIR) && python3 -m http.server $(FRONTEND_PORT); \
	else \
		echo "$(RED)❌ Frontend directory not found: $(FRONTEND_DIR)$(NC)"; \
		exit 1; \
	fi

frontend-serve-node: ## Serve frontend with Node.js serve
	@echo "$(YELLOW)🌐 Starting frontend with Node.js serve on port $(FRONTEND_PORT)...$(NC)"
	@if [ -d "$(FRONTEND_DIR)" ]; then \
		if command -v npx > /dev/null 2>&1; then \
			echo "$(BLUE)ℹ️  Frontend available at: http://localhost:$(FRONTEND_PORT)$(NC)"; \
			echo "$(BLUE)ℹ️  Press Ctrl+C to stop$(NC)"; \
			cd $(FRONTEND_DIR) && npx serve -s . -l $(FRONTEND_PORT); \
		else \
			echo "$(RED)❌ Node.js/npm not available. Use: make frontend-serve-python$(NC)"; \
			exit 1; \
		fi \
	else \
		echo "$(RED)❌ Frontend directory not found: $(FRONTEND_DIR)$(NC)"; \
		exit 1; \
	fi

frontend-dev: ## Start full development environment (backend + frontend)
	@echo "$(YELLOW)🚀 Starting full development environment...$(NC)"
	@echo "$(BLUE)ℹ️  Starting backend services...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(BLUE)ℹ️  Waiting for services to be ready...$(NC)"
	@sleep 5
	@echo "$(BLUE)ℹ️  Backend: http://localhost:8080$(NC)"
	@echo "$(BLUE)ℹ️  Frontend: http://localhost:$(FRONTEND_PORT)$(NC)"
	@echo "$(BLUE)ℹ️  pgAdmin: http://localhost:8082$(NC)"
	@echo "$(BLUE)ℹ️  Swagger: http://localhost:8083$(NC)"
	@echo "$(YELLOW)🌐 Starting frontend server...$(NC)"
	@$(MAKE) frontend-serve

frontend-test: ## Run frontend integration tests
	@echo "$(YELLOW)🧪 Running frontend integration tests...$(NC)"
	@echo "$(BLUE)ℹ️  Make sure backend is running: make up$(NC)"
	@if command -v curl > /dev/null 2>&1; then \
		echo "Testing backend connection..."; \
		curl -s http://localhost:8080/health > /dev/null && echo "$(GREEN)✅ Backend connected$(NC)" || echo "$(RED)❌ Backend not responding$(NC)"; \
		echo "Testing frontend files..."; \
		if [ -f "$(FRONTEND_DIR)/index.html" ] && [ -f "$(FRONTEND_DIR)/script.js" ]; then \
			echo "$(GREEN)✅ Frontend files present$(NC)"; \
		else \
			echo "$(RED)❌ Frontend files missing$(NC)"; \
		fi \
	else \
		echo "$(RED)❌ curl not available for testing$(NC)"; \
	fi

frontend-open: ## Open frontend in default browser
	@echo "$(YELLOW)🌐 Opening frontend in browser...$(NC)"
	@if command -v open > /dev/null 2>&1; then \
		open http://localhost:$(FRONTEND_PORT); \
	elif command -v xdg-open > /dev/null 2>&1; then \
		xdg-open http://localhost:$(FRONTEND_PORT); \
	elif command -v start > /dev/null 2>&1; then \
		start http://localhost:$(FRONTEND_PORT); \
	else \
		echo "$(BLUE)ℹ️  Please open http://localhost:$(FRONTEND_PORT) in your browser$(NC)"; \
	fi

##@ Testing

test: ## Run all tests
	@echo "$(YELLOW)🧪 Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "$(YELLOW)📊 Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report: coverage.html$(NC)"

test-race: ## Run tests with race detection
	@echo "$(YELLOW)🏃 Running tests with race detection...$(NC)"
	@go test -race -v ./...

benchmark: ## Run benchmarks
	@echo "$(YELLOW)⚡ Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

test-grpc: ## Test gRPC services with grpcurl
	@echo "$(YELLOW)🧪 Testing gRPC services...$(NC)"
	@echo "1. Testing Donation Service..."
	@grpcurl -plaintext localhost:9091 list || echo "Donation service not available"
	@echo "2. Testing Payment Service..."
	@grpcurl -plaintext localhost:9092 list || echo "Payment service not available"
	@echo "3. Testing Notification Service..."
	@grpcurl -plaintext localhost:9093 list || echo "Notification service not available"

##@ Docker Operations

docker-build: ## Build Docker containers
	@echo "$(YELLOW)🐳 Building Docker containers...$(NC)"
	@$(DOCKER_COMPOSE) build --no-cache

docker-up: up ## Start all Docker services (alias)

docker-down: down ## Stop Docker services (alias)

docker-restart: restart ## Restart Docker services (alias)

docker-clean: ## Clean Docker containers and volumes
	@echo "$(YELLOW)🧹 Cleaning Docker containers and volumes...$(NC)"
	@$(DOCKER_COMPOSE) down --volumes --remove-orphans
	@docker system prune -f

docker-logs: logs ## Show Docker application logs (alias)

docker-logs-all: logs ## Show all Docker services logs (alias)

docker-ps: ## Show Docker container status
	@echo "$(YELLOW)📊 Container status:$(NC)"
	@$(DOCKER_COMPOSE) ps

docker-exec: ## Execute bash in api-gateway container
	@echo "$(YELLOW)💻 Accessing api-gateway container...$(NC)"
	@$(DOCKER_COMPOSE) exec api-gateway sh

##@ Microservices Operations

up: ## Start microservices containers
	@echo "$(YELLOW)🚀 Starting microservices...$(NC)"
	@$(DOCKER_COMPOSE) up -d

down: ## Stop microservices containers
	@echo "$(YELLOW)🛑 Stopping microservices...$(NC)"
	@$(DOCKER_COMPOSE) down

logs: ## View microservices logs
	@echo "$(YELLOW)📋 Microservices logs:$(NC)"
	@$(DOCKER_COMPOSE) logs -f

logs-service: ## View specific service logs (usage: make logs-service SERVICE=donation-service)
	@echo "$(YELLOW)📋 $(SERVICE) logs:$(NC)"
	@$(DOCKER_COMPOSE) logs -f $(SERVICE)

clean: ## Clean microservices Docker resources
	@echo "$(YELLOW)🧹 Cleaning microservices resources...$(NC)"
	@$(DOCKER_COMPOSE) down -v --rmi all

restart: ## Restart microservices
	@echo "$(YELLOW)🔄 Restarting microservices...$(NC)"
	@$(DOCKER_COMPOSE) restart

rebuild: ## Rebuild and restart microservices
	@echo "$(YELLOW)🔨 Rebuilding and restarting microservices...$(NC)"
	@$(MAKE) down
	@$(MAKE) docker-build
	@$(MAKE) up

##@ Database

db-setup: ## Setup microservices databases
	@echo "$(YELLOW)🗄️  Setting up microservices databases...$(NC)"
	@$(DOCKER_COMPOSE) up -d gateway-db donation-db payment-db
	@echo "$(GREEN)✅ Databases are starting up. Please wait for health checks to pass.$(NC)"

db-connect-gateway: ## Connect to Gateway database
	@echo "$(YELLOW)🗄️  Connecting to Gateway database...$(NC)"
	@$(DOCKER_COMPOSE) exec gateway-db psql -U postgres -d gateway_db

db-connect-donation: ## Connect to Donation database
	@echo "$(YELLOW)🗄️  Connecting to Donation database...$(NC)"
	@$(DOCKER_COMPOSE) exec donation-db psql -U postgres -d donation_db

db-connect-payment: ## Connect to Payment database
	@echo "$(YELLOW)🗄️  Connecting to Payment database...$(NC)"
	@$(DOCKER_COMPOSE) exec payment-db psql -U postgres -d payment_db

##@ Health & Monitoring

health-check: ## Check microservices health
	@echo "$(YELLOW)🏥 Checking microservices health...$(NC)"
	@echo "API Gateway:"
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)❌ API Gateway health check failed$(NC)"
	@echo "Services Health:"
	@curl -s http://localhost:8080/services/health | jq . || echo "$(RED)❌ Services health check failed$(NC)"

status: ## Show complete system status
	@echo "$(YELLOW)📊 Complete System Status$(NC)"
	@echo "=========================="
	@echo ""
	@echo "$(BLUE)🐳 Docker Services:$(NC)"
	@$(MAKE) docker-ps
	@echo ""
	@echo "$(BLUE)🏥 Health Status:$(NC)"
	@$(MAKE) health-check
	@echo ""

##@ Protocol Buffers

proto-install: ## Install protobuf compiler and Go plugins
	@echo "$(YELLOW)📦 Installing protoc and Go plugins...$(NC)"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "$(GREEN)✅ Proto tools installed$(NC)"

proto-gen: ## Generate Go code from proto files
	@echo "$(YELLOW)🔧 Generating protobuf files...$(NC)"
	@mkdir -p pkg/pb
	@protoc --go_out=. --go_opt=module=github.com/rzfd/mediashar \
		   --go-grpc_out=. --go-grpc_opt=module=github.com/rzfd/mediashar \
		   proto/*.proto
	@echo "$(GREEN)✅ Protobuf files generated$(NC)"

proto-clean: ## Clean generated proto files
	@echo "$(YELLOW)🧹 Cleaning generated proto files...$(NC)"
	@rm -rf pkg/pb/*.pb.go
	@echo "$(GREEN)✅ Proto files cleaned$(NC)"

##@ Code Quality

deps: ## Download and tidy dependencies
	@echo "$(YELLOW)📦 Managing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@go mod verify

fmt: ## Format code
	@echo "$(YELLOW)🎨 Formatting code...$(NC)"
	@go fmt ./...
	@if command -v goimports > /dev/null 2>&1; then \
		goimports -w .; \
	fi

lint: ## Run linting
	@echo "$(YELLOW)🔍 Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(RED)❌ golangci-lint not installed. Run: make install-tools$(NC)"; \
	fi

vet: ## Run go vet
	@echo "$(YELLOW)🔬 Running go vet...$(NC)"
	@go vet ./...

security: ## Run security checks
	@echo "$(YELLOW)🔒 Running security checks...$(NC)"
	@if command -v gosec > /dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(RED)❌ gosec not installed. Run: make install-tools$(NC)"; \
	fi

check: fmt vet lint ## Run all code quality checks

##@ Tools & Installation

install-tools: ## Install development tools
	@echo "$(YELLOW)🛠️  Installing development tools...$(NC)"
	@echo "Installing air (hot reload)..."
	@go install github.com/cosmtrek/air@latest
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing gosec..."
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing grpcurl..."
	@go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	@echo "$(GREEN)✅ All tools installed$(NC)"

env-check: ## Check environment setup
	@echo "$(YELLOW)🔧 Checking environment...$(NC)"
	@echo "Go version: $$(go version)"
	@echo "Docker version: $$(docker --version 2>/dev/null || echo 'Not installed')"
	@echo "Docker Compose version: $$(docker-compose --version 2>/dev/null || echo 'Not installed')"
	@echo "Git version: $$(git --version 2>/dev/null || echo 'Not installed')"
	@echo "$(GREEN)✅ Environment check complete$(NC)"

##@ Production & Deployment

build-linux: ## Build for Linux production
	@echo "$(YELLOW)🐧 Building for Linux...$(NC)"
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(DONATION_SERVICE)-linux ./cmd/donation-service
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(PAYMENT_SERVICE)-linux ./cmd/payment-service
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(NOTIFICATION_SERVICE)-linux ./cmd/notification-service
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(API_GATEWAY)-linux ./cmd/api-gateway
	@echo "$(GREEN)✅ Linux builds completed$(NC)"

production: ## Build production Docker images
	@echo "$(YELLOW)🏭 Building production images...$(NC)"
	@$(DOCKER_COMPOSE) build
	@echo "$(GREEN)✅ Production images built$(NC)"

deploy-check: ## Check deployment readiness
	@echo "$(YELLOW)🚀 Checking deployment readiness...$(NC)"
	@echo "Version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@if [ -f .env.production ]; then \
		echo "$(GREEN)✅ Production environment file found$(NC)"; \
	else \
		echo "$(RED)❌ .env.production not found$(NC)"; \
	fi
	@echo "$(BLUE)ℹ️  Remember to update environment variables for production$(NC)"

##@ Cleanup

clean-build: ## Clean build artifacts
	@echo "$(YELLOW)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf bin
	@rm -f coverage.out coverage.html
	@go clean -cache
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

clean-all: clean-build clean ## Clean everything (build + Docker)
	@echo "$(GREEN)✅ Full cleanup completed$(NC)"

##@ Quick Commands

# Shortcuts for common operations
start: up ## Quick start (alias for up)
stop: down ## Quick stop (alias for down)
frontend: frontend-serve ## Quick frontend serve
web: frontend-dev ## Quick full stack development

# All-in-one commands
dev-full: ## Start complete microservices environment with frontend
	@echo "$(YELLOW)🚀 Starting complete microservices environment...$(NC)"
	@echo "$(BLUE)ℹ️  Step 1: Starting microservices...$(NC)"
	@$(MAKE) up
	@echo "$(BLUE)ℹ️  Step 2: Waiting for services...$(NC)"
	@sleep 8
	@echo "$(BLUE)ℹ️  Step 3: Checking health...$(NC)"
	@$(MAKE) health-check
	@echo "$(BLUE)ℹ️  Step 4: All services ready!$(NC)"
	@echo ""
	@echo "$(GREEN)✅ Microservices Environment Ready!$(NC)"
	@echo "$(BLUE)🌐 API Gateway: http://localhost:8080$(NC)"
	@echo "$(BLUE)🎨 Frontend UI: http://localhost:$(FRONTEND_PORT)$(NC)"
	@echo "$(BLUE)📊 pgAdmin: http://localhost:8082$(NC)"
	@echo "$(BLUE)📚 Swagger: http://localhost:8083$(NC)"
	@echo "$(BLUE)🔧 Donation Service: localhost:9091$(NC)"
	@echo "$(BLUE)💳 Payment Service: localhost:9092$(NC)"
	@echo "$(BLUE)🔔 Notification Service: localhost:9093$(NC)"
	@echo ""
	@echo "$(YELLOW)ℹ️  Run 'make frontend' in another terminal to start frontend$(NC)"
	@echo "$(YELLOW)ℹ️  Or run 'make frontend-open' to open browser$(NC)"

# MediaShar Monitoring Stack Management

.PHONY: help build up down logs monitoring-up monitoring-down monitoring-logs clean metrics-test dev-up dev-down health

# Default target
help:
	@echo "MediaShar Monitoring Stack Commands:"
	@echo "  make build           - Build all services"
	@echo "  make up              - Start all services including monitoring"
	@echo "  make down            - Stop all services"
	@echo "  make logs            - Show logs for all services"
	@echo "  make monitoring-up   - Start only monitoring services (Prometheus + Grafana)"
	@echo "  make monitoring-down - Stop only monitoring services"
	@echo "  make monitoring-logs - Show monitoring services logs"
	@echo "  make clean           - Clean up volumes and containers"
	@echo "  make metrics-test    - Test metrics endpoints"

# Build all services
build:
	@echo "🔨 Building MediaShar services..."
	docker-compose build

# Start all services including monitoring
up:
	@echo "🚀 Starting MediaShar with monitoring..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "📊 Prometheus: http://localhost:9090"
	@echo "📈 Grafana: http://localhost:3001 (admin/admin123)"
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🔧 API Gateway: http://localhost:8080"

# Stop all services
down:
	@echo "🛑 Stopping all services..."
	docker-compose down

# Show logs for all services
logs:
	docker-compose logs -f

# Start only monitoring services
monitoring-up:
	@echo "📊 Starting monitoring services..."
	docker-compose up -d prometheus grafana postgres-exporter node-exporter
	@echo "✅ Monitoring services started!"
	@echo "📊 Prometheus: http://localhost:9090"
	@echo "📈 Grafana: http://localhost:3001 (admin/admin123)"

# Stop only monitoring services
monitoring-down:
	@echo "🛑 Stopping monitoring services..."
	docker-compose stop prometheus grafana postgres-exporter node-exporter

# Show monitoring services logs
monitoring-logs:
	docker-compose logs -f prometheus grafana postgres-exporter node-exporter

# Clean up everything
clean:
	@echo "🧹 Cleaning up containers and volumes..."
	docker-compose down -v
	docker system prune -f
	@echo "✅ Cleanup completed!"

# Test metrics endpoints
metrics-test:
	@echo "🧪 Testing metrics endpoints..."
	@echo "API Gateway metrics:"
	@curl -s http://localhost:8080/metrics | head -10 || echo "❌ API Gateway metrics not available"
	@echo "\nPrometheus targets:"
	@curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets[].health' 2>/dev/null || echo "❌ Prometheus not available"

# Development helpers
dev-up:
	@echo "🔧 Starting development environment..."
	docker-compose up -d gateway-db donation-db payment-db prometheus grafana
	@echo "✅ Development databases and monitoring ready!"

dev-down:
	docker-compose stop gateway-db donation-db payment-db prometheus grafana

# Check service health
health:
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:8080/health | jq . || echo "❌ API Gateway not healthy"
	@curl -s http://localhost:9090/-/healthy || echo "❌ Prometheus not healthy"
	@curl -s http://localhost:3001/api/health || echo "❌ Grafana not healthy"

.PHONY: build clean test run deps \
	build-microservices build-donation-service build-payment-service build-notification-service build-api-gateway \
	run-donation-service run-payment-service run-notification-service run-api-gateway \
	docker-build docker-up docker-down docker-logs docker-clean \
	up down logs logs-service clean restart rebuild \
	proto-install proto-gen proto-clean \
	db-setup db-connect-gateway db-connect-donation db-connect-payment \
	health-check status \
	test-grpc benchmark \
	format lint security help 