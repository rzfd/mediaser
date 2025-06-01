# Docker Testing Guide untuk Midtrans Integration

## 🐳 Prerequisite

### 1. Install Docker dan Docker Compose
```bash
# Install Docker (Ubuntu/Debian)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt-get install docker-compose-plugin

# Atau install standalone docker-compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. Verifikasi Installation
```bash
docker --version
docker compose version
```

## 🚀 Quick Start Testing

### Automated Testing
```bash
# Make script executable
chmod +x scripts/test-docker.sh

# Run automated test
./scripts/test-docker.sh
```

### Manual Testing Steps

#### 1. Start Services
```bash
# Clean previous containers
docker-compose down --volumes --remove-orphans

# Build and start all services
docker-compose up -d --build

# Check container status
docker-compose ps
```

#### 2. Monitor Logs
```bash
# Watch all logs
docker-compose logs -f

# Watch specific service
docker-compose logs -f app
docker-compose logs -f postgres
```

#### 3. Health Check
```bash
# API Health Check
curl http://localhost:8080/health

# Readiness Check
curl http://localhost:8080/ready

# Database Check
curl http://localhost:8080/api/health
```

## 🧪 Manual API Testing

### 1. User Registration
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "teststreamer",
    "email": "streamer@mediashar.com",
    "password": "testpassword123",
    "is_streamer": true
  }'
```

### 2. User Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "streamer@mediashar.com",
    "password": "testpassword123"
  }'
```
**Save the JWT token dari response!**

### 3. Create Donation
```bash
# Replace YOUR_JWT_TOKEN dengan token dari login
curl -X POST http://localhost:8080/api/donations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "amount": 50000,
    "currency": "IDR",
    "message": "Test donation via Docker",
    "streamer_id": 1,
    "display_name": "Docker Tester",
    "is_anonymous": false
  }'
```
**Save the donation ID dari response!**

### 4. Create Midtrans Payment
```bash
# Replace DONATION_ID dengan ID dari step sebelumnya
curl -X POST http://localhost:8080/api/midtrans/payment/DONATION_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response akan berisi:
- `token`: Snap token untuk payment
- `redirect_url`: URL untuk redirect ke Midtrans
- `order_id`: Order ID untuk tracking

### 5. Test Webhook (Simulate Midtrans)
```bash
curl -X POST http://localhost:8080/api/midtrans/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "transaction_status": "settlement",
    "status_code": "200",
    "transaction_id": "test-transaction-123",
    "order_id": "DONATION-1-1642751234",
    "gross_amount": "50000.00",
    "payment_type": "bank_transfer",
    "signature_key": "test-signature"
  }'
```

### 6. Check Transaction Status
```bash
curl http://localhost:8080/api/midtrans/status/DONATION-1-1642751234
```

## 🎯 Frontend Testing dengan Snap.js

### 1. Basic HTML Test Page
Buat file `test-frontend.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Midtrans Test</title>
    <script src="https://app.sandbox.midtrans.com/snap/snap.js" 
            data-client-key="SB-Mid-client-Yy6kDu1A1cTYWiYy"></script>
</head>
<body>
    <h1>Midtrans Payment Test</h1>
    <button onclick="pay()">Test Payment</button>
    
    <script>
        // Replace with actual values from your API tests
        const token = 'YOUR_JWT_TOKEN';
        const donationId = 'YOUR_DONATION_ID';
        
        async function pay() {
            try {
                // Get snap token
                const response = await fetch(`http://localhost:8080/api/midtrans/payment/${donationId}`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`
                    }
                });
                
                const data = await response.json();
                console.log('Payment response:', data);
                
                // Open Snap payment
                snap.pay(data.data.token, {
                    onSuccess: function(result) {
                        console.log('Payment success:', result);
                        alert('Payment berhasil!');
                    },
                    onPending: function(result) {
                        console.log('Payment pending:', result);
                        alert('Payment pending, silakan selesaikan pembayaran');
                    },
                    onError: function(result) {
                        console.log('Payment error:', result);
                        alert('Payment error!');
                    },
                    onClose: function() {
                        console.log('Payment popup ditutup');
                    }
                });
            } catch (error) {
                console.error('Error:', error);
                alert('Error creating payment');
            }
        }
    </script>
