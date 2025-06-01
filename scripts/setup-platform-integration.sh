#!/bin/bash

# MediaShar Platform Integration Setup Script
# This script sets up YouTube and TikTok integration

set -e

echo "🎥 MediaShar Platform Integration Setup"
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

# Check if services are running
check_services() {
    print_info "Checking MediaShar services..."
    
    # Check if postgres is running
    if ! docker-compose ps postgres | grep -q "Up"; then
        print_warning "PostgreSQL is not running. Starting services..."
        docker-compose up -d postgres
        sleep 5
    fi
    
    print_status "Services are ready"
}

# Run database migration
run_migration() {
    print_info "Running platform integration migration..."
    
    # Check if migration file exists
    if [ ! -f "migrations/add_platform_tables.sql" ]; then
        print_error "Migration file not found: migrations/add_platform_tables.sql"
        exit 1
    fi
    
    # Run migration
    docker-compose exec -T postgres psql -U postgres -d donation_system -f /docker-entrypoint-initdb.d/../../../migrations/add_platform_tables.sql || {
        # If exec fails, try copying file and running
        docker cp migrations/add_platform_tables.sql mediashar_postgres:/tmp/
        docker-compose exec postgres psql -U postgres -d donation_system -f /tmp/add_platform_tables.sql
    }
    
    print_status "Database migration completed"
}

# Update Swagger UI
update_swagger() {
    print_info "Updating Swagger UI with new endpoints..."
    
    # Restart Swagger UI to pick up new endpoints
    docker-compose restart swagger-ui
    
    # Wait for Swagger UI to be ready
    timeout=30
    while [ $timeout -gt 0 ]; do
        if curl -s http://localhost:8081 > /dev/null 2>&1; then
            print_status "Swagger UI updated successfully"
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    if [ $timeout -le 0 ]; then
        print_warning "Swagger UI may not be ready yet"
    fi
}

# Test URL validation endpoint
test_url_validation() {
    print_info "Testing URL validation endpoint..."
    
    # Wait for API server
    timeout=30
    while [ $timeout -gt 0 ]; do
        if curl -s http://localhost:8080/api/streamers > /dev/null 2>&1; then
            break
        fi
        sleep 2
        timeout=$((timeout-2))
    done
    
    # Test YouTube URL validation
    print_info "Testing YouTube URL validation..."
    youtube_response=$(curl -s -X POST http://localhost:8080/api/content/validate \
        -H "Content-Type: application/json" \
        -d '{"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}' || echo "failed")
    
    if echo "$youtube_response" | grep -q "success"; then
        print_status "YouTube URL validation working"
    else
        print_warning "YouTube URL validation endpoint not yet implemented"
    fi
    
    # Test TikTok URL validation
    print_info "Testing TikTok URL validation..."
    tiktok_response=$(curl -s -X POST http://localhost:8080/api/content/validate \
        -H "Content-Type: application/json" \
        -d '{"url": "https://www.tiktok.com/@username/video/123456"}' || echo "failed")
    
    if echo "$tiktok_response" | grep -q "success"; then
        print_status "TikTok URL validation working"
    else
        print_warning "TikTok URL validation endpoint not yet implemented"
    fi
}

# Show integration examples
show_examples() {
    echo ""
    echo "🎯 Platform Integration Examples"
    echo "==============================="
    echo ""
    
    echo "📺 Supported URL Formats:"
    echo ""
    echo "YouTube:"
    echo "  • Videos: https://www.youtube.com/watch?v=VIDEO_ID"
    echo "  • Shorts: https://youtu.be/VIDEO_ID"
    echo "  • Live: https://www.youtube.com/live/VIDEO_ID"
    echo "  • Channel: https://www.youtube.com/@username"
    echo ""
    echo "TikTok:"
    echo "  • Videos: https://www.tiktok.com/@username/video/VIDEO_ID"
    echo "  • Short: https://vm.tiktok.com/SHORT_CODE"
    echo "  • Live: https://www.tiktok.com/@username/live"
    echo "  • Profile: https://www.tiktok.com/@username"
    echo ""
    
    echo "🧪 Testing Commands:"
    echo ""
    echo "# Test YouTube URL validation"
    echo "curl -X POST http://localhost:8080/api/content/validate \\"
    echo "  -H \"Content-Type: application/json\" \\"
    echo "  -d '{\"url\": \"https://www.youtube.com/watch?v=dQw4w9WgXcQ\"}'"
    echo ""
    echo "# Test TikTok URL validation"
    echo "curl -X POST http://localhost:8080/api/content/validate \\"
    echo "  -H \"Content-Type: application/json\" \\"
    echo "  -d '{\"url\": \"https://www.tiktok.com/@username/video/123456\"}'"
    echo ""
    echo "# Create donation to YouTube content"
    echo "curl -X POST http://localhost:8080/api/donations/to-content \\"
    echo "  -H \"Content-Type: application/json\" \\"
    echo "  -H \"Authorization: Bearer YOUR_JWT_TOKEN\" \\"
    echo "  -d '{"
    echo "    \"amount\": 25.50,"
    echo "    \"currency\": \"USD\","
    echo "    \"message\": \"Great content!\","
    echo "    \"content_url\": \"https://www.youtube.com/watch?v=dQw4w9WgXcQ\","
    echo "    \"display_name\": \"Anonymous Supporter\""
    echo "  }'"
    echo ""
}

# Show service URLs
show_urls() {
    echo ""
    echo "🌐 Service URLs"
    echo "==============="
    echo ""
    echo "📚 Services Available:"
    echo "  • API Server:    http://localhost:8080"
    echo "  • Swagger UI:    http://localhost:8081"
    echo "  • PgAdmin:       http://localhost:8082"
    echo "  • PostgreSQL:    localhost:5432"
    echo ""
    echo "🔗 New Platform Endpoints:"
    echo "  • URL Validation:     POST /api/content/validate"
    echo "  • Connect Platform:   POST /api/platforms/connect"
    echo "  • Content Donation:   POST /api/donations/to-content"
    echo "  • List Platforms:     GET /api/platforms"
    echo ""
    echo "📖 Documentation:"
    echo "  • Platform Integration: docs/SOCIAL_MEDIA_INTEGRATION.md"
    echo "  • Swagger UI: http://localhost:8081"
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
    check_services
    run_migration
    update_swagger
    test_url_validation
    show_examples
    show_urls
    open_browser "$1"
    
    echo "🎉 Platform Integration Setup Complete!"
    echo ""
    echo "Next Steps:"
    echo "1. Open Swagger UI: http://localhost:8081"
    echo "2. Test the new Platform Integration endpoints"
    echo "3. Try Content Management endpoints"
    echo "4. Implement the Go handlers for full functionality"
    echo ""
}

# Handle script arguments
case "$1" in
    --help|-h)
        echo "MediaShar Platform Integration Setup"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  -h, --help     Show this help message"
        echo "  -o, --open     Open Swagger UI in browser after setup"
        echo "  --migrate      Only run database migration"
        echo "  --test         Only run endpoint tests"
        echo "  --examples     Show usage examples"
        echo ""
        exit 0
        ;;
    --migrate)
        check_docker
        check_services
        run_migration
        exit 0
        ;;
    --test)
        test_url_validation
        exit 0
        ;;
    --examples)
        show_examples
        exit 0
        ;;
    *)
        main "$1"
        ;;
esac 