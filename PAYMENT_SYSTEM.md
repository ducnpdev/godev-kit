# Auto Payment System - Thanh toán điện tự động

Hệ thống auto payment cho thanh toán điện được xây dựng với kiến trúc microservices, sử dụng Kafka để xử lý bất đồng bộ.

## Kiến trúc hệ thống

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │    │   Kafka         │    │   Consumer      │
│   (Register)    │───▶│   (Events)      │───▶│   (Process)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Database      │    │   Database      │    │   Database      │
│   (Insert)      │    │   (Update)      │    │   (History)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Các thành phần chính

### 1. Entity (`internal/entity/payment.go`)
- `Payment`: Entity chính cho payment
- `PaymentEvent`: Event cho Kafka
- `PaymentRequest`: Request từ API
- `PaymentResponse`: Response cho API

### 2. Repository (`internal/repo/persistent/payment_postgres.go`)
- `PaymentRepo`: Repository cho payment
- Các method: Create, GetByID, GetByUserID, UpdateStatus, CreateHistory

### 3. Use Case (`internal/usecase/payment/`)
- `PaymentUseCase`: Business logic cho payment
- `PaymentConsumer`: Kafka consumer cho payment processing

### 4. Controller (`internal/controller/http/v1/payment.go`)
- `PaymentController`: HTTP controller cho payment API

## API Endpoints

### 1. Register Payment
```http
POST /api/v1/payments
Content-Type: application/json

{
  "user_id": 1,
  "amount": 500000,
  "currency": "VND",
  "payment_type": "electric",
  "meter_number": "EVN001234567",
  "customer_code": "CUST001",
  "description": "Thanh toán tiền điện tháng 12/2024",
  "payment_method": "bank_transfer"
}
```

Response:
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
```http
GET /api/v1/payments/{id}
```

### 3. Get Payments by User ID
```http
GET /api/v1/users/{user_id}/payments
```

## Luồng xử lý

### 1. Register Payment
1. Client gọi API `POST /payments`
2. Controller validate request
3. Use case tạo payment entity với status "pending"
4. Lưu vào database
5. Tạo PaymentEvent và gửi đến Kafka topic "payment-events"
6. Trả về response với payment ID

### 2. Process Payment (Kafka Consumer)
1. Consumer nhận message từ Kafka topic "payment-events"
2. Parse PaymentEvent từ JSON
3. Update status thành "processing"
4. Simulate payment processing (trong thực tế sẽ gọi payment gateway)
5. Update status thành "completed" hoặc "failed"
6. Tạo payment history record

## Database Schema

### Payments Table
```sql
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'VND',
    payment_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    meter_number VARCHAR(50) NOT NULL,
    customer_code VARCHAR(50) NOT NULL,
    description TEXT,
    transaction_id VARCHAR(100) UNIQUE NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Payment History Table
```sql
CREATE TABLE payment_history (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    payment_type VARCHAR(20) NOT NULL,
    meter_number VARCHAR(50) NOT NULL,
    customer_code VARCHAR(50) NOT NULL,
    description TEXT,
    transaction_id VARCHAR(100) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Payment Status

- `pending`: Chờ xử lý
- `processing`: Đang xử lý
- `completed`: Hoàn thành
- `failed`: Thất bại
- `cancelled`: Đã hủy

## Payment Types

- `electric`: Thanh toán điện
- `water`: Thanh toán nước
- `gas`: Thanh toán gas

## Cấu hình

### Kafka Configuration
```yaml
kafka:
  brokers:
    - localhost:9092
  topics:
    payment-events: payment-events
  consumer:
    group-id: payment-processor
```

### Database Configuration
```yaml
database:
  host: localhost
  port: 5432
  name: godev_kit
  user: postgres
  password: password
```

## Chạy hệ thống

### 1. Setup Database
```bash
# Chạy migration
make migrate-up
```

### 2. Setup Kafka
```bash
# Tạo topic
kafka-topics.sh --create --topic payment-events --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
```

### 3. Chạy ứng dụng
```bash
# Chạy main application
go run main.go
```

### 4. Test API
```bash
# Register payment
curl -X POST http://localhost:8080/api/v1/payments \
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

# Get payment by ID
curl http://localhost:8080/api/v1/payments/1

# Get payments by user ID
curl http://localhost:8080/api/v1/users/1/payments
```

## Monitoring

### Logs
Hệ thống sử dụng structured logging với zerolog:
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

### Metrics
- Số lượng payment được tạo
- Số lượng payment được xử lý thành công/thất bại
- Thời gian xử lý payment
- Kafka lag

## Security

- Validate input data
- Rate limiting
- Authentication/Authorization (cần implement)
- Audit logging
- Encryption cho sensitive data

## Scaling

- Horizontal scaling cho API servers
- Multiple Kafka consumers
- Database read replicas
- Caching với Redis
- Load balancing

## Error Handling

- Retry mechanism cho Kafka messages
- Dead letter queue cho failed messages
- Circuit breaker cho external services
- Graceful degradation 