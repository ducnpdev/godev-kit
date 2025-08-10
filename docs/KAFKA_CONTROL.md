# Kafka Producer & Consumer Control

## 🎯 Overview

This implementation provides on/off control functionality for Kafka producers and consumers through HTTP endpoints. This allows you to:

- **Enable/Disable Kafka Producer**: Control message sending capability
- **Enable/Disable Kafka Consumer**: Control message receiving capability  
- **Monitor Status**: Get real-time status of both producer and consumer
- **Graceful Control**: Safe enable/disable without losing connections

## 🚀 API Endpoints

### **Producer Control**

#### **Enable Producer**
```bash
POST /api/v1/kafka/producer/enable
```
**Response:**
```json
{
  "status": "producer enabled"
}
```

#### **Disable Producer**
```bash
POST /api/v1/kafka/producer/disable
```
**Response:**
```json
{
  "status": "producer disabled"
}
```

### **Consumer Control**

#### **Enable Consumer**
```bash
POST /api/v1/kafka/consumer/enable
```
**Response:**
```json
{
  "status": "consumer enabled"
}
```

#### **Disable Consumer**
```bash
POST /api/v1/kafka/consumer/disable
```
**Response:**
```json
{
  "status": "consumer disabled"
}
```

### **Status Monitoring**

#### **Get Kafka Status**
```bash
GET /api/v1/kafka/status
```
**Response:**
```json
{
  "status": "success",
  "data": {
    "producer_enabled": true,
    "consumer_enabled": false,
    "consumer_count": 2,
    "brokers": ["localhost:9092"]
  }
}
```

## 🔧 Usage Examples

### **1. Disable Producer (Stop Sending Messages)**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable
```

**When disabled:**
- `POST /api/v1/kafka/producer/request` will return error: `"kafka producer is disabled"`
- No messages will be sent to Kafka topics
- Existing connections remain intact

### **2. Enable Producer (Resume Sending Messages)**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/producer/enable
```

**When enabled:**
- Producer can send messages normally
- All existing functionality restored

### **3. Disable Consumer (Stop Receiving Messages)**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable
```

**When disabled:**
- `GET /api/v1/kafka/consumer/receiver` will return error: `"kafka consumer is disabled"`
- No new messages will be consumed
- Running consumers continue but new ones won't start

### **4. Enable Consumer (Resume Receiving Messages)**
```bash
curl -X POST http://localhost:10000/api/v1/kafka/consumer/enable
```

**When enabled:**
- Consumers can receive messages normally
- All existing functionality restored

### **5. Check Current Status**
```bash
curl http://localhost:10000/api/v1/kafka/status
```

## 🛡️ Safety Features

### **Thread-Safe Operations**
- All control operations are protected by mutex locks
- No race conditions during enable/disable operations
- Safe concurrent access from multiple HTTP requests

### **Graceful State Management**
- Producer/Consumer state is maintained independently
- Enabling one doesn't affect the other
- State persists across HTTP requests

### **Error Handling**
- Clear error messages when operations are blocked
- Status endpoint shows current state
- Logging for all state changes

## 🔍 Implementation Details

### **Manager Level Control**
```go
// Enable/Disable with thread safety
func (m *Manager) EnableProducer() {
    m.controlMu.Lock()
    defer m.controlMu.Unlock()
    m.producerEnabled = true
    m.logger.Info().Msg("kafka producer enabled")
}
```

### **Request-Level Checks**
```go
func (m *Manager) SendMessage(ctx context.Context, topic string, key []byte, value interface{}) error {
    m.controlMu.RLock()
    defer m.controlMu.RUnlock()
    
    if !m.producerEnabled {
        return fmt.Errorf("kafka producer is disabled")
    }
    
    return m.producer.SendMessage(ctx, topic, key, value)
}
```

### **Status Reporting**
```go
func (m *Manager) GetStatus() map[string]interface{} {
    return map[string]interface{}{
        "producer_enabled":  m.producerEnabled,
        "consumer_enabled":  m.consumerEnabled,
        "consumer_count":    len(m.consumers),
        "brokers":           m.brokers,
    }
}
```

## 🎯 Use Cases

### **1. Maintenance Mode**
```bash
# Disable both during maintenance
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable

# Resume after maintenance
curl -X POST http://localhost:10000/api/v1/kafka/producer/enable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/enable
```

### **2. Debugging**
```bash
# Disable producer to test consumer only
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable

# Check status
curl http://localhost:10000/api/v1/kafka/status
```

### **3. Load Testing**
```bash
# Disable consumer to test producer performance
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable

# Send test messages
for i in {1..1000}; do
  curl -X POST http://localhost:10000/api/v1/kafka/producer/request \
    -H "Content-Type: application/json" \
    -d '{"topic":"test","key":"key'$i'","value":"message'$i'"}'
done
```

### **4. Emergency Stop**
```bash
# Quick emergency stop
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable

# Verify stopped
curl http://localhost:10000/api/v1/kafka/status
```

## 🔧 Configuration

### **Default State**
- **Producer**: Enabled by default
- **Consumer**: Enabled by default
- **Thread Safety**: Enabled
- **Logging**: All state changes logged

### **Environment Variables**
```yaml
KAFKA:
  BROKERS:
    - localhost:9092
  GROUP_ID: godev-kit-group
  # Control can be managed via API endpoints
```

## 🚨 Error Scenarios

### **Producer Disabled**
```json
{
  "error": "kafka producer is disabled"
}
```

### **Consumer Disabled**
```json
{
  "error": "kafka consumer is disabled"
}
```

### **Invalid Topic**
```json
{
  "error": "consumer for topic test-topic not found"
}
```

## 📊 Monitoring

### **Status Endpoint Response**
```json
{
  "status": "success",
  "data": {
    "producer_enabled": true,
    "consumer_enabled": true,
    "consumer_count": 3,
    "brokers": ["localhost:9092", "localhost:9093"]
  }
}
```

### **Log Messages**
```
INFO kafka producer enabled
INFO kafka consumer disabled
WARN kafka consumer is disabled, skipping start all consumers
```

## 🔄 State Transitions

### **Producer States**
```
ENABLED → DISABLED: POST /api/v1/kafka/producer/disable
DISABLED → ENABLED: POST /api/v1/kafka/producer/enable
```

### **Consumer States**
```
ENABLED → DISABLED: POST /api/v1/kafka/consumer/disable
DISABLED → ENABLED: POST /api/v1/kafka/consumer/enable
```

## 🎯 Best Practices

1. **Check Status Before Operations**: Always verify current state
2. **Use in Pairs**: Enable/disable both for maintenance
3. **Monitor Logs**: Watch for state change confirmations
4. **Test in Development**: Verify behavior before production use
5. **Document State**: Keep track of current Kafka state

## 🔧 Integration

### **With CI/CD**
```bash
# Pre-deployment: Disable Kafka
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable

# Deploy application
# ...

# Post-deployment: Enable Kafka
curl -X POST http://localhost:10000/api/v1/kafka/producer/enable
curl -X POST http://localhost:10000/api/v1/kafka/consumer/enable
```

### **With Monitoring**
```bash
# Health check script
status=$(curl -s http://localhost:10000/api/v1/kafka/status | jq -r '.data.producer_enabled')
if [ "$status" = "false" ]; then
    echo "Kafka producer is disabled!"
    exit 1
fi
``` 