# gRPC Implementation Guide

## 📋 Overview

MediaShar donation system now supports **gRPC** alongside REST API untuk persiapan migrasi ke **microservices architecture**. Implementasi ini memungkinkan:

- **High-performance** inter-service communication
- **Type-safe** service contracts dengan Protocol Buffers
- **Real-time streaming** untuk donation notifications
- **Future-ready** untuk microservices scaling

## 🏗️ Architecture

### Current Implementation (Hybrid)
```
┌─────────────────┐    ┌─────────────────┐
│   REST API      │    │   gRPC API      │
│   (Port 8080)   │    │   (Port 9090)   │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────┬───────────────┘
                 │
    ┌─────────────────────────┐
    │     Service Layer       │
    │  (Shared Business Logic)│
    └─────────────────────────┘
                 │
    ┌─────────────────────────┐
    │    Repository Layer     │
    │      (Database)         │
    └─────────────────────────┘
```

### Future Microservices Architecture
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Donation   │    │  Payment    │    │Notification │
│  Service    │    │  Service    │    │  Service    │
│ (gRPC:9091) │    │(gRPC:9092)  │    │(gRPC:9093)  │
└─────────────┘    └─────────────┘    └─────────────┘
       │                   │                   │
       └──────────┬────────────────┬──────────┘
                  │                │
         ┌─────────────────┐       │
         │   API Gateway   │       │
         │  (HTTP & gRPC)  │       │
         └─────────────────┘       │
                  │                │
         ┌─────────────────────────┴─────┐
         │     Message Bus (NATS/Kafka)  │
         └───────────────────────────────┘
```

## 🚀 Services Implemented

### 1. DonationService
**Purpose**: Manages donation operations
**Endpoints**:
- `CreateDonation` - Create new donation
- `GetDonation` - Get donation by ID
- `GetDonationsByStreamer` - Get paginated donations
- `UpdateDonationStatus` - Update payment status
- `StreamDonationEvents` - Real-time event streaming
- `GetDonationStats` - Donation statistics

### 2. PaymentService  
**Purpose**: Handles payment processing
**Endpoints**:
- `ProcessPayment` - Process donation payment
- `VerifyPayment` - Verify payment status
- `HandleWebhook` - Process payment webhooks

### 3. NotificationService
**Purpose**: Real-time notifications and events
**Endpoints**:
- `SendDonationNotification` - Send notifications
- `SubscribeDonationEvents` - Stream subscription

## 📁 File Structure

```
├── proto/
│   └── donation.proto              # Protocol Buffer definitions
├── pkg/pb/                         # Generated Go code (auto)
│   ├── donation.pb.go
│   └── donation_grpc.pb.go
├── internal/grpc/
│   ├── server.go                   # Main gRPC server
│   ├── donation_server.go          # Donation service implementation
│   ├── payment_server.go           # Payment service implementation
│   └── notification_server.go      # Notification service implementation
└── cmd/api/main.go                 # Updated with gRPC support
```

## 🛠️ Setup & Usage

### 1. Install protoc tools
```bash
make proto-install
```

### 2. Generate proto files
```bash
make proto-gen
```

### 3. Start application
```bash
make up
```

**Services Available**:
- **REST API**: http://localhost:8080
- **gRPC API**: localhost:9090

### 4. Test gRPC endpoints

#### Using grpcurl
```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:9090 list

# Create donation
grpcurl -plaintext -d '{
  "amount": 50000,
  "currency": "IDR", 
  "message": "Keep up the good work!",
  "streamer_id": 1,
  "display_name": "Anonymous",
  "is_anonymous": false,
  "payment_method": "qris"
}' localhost:9090 donation.DonationService/CreateDonation

# Get donation
grpcurl -plaintext -d '{"donation_id": 1}' \
  localhost:9090 donation.DonationService/GetDonation

# Subscribe to events (streaming)
grpcurl -plaintext -d '{"user_id": 1}' \
  localhost:9090 donation.NotificationService/SubscribeDonationEvents
```

#### Using Go client
```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "github.com/rzfd/mediashar/pkg/pb"
)

