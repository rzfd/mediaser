# MediaShar Monitoring Guide

## Overview

MediaShar menggunakan **Prometheus** dan **Grafana** untuk monitoring dan observability yang comprehensive. Stack monitoring ini memberikan insight real-time tentang performa aplikasi, kesehatan sistem, dan business metrics.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Application   │───▶│   Prometheus    │───▶│     Grafana     │
│   (Metrics)     │    │   (Collection)  │    │  (Visualization)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Node Exporter  │    │ Postgres Export │    │   Dashboards    │
│ (System Metrics)│    │  (DB Metrics)   │    │   & Alerts      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Quick Start

### 1. Start Monitoring Stack

```bash
# Start semua services termasuk monitoring
make up

# Atau hanya monitoring services
make monitoring-up
```

### 2. Access Dashboards

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001
  - Username: `admin`
  - Password: `admin123`

### 3. Test Metrics

```bash
# Test metrics endpoints
make metrics-test

# Check service health
make health
```

## Metrics Overview

### HTTP Metrics
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request duration
- `http_response_size_bytes` - Response size

### gRPC Metrics
- `grpc_requests_total` - Total gRPC requests
- `grpc_request_duration_seconds` - gRPC request duration

### Database Metrics
- `db_connections_active` - Active database connections
- `db_query_duration_seconds` - Database query duration
- `db_queries_total` - Total database queries

### Business Metrics
- `donations_total` - Total donations
- `donation_amount` - Donation amounts
- `payments_processed_total` - Payment processing

### System Metrics
- `go_routines` - Number of goroutines
- `memory_usage_bytes` - Memory usage
- `cpu_usage_percent` - CPU usage

## Grafana Dashboards

### MediaShar Overview Dashboard

Dashboard utama yang menampilkan:

1. **Service Health** - Status UP/DOWN semua services
2. **HTTP Request Rate** - Request rate per service
3. **Response Time** - 95th percentile response time
4. **Error Rate** - Error rate by service
5. **Database Performance** - Query performance metrics
6. **Business Metrics** - Donations dan payments
7. **System Resources** - Memory dan CPU usage

### Custom Queries

Beberapa query Prometheus yang berguna:

```promql
# Request rate per service
sum(rate(http_requests_total[5m])) by (service)

# Error rate
rate(http_requests_total{status_code=~"5.."}[5m]) / rate(http_requests_total[5m])

# 95th percentile response time
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Database query performance
histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))
```

## Alerting Rules

### Critical Alerts

1. **HighErrorRate** - Error rate > 5% for 5 minutes
2. **ServiceDown** - Service unavailable for 1 minute
3. **DatabaseConnectionFailure** - DB connection issues
4. **PaymentFailures** - Payment failure rate > 2%

### Warning Alerts

1. **HighResponseTime** - 95th percentile > 2 seconds
2. **HighMemoryUsage** - Memory usage > 500MB
3. **SlowDatabaseQueries** - Query time > 1 second
4. **HighGoRoutines** - Goroutines > 1000

## Configuration

### Prometheus Configuration

File: `monitoring/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'api-gateway'
    static_configs:
      - targets: ['api-gateway:8080']
    metrics_path: /metrics
    scrape_interval: 10s
```

### Grafana Configuration

- **Datasource**: `monitoring/grafana/provisioning/datasources/prometheus.yml`
- **Dashboards**: `monitoring/grafana/provisioning/dashboards/`

## Development

### Adding Custom Metrics

1. **Import metrics package**:
```go
import "github.com/rzfd/mediashar/pkg/metrics"
```

2. **Record business metrics**:
```go
// Record donation
metrics.GetMetrics().RecordDonation("api-gateway", "USD", "success", 100.0)

// Record payment
metrics.GetMetrics().RecordPayment("api-gateway", "midtrans", "success")

// Record database query
metrics.GetMetrics().RecordDBQuery("api-gateway", "SELECT", "success", duration)
```

3. **Add HTTP middleware**:
```go
e.Use(middleware.MetricsMiddleware("service-name"))
```

### Custom Dashboards

1. Create dashboard JSON in `monitoring/grafana/dashboards/`
2. Restart Grafana container
3. Dashboard akan auto-load

## Troubleshooting

### Common Issues

1. **Metrics endpoint tidak accessible**
   ```bash
   # Check if service is running
   curl http://localhost:8080/metrics
   
   # Check container logs
   docker-compose logs api-gateway
   ```

2. **Prometheus tidak scrape targets**
   ```bash
   # Check Prometheus targets
   curl http://localhost:9090/api/v1/targets
   
   # Check Prometheus config
   docker-compose logs prometheus
   ```

3. **Grafana tidak show data**
   - Verify datasource connection
   - Check time range
   - Verify query syntax

### Logs

```bash
# All monitoring logs
make monitoring-logs

# Specific service logs
docker-compose logs prometheus
docker-compose logs grafana
```

## Production Considerations

### Security

1. **Change default passwords**:
```yaml
environment:
  - GF_SECURITY_ADMIN_PASSWORD=your-secure-password
```

2. **Enable HTTPS**:
```yaml
environment:
  - GF_SERVER_PROTOCOL=https
  - GF_SERVER_CERT_FILE=/etc/ssl/certs/grafana.crt
  - GF_SERVER_CERT_KEY=/etc/ssl/private/grafana.key
```

### Scaling

1. **Prometheus retention**:
```yaml
command:
  - '--storage.tsdb.retention.time=200h'
```

2. **Resource limits**:
```yaml
deploy:
  resources:
    limits:
      memory: 1G
      cpus: '0.5'
```

### Backup

```bash
# Backup Prometheus data
docker run --rm -v mediashar_prometheus_data:/data -v $(pwd):/backup alpine tar czf /backup/prometheus-backup.tar.gz /data

# Backup Grafana data
docker run --rm -v mediashar_grafana_data:/data -v $(pwd):/backup alpine tar czf /backup/grafana-backup.tar.gz /data
```

## Useful Commands

```bash
# Start monitoring only
make monitoring-up

# Stop monitoring
make monitoring-down

# View monitoring logs
make monitoring-logs

# Test metrics endpoints
make metrics-test

# Check service health
make health

# Clean up everything
make clean
```

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Go Prometheus Client](https://github.com/prometheus/client_golang)
- [Echo Middleware](https://echo.labstack.com/middleware/) 