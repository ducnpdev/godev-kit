# Payment System Setup Guide

## Prerequisites

1. **Go 1.24+**
2. **PostgreSQL**
3. **Kafka**
4. **Redis** (optional, for caching)

## Setup Steps

### 1. Database Setup

```bash
# Create database
createdb godev_kit

# Run migrations
make migrate-up
```

### 2. Kafka Setup

```bash
# Start Kafka (if using Docker)
docker-compose up -d kafka

# Create payment-events topic
kafka-topics.sh --create \
  --topic payment-events \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1
```

### 3. Environment Configuration

Create `.env` file or set environment variables:

```bash
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/godev_kit?sslmode=disable

# Kafka
KAFKA_BROKERS=localhost:9092

# HTTP Server
HTTP_PORT=8080

# Log Level
LOG_LEVEL=info
```

### 4. Build and Run

```bash
# Build the application
go build -o bin/app cmd/app/main.go

# Run the application
./bin/app
```

## API Testing

### 1. Register Payment

```bash
curl -X POST http://localhost:8080/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "amount": 500000,
    "currency": "VND",
    "payment_type": "electric",
    "meter_number": "EVN001234567",
    "customer_code": "CUST001",
    "description": "Thanh toán tiền điện tháng 12/2024",
    "payment_method": "bank_transfer"
  }'
```

Expected Response:
```json
{
  "id": 1,
  "user_id": 1,
  "amount": 500000,
  "currency": "VND",
  "payment_type": "electric",
  "status": "pending",
  "meter_number": "EVN001234567",
  "customer_code": "CUST001",
  "description": "Thanh toán tiền điện tháng 12/2024",
  "transaction_id": "uuid-here",
  "payment_method": "bank_transfer",
  "created_at": "2024-12-20T10:30:00Z"
}
```

### 2. Get Payment by ID

```bash
curl http://localhost:8080/v1/payments/1
```

### 3. Get Payments by User ID

```bash
curl http://localhost:8080/v1/users/1/payments
```

## Monitoring

### 1. Check Logs

The application uses structured logging. Look for payment-related logs:

```json
{
  "level": "info",
  "payment_id": 1,
  "user_id": 1,
  "amount": 500000,
  "status": "pending",
  "message": "Payment registered successfully"
}
```

### 2. Check Kafka Messages

```bash
# Consume messages from payment-events topic
kafka-console-consumer.sh \
  --topic payment-events \
  --bootstrap-server localhost:9092 \
  --from-beginning
```

### 3. Check Database

```sql
-- Check payments table
SELECT * FROM payments ORDER BY created_at DESC;

-- Check payment history
SELECT * FROM payment_history ORDER BY created_at DESC;
```

## Troubleshooting

### 1. Database Connection Issues

```bash
# Check PostgreSQL connection
psql -h localhost -U postgres -d godev_kit

# Check if tables exist
\dt payments
\dt payment_history
```

### 2. Kafka Connection Issues

```bash
# Check Kafka brokers
kafka-broker-api-versions.sh --bootstrap-server localhost:9092

# Check topics
kafka-topics.sh --list --bootstrap-server localhost:9092
```

### 3. Application Issues

```bash
# Check application logs
tail -f logs/app.log

# Check if all dependencies are resolved
go mod tidy
```

## Development

### 1. Run Tests

```bash
# Run all tests
go test ./...

# Run payment tests specifically
go test ./internal/usecase/payment/...
```

### 2. Code Generation

```bash
# Generate swagger docs
swag init -g cmd/app/main.go

# Generate mocks (if using mockgen)
mockgen -source=internal/usecase/payment/payment.go -destination=internal/usecase/payment/mocks/payment.go
```

### 3. Linting

```bash
# Run linter
golangci-lint run

# Fix formatting
go fmt ./...
```

## Production Deployment

### 1. Environment Variables

```bash
# Production database
DATABASE_URL=postgres://user:pass@prod-db:5432/godev_kit?sslmode=require

# Production Kafka
KAFKA_BROKERS=prod-kafka-1:9092,prod-kafka-2:9092,prod-kafka-3:9092

# Production settings
HTTP_PORT=8080
LOG_LEVEL=info
ENVIRONMENT=production
```

### 2. Docker Deployment

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### 3. Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payment-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: payment-service
  template:
    metadata:
      labels:
        app: payment-service
    spec:
      containers:
      - name: payment-service
        image: payment-service:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
        - name: KAFKA_BROKERS
          value: "kafka-service:9092"
```

## Security Considerations

1. **Authentication**: Implement JWT authentication
2. **Authorization**: Add role-based access control
3. **Rate Limiting**: Implement rate limiting for API endpoints
4. **Input Validation**: Validate all input data
5. **Encryption**: Encrypt sensitive data at rest
6. **Audit Logging**: Log all payment operations
7. **HTTPS**: Use HTTPS in production

## Performance Optimization

1. **Database Indexing**: Ensure proper indexes on payment tables
2. **Caching**: Use Redis for caching frequently accessed data
3. **Connection Pooling**: Configure proper database connection pools
4. **Kafka Partitioning**: Use appropriate number of partitions
5. **Monitoring**: Implement metrics and monitoring
6. **Load Balancing**: Use load balancers for high availability 