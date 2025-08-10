# Kafka Control - Quick Start Guide

## 🚀 Quick Start

### **1. Start Your Application**
```bash
go run cmd/app/main.go
```

### **2. Test Kafka Control**

#### **Check Status**
```bash
curl http://localhost:10000/api/v1/kafka/status
```

#### **Disable Producer**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable
```

#### **Try to Send Message (Should Fail)**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/producer/request \
  -H "Content-Type: application/json" \
  -d '{"topic":"test","key":"test-key","value":"test-message"}'
```

#### **Enable Producer**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/producer/enable
```

#### **Send Message (Should Work)**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/producer/request \
  -H "Content-Type: application/json" \
  -d '{"topic":"test","key":"test-key","value":"test-message"}'
```

### **3. Run Demo**
```bash
go test -v examples/kafka_control_test.go -run TestKafkaControl
```

## 📋 Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/kafka/status` | Get Kafka status |
| `POST` | `/api/v1/kafka/producer/enable` | Enable producer |
| `POST` | `/api/v1/kafka/producer/disable` | Disable producer |
| `POST` | `/api/v1/kafka/consumer/enable` | Enable consumer |
| `POST` | `/api/v1/kafka/consumer/disable` | Disable consumer |
| `POST` | `/api/v1/kafka/producer/request` | Send message |
| `GET` | `/api/v1/kafka/consumer/receiver` | Receive message |

## 🎯 Use Cases

### **Maintenance Mode**
```bash
# Disable both
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable

# Do maintenance...

# Re-enable both
curl -X POST http://localhost:10000/api/v1/kafka/producer/enable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/enable
```

### **Emergency Stop**
```bash
# Quick stop
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable

# Verify
curl http://localhost:10000/api/v1/kafka/status
```

## 🔧 Configuration

The Kafka control is enabled by default. Both producer and consumer start in **enabled** state.

## 📊 Status Response

```json
{
  "status": "success",
  "data": {
    "producer_enabled": true,
    "consumer_enabled": true,
    "consumer_count": 2,
    "brokers": ["localhost:9092"]
  }
}
```

## 🚨 Error Messages

- **Producer Disabled**: `"kafka producer is disabled"`
- **Consumer Disabled**: `"kafka consumer is disabled"`
- **Topic Not Found**: `"consumer for topic test-topic not found"`

## 🎉 That's It!

Your Kafka producer and consumer can now be controlled via HTTP endpoints. Perfect for:
- Maintenance operations
- Debugging
- Load testing
- Emergency situations 