#!/bin/bash

# MediaShar Swagger Integration Setup Script
# This script sets up Swagger UI integration with Docker Compose

set -e

echo "🚀 MediaShar Swagger Integration Setup"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Docker is running
check_docker() {
    print_info "Checking Docker..."
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_status "Docker is running"
}

# Check if docker-compose is available
check_docker_compose() {
    print_info "Checking Docker Compose..."
    if ! command -v docker-compose &> /dev/null; then
        print_error "docker-compose is not installed. Please install it and try again."
        exit 1
    fi
    print_status "Docker Compose is available"
}

# Validate swagger.yaml
validate_swagger() {
    print_info "Validating OpenAPI specification..."
    if [ -f "docs/swagger.yaml" ]; then
        # Try to validate with yq if available
        if command -v yq &> /dev/null; then
            if yq eval docs/swagger.yaml > /dev/null 2>&1; then
                print_status "swagger.yaml is valid"
            else
                print_warning "swagger.yaml may have syntax issues"
            fi
        else
            print_info "yq not available, skipping validation"
        fi
    else
        print_error "docs/swagger.yaml not found!"
        exit 1
    fi
}

# Start services
start_services() {
    print_info "Starting MediaShar services..."
    
    # Stop any existing services
    docker-compose down > /dev/null 2>&1 || true
    
    # Start all services
    docker-compose up -d
    
    print_status "Services started successfully"
}

# Wait for services to be ready
wait_for_services() {
    print_info "Waiting for services to be ready..."
    
    # Wait for API server
    print_info "Waiting for API server (localhost:8080)..."
    timeout=60
    while [ $timeout -gt 0 ]; do
        if curl -s http://localhost:8080/api/streamers > /dev/null 2>&1; then
            print_status "API server is ready"
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    if [ $timeout -le 0 ]; then
        print_warning "API server may not be ready yet"
    fi
    
    # Wait for Swagger UI
    print_info "Waiting for Swagger UI (localhost:8081)..."
    timeout=30
    while [ $timeout -gt 0 ]; do
        if curl -s http://localhost:8081 > /dev/null 2>&1; then
            print_status "Swagger UI is ready"
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    if [ $timeout -le 0 ]; then
        print_warning "Swagger UI may not be ready yet"
    fi
}

# Display service URLs
show_urls() {
    echo ""
    echo "🎉 Setup Complete!"
    echo "=================="
    echo ""
    echo "📚 Services Available:"
    echo "  • API Server:    http://localhost:8080"
    echo "  • Swagger UI:    http://localhost:8081"
    echo "  • PgAdmin:       http://localhost:8082"
    echo "  • PostgreSQL:    localhost:5432"
    echo ""
    echo "🔐 Default Credentials:"
    echo "  • PgAdmin: admin@mediashar.com / admin123"
    echo "  • PostgreSQL: postgres / password"
    echo ""
    echo "🧪 Quick Test:"
    echo "  1. Open Swagger UI: http://localhost:8081"
    echo "  2. Try POST /auth/register to create a user"
    echo "  3. Try POST /auth/login to get JWT token"
    echo "  4. Click 'Authorize' and paste the token"
    echo "  5. Test protected endpoints"
    echo ""
    echo "📖 Documentation:"
    echo "  • Swagger Integration: docs/SWAGGER_INTEGRATION.md"
    echo "  • Postman Collection: postman/README.md"
    echo ""
}

# Open browser (optional)
open_browser() {
    if [ "$1" = "--open" ] || [ "$1" = "-o" ]; then
        print_info "Opening Swagger UI in browser..."
        if command -v xdg-open > /dev/null; then
            xdg-open http://localhost:8081
        elif command -v open > /dev/null; then
            open http://localhost:8081
        else
            print_info "Please open http://localhost:8081 in your browser"
        fi
    fi
}

# Main execution
main() {
    echo ""
    check_docker
    check_docker_compose
    validate_swagger
    start_services
    wait_for_services
    show_urls
    open_browser "$1"
    
    echo "🚀 Ready to test your API with Swagger UI!"
    echo ""
}

# Handle script arguments
case "$1" in
    --help|-h)
        echo "MediaShar Swagger Integration Setup"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  -h, --help     Show this help message"
        echo "  -o, --open     Open Swagger UI in browser after setup"
        echo "  --validate     Only validate swagger.yaml"
        echo "  --status       Show service status"
        echo ""
        exit 0
        ;;
    --validate)
        validate_swagger
        exit 0
        ;;
    --status)
        echo "Service Status:"
        docker-compose ps
        exit 0
        ;;
    *)
        main "$1"
        ;;
esac 