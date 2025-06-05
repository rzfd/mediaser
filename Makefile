# MediaShar Project Makefile
# Comprehensive task automation for Go backend with Docker and Midtrans integration

.PHONY: help build test clean dev docker-build docker-up docker-down docker-test docker-logs swagger platform deps lint fmt vet security install-tools migrate-up migrate-down backup restore quick-test midtrans-test production check frontend-serve frontend-serve-python frontend-serve-node frontend-dev frontend-test frontend-open proto-install proto-gen proto-clean

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
MAIN_PATH := cmd/api/main.go
BUILD_DIR := bin
FRONTEND_DIR := frontend
FRONTEND_PORT := 8000
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

##@ Help

help: ## Display available commands
	@echo ""
	@echo "$(GREEN)🚀 MediaShar Development Commands$(NC)"
	@echo "=================================="
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(GREEN)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""

##@ Development

build: ## Build the application binary
	@echo "$(YELLOW)📦 Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build completed: $(BUILD_DIR)/$(APP_NAME)$(NC)"

dev: ## Run application in development mode
	@echo "$(YELLOW)🔧 Starting development server...$(NC)"
	@go run $(MAIN_PATH)

watch: ## Run with file watching (requires air)
	@echo "$(YELLOW)👀 Starting with hot reload...$(NC)"
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "$(RED)❌ Air not installed. Run: make install-tools$(NC)"; \
		exit 1; \
	fi

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
	@echo "$(BLUE)ℹ️  pgAdmin: http://localhost:5050$(NC)"
	@echo "$(BLUE)ℹ️  Swagger: http://localhost:8080/swagger/index.html$(NC)"
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

frontend-docker-build: ## Build frontend Docker image
	@echo "$(YELLOW)🐳 Building frontend Docker image...$(NC)"
	@cd $(FRONTEND_DIR) && docker build -t $(APP_NAME)-frontend:latest .
	@echo "$(GREEN)✅ Frontend Docker image built$(NC)"

frontend-docker-run: ## Run frontend container standalone
	@echo "$(YELLOW)🚀 Running frontend container...$(NC)"
	@docker run -d -p $(FRONTEND_PORT):80 --name $(APP_NAME)-frontend $(APP_NAME)-frontend:latest
	@echo "$(GREEN)✅ Frontend container running at http://localhost:$(FRONTEND_PORT)$(NC)"

frontend-docker-stop: ## Stop frontend container
	@echo "$(YELLOW)🛑 Stopping frontend container...$(NC)"
	@docker stop $(APP_NAME)-frontend 2>/dev/null || true
	@docker rm $(APP_NAME)-frontend 2>/dev/null || true
	@echo "$(GREEN)✅ Frontend container stopped$(NC)"

frontend-docker-logs: ## View frontend container logs
	@echo "$(YELLOW)📋 Frontend container logs:$(NC)"
	@docker logs -f $(APP_NAME)-frontend

frontend-docker-shell: ## Access frontend container shell
	@echo "$(YELLOW)💻 Accessing frontend container...$(NC)"
	@docker exec -it $(APP_NAME)-frontend sh

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

quick-test: ## Run quick integration test
	@echo "$(YELLOW)⚡ Running quick test...$(NC)"
	@if [ -f quick_test.sh ]; then \
		chmod +x quick_test.sh && ./quick_test.sh; \
	else \
		echo "$(RED)❌ quick_test.sh not found$(NC)"; \
	fi

##@ Docker Operations

docker-build: ## Build Docker containers
	@echo "$(YELLOW)🐳 Building Docker containers...$(NC)"
	@$(DOCKER_COMPOSE) build --no-cache

docker-up: ## Start all Docker services
	@echo "$(YELLOW)🚀 Starting Docker services...$(NC)"
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker services
	@echo "$(YELLOW)🛑 Stopping Docker services...$(NC)"
	@$(DOCKER_COMPOSE) down

docker-restart: ## Restart Docker services
	@echo "$(YELLOW)🔄 Restarting Docker services...$(NC)"
	@$(DOCKER_COMPOSE) restart

docker-clean: ## Clean Docker containers and volumes
	@echo "$(YELLOW)🧹 Cleaning Docker containers and volumes...$(NC)"
	@$(DOCKER_COMPOSE) down --volumes --remove-orphans
	@docker system prune -f

docker-logs: ## Show Docker application logs
	@echo "$(YELLOW)📋 Application logs:$(NC)"
	@$(DOCKER_COMPOSE) logs -f app

