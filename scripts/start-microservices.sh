#!/bin/bash

# MediaShar Microservices Startup Script
# This script helps start the microservices architecture for development

set -e

echo "🚀 Starting MediaShar Microservices..."

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if a port is in use
port_in_use() {
    lsof -i :"$1" >/dev/null 2>&1
}

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command_exists docker; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command_exists docker-compose; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

if ! command_exists make; then
    echo "❌ Make is not installed. Please install Make first."
    exit 1
fi

echo "✅ All prerequisites are installed."

# Function to wait for service to be healthy
wait_for_service() {
    local service_name=$1
    local port=$2
    local max_attempts=30
    local attempt=1
    
    echo "⏳ Waiting for $service_name to be healthy on port $port..."
    
    while [ $attempt -le $max_attempts ]; do
        if port_in_use $port; then
            echo "✅ $service_name is running on port $port"
            return 0
        fi
        
        echo "   Attempt $attempt/$max_attempts: $service_name not ready yet..."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    echo "❌ $service_name failed to start after $max_attempts attempts"
    return 1
}

# Function to stop existing services
stop_existing() {
    echo "🛑 Stopping any existing services..."
    
    # Stop services if running
    if [ -f docker-compose.yml ]; then
        docker-compose down >/dev/null 2>&1 || true
    fi
    
    echo "✅ Stopped existing services"
}

# Function to build services
build_services() {
    echo "🔨 Building microservices..."
    
    # Generate protobuf files
    echo "  📝 Generating protobuf files..."
    make proto-gen
    
    # Build Docker images
    echo "  🐳 Building Docker images..."
    make docker-build
    
    echo "✅ Services built successfully"
}

# Function to start databases first
start_databases() {
    echo "🗄️  Starting databases..."
    
    docker-compose up -d gateway-db donation-db payment-db
    
    # Wait for databases to be ready
    echo "⏳ Waiting for databases to be ready..."
    sleep 10
    
    # Check database health
    if docker-compose ps gateway-db | grep -q "healthy"; then
        echo "✅ Gateway database is healthy"
    else
        echo "⚠️  Gateway database might not be fully ready"
    fi
    
    if docker-compose ps donation-db | grep -q "healthy"; then
        echo "✅ Donation database is healthy"
    else
        echo "⚠️  Donation database might not be fully ready"
    fi
    
    if docker-compose ps payment-db | grep -q "healthy"; then
        echo "✅ Payment database is healthy"
    else
        echo "⚠️  Payment database might not be fully ready"
    fi
}

# Function to start microservices
start_microservices() {
    echo "🚀 Starting microservices..."
    
    # Start all services
    docker-compose up -d
    
    # Wait for each service
    wait_for_service "Donation Service" 9091
    wait_for_service "Payment Service" 9092
    wait_for_service "Notification Service" 9093
    wait_for_service "API Gateway" 8080
    
    echo "✅ All microservices are running"
}

# Function to show service status
show_status() {
    echo ""
    echo "📊 Service Status:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    
    docker-compose ps
    
    echo ""
    echo "🌐 Service URLs:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  🌐 API Gateway:        http://localhost:8080"
    echo "  📱 Frontend:           http://localhost:8000"
    echo "  🗄️  pgAdmin:            http://localhost:8082 (admin@mediashar.com / admin123)"
    echo "  📚 Swagger UI:         http://localhost:8083"
    echo ""
    echo "  🔧 gRPC Services:"
    echo "    🎁 Donation Service:   localhost:9091"
    echo "    💳 Payment Service:    localhost:9092"
    echo "    🔔 Notification Service: localhost:9093"
    echo ""
}

# Function to test services
test_services() {
    echo "🧪 Testing services..."
    
    # Test API Gateway health
    if curl -s http://localhost:8080/health >/dev/null; then
        echo "✅ API Gateway health check passed"
    else
        echo "❌ API Gateway health check failed"
    fi
    
    # Test services health
    if curl -s http://localhost:8080/services/health >/dev/null; then
        echo "✅ Services health check passed"
    else
        echo "⚠️  Some services might not be fully ready"
    fi
    
    echo "✅ Service tests completed"
}

# Main execution
main() {
    local action=${1:-"start"}
    
    case $action in
        "start")
            echo "🎯 Starting microservices architecture..."
            stop_existing
            build_services
            start_databases
            sleep 5  # Give databases more time
            start_microservices
            sleep 10 # Give services time to fully start
            test_services
            show_status
            
            echo ""
            echo "🎉 Microservices started successfully!"
            echo "🔗 API Gateway is available at: http://localhost:8080"
            echo ""
            echo "💡 Useful commands:"
            echo "   📊 Check logs:          make logs"
            echo "   🔄 Restart services:   make restart"
            echo "   🛑 Stop services:      make down"
            echo "   🧪 Test gRPC services: make test-grpc"
            echo "   ❤️  Health check:       curl http://localhost:8080/services/health"
            ;;
        "stop")
            echo "🛑 Stopping microservices..."
            make down
            echo "✅ Microservices stopped"
            ;;
        "restart")
            echo "🔄 Restarting microservices..."
            make restart
            echo "✅ Microservices restarted"
            ;;
        "status")
            show_status
            ;;
        "logs")
            echo "📋 Showing logs..."
            make logs
            ;;
        "test")
            test_services
            ;;
        *)
            echo "Usage: $0 {start|stop|restart|status|logs|test}"
            echo ""
            echo "Commands:"
            echo "  start    - Start the microservices architecture"
            echo "  stop     - Stop all microservices"
            echo "  restart  - Restart all microservices"
            echo "  status   - Show service status"
            echo "  logs     - Show service logs"
            echo "  test     - Test service health"
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@" 