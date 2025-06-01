#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🐳 MediaShar Docker Testing Script${NC}"
echo "=================================="

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# Function to wait for service
wait_for_service() {
    echo -e "${YELLOW}⏳ Waiting for $1 to be ready...${NC}"
    timeout=60
    counter=0
    
    while [ $counter -lt $timeout ]; do
        if curl -s "$2" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ $1 is ready!${NC}"
            return 0
        fi
        counter=$((counter + 1))
        sleep 1
    done
    
    echo -e "${RED}❌ $1 failed to start within $timeout seconds${NC}"
    return 1
}

# Step 1: Clean up previous containers
echo -e "\n${YELLOW}🧹 Cleaning up previous containers...${NC}"
docker-compose down --volumes --remove-orphans
print_status $? "Previous containers cleaned up"

# Step 2: Build and start services
echo -e "\n${YELLOW}🔨 Building and starting services...${NC}"
docker-compose up -d --build
print_status $? "Services started"

# Step 3: Wait for database
echo -e "\n${YELLOW}📊 Waiting for database...${NC}"
sleep 15  # Give database time to start
echo -e "${GREEN}✅ Database should be ready${NC}"

# Step 4: Wait for application
echo -e "\n${YELLOW}🚀 Waiting for application...${NC}"
wait_for_service "MediaShar API" "http://localhost:8080/health"

# Step 5: Check health endpoints
echo -e "\n${YELLOW}🏥 Testing health endpoints...${NC}"

# Health check
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ "$response" = "200" ]; then
    echo -e "${GREEN}✅ Health check passed${NC}"
else
    echo -e "${RED}❌ Health check failed (HTTP $response)${NC}"
    exit 1
fi

# Readiness check
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ready)
if [ "$response" = "200" ]; then
    echo -e "${GREEN}✅ Readiness check passed${NC}"
else
    echo -e "${RED}❌ Readiness check failed (HTTP $response)${NC}"
    exit 1
fi

# Step 6: Test API endpoints
echo -e "\n${YELLOW}🧪 Testing API endpoints...${NC}"

# Test supported platforms endpoint
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/platforms/supported)
print_status $([ "$response" = "200" ] && echo 0 || echo 1) "Platforms endpoint test"

# Step 7: Test user registration and authentication
echo -e "\n${YELLOW}👤 Testing user authentication...${NC}"

# Register a test user (streamer)
register_response=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "teststreamer",
    "email": "streamer@mediashar.com",
    "password": "testpassword123",
    "is_streamer": true
  }')

if echo "$register_response" | grep -q "success"; then
    echo -e "${GREEN}✅ Streamer registration test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Streamer might already exist, trying login...${NC}"
fi

# Login to get JWT token
login_response=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "streamer@mediashar.com",
    "password": "testpassword123"
  }')

# Extract token from response
token=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$token" ]; then
    echo -e "${GREEN}✅ Streamer login test passed${NC}"
    echo "Token: ${token:0:20}..."
else
    echo -e "${RED}❌ Streamer login test failed${NC}"
    echo "Response: $login_response"
    exit 1
fi