docker-logs-all: ## Show all Docker services logs
	@echo "$(YELLOW)📋 All services logs:$(NC)"
	@$(DOCKER_COMPOSE) logs -f

docker-ps: ## Show Docker container status
	@echo "$(YELLOW)📊 Container status:$(NC)"
	@$(DOCKER_COMPOSE) ps

docker-exec: ## Execute bash in app container
	@echo "$(YELLOW)💻 Accessing app container...$(NC)"
	@$(DOCKER_COMPOSE) exec app sh

##@ Testing & Integration

docker-test: ## Run comprehensive Docker integration tests
	@echo "$(YELLOW)🧪 Running Docker integration tests...$(NC)"
	@if [ -f scripts/test-docker.sh ]; then \
		chmod +x scripts/test-docker.sh && ./scripts/test-docker.sh; \
	else \
		echo "$(RED)❌ scripts/test-docker.sh not found$(NC)"; \
		exit 1; \
	fi

midtrans-test: ## Test Midtrans integration specifically
	@echo "$(YELLOW)💳 Testing Midtrans integration...$(NC)"
	@curl -s http://localhost:8080/health || echo "$(RED)❌ App not running. Run: make docker-up$(NC)"
	@echo "$(BLUE)ℹ️  Check Midtrans configuration in docker-compose.yml$(NC)"

health-check: ## Check service health
	@echo "$(YELLOW)🏥 Checking service health...$(NC)"
	@curl -s http://localhost:8080/health | jq . || echo "$(RED)❌ Health check failed$(NC)"
	@curl -s http://localhost:8080/ready | jq . || echo "$(RED)❌ Readiness check failed$(NC)"

##@ Database

db-connect: ## Connect to PostgreSQL database
	@echo "$(YELLOW)🗄️  Connecting to database...$(NC)"
	@$(DOCKER_COMPOSE) exec postgres psql -U postgres -d donation_system

db-migrate-up: ## Run database migrations up
	@echo "$(YELLOW)⬆️  Running migrations up...$(NC)"
	@if command -v migrate > /dev/null 2>&1; then \
		migrate -path migrations -database "postgres://postgres:password@localhost:5432/donation_system?sslmode=disable" up; \
	else \
		echo "$(RED)❌ migrate tool not installed. Run: make install-tools$(NC)"; \
	fi

db-migrate-down: ## Run database migrations down
	@echo "$(YELLOW)⬇️  Running migrations down...$(NC)"
	@if command -v migrate > /dev/null 2>&1; then \
		migrate -path migrations -database "postgres://postgres:password@localhost:5432/donation_system?sslmode=disable" down; \
	else \
		echo "$(RED)❌ migrate tool not installed. Run: make install-tools$(NC)"; \
	fi

db-backup: ## Backup database
	@echo "$(YELLOW)💾 Backing up database...$(NC)"
	@mkdir -p backups
	@$(DOCKER_COMPOSE) exec postgres pg_dump -U postgres donation_system > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "$(GREEN)✅ Database backed up to backups/$(NC)"

db-restore: ## Restore database (requires BACKUP_FILE=filename)
	@echo "$(YELLOW)📥 Restoring database...$(NC)"
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "$(RED)❌ Please specify BACKUP_FILE=filename$(NC)"; \
		exit 1; \
	fi
	@$(DOCKER_COMPOSE) exec -T postgres psql -U postgres -d donation_system < $(BACKUP_FILE)
	@echo "$(GREEN)✅ Database restored$(NC)"

##@ Documentation & Setup

swagger: ## Setup Swagger documentation
	@echo "$(YELLOW)📚 Setting up Swagger documentation...$(NC)"
	@if [ -f scripts/setup-swagger.sh ]; then \
		chmod +x scripts/setup-swagger.sh && ./scripts/setup-swagger.sh; \
	else \
		echo "$(RED)❌ scripts/setup-swagger.sh not found$(NC)"; \
	fi

platform: ## Setup platform integration
	@echo "$(YELLOW)🔗 Setting up platform integration...$(NC)"
	@if [ -f scripts/setup-platform-integration.sh ]; then \
		chmod +x scripts/setup-platform-integration.sh && ./scripts/setup-platform-integration.sh; \
	else \
		echo "$(RED)❌ scripts/setup-platform-integration.sh not found$(NC)"; \
	fi

