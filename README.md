# Message Service

Automatic message sending service that sends SMS messages every 2 minutes via webhook using a custom Go scheduler.

**Note:** Requires **Go 1.23+** and **Docker**.

---

## Quick Start

### Option 1: Using Docker (Recommended)

```bash
git clone https://github.com/srcnDev/message-service.git
cd message-service
make docker-prod
```

API: http://localhost:8080  
Swagger: http://localhost:8080/swagger/

### Option 2: Local Development

```bash
# 1. Clone repository
git clone https://github.com/srcnDev/message-service.git
cd message-service

# 2. Start infrastructure (PostgreSQL + Redis)
docker-compose up -d

# 3. Run migrations and seed data
go run cmd/migrate/main.go -seed

# 4. Start application
go run cmd/api/main.go
```

**With Makefile:**

```bash
make docker-up    # Start infrastructure
make migrate      # Run migrations + seed
make run          # Start application
```

---

## Project Structure

```
message-service/
├── cmd/
│   ├── api/              # Main application
│   └── migrate/          # Database migration tool
├── internal/
│   ├── domain/           # Business entities
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   ├── handler/          # HTTP handlers
│   └── job/              # Background jobs (message sender)
├── pkg/
│   ├── scheduler/        # Custom Go scheduler (no cron)
│   ├── webhook/          # Webhook client
│   ├── database/         # PostgreSQL client
│   └── health/           # Health check
├── test/
│   └── e2e/              # End-to-end tests
├── config/               # Configuration management
├── docs/                 # Swagger documentation
├── docker-compose.yaml   # Infrastructure (PostgreSQL + Redis)
├── docker-compose.prod.yaml  # Production (+ App)
├── Dockerfile            # Application image
└── Makefile              # Build commands
```

---

## Available Commands

### Makefile

```bash
make build        # Build application binary
make run          # Run application locally
make migrate      # Run database migrations and seed data
make test         # Run all tests

# Infrastructure
make docker-up    # Start PostgreSQL + Redis
make docker-down  # Stop infrastructure

# Production
make docker-prod      # Start full stack (PostgreSQL + Redis + App)
make docker-prod-down # Stop production
```

### Manual Commands

```bash
# Build
go build -o bin/message-service cmd/api/main.go

# Run
go run cmd/api/main.go

# Test
go test ./... -v

# Migration
go run cmd/migrate/main.go -seed

# Docker
docker-compose up -d
docker-compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d
```

---

## API Endpoints

```bash
GET  /health                   # Health check

GET  /api/v1/messages          # List all messages
GET  /api/v1/messages/:id      # Get single message
POST /api/v1/messages          # Create message
PUT  /api/v1/messages/:id      # Update message
DELETE /api/v1/messages/:id    # Delete message

GET  /api/v1/sender/status     # Job status
POST /api/v1/sender/start      # Start job
POST /api/v1/sender/stop       # Stop job
```

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "+905551234567",
    "content": "Hello World"
  }'
```

---

## Configuration

Edit `.env` file:

```env
# Application
APP_PORT=8080

# Database
POSTGRES_DB_HOST=localhost
POSTGRES_DB_PORT=5432
POSTGRES_DB_USERNAME=postgres
POSTGRES_DB_PASSWORD=postgres
POSTGRES_DB_NAME=message_service

# Redis
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379

# Message Sender (Case Study Requirements)
MESSAGE_SENDER_INTERVAL=120    # seconds (2 minutes)
MESSAGE_SENDER_BATCH_SIZE=2    # messages per cycle

# Webhook
WEBHOOK_URL=https://webhook.site/7d2fa94f-bb3c-47d7-b787-8aaacbd5097d
WEBHOOK_AUTH_KEY=INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo
```

---

## Key Features

- ⏰ **Custom Scheduler**: Pure Go implementation using `time.Ticker`
- 📦 **Batch Processing**: Sends exactly 2 messages per cycle
- 🔄 **Idempotency**: Messages are never sent twice (status tracking)
- 💾 **PostgreSQL**: Message storage with GORM
- 🔴 **Redis Cache**: Stores sent message history (bonus feature)
- 📝 **Swagger**: Auto-generated API documentation
- ✅ **Tests**: 169 tests with 95%+ coverage

---

## Testing

```bash
# Run unit tests
go test ./... -v

# With Makefile
make test

# With coverage
go test ./... -cover
```

**Test Coverage:** 169 tests (163 unit + 6 e2e), 95-100% on core modules

---

_Sercan Yılmaz_

Email: [sercanyilmaz.dev@gmail.com](mailto:sercanyilmaz.dev@gmail.com)  
Phone: [+90 538 786 3537](tel:+905387863537)