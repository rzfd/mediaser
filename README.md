# MediaShar - Microservices Donation Platform

MediaShar adalah platform donasi modern untuk content creator yang dibangun dengan **microservices architecture** untuk scalability dan maintainability yang optimal.

## 🚀 Microservices Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           API GATEWAY                                      │
│                         (REST API + gRPC Client)                          │
│                              Port: 8080                                    │
└─┬─────────┬─────────────┬───────────────────────────────────────────────────┘
  │         │             │
┌─▼─────┐ ┌─▼─────────┐ ┌─▼─────────────┐
│Donation│ │  Payment  │ │ Notification  │
│Service │ │  Service  │ │   Service     │
│:9091   │ │   :9092   │ │     :9093     │
└─┬─────┘ └─┬─────────┘ └───────────────┘
  │         │
┌─▼─────┐ ┌─▼─────────┐     ┌─────────────┐
│Donation│ │  Payment  │     │ Gateway DB  │
│   DB   │ │    DB     │     │   :5432     │
│ :5433  │ │   :5434   │     └─────────────┘
└───────┘ └───────────┘
```

## ⚡ Quick Start

```bash
# Clone repository
git clone https://github.com/rzfd/mediashar.git
cd mediashar

# Start microservices
./scripts/start-microservices.sh start

# Or using Makefile
make up

# Access application
open http://localhost:8080
```

## 🚀 Features

### Core Features
- 🎁 **Donation Management**: Create, track, and manage donations
- 💳 **Multi-Payment Support**: Midtrans, PayPal, Stripe, QRIS
- 👤 **User Management**: Registration, authentication, profiles
- 📱 **Platform Integration**: YouTube, TikTok content sync
- 🔔 **Real-time Notifications**: Live donation alerts
- 📊 **Analytics Dashboard**: Donation statistics and reports

### Technical Features  
- 🏗️ **Microservices Architecture**: Scalable distributed system
- 🔄 **gRPC Communication**: High-performance inter-service calls
- 🗄️ **Database Per Service**: Isolated databases for each microservice
- 🔒 **JWT Authentication**: Secure API access
- 🐳 **Docker Support**: Containerized deployment
- 📚 **API Documentation**: Swagger/OpenAPI specs
- 🧪 **Health Checks**: Service monitoring and status
- ⚖️ **Load Balancing**: Horizontal scaling capability

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Echo (REST API)
- **Communication**: gRPC
- **Database**: PostgreSQL 15 (per service)
- **Authentication**: JWT
- **ORM**: GORM

### Infrastructure
- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx
- **Database Admin**: pgAdmin
- **API Docs**: Swagger UI
- **Monitoring**: Built-in health checks

### External Services
- **Payment**: Midtrans, PayPal, Stripe
- **Media**: YouTube API, TikTok API
- **QRIS**: Bank payment integration

## 📁 Project Structure

```
mediashar/
├── cmd/
│   ├── api-gateway/            # API Gateway (main entry point)
│   ├── donation-service/       # Donation microservice
│   ├── payment-service/        # Payment microservice
│   └── notification-service/   # Notification microservice
├── internal/
│   ├── handler/               # HTTP handlers
│   ├── service/               # Business logic
│   ├── repository/            # Data access layer
│   ├── models/                # Data models
│   ├── grpc/                  # gRPC servers
│   └── routes/                # Route definitions
├── proto/                     # Protocol buffer definitions
├── pkg/
│   ├── pb/                    # Generated protobuf files
│   └── utils/                 # Utilities
├── configs/                   # Configuration files
├── scripts/                   # Deployment scripts
├── docs/                      # Documentation
├── docker-compose.microservices.yml  # Microservices setup
└── Makefile                   # Build commands
```

## 🔧 Development Setup

### Prerequisites
```bash
# Install Go 1.21+
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Install Docker & Docker Compose
sudo apt install docker.io docker-compose

# Install development tools
sudo apt install make curl jq

# Install protobuf compiler
sudo apt install protobuf-compiler
```

### Local Development

```bash
# Install protobuf tools
make proto-install

# Generate protobuf files
make proto-gen

# Start all microservices
make up

# Or step by step
make db-setup      # Start databases first
make docker-build  # Build services
make up           # Start services
```

## 🐳 Docker Commands

```bash
make up          # Start microservices
make down        # Stop microservices
make logs        # View all logs
make logs-service SERVICE=donation-service  # Specific service
make restart     # Restart all services
make rebuild     # Rebuild and restart
make status      # Show system status
```

## 🌐 Service URLs

| Service | Port | URL | Protocol | Description |
|---------|------|-----|----------|-------------|
| **API Gateway** | 8080 | `http://localhost:8080` | HTTP/REST | Main API endpoint |
| **Donation Service** | 9091 | `localhost:9091` | gRPC | Donation management |
| **Payment Service** | 9092 | `localhost:9092` | gRPC | Payment processing |
| **Notification Service** | 9093 | `localhost:9093` | gRPC | Real-time notifications |
| **Frontend** | 8000 | `http://localhost:8000` | HTTP | Web interface |
| **pgAdmin** | 8082 | `http://localhost:8082` | HTTP | Database admin |
| **Swagger UI** | 8083 | `http://localhost:8083` | HTTP | API documentation |

### Database Ports
| Database | Port | Connection |
|----------|------|------------|
| **Gateway DB** | 5432 | `postgresql://postgres:password@localhost:5432/gateway_db` |
| **Donation DB** | 5433 | `postgresql://postgres:password@localhost:5433/donation_db` |
| **Payment DB** | 5434 | `postgresql://postgres:password@localhost:5434/payment_db` |

