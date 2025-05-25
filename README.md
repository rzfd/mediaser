# Donation System for Streamers

A backend service built with Go, Echo, and GORM to handle donations to streamers with JWT authentication.

## Features

- User management (streamers and donators)
- JWT-based authentication and authorization
- Donation processing with multiple payment providers
- **QRIS (Quick Response Code Indonesian Standard) integration**
- Payment webhooks handling
- Donation analytics and statistics
- Role-based access control (Streamer vs Donator)
- QR Code generation for instant payments

## Architecture

The application follows a clean architecture pattern:

- **Models**: Data structures and database schema
- **Repository**: Database access layer
- **Service**: Business logic layer
- **Handler**: HTTP request handlers
- **Routes**: API routing configuration
- **Middleware**: JWT authentication and authorization
- **Config**: Application configuration

## Project Structure

```
├── cmd/api/                 # Application entry point
├── configs/                 # Configuration files
├── internal/
│   ├── models/             # Data models and database schema
│   ├── repository/         # Database access layer
│   ├── service/            # Business logic layer
│   ├── handler/            # HTTP request handlers
│   ├── middleware/         # JWT authentication middleware
│   └── routes/             # API routing configuration
├── pkg/utils/              # Utility functions
├── docker-compose.yml      # Docker setup for PostgreSQL
└── init.sql               # Database initialization script
```

## Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Docker & Docker Compose (recommended, for easy setup)

### Option 1: Full Docker Setup (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/rzfd/mediashar.git
cd mediashar
```

2. Start all services with Docker:
```bash
make docker-setup
```

3. Access the application:
   - **API**: http://localhost:8080
   - **pgAdmin**: http://localhost:8082 (admin@mediashar.com / admin123)
   - **Adminer**: http://localhost:8081

### Option 2: Development Setup (Database Only)

1. Clone the repository:
```bash
git clone https://github.com/rzfd/mediashar.git
cd mediashar
```

2. Start database services:
```bash
make dev-setup
```

3. Run the application locally:
```bash
make run
```

4. Access database admin panels:
   - **pgAdmin**: http://localhost:8082 (admin@mediashar.com / admin123)
   - **Adminer**: http://localhost:8081

### Option 3: Manual PostgreSQL Setup

1. Clone the repository:
```bash
git clone https://github.com/rzfd/mediashar.git
cd mediashar
```

2. Install dependencies:
```bash
go mod download
```

3. Setup PostgreSQL database:
```bash
# Create database
createdb donation_system

# Or using psql
psql -U postgres
CREATE DATABASE donation_system;
```

4. Configure the application:

Edit `configs/config.yaml` with your database settings:
```yaml
db:
  host: "localhost"
  port: "5432"
  username: "postgres"
  password: "your_password"
  name: "donation_system"

auth:
  jwtSecret: "your-super-secret-jwt-key"
  tokenExpiry: 86400  # 24 hours in seconds
