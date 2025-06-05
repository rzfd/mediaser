# MediaShar Microservices Architecture

## 🏗️ Architecture Overview

MediaShar telah berhasil dikonversi dari **monolithic architecture** menjadi **true microservices architecture** dengan pemisahan yang jelas antar service dan database.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                FRONTEND                                     │
│                            (React/Vue)                                     │
│                           Port: 8000                                       │
└─────────────────────────┬───────────────────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────────────────┐
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
┌─▼─────┐ ┌─▼─────────┐
│Donation│ │  Payment  │
│   DB   │ │    DB     │
│ :5433  │ │   :5434   │
└───────┘ └───────────┘

     ┌─────────────┐
     │ Gateway DB  │
     │   :5432     │
     └─────────────┘
```

## 🚀 Services Overview

### 1. **API Gateway** (Port: 8080)
- **Tanggung jawab**: Entry point untuk semua REST API requests
- **Database**: Gateway DB (User management, Auth, Platform data)
- **Teknologi**: Echo Framework + gRPC Client
- **Features**:
  - Authentication & Authorization
  - Request routing ke microservices
  - User management
  - Platform management
  - Health checks untuk semua services

### 2. **Donation Service** (Port: 9091)
- **Tanggung jawab**: Mengelola semua operasi donation
- **Database**: Donation DB (Isolated)
- **Teknologi**: gRPC Server
- **Features**:
  - Create/Read donations
  - Donation statistics
  - Real-time donation streaming
  - Donation history

### 3. **Payment Service** (Port: 9092)
- **Tanggung jawab**: Mengelola payment processing
- **Database**: Payment DB (Isolated)
- **Teknologi**: gRPC Server
- **Features**:
  - Payment processing (Midtrans, PayPal, Stripe)
  - Payment verification
  - Webhook handling
  - Transaction management

### 4. **Notification Service** (Port: 9093)
- **Tanggung jawab**: Real-time notifications
- **Database**: None (stateless)
- **Teknologi**: gRPC Server
- **Features**:
  - Real-time notifications
  - Event streaming
  - Push notifications
  - Email notifications

## 🗃️ Database Architecture

### **Database Per Service (Database-per-Service Pattern)**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Gateway DB    │    │  Donation DB    │    │   Payment DB    │
│    Port: 5432   │    │   Port: 5433    │    │   Port: 5434    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • users         │    │ • donations     │    │ • donations     │
│ • platforms     │    │ • users (ref)   │    │   (payment_view)│
│ • content       │    │                 │    │ • users (ref)   │
│ • sessions      │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

**Keuntungan Database Isolation:**
- ✅ **Independent scaling** per service
- ✅ **Data isolation** dan security
- ✅ **Technology diversity** (bisa pakai database berbeda)
- ✅ **Fault isolation** (failure di satu DB tidak mempengaruhi yang lain)

## 🔄 Communication Patterns

### **Synchronous Communication (gRPC)**
```
API Gateway ←→ Donation Service   (gRPC)
API Gateway ←→ Payment Service    (gRPC)
API Gateway ←→ Notification Service (gRPC)
```

### **Asynchronous Communication (Future: Event-Driven)**
```
Payment Service → Event Bus → Notification Service
Donation Service → Event Bus → Analytics Service
```

## 🚀 Getting Started

### **1. Prerequisites**
```bash
# Install required tools
sudo apt update
sudo apt install docker.io docker-compose make curl jq

# Install gRPC tools (optional)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### **2. Start Microservices**
```bash
# Easy way - using script
./scripts/start-microservices.sh start

# Manual way - using Makefile
make up

# Or step by step
make proto-gen        # Generate protobuf files
make docker-build     # Build Docker images
make up              # Start all services
```

### **3. Verify Services**
```bash
# Check all services health
curl http://localhost:8080/services/health | jq

# Test individual services
curl http://localhost:8080/health                    # API Gateway
grpcurl -plaintext localhost:9091 list              # Donation Service
grpcurl -plaintext localhost:9092 list              # Payment Service
grpcurl -plaintext localhost:9093 list              # Notification Service
```