</body>
</html>
```

### 2. Test dengan Browser
1. Buka file HTML di browser
2. Klik tombol "Test Payment"
3. Akan muncul Snap popup
4. Gunakan test credentials Midtrans

## 🔧 Troubleshooting

### Common Issues

#### 1. Container tidak start
```bash
# Check container status
docker-compose ps

# Check logs
docker-compose logs app

# Restart specific service
docker-compose restart app
```

#### 2. Database connection error
```bash
# Check database logs
docker-compose logs postgres

# Check database health
docker exec mediashar_postgres pg_isready -U postgres

# Connect to database manually
docker exec -it mediashar_postgres psql -U postgres -d donation_system
```

#### 3. Midtrans error responses
```bash
# Check app logs for Midtrans errors
docker-compose logs app | grep -i midtrans

# Verify environment variables
docker exec mediashar_app env | grep MIDTRANS
```

#### 4. Network issues
```bash
# Check network
docker network ls
docker network inspect mediashar_mediashar_network

# Test internal connectivity
docker exec mediashar_app curl http://postgres:5432
```

### Debug Commands

```bash
# Enter app container
docker exec -it mediashar_app sh

# Check app configuration
docker exec mediashar_app cat configs/config.yaml

# Test database from app container
docker exec mediashar_app wget --spider http://postgres:5432

# Check environment variables
docker exec mediashar_app printenv | grep -E "(DB_|MIDTRANS_)"
```

## 📊 Monitoring

### Container Resources
```bash
# Check resource usage
docker stats

# Check container health
docker inspect --format='{{.State.Health.Status}}' mediashar_app
```

### Application Metrics
```bash
# Health endpoint
curl http://localhost:8080/health

# Database status
curl http://localhost:8080/ready
```

## 🛑 Cleanup

### Stop Services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down --volumes

# Remove images as well
docker-compose down --rmi all --volumes --remove-orphans
```

### Complete Cleanup
```bash
# Remove all MediaShar related containers, networks, images
docker system prune -f
docker volume prune -f
docker network prune -f
```

## 🚀 Production-like Testing

### Environment Variables Testing
Buat file `.env.test`:

```bash
# Test with different configurations
MIDTRANS_ENVIRONMENT=sandbox
DB_NAME=test_donation_system
JWT_SECRET=test-secret-key
```

### Run with test environment
```bash
# Use test environment
docker-compose --env-file .env.test up -d

# Test with production-like settings
MIDTRANS_ENVIRONMENT=production docker-compose up -d
```

## 📈 Performance Testing

### Load Testing dengan curl
```bash
# Create multiple concurrent requests
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/health &
done
wait
```

### Database Performance
```bash
# Check database performance
docker exec mediashar_postgres psql -U postgres -d donation_system -c "
SELECT 
  schemaname,
  tablename,
  attname,
  n_distinct,
  correlation
FROM pg_stats 
WHERE schemaname = 'public';
"
```

## 💡 Tips

1. **Always check logs first** saat ada masalah
2. **Use health checks** untuk memantau service
3. **Test incrementally** - mulai dari health check, lalu API, terakhir payment
4. **Save tokens dan IDs** untuk testing manual
5. **Use proper test data** yang realistis
6. **Monitor resource usage** jika testing di mesin terbatas

## 🔗 Useful URLs

Saat services berjalan:
- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API Documentation**: http://localhost:8083
- **Database Admin**: http://localhost:8082 (admin@mediashar.com / admin123)
- **Midtrans Sandbox**: https://dashboard.sandbox.midtrans.com 