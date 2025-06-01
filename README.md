# MediaShar - Donation System Backend

🚀 **Modern Go backend** untuk sistem donasi dengan integrasi **Midtrans payment gateway**, **Docker containerization**, dan **comprehensive testing**.

## 📋 Table of Contents

- [Features](#-features)
- [Quick Start](#-quick-start)
- [Development Commands](#-development-commands)
- [Project Structure](#-project-structure)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)

## ✨ Features

### 🏗️ **Architecture**
- **Clean Architecture** dengan separation of concerns
- **Repository Pattern** untuk data access
- **Service Layer** untuk business logic
- **JWT Authentication** untuk security
- **Middleware** untuk cross-cutting concerns

### 💳 **Payment Integration**
- **Midtrans Snap** payment gateway
- **Webhook handling** untuk payment notifications
- **Multiple payment methods** support
- **Transaction status tracking**

### 🐳 **DevOps & Infrastructure**
- **Docker containerization** dengan multi-stage builds
- **Docker Compose** untuk local development
- **Health checks** dan monitoring
- **Automated testing** dengan comprehensive test suite

### 📊 **Database & Storage**
- **PostgreSQL** sebagai primary database
- **GORM** untuk ORM dan migrations
- **Database backup/restore** utilities
- **Connection pooling** dan optimization

## 🚀 Quick Start

### Prerequisites

```bash
# Check requirements
make env-check
```

Required:
- **Go 1.21+**
- **Docker & Docker Compose**
- **Git**

### 1. Clone & Setup

```bash
git clone <repository-url>
cd mediashar
```

### 2. Install Development Tools

```bash
make install-tools
```

### 3. Start Services

```bash
# Start all services with Docker
make up

# Check status
make status
```

### 4. Run Tests

```bash
# Quick integration test
make docker-test
```

## 🛠️ Development Commands

### **Quick Commands**
```bash
make up          # Start all services
make down        # Stop all services  
make logs        # View application logs
make status      # Check service status
make restart     # Restart services
```

### **Development**
```bash
make build       # Build application binary
make dev         # Run in development mode
make watch       # Run with hot reload (requires air)
make test        # Run all tests
make fmt         # Format code
make lint        # Run linter
```

### **Docker Operations**
```bash
make docker-build    # Build Docker containers
make docker-up       # Start Docker services
make docker-down     # Stop Docker services
make docker-clean    # Clean containers & volumes
make docker-logs     # View app logs
make docker-exec     # Access app container
```

### **Testing & Integration**
```bash
make test            # Unit tests
make test-coverage   # Tests with coverage
make test-race       # Race condition tests
make docker-test     # Full integration tests
make midtrans-test   # Test Midtrans integration
make health-check    # Check service health
```

### **Database**
```bash
make db-connect      # Connect to PostgreSQL
make db-backup       # Backup database
make db-restore      # Restore database
make db-migrate-up   # Run migrations up
make db-migrate-down # Run migrations down
```

### **Documentation & Setup**
```bash
make swagger         # Setup Swagger docs
make platform        # Setup platform integration
make docs            # Generate Go docs
```

### **Production**
```bash
make build-linux     # Build for Linux
make production      # Build production image
make deploy-check    # Check deployment readiness
```

### **Cleanup**
```bash
make clean           # Clean build artifacts
make clean-all       # Clean everything
```

## 📁 Project Structure

```
mediashar/
├── cmd/api/                 # Application entry point
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic interfaces
│   │   └── serviceImpl/    # Service implementations
│   ├── repository/         # Data access interfaces  
│   │   └── repositoryImpl/ # Repository implementations
│   ├── models/             # Domain models
│   ├── middleware/         # HTTP middleware
│   └── routes/             # Route definitions
├── pkg/
│   ├── config/             # Configuration management
│   ├── database/           # Database utilities
│   └── utils/              # Common utilities
├── scripts/                # Automation scripts
├── docs/                   # Documentation
├── configs/                # Configuration files
├── Dockerfile              # Container definition
├── docker-compose.yml      # Multi-service setup
└── Makefile               # Task automation
```

## 📚 API Documentation

### **Available Services**
- **API Server**: http://localhost:8080
- **Swagger UI**: http://localhost:8083  
- **PgAdmin**: http://localhost:8082
- **Health Check**: http://localhost:8080/health

### **Key Endpoints**

#### Authentication
```bash
POST /api/auth/register     # User registration
POST /api/auth/login        # User login
```

#### Users
```bash
GET  /api/users/profile     # Get current user profile
GET  /api/users/:id         # Get user by ID
GET  /api/streamers         # List streamers
```

#### Donations
```bash
POST /api/donations         # Create donation
GET  /api/donations/:id     # Get donation
GET  /api/donations         # List donations
```

#### Midtrans Payment
```bash
POST /api/midtrans/payment/:donationId  # Create payment
POST /api/midtrans/webhook              # Payment webhook
GET  /api/midtrans/status/:orderId      # Check status
```

## 🧪 Testing

### **Test Types**

1. **Unit Tests**
   ```bash
   make test
   ```

2. **Integration Tests**
   ```bash
   make docker-test
   ```

3. **Frontend Testing**
   ```bash
   make frontend-test      # Check frontend integration
   make test-ui           # Quick frontend test
   ```

4. **Coverage Report**
   ```bash
   make test-coverage
   # Opens coverage.html
   ```

5. **Race Condition Tests**
   ```bash
   make test-race
   ```

### **Visual Testing dengan Frontend Interface**

#### **Quick Start Testing**
```bash
# Start complete environment
make dev-full

# Open frontend di browser
make frontend-open
# or manually: http://localhost:8000
```

#### **Automated Testing Workflow**
1. **Open Frontend Interface**: http://localhost:8000
2. **Click "Full Flow Test"** atau tekan `Ctrl+Enter`
3. **Watch Automated Process**:
   - ✅ Health check
   - ✅ User registration & login
   - ✅ Donation creation (75,000 IDR)
   - ✅ Payment token generation
4. **Manual Payment Test**: Click "Open Snap Payment"
5. **Test Scenarios**:
   - **Success**: Card `4811 1111 1111 1114`
   - **Pending**: Card `4911 1111 1111 1113`
   - **Failed**: Card `4411 1111 1111 1118`

#### **Frontend Features**
- **Real-time API monitoring** dengan colored logs
- **Session management** dengan JWT token tracking
- **Interactive forms** untuk testing scenarios
- **Payment integration** dengan Midtrans Snap
- **Health status dashboard** untuk service monitoring

### **Manual Testing**

```bash
# Start services
make up

# Run comprehensive tests
make docker-test

# Test specific features
make midtrans-test
make health-check

# Frontend integration test
make frontend-test
```

### **Test Data**

#### **Backend Test Data**
Test script creates:
- **Streamer user**: `streamer@mediashar.com`
- **Donator user**: `donator@mediashar.com`
- **Test donation**: 50,000 IDR
- **Midtrans payment**: Sandbox environment

#### **Frontend Test Accounts**
Quick login credentials:
- **Streamer**: `streamer@test.com` / `password123`
- **Donator**: `donator@test.com` / `password123`

#### **Midtrans Test Cards**
```
✅ SUCCESS: 4811 1111 1111 1114
⏳ PENDING: 4911 1111 1111 1113  
❌ FAILED:  4411 1111 1111 1118
CVV: 123, Exp: 12/25
```

### **Testing Best Practices**

1. **Start with Frontend Automated Test**:
   ```bash
   make frontend
   # Open browser → Click "Full Flow Test"
   ```

2. **Verify Each Component**:
   - API health check
   - Database connectivity
   - Midtrans configuration
   - Authentication flow
   - Payment processing

3. **Test Edge Cases**:
   - Invalid input data
   - Network failures
   - Payment failures
   - Session expiration

4. **Monitor Logs**:
   - Frontend API response log
   - Backend application logs: `make logs`
   - Database connection status

## 🚀 Deployment

### **Environment Setup**

1. **Production Environment**
   ```bash
   cp .env.example .env.production
   # Edit production values
   ```

2. **Build Production Image**
   ```bash
   make production
   ```

3. **Deployment Check**
   ```bash
   make deploy-check
   ```

### **Environment Variables**

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=password
DB_NAME=donation_system

# Server
SERVER_PORT=8080
SERVER_ENV=production

# Authentication
JWT_SECRET=your-secret-key
JWT_TOKEN_EXPIRY_HOURS=24

# Midtrans (Production)
MIDTRANS_MERCHANT_ID=your-merchant-id
MIDTRANS_CLIENT_KEY=your-client-key
MIDTRANS_SERVER_KEY=your-server-key
MIDTRANS_ENVIRONMENT=production
```

### **Production Deployment**

```bash
# Build for production
make build-linux

# Or use Docker
docker build -t mediashar:latest .
docker run -p 8080:8080 --env-file .env.production mediashar:latest
```

## 🤝 Contributing

### **Development Workflow**

1. **Setup Development Environment**
   ```bash
   make install-tools
   make up
   ```

2. **Make Changes**
   ```bash
   # Format code
   make fmt
   
   # Run linter
   make lint
   
   # Run tests
   make test
   ```

3. **Test Integration**
   ```bash
   make docker-test
   ```

4. **Code Quality Checks**
   ```bash
   make check  # Runs fmt, vet, lint
   ```

### **Code Standards**

- **Go formatting**: `gofmt` + `goimports`
- **Linting**: `golangci-lint`
- **Security**: `gosec`
- **Testing**: Minimum 80% coverage
- **Documentation**: Godoc comments

### **Commit Guidelines**

```bash
feat: add Midtrans payment integration
fix: resolve database connection issue  
docs: update API documentation
test: add integration tests for payments
refactor: improve service layer structure
```

## 📞 Support

### **Useful Commands**

```bash
# Show all available commands
make help

# Check environment
make env-check

# View logs
make logs

# Access database
make db-connect

# Backup database
make db-backup
```

### **Troubleshooting**

1. **Services not starting**
   ```bash
   make docker-clean
   make up
   ```

2. **Database issues**
   ```bash
   make db-connect
   # Check database manually
   ```

3. **Build issues**
   ```bash
   make clean
   make deps
   make build
   ```

### **Documentation**

- **API Docs**: http://localhost:8083 (Swagger)
- **Code Docs**: `make docs` (Godoc)
- **Docker Guide**: [docs/DOCKER_TESTING.md](docs/DOCKER_TESTING.md)
- **Midtrans Guide**: [docs/MIDTRANS_INTEGRATION.md](docs/MIDTRANS_INTEGRATION.md)

---

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- **Midtrans** untuk payment gateway integration
- **Go community** untuk excellent tooling
- **Docker** untuk containerization platform