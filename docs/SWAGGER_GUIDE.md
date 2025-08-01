# Swagger UI Guide

## Overview

Swagger UI đã được cập nhật với đầy đủ payment endpoints. Bạn có thể truy cập Swagger UI để xem và test các API endpoints.

## Access Swagger UI

### 1. Start the Application

```bash
# Run the application
go run cmd/app/main.go
```

### 2. Access Swagger UI

Mở trình duyệt và truy cập:
```
http://localhost:10000/swagger/index.html
```

## Payment Endpoints in Swagger

### 1. Register Payment
- **Endpoint**: `POST /v1/payments`
- **Description**: Register a new payment for electric bill and send to Kafka for processing
- **Tags**: payments
- **Request Body**: PaymentRequest
- **Response**: PaymentResponse (201 Created)

### 2. Get Payment by ID
- **Endpoint**: `GET /v1/payments/{id}`
- **Description**: Get payment details by ID
- **Tags**: payments
- **Parameters**: id (integer, required)
- **Response**: PaymentResponse (200 OK)

### 3. Get Payments by User ID
- **Endpoint**: `GET /v1/users/{user_id}/payments`
- **Description**: Get all payments for a specific user
- **Tags**: payments
- **Parameters**: user_id (integer, required)
- **Response**: Array of PaymentResponse (200 OK)

## Request/Response Models

### PaymentRequest
```json
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

### PaymentResponse
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

## Testing with Swagger UI

### 1. Register a Payment

1. Mở Swagger UI
2. Tìm endpoint `POST /v1/payments`
3. Click "Try it out"
4. Nhập request body:
```json
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
5. Click "Execute"
6. Xem response với payment ID và transaction ID

### 2. Get Payment by ID

1. Tìm endpoint `GET /v1/payments/{id}`
2. Click "Try it out"
3. Nhập payment ID (ví dụ: 1)
4. Click "Execute"
5. Xem payment details

### 3. Get Payments by User ID

1. Tìm endpoint `GET /v1/users/{user_id}/payments`
2. Click "Try it out"
3. Nhập user ID (ví dụ: 1)
4. Click "Execute"
5. Xem danh sách payments của user

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request",
  "message": "validation error details"
}
```

### 404 Not Found
```json
{
  "error": "Payment not found",
  "message": "Payment with the specified ID was not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "message": "error details"
}
```

## Regenerate Swagger Docs

Nếu bạn thay đổi code và muốn cập nhật Swagger docs:

```bash
# Regenerate swagger docs
swag init -g cmd/app/main.go

# Hoặc sử dụng make command
make swag-v1
```

## Swagger Configuration

### Main API Info
- **Title**: Go Dev Kit Template API
- **Version**: 1.0
- **Description**: A comprehensive API template with translation, user management, Kafka, Redis, NATS, VietQR, and Payment services
- **Host**: localhost:10000
- **Base Path**: /

### Security
- **BearerAuth**: JWT token authentication
- **Header**: Authorization
- **Format**: "Bearer {token}"

## Other Available Endpoints

Swagger UI cũng hiển thị các endpoints khác:

### Translation
- `POST /translation/do-translate` - Translate text
- `GET /translation/history` - Get translation history

### User Management
- `GET /v1/user` - List users
- `POST /v1/user` - Create user
- `GET /v1/user/{id}` - Get user by ID
- `PUT /v1/user/{id}` - Update user
- `DELETE /v1/user/{id}` - Delete user

### Kafka
- `POST /v1/kafka/producer/request` - Send Kafka message
- `GET /v1/kafka/consumer/receiver` - Receive Kafka message

### Redis
- `POST /v1/redis/set` - Set Redis value
- `GET /v1/redis/get/{key}` - Get Redis value
- `POST /v1/redis/shipper/location` - Update shipper location
- `GET /v1/redis/shipper/location/{shipper_id}` - Get shipper location

### NATS
- `POST /v1/nats/publish/{subject}` - Publish NATS message
- `GET /v1/nats/subscribe/{subject}` - Subscribe to NATS subject

### VietQR
- `POST /v1/vietqr/gen` - Generate VietQR code
- `GET /v1/vietqr/inquiry/{id}` - Inquiry QR status
- `PUT /v1/vietqr/update/{id}` - Update QR status

### Billing
- `POST /v1/billing/invoice` - Generate invoice PDF

## Tips for Using Swagger UI

1. **Authentication**: Nếu endpoint yêu cầu authentication, click "Authorize" button và nhập Bearer token
2. **Request Body**: Sử dụng example values được cung cấp
3. **Response**: Xem response schema để hiểu cấu trúc data
4. **Error Handling**: Test các error cases để hiểu error responses
5. **Try it out**: Luôn sử dụng "Try it out" để test thực tế

## Troubleshooting

### Swagger UI không load
- Kiểm tra application đã start chưa
- Kiểm tra port 10000 có đúng không
- Kiểm tra Swagger được enable trong config

### Endpoints không hiển thị
- Regenerate swagger docs: `swag init -g cmd/app/main.go`
- Kiểm tra annotations trong code
- Restart application

### Authentication issues
- Kiểm tra Bearer token format
- Đảm bảo token còn hiệu lực
- Kiểm tra JWT secret configuration 