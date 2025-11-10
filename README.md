# Message Service

Automatic message sending service that sends SMS messages every 2 minutes via webhook using a custom Go scheduler.

**Note:** Requires **Go 1.25.4+** and **Docker**.

---

## Quick Start

### Option 1: Using Docker (Recommended)

```bash
git clone https://github.com/srcnDev/message-service.git
cd message-service
make docker-prod
```

API: http://localhost:8080  
Swagger: http://localhost:8080/swagger/index.html

**Note:** Message sender job starts automatically on application startup.

### Option 2: Local Development

```bash
# 1. Clone repository
git clone https://github.com/srcnDev/message-service.git
cd message-service

# 2. Configure environment (optional - .env file is already included)
# Edit .env file if you need to change default settings

# 3. Start infrastructure (PostgreSQL + Redis)
docker-compose up -d

# 4. Wait for database to be ready (5-10 seconds)
sleep 10

# 5. Run migrations and seed data
go run cmd/migrate/main.go -seed

# 6. Start application
go run cmd/api/main.go
```

**With Makefile:**

```bash
make docker-up    # Start infrastructure
make migrate      # Run migrations + seed
make run          # Start application
```

---

## Technology Stack

- **Go 1.25.4+** - Programming language
- **Gin** - High-performance HTTP web framework
- **PostgreSQL** - Primary database for message storage and status tracking
- **Redis** - Caching layer for sent messages (bonus feature)

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
│   ├── handler/          # HTTP handlers (Gin)
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
make build         # Build application binary
make run           # Run application locally
make migrate       # Run database migrations and seed data
make test          # Run unit tests (short mode)
make test-e2e      # Run end-to-end tests
make test-coverage # Run tests with coverage report

# Infrastructure
make docker-up    # Start PostgreSQL + Redis
make docker-down  # Stop infrastructure

# Production
make docker-prod      # Start full stack (PostgreSQL + Redis + App)
make docker-prod-down # Stop production

make clean        # Clean build artifacts
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

### Messages

```bash
GET  /health                      # Health check

GET  /api/v1/messages             # List all messages (with pagination)
GET  /api/v1/messages/sent        # List only sent messages (with pagination)
GET  /api/v1/messages/:id         # Get single message by ID
POST /api/v1/messages             # Create new message
PUT  /api/v1/messages/:id         # Update message
DELETE /api/v1/messages/:id       # Soft delete message
```

### Message Sender Job

```bash
GET  /api/v1/sender/status        # Get job status
POST /api/v1/sender/start         # Start sending job
POST /api/v1/sender/stop          # Stop sending job
```

**Note:** Job starts automatically on application startup.

**Example - Create Message:**

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "phoneNumber": "+905551234567",
    "content": "Hello World"
  }'
```

**Example - List Sent Messages:**

```bash
curl http://localhost:8080/api/v1/messages/sent?limit=10&offset=0
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

## Testing

### Unit Tests

```bash
# Run all unit tests
go test ./... -v -short

# With Makefile
make test

# With coverage report
go test ./... -cover -short

# Test specific package
go test ./internal/handler/... -v
```

### E2E Tests

```bash
# Run end-to-end tests (requires running database)
go test ./test/e2e/... -v -tags=e2e

# Or with Makefile
make test-e2e
```

**Test Coverage:** 225+ tests covering:
- Unit tests: Handler, Service, Repository layers
- Integration tests: Database operations
- E2E tests: Full workflow testing

---

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# View logs
docker-compose logs postgres

# Restart infrastructure
make docker-down && make docker-up
```

### Application Won't Start

```bash
# Check if port 8080 is available
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Check .env configuration
cat .env

# Verify database migrations
go run cmd/migrate/main.go
```

### Message Sender Not Working

1. Check job status: \`GET /api/v1/sender/status\`
2. Verify webhook URL is accessible
3. Check logs for errors
4. Ensure database has pending messages

---

_Sercan Yılmaz_

Email: [sercanyilmaz.dev@gmail.com](mailto:sercanyilmaz.dev@gmail.com)  
Phone: [+90 538 786 3537](tel:+905387863537)