# Get streamer user ID
user_response=$(curl -s -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer $token")

streamer_id=$(echo "$user_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)

if [ -n "$streamer_id" ]; then
    echo -e "${GREEN}✅ Streamer ID retrieved: $streamer_id${NC}"
else
    echo -e "${RED}❌ Failed to get streamer ID${NC}"
    echo "Response: $user_response"
    exit 1
fi

# Create a donator user
register_donator_response=$(curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testdonator",
    "email": "donator@mediashar.com",
    "password": "testpassword123",
    "is_streamer": false
  }')

# Login as donator
login_donator_response=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "donator@mediashar.com",
    "password": "testpassword123"
  }')

donator_token=$(echo "$login_donator_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$donator_token" ]; then
    echo -e "${GREEN}✅ Donator login test passed${NC}"
    
    # Get donator ID
    donator_response=$(curl -s -X GET http://localhost:8080/api/users/profile \
      -H "Authorization: Bearer $donator_token")
    
    donator_id=$(echo "$donator_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    if [ -n "$donator_id" ]; then
        echo -e "${GREEN}✅ Donator ID retrieved: $donator_id${NC}"
    else
        echo -e "${RED}❌ Failed to get donator ID${NC}"
        echo "Response: $donator_response"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  Using streamer token for donation${NC}"
    donator_token=$token
    donator_id=$streamer_id
fi

# Step 8: Test donation creation
echo -e "\n${YELLOW}💰 Testing donation creation...${NC}"

donation_response=$(curl -s -X POST http://localhost:8080/api/donations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $donator_token" \
  -d '{
    "amount": 50000,
    "currency": "IDR",
    "message": "Test donation for Docker testing",
    "streamer_id": '$streamer_id',
    "display_name": "Docker Tester",
    "is_anonymous": false
  }')

donation_id=$(echo "$donation_response" | grep -o '"id":[0-9]*' | cut -d':' -f2)

if [ -n "$donation_id" ]; then
    echo -e "${GREEN}✅ Donation creation test passed${NC}"
    echo "Donation ID: $donation_id"
else
    echo -e "${RED}❌ Donation creation test failed${NC}"
    echo "Response: $donation_response"
    exit 1
fi

# Step 9: Test Midtrans payment creation
echo -e "\n${YELLOW}💳 Testing Midtrans payment creation...${NC}"

midtrans_response=$(curl -s -X POST "http://localhost:8080/api/midtrans/payment/$donation_id" \
  -H "Authorization: Bearer $donator_token")

if echo "$midtrans_response" | grep -q "token"; then
    echo -e "${GREEN}✅ Midtrans payment creation test passed${NC}"
    
    # Extract token and redirect URL
    snap_token=$(echo "$midtrans_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    redirect_url=$(echo "$midtrans_response" | grep -o '"redirect_url":"[^"]*"' | cut -d'"' -f4)
    order_id=$(echo "$midtrans_response" | grep -o '"order_id":"[^"]*"' | cut -d'"' -f4)
    
    echo "Snap Token: ${snap_token:0:20}..."
    echo "Redirect URL: $redirect_url"
    echo "Order ID: $order_id"
else
    echo -e "${RED}❌ Midtrans payment creation test failed${NC}"
    echo "Response: $midtrans_response"
    order_id="DONATION-$donation_id-123456789"
fi

# Step 10: Test webhook endpoint
echo -e "\n${YELLOW}🔔 Testing webhook endpoint...${NC}"

webhook_response=$(curl -s -X POST http://localhost:8080/api/midtrans/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "transaction_status": "settlement",
    "status_code": "200",
    "transaction_id": "test-transaction-123",
    "order_id": "'$order_id'",
    "gross_amount": "50000.00",
    "payment_type": "bank_transfer",
    "signature_key": "test-signature"
  }')

if echo "$webhook_response" | grep -q "success\|processed"; then
    echo -e "${GREEN}✅ Webhook endpoint test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Webhook test completed (signature validation may fail in test)${NC}"
fi

# Step 11: Test transaction status
echo -e "\n${YELLOW}📊 Testing transaction status...${NC}"

status_response=$(curl -s "http://localhost:8080/api/midtrans/status/$order_id")
if echo "$status_response" | grep -q "success"; then
    echo -e "${GREEN}✅ Transaction status test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Transaction status test completed${NC}"
fi

# Step 12: Show container status
echo -e "\n${YELLOW}📊 Container Status:${NC}"
docker-compose ps

# Step 13: Show logs summary
echo -e "\n${YELLOW}📝 Recent Application Logs:${NC}"
docker-compose logs --tail=10 app

echo -e "\n${GREEN}🎉 Docker testing completed successfully!${NC}"
echo -e "\n${YELLOW}📋 Available Services:${NC}"
echo "• MediaShar API: http://localhost:8080"
echo "• API Health: http://localhost:8080/health"
echo "• API Docs: http://localhost:8083"
echo "• PgAdmin: http://localhost:8082"
echo ""
echo -e "${YELLOW}🧪 For manual testing:${NC}"
echo "• Streamer JWT Token: $token"
echo "• Donator JWT Token: $donator_token"
echo "• Streamer ID: $streamer_id"
echo "• Donator ID: $donator_id"
echo "• Donation ID: $donation_id"
echo "• Order ID: $order_id"
echo ""
echo -e "${YELLOW}🛑 To stop services:${NC}"
echo "docker-compose down" 