docs: ## Generate documentation
	@echo "$(YELLOW)📖 Generating documentation...$(NC)"
	@if command -v godoc > /dev/null 2>&1; then \
		echo "$(BLUE)ℹ️  Starting godoc server at http://localhost:6060$(NC)"; \
		godoc -http=:6060; \
	else \
		echo "$(RED)❌ godoc not installed. Run: make install-tools$(NC)"; \
	fi

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
	@echo "Installing migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
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
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux $(MAIN_PATH)
	@echo "$(GREEN)✅ Linux build completed$(NC)"

production: ## Build production Docker image
	@echo "$(YELLOW)🏭 Building production image...$(NC)"
	@docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .
	@echo "$(GREEN)✅ Production image built$(NC)"

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

clean: ## Clean build artifacts
	@echo "$(YELLOW)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean -cache
	@echo "$(GREEN)✅ Cleanup completed$(NC)"

clean-all: clean docker-clean ## Clean everything (build + Docker)
	@echo "$(GREEN)✅ Full cleanup completed$(NC)"

##@ Quick Commands

up: docker-up ## Quick start (alias for docker-up)

down: docker-down ## Quick stop (alias for docker-down)

logs: docker-logs ## Quick logs (alias for docker-logs)

restart: docker-restart ## Quick restart (alias for docker-restart)

status: docker-ps health-check ## Show status and health

# Frontend shortcuts
frontend: frontend-serve ## Quick frontend serve (alias for frontend-serve)

web: frontend-dev ## Quick full stack development (alias for frontend-dev)

test-ui: frontend-test ## Quick frontend integration test (alias for frontend-test)

# All-in-one commands
dev-full: ## Start complete development environment with frontend
	@echo "$(YELLOW)🚀 Starting complete development environment...$(NC)"
	@echo "$(BLUE)ℹ️  Step 1: Starting backend services...$(NC)"
	@$(MAKE) up
	@echo "$(BLUE)ℹ️  Step 2: Waiting for services...$(NC)"
	@sleep 8
	@echo "$(BLUE)ℹ️  Step 3: Checking health...$(NC)"
	@$(MAKE) health-check
	@echo "$(BLUE)ℹ️  Step 4: All services ready!$(NC)"
	@echo ""
	@echo "$(GREEN)✅ Development Environment Ready!$(NC)"
	@echo "$(BLUE)🌐 Backend API: http://localhost:8080$(NC)"
	@echo "$(BLUE)🎨 Frontend UI: http://localhost:$(FRONTEND_PORT)$(NC)"
	@echo "$(BLUE)📊 pgAdmin: http://localhost:5050$(NC)"
	@echo "$(BLUE)📚 Swagger: http://localhost:8080/swagger/index.html$(NC)"
	@echo ""
	@echo "$(YELLOW)ℹ️  Run 'make frontend' in another terminal to start frontend$(NC)"
	@echo "$(YELLOW)ℹ️  Or run 'make frontend-open' to open browser$(NC)"

status-full: ## Show complete system status
	@echo "$(YELLOW)📊 Complete System Status$(NC)"
	@echo "=========================="
	@echo ""
	@echo "$(BLUE)🐳 Docker Services:$(NC)"
	@$(MAKE) docker-ps
	@echo ""
	@echo "$(BLUE)🏥 Health Status:$(NC)"
	@$(MAKE) health-check
	@echo ""
	@echo "$(BLUE)🌐 Frontend Status:$(NC)"
	@$(MAKE) frontend-test
	@echo ""
	@echo "$(BLUE)💾 Database Connection:$(NC)"
	@$(DOCKER_COMPOSE) exec postgres pg_isready -U postgres 2>/dev/null && echo "$(GREEN)✅ Database ready$(NC)" || echo "$(RED)❌ Database not ready$(NC)"
	@echo ""
	@echo "$(BLUE)📋 Environment Info:$(NC)"
	@$(MAKE) env-check

# Proto generation
proto-install:
	@echo "📦 Installing protoc and Go plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Proto tools installed"

proto-gen:
	@echo "🔧 Generating protobuf files..."
	@mkdir -p pkg/pb
	protoc --go_out=. --go_opt=module=github.com/rzfd/mediashar \
		   --go-grpc_out=. --go-grpc_opt=module=github.com/rzfd/mediashar \
		   proto/*.proto
	@echo "✅ Protobuf files generated"

proto-clean:
	@echo "🧹 Cleaning generated proto files..."
	@rm -rf pkg/pb/*.pb.go
	@echo "✅ Proto files cleaned" 