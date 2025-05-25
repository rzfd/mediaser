# Docker Setup Guide

Panduan lengkap untuk menjalankan MediaShar Donation System menggunakan Docker.

## 🐳 **Docker Services**

Aplikasi ini menggunakan beberapa service Docker:

| Service | Port | Description |
|---------|------|-------------|
| **app** | 8080 | Go application (MediaShar API) |
| **postgres** | 5432 | PostgreSQL database |
| **pgadmin** | 8082 | pgAdmin web interface |
| **adminer** | 8081 | Adminer database admin |

## 🚀 **Quick Start**

### **1. Full Docker Setup (Recommended)**
Menjalankan semua service termasuk aplikasi:

```bash
# Build dan start semua service
make docker-setup

# Atau manual
docker-compose up -d
```

**Access URLs:**
- **API**: http://localhost:8080
- **pgAdmin**: http://localhost:8082
- **Adminer**: http://localhost:8081

### **2. Development Setup**
Menjalankan hanya database, aplikasi dijalankan secara lokal:

```bash
# Start database services only
make dev-setup

# Run aplikasi secara lokal
make run
```

## 🔧 **Available Commands**

### **Build Commands**
```bash
# Build Docker image
make docker-build

# Build aplikasi Go
make build
```

### **Service Management**
```bash
# Start all services
make docker-up

# Start only database services
make docker-db

# Stop all services
make docker-down

# Stop and remove volumes
make docker-clean

# Rebuild and restart
make docker-rebuild
```

### **Monitoring & Debugging**
```bash
# View all logs
make docker-logs

# View app logs only
make docker-logs-app

# View database logs only
make docker-logs-db

# Execute shell in app container
make docker-shell

# Execute psql in postgres
make docker-psql
```

## 🗄️ **Database Access**

### **pgAdmin (Recommended)**
- **URL**: http://localhost:8082
- **Email**: admin@mediashar.com
- **Password**: admin123

**Pre-configured server:**
- **Name**: MediaShar PostgreSQL
- **Host**: postgres
- **Port**: 5432
- **Database**: donation_system
- **Username**: postgres
- **Password**: password

### **Adminer (Alternative)**
- **URL**: http://localhost:8081
- **System**: PostgreSQL
- **Server**: postgres
- **Username**: postgres
- **Password**: password
- **Database**: donation_system

### **Direct psql Access**
```bash
# Via Docker
make docker-psql

# Via local psql (if installed)
psql -h localhost -U postgres -d donation_system
```

## 📁 **Docker Files Structure**

```
├── Dockerfile              # Multi-stage build untuk Go app
├── docker-compose.yml      # Service orchestration
├── .dockerignore           # Files to exclude from build
├── pgadmin/
│   └── servers.json        # pgAdmin server configuration
└── init.sql               # Database initialization
```

## 🏗️ **Dockerfile Explanation**

### **Multi-stage Build**
```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
# ... build process ...

# Stage 2: Runtime
FROM alpine:latest
# ... minimal runtime image ...
```

**Benefits:**
- ✅ Smaller final image size
- ✅ No build dependencies in production
- ✅ Better security (minimal attack surface)

## 🔧 **Environment Variables**

### **Application Environment**
```bash
# Database connection
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=donation_system

# JWT configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Server configuration
SERVER_PORT=8080
```

### **pgAdmin Environment**
```bash
PGADMIN_DEFAULT_EMAIL=admin@mediashar.com
PGADMIN_DEFAULT_PASSWORD=admin123
PGADMIN_CONFIG_SERVER_MODE=False
PGADMIN_CONFIG_MASTER_PASSWORD_REQUIRED=False
```

## 🔒 **Production Considerations**

### **1. Security**
```bash
# Change default passwords
POSTGRES_PASSWORD=your-secure-password
PGADMIN_DEFAULT_PASSWORD=your-secure-password
JWT_SECRET=your-super-secure-jwt-secret
```

### **2. Volumes**
```yaml
volumes:
  postgres_data:     # Database persistence
  pgadmin_data:      # pgAdmin settings persistence
```

### **3. Networks**
```yaml
networks:
  mediashar_network:  # Isolated network for services
    driver: bridge
```

### **4. Health Checks**
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres"]
  interval: 10s
  timeout: 5s
  retries: 5
```

## 🐛 **Troubleshooting**

### **Common Issues**

#### **1. Port Already in Use**
```bash
# Check what's using the port
lsof -i :8080
lsof -i :5432

# Kill the process or change port in docker-compose.yml
```

#### **2. Database Connection Failed**
```bash
# Check if postgres is running
docker-compose ps

# Check postgres logs
make docker-logs-db

# Restart postgres
docker-compose restart postgres
```

#### **3. App Container Fails to Start**
```bash
# Check app logs
make docker-logs-app

# Rebuild image
make docker-rebuild

# Check if database is ready
docker-compose exec postgres pg_isready -U postgres
```

#### **4. pgAdmin Can't Connect**
```bash
# Ensure postgres is running
docker-compose ps postgres

# Check network connectivity
docker-compose exec pgadmin ping postgres

# Reset pgAdmin data
docker-compose down
docker volume rm mediashar_pgadmin_data
docker-compose up -d
```

### **Reset Everything**
```bash
# Stop all services and remove volumes
make docker-clean

# Remove all images
docker-compose down --rmi all

# Start fresh
make docker-setup
```

## 📊 **Monitoring**

### **Container Stats**
```bash
# View resource usage
docker stats

# View specific container
docker stats mediashar_app
```

### **Logs Management**
```bash
# Follow logs with timestamps
docker-compose logs -f -t

# Limit log lines
docker-compose logs --tail=100

# Filter by service
docker-compose logs postgres
```

## 🚀 **Deployment**

### **Production Deployment**
```bash
# Build for production
docker build -t mediashar:prod .

# Run with production config
docker-compose -f docker-compose.prod.yml up -d
```

### **Docker Registry**
```bash
# Tag for registry
docker tag mediashar:latest your-registry/mediashar:latest

# Push to registry
docker push your-registry/mediashar:latest
```

## 📝 **Development Workflow**

### **1. Code Changes**
```bash
# For local development
make dev-setup
make run

# For Docker development
make docker-rebuild
```

### **2. Database Changes**
```bash
# Reset database
make docker-clean
make docker-db

# Or manually
docker-compose down -v
docker-compose up -d postgres
```

### **3. Testing**
```bash
# Run tests locally
make test

# Run tests in container
docker-compose exec app go test -v ./...
```

Dengan setup Docker ini, Anda dapat dengan mudah menjalankan dan mengembangkan aplikasi MediaShar! 🎯 