## 🛠️ Development Workflow

### **Service Development**
```bash
# Build specific service
make build-donation-service
make build-payment-service
make build-notification-service
make build-api-gateway

# Run service locally (for development)
make run-donation-service
make run-payment-service
make run-notification-service
make run-api-gateway
```

### **Database Operations**
```bash
# Setup databases
make db-setup

# Access databases
docker exec -it mediashar_gateway_db psql -U postgres -d gateway_db
docker exec -it mediashar_donation_db psql -U postgres -d donation_db
docker exec -it mediashar_payment_db psql -U postgres -d payment_db
```

### **Monitoring & Debugging**
```bash
# View logs
make logs                                         # All services
make logs-service SERVICE=donation-service        # Specific service

# Check service status
docker-compose ps

# Restart services
make ms-restart
```

## 🔧 Service URLs & Ports

| Service | Type | Port | URL | Protocol |
|---------|------|------|-----|----------|
| **API Gateway** | REST API | 8080 | `http://localhost:8080` | HTTP/REST |
| **Donation Service** | gRPC | 9091 | `localhost:9091` | gRPC |
| **Payment Service** | gRPC | 9092 | `localhost:9092` | gRPC |
| **Notification Service** | gRPC | 9093 | `localhost:9093` | gRPC |
| **Frontend** | Web App | 8000 | `http://localhost:8000` | HTTP |
| **pgAdmin** | Database UI | 8082 | `http://localhost:8082` | HTTP |
| **Swagger UI** | API Docs | 8083 | `http://localhost:8083` | HTTP |

### **Database Ports**
| Database | Port | Connection |
|----------|------|------------|
| **Gateway DB** | 5432 | `postgresql://postgres:password@localhost:5432/gateway_db` |
| **Donation DB** | 5433 | `postgresql://postgres:password@localhost:5433/donation_db` |
| **Payment DB** | 5434 | `postgresql://postgres:password@localhost:5434/payment_db` |

## 🔒 Security Considerations

### **Authentication Flow**
```
1. Client → API Gateway (JWT Token)
2. API Gateway validates JWT
3. API Gateway → Microservice (Authenticated Request)
4. Microservice processes request
5. Response back through API Gateway
```

### **Service-to-Service Communication**
- **Internal Network**: Services communicate dalam Docker network
- **gRPC Security**: TLS encryption untuk production
- **Service Discovery**: Static configuration (dapat diperkuat dengan Consul/Eureka)

## 📊 Monitoring & Health Checks

### **Health Check Endpoints**
```bash
# API Gateway health
curl http://localhost:8080/health

# All services health status
curl http://localhost:8080/services/health

# Individual service reflection
grpcurl -plaintext localhost:9091 grpc.health.v1.Health/Check
grpcurl -plaintext localhost:9092 grpc.health.v1.Health/Check
grpcurl -plaintext localhost:9093 grpc.health.v1.Health/Check
```

### **Performance Testing**
```bash
# Load testing API Gateway
ab -n 1000 -c 10 http://localhost:8080/health

# gRPC performance testing
ghz --insecure --proto proto/donation.proto \
    --call pb.DonationService.GetDonationStats \
    -d '{"streamer_id": 1}' \
    -n 1000 -c 10 \
    localhost:9091
```

## 🔗 Useful Links

- [gRPC Documentation](https://grpc.io/docs/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [Microservices Patterns](https://microservices.io/patterns/)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [PostgreSQL Docker Hub](https://hub.docker.com/_/postgres)

## 🤝 Contributing

Untuk berkontribusi pada microservices architecture:

1. **Fork repository** ini
2. **Create feature branch** untuk service baru
3. **Follow** service design patterns yang ada
4. **Add tests** untuk service baru
5. **Update documentation** sesuai changes
6. **Submit pull request** dengan deskripsi yang jelas

---

**🎉 Selamat! Anda sekarang memiliki true microservices architecture yang scalable dan maintainable!** 