```

5. Run the application:
```bash
go run cmd/api/main.go
```

## Database

The application uses **PostgreSQL** as the primary database with the following configuration:

- **Driver**: `gorm.io/driver/postgres`
- **ORM**: GORM (Go Object-Relational Mapping)
- **Default Port**: 5432
- **Connection**: Uses pgx driver with SSL disabled for development

### Database Schema

The application will automatically create the following tables:
- `users` - User information (streamers and donators)
- `donations` - Donation records with payment information

### Docker Services

- **app**: Go application (port 8080)
- **postgres**: PostgreSQL database (port 5432)
- **pgadmin**: pgAdmin web interface (port 8082)
- **adminer**: Adminer database admin (port 8081)

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Token Claims

JWT tokens contain the following claims:
- `user_id`: User's unique identifier
- `email`: User's email address
- `is_streamer`: Boolean indicating if user is a streamer
- `exp`: Token expiration time
- `iat`: Token issued at time

## API Endpoints

### Authentication Endpoints (Public)

- `POST /api/auth/register`: Register a new user account
- `POST /api/auth/login`: Login and receive JWT token
- `POST /api/auth/refresh`: Refresh an existing JWT token

### Profile Management Endpoints (Protected)

- `GET /api/auth/profile`: Get current user's profile
- `PUT /api/auth/profile`: Update current user's profile
- `POST /api/auth/change-password`: Change current user's password
- `POST /api/auth/logout`: Logout (client-side token invalidation)

### User Endpoints

**Public:**
- `GET /api/users/:id`: Get user by ID
- `GET /api/streamers`: List all streamers

**Protected:**
- `POST /api/users`: Create a new user (admin)
- `PUT /api/users/:id`: Update user information (self or admin)
- `GET /api/users/:id/donations`: Get donations by user (as donator)

### Donation Endpoints (Protected)

- `POST /api/donations`: Create a new donation
- `GET /api/donations`: List all donations
- `GET /api/donations/:id`: Get donation by ID
- `GET /api/donations/latest`: Get latest donations

### Streamer-Only Endpoints (Protected + Streamer Role)

- `GET /api/streamers/:id/donations`: Get donations for a streamer
- `GET /api/streamers/:id/total`: Get total donation amount for a streamer

### Payment Endpoints (Protected)

- `POST /api/payments/process`: Process payment for a donation

### QRIS Payment Endpoints

**Public:**
- `POST /api/qris/donate`: Create donation with QRIS QR code (anonymous or authenticated)

**Protected:**
- `POST /api/qris/donations/:id/generate`: Generate QRIS for existing donation
- `GET /api/qris/status/:transaction_id`: Check QRIS payment status

### Payment Webhook Endpoints (Public)

- `POST /api/webhooks/paypal`: PayPal webhook handler
- `POST /api/webhooks/stripe`: Stripe webhook handler
- `POST /api/webhooks/crypto`: Crypto payment webhook handler
- `POST /api/webhooks/qris`: QRIS payment webhook handler

## Usage Examples

### Register a new user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "streamer1",
    "email": "streamer@example.com",
    "password": "password123",
    "full_name": "John Streamer",
    "is_streamer": true,
    "description": "Gaming streamer"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "streamer@example.com",
    "password": "password123"
  }'
```

### Access protected endpoint
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Create a donation
```bash
curl -X POST http://localhost:8080/api/donations \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 10.00,
    "currency": "USD",
    "message": "Great stream!",
    "streamer_id": 1,
    "display_name": "Anonymous"
  }'
```

### Create donation with QRIS
```bash
curl -X POST http://localhost:8080/api/qris/donate \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 50000,
    "currency": "IDR",
    "message": "Semangat streaming!",
    "streamer_id": 1,
    "display_name": "Anonymous Supporter",
    "is_anonymous": true
  }'
```

## Development

### Adding New Payment Providers

To add a new payment provider:

1. Implement the `PaymentProcessor` interface in a new file
2. Register the processor in the `PaymentService`
3. Add appropriate webhook handling in the `WebhookHandler`
4. Add the webhook route in `internal/routes/routes.go`

### Database Migrations

Database migrations are handled automatically when the application starts using GORM's AutoMigrate feature.

### Routes Organization

All API routes are organized in the `internal/routes` package. Routes are grouped by:
- **Public routes**: No authentication required
- **Protected routes**: JWT authentication required
- **Streamer-only routes**: JWT authentication + streamer role required
- **Webhook routes**: No authentication (secured by webhook secrets)

### JWT Middleware

The application includes several middleware functions:
- `JWTMiddleware`: Validates JWT tokens and extracts user info
- `OptionalJWTMiddleware`: Optional JWT validation (doesn't fail if no token)
- `StreamerOnlyMiddleware`: Ensures only streamers can access certain endpoints

### Docker Commands

```bash
# Full setup (all services)
make docker-setup

# Development setup (database only)
make dev-setup

# Build Docker image
make docker-build

# Start all services
make docker-up

# Start only database services
make docker-db

# Stop services
make docker-down

# Stop and remove volumes
make docker-clean

# View logs
make docker-logs

# Execute shell in app container
make docker-shell

# Execute psql in postgres
make docker-psql
```

For detailed Docker documentation, see [docs/DOCKER_SETUP.md](docs/DOCKER_SETUP.md).

## Security Considerations

- JWT tokens expire after 24 hours (configurable)
- Passwords are hashed using bcrypt
- CORS is enabled for cross-origin requests
- SQL injection protection via GORM
- Input validation on all endpoints

## License

MIT