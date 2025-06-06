#!/bin/bash

# MediaShar Currency & Language System Test Script
# This script tests the currency conversion and translation features

set -e

echo "🧪 MediaShar Currency & Language System Test"
echo "============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="http://localhost:8080"
CURRENCY_API_URL="$API_BASE_URL/api/currency"
LANGUAGE_API_URL="$API_BASE_URL/api/language"

# Test functions
test_currency_conversion() {
    echo -e "\n${BLUE}Testing Currency Conversion...${NC}"
    
    # Test 1: Convert USD to IDR
    echo "Test 1: Converting 100 USD to IDR"
    response=$(curl -s -X POST "$CURRENCY_API_URL/convert" \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 100,
            "from_currency": "USD",
            "to_currency": "IDR"
        }' || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Currency conversion test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Currency conversion successful${NC}"
        echo "Response: $response"
    fi
    
    # Test 2: Get exchange rate
    echo -e "\nTest 2: Getting USD to IDR exchange rate"
    response=$(curl -s "$CURRENCY_API_URL/rate?from=USD&to=IDR" || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Exchange rate test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Exchange rate retrieval successful${NC}"
        echo "Response: $response"
    fi
    
    # Test 3: List supported currencies
    echo -e "\nTest 3: Listing supported currencies"
    response=$(curl -s "$CURRENCY_API_URL/list" || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Currency list test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Currency list retrieval successful${NC}"
        echo "Response: $response"
    fi
}

test_language_translation() {
    echo -e "\n${BLUE}Testing Language Translation...${NC}"
    
    # Test 1: Translate text
    echo "Test 1: Translating 'Hello World' from English to Indonesian"
    response=$(curl -s -X POST "$LANGUAGE_API_URL/translate" \
        -H "Content-Type: application/json" \
        -d '{
            "text": "Hello World",
            "from_language": "en",
            "to_language": "id"
        }' || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Translation test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Translation successful${NC}"
        echo "Response: $response"
    fi
    
    # Test 2: Detect language
    echo -e "\nTest 2: Detecting language of 'Selamat pagi'"
    response=$(curl -s -X POST "$LANGUAGE_API_URL/detect" \
        -H "Content-Type: application/json" \
        -d '{
            "text": "Selamat pagi"
        }' || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Language detection test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Language detection successful${NC}"
        echo "Response: $response"
    fi
    
    # Test 3: List supported languages
    echo -e "\nTest 3: Listing supported languages"
    response=$(curl -s "$LANGUAGE_API_URL/list" || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Language list test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Language list retrieval successful${NC}"
        echo "Response: $response"
    fi
    
    # Test 4: Bulk translation
    echo -e "\nTest 4: Bulk translation"
    response=$(curl -s -X POST "$LANGUAGE_API_URL/bulk-translate" \
        -H "Content-Type: application/json" \
        -d '{
            "texts": ["Hello", "Thank you", "Good morning"],
            "from_language": "en",
            "to_language": "id"
        }' || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Bulk translation test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Bulk translation successful${NC}"
        echo "Response: $response"
    fi
}

test_integration() {
    echo -e "\n${BLUE}Testing Integration Scenarios...${NC}"
    
    # Test 1: Multi-currency donation with translation
    echo "Test 1: Creating donation with currency conversion and message translation"
    
    # First convert currency
    currency_response=$(curl -s -X POST "$CURRENCY_API_URL/convert" \
        -H "Content-Type: application/json" \
        -d '{
            "amount": 50,
            "from_currency": "USD",
            "to_currency": "IDR"
        }' || echo "ERROR")
    
    # Then translate message
    translation_response=$(curl -s -X POST "$LANGUAGE_API_URL/translate" \
        -H "Content-Type: application/json" \
        -d '{
            "text": "Thank you for the great stream!",
            "from_language": "en",
            "to_language": "id"
        }' || echo "ERROR")
    
    if [[ "$currency_response" == *"ERROR"* ]] || [[ "$translation_response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ Integration test failed${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Integration test successful${NC}"
        echo "Currency conversion: $currency_response"
        echo "Message translation: $translation_response"
    fi
}

check_service_health() {
    echo -e "\n${BLUE}Checking Service Health...${NC}"
    
    # Check API Gateway health
    echo "Checking API Gateway health..."
    response=$(curl -s "$API_BASE_URL/health" || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${RED}❌ API Gateway is not healthy${NC}"
        return 1
    else
        echo -e "${GREEN}✅ API Gateway is healthy${NC}"
    fi
    
    # Check services health
    echo "Checking microservices health..."
    response=$(curl -s "$API_BASE_URL/health/services" || echo "ERROR")
    
    if [[ "$response" == *"ERROR"* ]]; then
        echo -e "${YELLOW}⚠️  Some services may not be available${NC}"
    else
        echo -e "${GREEN}✅ All services are healthy${NC}"
        echo "Response: $response"
    fi
}

# Main test execution
main() {
    echo -e "${YELLOW}Starting MediaShar Currency & Language System Tests...${NC}"
    
    # Wait for services to be ready
    echo "Waiting for services to be ready..."
    sleep 5
    
    # Run health checks first
    if ! check_service_health; then
        echo -e "${RED}❌ Services are not ready. Please start the system first.${NC}"
        echo "Run: docker-compose up -d"
        exit 1
    fi
    
    # Run currency tests
    if test_currency_conversion; then
        echo -e "${GREEN}✅ All currency tests passed${NC}"
    else
        echo -e "${RED}❌ Some currency tests failed${NC}"
    fi
    
    # Run language tests
    if test_language_translation; then
        echo -e "${GREEN}✅ All language tests passed${NC}"
    else
        echo -e "${RED}❌ Some language tests failed${NC}"
    fi
    
    # Run integration tests
    if test_integration; then
        echo -e "${GREEN}✅ All integration tests passed${NC}"
    else
        echo -e "${RED}❌ Some integration tests failed${NC}"
    fi
    
    echo -e "\n${GREEN}🎉 Test execution completed!${NC}"
    echo -e "${BLUE}Check the results above for any failures.${NC}"
}

# Help function
show_help() {
    echo "MediaShar Currency & Language System Test Script"
    echo ""
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -c, --currency      Test only currency features"
    echo "  -l, --language      Test only language features"
    echo "  -i, --integration   Test only integration scenarios"
    echo "  -s, --health        Check only service health"
    echo ""
    echo "Examples:"
    echo "  $0                  Run all tests"
    echo "  $0 -c               Test only currency conversion"
    echo "  $0 -l               Test only language translation"
    echo "  $0 -s               Check service health only"
}

# Parse command line arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    -c|--currency)
        check_service_health && test_currency_conversion
        ;;
    -l|--language)
        check_service_health && test_language_translation
        ;;
    -i|--integration)
        check_service_health && test_integration
        ;;
    -s|--health)
        check_service_health
        ;;
    "")
        main
        ;;
    *)
        echo "Unknown option: $1"
        show_help
        exit 1
        ;;
esac 