## 🧪 Testing

### API Testing
```bash
# Health check
curl http://localhost:8080/health

# Service health
curl http://localhost:8080/services/health

# API endpoints
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'
```

### gRPC Testing
```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Test services
make test-grpc

# Manual testing
grpcurl -plaintext localhost:9091 list
grpcurl -plaintext localhost:9091 pb.DonationService.GetDonationStats
```

### Load Testing
```bash
# Test API Gateway
ab -n 1000 -c 10 http://localhost:8080/health

# Test gRPC service
ghz --insecure --proto proto/donation.proto \
    --call pb.DonationService.GetDonationStats \
    -d '{"streamer_id": 1}' -n 1000 -c 10 localhost:9091
```

## 🏗️ Microservices Details

### 1. **API Gateway** (Port: 8080)
- **Role**: Main entry point for all REST API requests
- **Database**: Gateway DB (User management, Auth, Platform data)
- **Technology**: Echo Framework + gRPC Client
- **Features**:
  - Authentication & Authorization
  - Request routing to microservices
  - User management
  - Platform management
  - Health checks for all services

### 2. **Donation Service** (Port: 9091)
- **Role**: Handle all donation operations
- **Database**: Donation DB (Isolated)
- **Technology**: gRPC Server
- **Features**:
  - Create/Read donations
  - Donation statistics
  - Real-time donation streaming
  - Donation history

### 3. **Payment Service** (Port: 9092)
- **Role**: Payment processing
- **Database**: Payment DB (Isolated)
- **Technology**: gRPC Server
- **Features**:
  - Payment processing (Midtrans, PayPal, Stripe)
  - Payment verification
  - Webhook handling
  - Transaction management

### 4. **Notification Service** (Port: 9093)
- **Role**: Real-time notifications
- **Database**: None (stateless)
- **Technology**: gRPC Server
- **Features**:
  - Real-time notifications
  - Event streaming
  - Push notifications
  - Email notifications

## 🔄 Communication Patterns

### **Synchronous Communication (gRPC)**
```
API Gateway ←→ Donation Service   (gRPC)
API Gateway ←→ Payment Service    (gRPC)
API Gateway ←→ Notification Service (gRPC)
```

### **Database Architecture (Database-per-Service)**
- ✅ **Independent scaling** per service
- ✅ **Data isolation** and security
- ✅ **Technology diversity** support
- ✅ **Fault isolation**

## 📈 Scaling Strategy

### **Horizontal Scaling**
```bash
# Scale specific service
docker-compose -f docker-compose.microservices.yml up -d --scale donation-service=3
docker-compose -f docker-compose.microservices.yml up -d --scale payment-service=2
```

### **Production Considerations**
- **Load Balancing**: Nginx or HAProxy
- **Service Discovery**: Consul, Eureka, or Kubernetes
- **Circuit Breaker**: For fault tolerance
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK Stack
- **Security**: TLS for gRPC, API rate limiting

## 📚 Documentation

- 📖 **[API Documentation](docs/API.md)** - REST API reference
- 🏗️ **[Microservices Guide](docs/MICROSERVICES_ARCHITECTURE.md)** - Detailed microservices documentation  
- 🔧 **[Development Guide](docs/DEVELOPMENT.md)** - Development setup and workflow
- 🐳 **[Deployment Guide](docs/DEPLOYMENT.md)** - Production deployment
- 🔒 **[Security Guide](docs/SECURITY.md)** - Security considerations

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch
3. **Follow** the microservices patterns
4. **Add** tests for new features
5. **Update** documentation
6. **Submit** a pull request

### Development Workflow
```bash
# Create feature branch
git checkout -b feature/new-payment-method

# Make changes and test
make test
make lint

# Generate protobuf if needed
make proto-gen
make rebuild

# Commit and push
git commit -m "feat: add new payment method"
git push origin feature/new-payment-method
```

## 🔧 Environment Variables

```bash
# API Gateway
SERVER_PORT=8080
JWT_SECRET=your-secret-key
DB_HOST=gateway-db

# Microservice URLs
DONATION_SERVICE_URL=donation-service:9091
PAYMENT_SERVICE_URL=payment-service:9092
NOTIFICATION_SERVICE_URL=notification-service:9093

# Database configurations per service
DONATION_DB_HOST=donation-db
PAYMENT_DB_HOST=payment-db
```

## 🚀 Production Deployment

```bash
# Build production images
make production

# Deploy to production
docker-compose -f docker-compose.microservices.yml up -d --scale donation-service=3

# Check deployment
make health-check
make status
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📧 **Email**: support@mediashar.com
- 💬 **Discord**: [MediaShar Community](https://discord.gg/mediashar)
- 🐛 **Issues**: [GitHub Issues](https://github.com/rzfd/mediashar/issues)
- 📖 **Wiki**: [Documentation Wiki](https://github.com/rzfd/mediashar/wiki)

## 🏆 Acknowledgments

- **Go Community** for amazing tools and libraries
- **gRPC Team** for high-performance communication
- **Echo Framework** for excellent REST API support
- **Docker** for containerization made easy
- **PostgreSQL** for reliable database management

---

**🎉 Ready to build the next-generation donation platform with microservices!**

```bash
# Quick start
./scripts/start-microservices.sh start

# Or manual start
make up
```