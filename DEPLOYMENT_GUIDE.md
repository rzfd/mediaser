# 🚀 MediaShar Currency & Language System Deployment Guide

This guide covers the complete deployment of MediaShar's enhanced currency conversion and language translation system.

## 📋 Table of Contents

1. [System Overview](#system-overview)
2. [Prerequisites](#prerequisites)
3. [Quick Start](#quick-start)
4. [Detailed Setup](#detailed-setup)
5. [Configuration](#configuration)
6. [Testing](#testing)
7. [Monitoring](#monitoring)
8. [Troubleshooting](#troubleshooting)

## 🏗️ System Overview

### New Features Added
- **Currency System**: Real-time exchange rates with 7 supported currencies
- **Language System**: Free translation with 3 supported languages
- **React Components**: Beautiful UI for currency conversion and translation
- **External APIs**: Integration with ExchangeRate-API and LibreTranslate

### Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │  Microservices  │
│   (React)       │◄──►│   (Port 8080)   │◄──►│                 │
│                 │    │                 │    │ • Currency      │
│ • CurrencyConv  │    │ • Routes        │    │ • Language      │
│ • Translation   │    │ • Auth          │    │ • Donation      │
│ • Modern UI     │    │ • Load Balance  │    │ • Payment       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Supported Features

#### 💱 Currency System
- **Currencies**: IDR, USD, CNY, EUR, JPY, SGD, MYR
- **Real-time Rates**: ExchangeRate-API integration
- **Smart Caching**: 1-hour cache with auto-refresh
- **Auto-formatting**: Region-specific number formatting

#### 🌐 Language System
- **Languages**: Indonesian, English, Mandarin
- **Free Translation**: LibreTranslate API
- **Auto-detection**: Language detection with confidence
- **Bulk Support**: Translate up to 100 texts at once

## 🔧 Prerequisites

### System Requirements
- **OS**: Linux/macOS/Windows with WSL2
- **Docker**: Version 20.10+
- **Docker Compose**: Version 2.0+
- **Memory**: Minimum 4GB RAM
- **Storage**: 10GB free space

### Development Tools (Optional)
- **Go**: Version 1.21+ for backend development
- **Node.js**: Version 18+ for frontend development
- **Git**: For version control

## ⚡ Quick Start

### 1. Clone and Setup
```bash
# Clone repository
git clone https://github.com/rzfd/mediashar.git
cd mediashar

# Make scripts executable
chmod +x scripts/*.sh

# Create required directories
mkdir -p build logs data
```

### 2. Environment Configuration
```bash
# Copy environment template
cp .env.example .env

# Edit configuration (optional)
nano .env
```

### 3. Start Services
```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f api-gateway
```

### 4. Verify Installation
```bash
# Run health check
curl http://localhost:8080/health

# Run comprehensive tests
./scripts/test-currency-language.sh
```

## 🔧 Detailed Setup

### Database Initialization

The system automatically creates and initializes databases:

#### Currency Database (Port 5435)
- **Tables**: currency_rates, currency_info, user_currency_preferences
- **Initial Data**: 7 currencies with basic exchange rates
- **Indexes**: Optimized for fast lookups

#### Language Database (Port 5436)
- **Tables**: language_configs, language_info, translation_cache
- **Initial Data**: 3 languages with system translations
- **Cleanup**: Auto-cleanup of expired translations

### Service Configuration

#### API Gateway (Port 8080)
```yaml
# config/api-gateway.yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  host: gateway_db
  port: 5432
  user: gateway_user
  password: gateway_password
  name: gateway_db

services:
  currency_service: "currency-service:8084"
  language_service: "language-service:8085"
```

#### Currency Service (Port 8084)
```yaml
# config/currency-service.yaml
grpc:
  port: 8084

http:
  port: 8094

database:
  host: currency_db
  port: 5432
  user: currency_user
  password: currency_password
  name: currency_db

external_apis:
  exchange_rate_api: "https://v6.exchangerate-api.com/v6/latest"
  cache_ttl: 3600  # 1 hour
```

#### Language Service (Port 8085)
```yaml
# config/language-service.yaml
grpc:
  port: 8085

http:
  port: 8095

database:
  host: language_db
  port: 5432
  user: language_user
  password: language_password
  name: language_db

external_apis:
  libretranslate_api: "https://libretranslate.de/translate"
  translation_cache_ttl: 86400  # 24 hours
```

## ⚙️ Configuration

### Environment Variables

```bash
# .env file
COMPOSE_PROJECT_NAME=mediashar

# Database Configuration
POSTGRES_VERSION=15-alpine
DB_MAX_CONNECTIONS=100

# API Configuration
API_GATEWAY_PORT=8080
CURRENCY_SERVICE_PORT=8084
LANGUAGE_SERVICE_PORT=8085

# External APIs
EXCHANGE_RATE_API_URL=https://v6.exchangerate-api.com/v6/latest
LIBRETRANSLATE_API_URL=https://libretranslate.de/translate

# Cache Configuration
CURRENCY_CACHE_TTL=3600
TRANSLATION_CACHE_TTL=86400

# Security
JWT_SECRET=your-super-secret-jwt-key-change-in-production
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
```

### Frontend Configuration

```javascript
// frontend/src/config/api.js
export const API_CONFIG = {
  BASE_URL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  CURRENCY_ENDPOINT: '/api/currency',
  LANGUAGE_ENDPOINT: '/api/language',
  TIMEOUT: 10000,
  RETRY_ATTEMPTS: 3
};

export const CURRENCY_CONFIG = {
  SUPPORTED_CURRENCIES: ['IDR', 'USD', 'CNY', 'EUR', 'JPY', 'SGD', 'MYR'],
  DEFAULT_FROM: 'USD',
  DEFAULT_TO: 'IDR',
  CACHE_DURATION: 300000 // 5 minutes
};

export const LANGUAGE_CONFIG = {
  SUPPORTED_LANGUAGES: [
    { code: 'id', name: 'Indonesian', flag: '🇮🇩' },
    { code: 'en', name: 'English', flag: '🇺🇸' },
    { code: 'zh', name: 'Chinese', flag: '🇨🇳' }
  ],
  DEFAULT_FROM: 'en',
  DEFAULT_TO: 'id',
  MAX_TEXT_LENGTH: 5000,
  BULK_LIMIT: 100
};
```

## 🧪 Testing

### Automated Testing

```bash
# Run all tests
./scripts/test-currency-language.sh

# Test specific features
./scripts/test-currency-language.sh --currency
./scripts/test-currency-language.sh --language
./scripts/test-currency-language.sh --integration

# Check service health
./scripts/test-currency-language.sh --health
```

### Manual Testing

#### Currency Conversion
```bash
# Convert 100 USD to IDR
curl -X POST http://localhost:8080/api/currency/convert \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100,
    "from_currency": "USD",
    "to_currency": "IDR"
  }'

# Get exchange rate
curl "http://localhost:8080/api/currency/rate?from=USD&to=IDR"

# List supported currencies
curl http://localhost:8080/api/currency/list
```

#### Language Translation
```bash
# Translate text
curl -X POST http://localhost:8080/api/language/translate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello World",
    "from_language": "en",
    "to_language": "id"
  }'

# Detect language
curl -X POST http://localhost:8080/api/language/detect \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Selamat pagi"
  }'

# Bulk translate
curl -X POST http://localhost:8080/api/language/bulk-translate \
  -H "Content-Type: application/json" \
  -d '{
    "texts": ["Hello", "Thank you", "Good morning"],
    "from_language": "en",
    "to_language": "id"
  }'
```

### Frontend Testing

```bash
# Install dependencies
cd frontend
npm install

# Run development server
npm start

# Run tests
npm test

# Build for production
npm run build
```

## 📊 Monitoring

### Health Checks

```bash
# API Gateway health
curl http://localhost:8080/health

# Services health
curl http://localhost:8080/health/services

# Database health
docker-compose exec currency_db pg_isready
docker-compose exec language_db pg_isready
```

### Logs Monitoring

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f api-gateway
docker-compose logs -f currency-service
docker-compose logs -f language-service

# View database logs
docker-compose logs -f currency_db
docker-compose logs -f language_db
```

### Performance Metrics

```bash
# Check container stats
docker stats

# Check database connections
docker-compose exec currency_db psql -U currency_user -d currency_db -c "SELECT count(*) FROM pg_stat_activity;"

# Check cache hit rates
docker-compose exec currency_db psql -U currency_user -d currency_db -c "SELECT count(*) FROM currency_rates WHERE last_updated > EXTRACT(EPOCH FROM NOW()) - 3600;"
```

## 🔧 Troubleshooting

### Common Issues

#### 1. Services Not Starting
```bash
# Check Docker status
docker --version
docker-compose --version

# Check port conflicts
netstat -tulpn | grep :8080
netstat -tulpn | grep :5432

# Restart services
docker-compose down
docker-compose up -d
```

#### 2. Database Connection Issues
```bash
# Check database status
docker-compose ps

# Reset databases
docker-compose down -v
docker-compose up -d

# Check database logs
docker-compose logs currency_db
docker-compose logs language_db
```

#### 3. External API Issues
```bash
# Test ExchangeRate-API
curl "https://v6.exchangerate-api.com/v6/latest/USD"

# Test LibreTranslate
curl -X POST "https://libretranslate.de/translate" \
  -H "Content-Type: application/json" \
  -d '{
    "q": "Hello",
    "source": "en",
    "target": "id"
  }'
```

#### 4. Frontend Issues
```bash
# Check API connectivity
curl http://localhost:8080/health

# Clear browser cache
# Check browser console for errors
# Verify API endpoints in Network tab
```

### Performance Optimization

#### Database Optimization
```sql
-- Check slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;

-- Analyze table statistics
ANALYZE currency_rates;
ANALYZE language_configs;
```

#### Cache Optimization
```bash
# Monitor cache hit rates
docker-compose exec currency_db psql -U currency_user -d currency_db -c "
SELECT 
  'currency_rates' as table_name,
  count(*) as total_rates,
  count(*) FILTER (WHERE last_updated > EXTRACT(EPOCH FROM NOW()) - 3600) as fresh_rates
FROM currency_rates;
"
```

### Scaling Considerations

#### Horizontal Scaling
```yaml
# docker-compose.override.yml
version: '3.8'
services:
  currency-service:
    deploy:
      replicas: 3
  
  language-service:
    deploy:
      replicas: 2
```

#### Load Balancing
```nginx
# nginx.conf
upstream currency_backend {
    server currency-service-1:8084;
    server currency-service-2:8084;
    server currency-service-3:8084;
}

upstream language_backend {
    server language-service-1:8085;
    server language-service-2:8085;
}
```

## 🎯 Next Steps

### Production Deployment
1. **Security**: Implement proper authentication and authorization
2. **SSL/TLS**: Add HTTPS certificates
3. **Monitoring**: Set up Prometheus and Grafana
4. **Backup**: Implement database backup strategy
5. **CI/CD**: Set up automated deployment pipeline

### Feature Enhancements
1. **More Currencies**: Add support for additional currencies
2. **More Languages**: Expand language support
3. **Caching**: Implement Redis for better caching
4. **Rate Limiting**: Add API rate limiting
5. **Analytics**: Add usage analytics and reporting

### Maintenance
1. **Updates**: Regular dependency updates
2. **Monitoring**: Set up alerting for failures
3. **Cleanup**: Automated cleanup of old cache data
4. **Optimization**: Regular performance tuning

---

## 📞 Support

For issues and questions:
- **GitHub Issues**: [Create an issue](https://github.com/rzfd/mediashar/issues)
- **Documentation**: Check the `/docs` directory
- **Logs**: Always include relevant logs when reporting issues

---

**Happy Deploying! 🚀** 