func main() {
    // Connect to gRPC server
    conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    // Create client
    client := pb.NewDonationServiceClient(conn)
    
    // Create donation
    resp, err := client.CreateDonation(context.Background(), &pb.CreateDonationRequest{
        Amount:      50000,
        Currency:    "IDR",
        Message:     "Test donation",
        StreamerId:  1,
        DisplayName: "Test User",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Donation created: %v", resp)
}
```

## 🔄 Migration Strategy

### Phase 1: Hybrid (Current)
- ✅ REST API + gRPC API coexist
- ✅ Shared business logic
- ✅ Single database
- ✅ Easy testing and comparison

### Phase 2: Service Separation
```bash
# Extract donation service
docker run -p 9091:9091 mediashar-donation-service

# Extract payment service  
docker run -p 9092:9092 mediashar-payment-service

# Extract notification service
docker run -p 9093:9093 mediashar-notification-service
```

### Phase 3: Full Microservices
- Separate databases per service
- API Gateway for routing
- Service discovery
- Message bus for async communication

## 📊 Benefits of gRPC

### Performance
- **HTTP/2**: Multiplexing, binary protocol
- **Protobuf**: Compact serialization
- **Streaming**: Real-time bidirectional communication

### Developer Experience
- **Type Safety**: Generated client/server code
- **Documentation**: Self-documenting proto files
- **Tooling**: Rich ecosystem (grpcurl, grpcui, etc.)

### Scalability
- **Load Balancing**: Built-in client-side load balancing
- **Connection Pooling**: Efficient connection reuse
- **Backpressure**: Flow control for streaming

## 🔧 Configuration

### Environment Variables
```bash
# gRPC Server
GRPC_PORT=9090
GRPC_REFLECTION_ENABLED=true

# Service Discovery (Future)
SERVICE_REGISTRY_URL=consul://localhost:8500
GRPC_HEALTH_CHECK_ENABLED=true
```

### Docker Compose (Future)
```yaml
services:
  donation-service:
    image: mediashar-donation-service
    ports:
      - "9091:9091"
    environment:
      - GRPC_PORT=9091
      - DB_URL=postgres://...
      
  payment-service:
    image: mediashar-payment-service  
    ports:
      - "9092:9092"
    depends_on:
      - donation-service
```

## 🧪 Testing

### Unit Tests
```go
func TestDonationGRPCServer_CreateDonation(t *testing.T) {
    // Setup mock service
    mockService := &MockDonationService{}
    server := NewDonationGRPCServer(mockService)
    
    // Test request
    req := &pb.CreateDonationRequest{
        Amount:     50000,
        Currency:   "IDR",
        StreamerId: 1,
    }
    
    resp, err := server.CreateDonation(context.Background(), req)
    assert.NoError(t, err)
    assert.NotNil(t, resp)
}
```

### Integration Tests
```bash
# Start test environment
make docker-test

# Test gRPC endpoints
grpcurl -plaintext localhost:9090 grpc.health.v1.Health/Check
```

## 🚨 When to Use gRPC vs REST

### Use gRPC for:
- ✅ **Internal service communication**
- ✅ **High-performance requirements**
- ✅ **Real-time streaming** (donations, notifications)
- ✅ **Type-safe contracts** between services
- ✅ **Binary data transfer**

### Use REST for:
- ✅ **Public APIs** (Web, Mobile apps)
- ✅ **Simple CRUD operations**
- ✅ **Browser compatibility**
- ✅ **Caching** (HTTP caching)
- ✅ **Third-party integrations**

## 🔮 Future Enhancements

### 1. Service Mesh (Istio)
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: donation-service
spec:
  http:
  - route:
    - destination:
        host: donation-service
        port:
          number: 9091
```

### 2. Distributed Tracing
```go
import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// Add tracing to gRPC server
s := grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
    grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
)
```

### 3. Circuit Breaker
```go
import "github.com/sony/gobreaker"

// Add circuit breaker for service calls
cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "donation-service",
    MaxRequests: 3,
    Timeout:     60 * time.Second,
})
```

## 📈 Monitoring & Metrics

### Prometheus Metrics
```go
import "github.com/grpc-ecosystem/go-grpc-prometheus"

// Add metrics to gRPC server
s := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
    grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
)
grpc_prometheus.Register(s)
```

### Health Checks
```go
import "google.golang.org/grpc/health"

// Register health service
healthServer := health.NewServer()
grpc_health_v1.RegisterHealthServer(s, healthServer)
healthServer.SetServingStatus("donation", grpc_health_v1.HealthCheckResponse_SERVING)
```

## 🎯 Conclusion

gRPC implementation di MediaShar memberikan:

1. **Future-Proofing**: Siap untuk microservices migration
2. **Performance**: High-throughput internal communication  
3. **Type Safety**: Reduced integration errors
4. **Real-time**: Streaming capabilities untuk notifications
5. **Flexibility**: Choice between REST dan gRPC sesuai use case

**Recommendation**: 
- Keep REST API untuk client-facing applications
- Use gRPC untuk internal service communication
- Implement gradually service by service 