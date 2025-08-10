# Kafka Configuration from YAML

## 🎯 Overview

Your Kafka producer and consumer can now be configured directly from YAML configuration files, allowing you to:

- **Set Initial State**: Configure producer/consumer enabled/disabled state at startup
- **Environment-Specific Settings**: Different configs for dev/staging/production
- **Runtime Status Checking**: Check current on/off status via HTTP endpoints
- **Dynamic Control**: Override config settings at runtime via API

## 📁 Configuration Files

### **1. YAML Configuration**

#### **`config/config.yaml`**
```yaml
KAFKA:
  BROKERS:
    - localhost:9092
  GROUP_ID: godev-kit-group
  TOPICS:
    USER_EVENTS: dev.user-events
    TRANSLATION_EVENTS: dev.translation-events
  CONTROL:
    PRODUCER_ENABLED: true   # Enable/disable Kafka producer
    CONSUMER_ENABLED: true   # Enable/disable Kafka consumer
```

#### **Environment Variables Override**
```bash
# Override YAML settings with environment variables
export KAFKA_CONTROL_PRODUCER_ENABLED=false
export KAFKA_CONTROL_CONSUMER_ENABLED=false
```

### **2. Configuration Structure**

#### **`config/config.go`**
```go
// Kafka configuration structure
type Kafka struct {
    Brokers []string `mapstructure:"BROKERS"`
    GroupID string   `mapstructure:"GROUP_ID"`
    Topics  Topics   `mapstructure:"TOPICS"`
    Control Control  `mapstructure:"CONTROL"`
}

// Control settings
type Control struct {
    ProducerEnabled bool `mapstructure:"PRODUCER_ENABLED"`
    ConsumerEnabled bool `mapstructure:"CONSUMER_ENABLED"`
}
```

## 🔧 Usage Examples

### **1. Check Initial Configuration**

#### **Start Application with Config**
```bash
# Start with default config (both enabled)
go run cmd/app/main.go

# Start with custom config
KAFKA_CONTROL_PRODUCER_ENABLED=false go run cmd/app/main.go
```

#### **Check Status from YAML Config**
```bash
# Check overall status
curl http://localhost:10000/api/v1/kafka/status

# Check producer status specifically
curl http://localhost:10000/api/v1/kafka/producer/status

# Check consumer status specifically
curl http://localhost:10000/api/v1/kafka/consumer/status
```

### **2. Configuration Scenarios**

#### **Development Environment**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: true   # Enable for testing
    CONSUMER_ENABLED: true   # Enable for testing
```

#### **Production Environment**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: true   # Enable for production
    CONSUMER_ENABLED: false  # Disable during deployment
```

#### **Maintenance Mode**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: false  # Disable during maintenance
    CONSUMER_ENABLED: false  # Disable during maintenance
```

### **3. Runtime Status Checking**

#### **Check Producer Status**
```bash
curl http://localhost:10000/api/v1/kafka/producer/status
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "producer_enabled": true,
    "message": "Kafka producer is enabled and can send messages"
  }
}
```

#### **Check Consumer Status**
```bash
curl http://localhost:10000/api/v1/kafka/consumer/status
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "consumer_enabled": false,
    "message": "Kafka consumer is disabled and cannot receive messages"
  }
}
```

#### **Check Overall Status**
```bash
curl http://localhost:10000/api/v1/kafka/status
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

## 🚀 API Endpoints

### **Status Check Endpoints**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/kafka/status` | Overall Kafka status |
| `GET` | `/api/v1/kafka/producer/status` | Producer status only |
| `GET` | `/api/v1/kafka/consumer/status` | Consumer status only |

### **Control Endpoints**

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/kafka/producer/enable` | Enable producer |
| `POST` | `/api/v1/kafka/producer/disable` | Disable producer |
| `POST` | `/api/v1/kafka/consumer/enable` | Enable consumer |
| `POST` | `/api/v1/kafka/consumer/disable` | Disable consumer |

## 🔍 Implementation Details

### **1. Configuration Loading**

#### **App Initialization**
```go
// In internal/app/app.go
kafkaRepo := persistent.NewKafkaRepoWithConfig(
    cfg.Kafka.Brokers, 
    l.Zerolog(), 
    cfg.Kafka.Control.ProducerEnabled,  // From YAML
    cfg.Kafka.Control.ConsumerEnabled,  // From YAML
)
```

#### **Manager Creation**
```go
// In pkg/kafka/manager.go
func NewManagerWithConfig(brokers []string, logger zerolog.Logger, producerEnabled, consumerEnabled bool) *Manager {
    return &Manager{
        producer:        NewProducer(brokers, logger),
        consumers:       make(map[string]*Consumer),
        logger:          logger,
        brokers:         brokers,
        producerEnabled: producerEnabled,  // From config
        consumerEnabled: consumerEnabled,  // From config
    }
}
```

### **2. Configuration Validation**

#### **Validation Function**
```go
func (c *Config) validateKafkaConfig() error {
    if len(c.Kafka.Brokers) == 0 {
        return errors.New("kafka brokers are required")
    }
    
    if c.Kafka.GroupID == "" {
        return errors.New("kafka group ID is required")
    }
    
    // Log Kafka control settings
    log.Printf("Kafka Control Settings:")
    log.Printf("  Producer Enabled: %v", c.Kafka.Control.ProducerEnabled)
    log.Printf("  Consumer Enabled: %v", c.Kafka.Control.ConsumerEnabled)
    
    return nil
}
```

### **3. Status Checking**

#### **Producer Status Check**
```go
func (h *V1) CheckProducerStatus(c *gin.Context) {
    isEnabled := h.kafka.IsProducerEnabled()
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data": gin.H{
            "producer_enabled": isEnabled,
            "message":          getProducerStatusMessage(isEnabled),
        },
    })
}
```

## 🎯 Use Cases

### **1. Environment-Specific Configuration**

#### **Development**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: true
    CONSUMER_ENABLED: true
```

#### **Staging**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: true
    CONSUMER_ENABLED: false  # Disable to test producer only
```

#### **Production**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: true
    CONSUMER_ENABLED: true
```

### **2. CI/CD Integration**

#### **Pre-Deployment Script**
```bash
#!/bin/bash
# Disable consumers before deployment
curl -X POST http://localhost:10000/api/v1/kafka/consumer/disable

# Verify status
curl http://localhost:10000/api/v1/kafka/consumer/status

# Deploy application
# ...

# Re-enable consumers after deployment
curl -X POST http://localhost:10000/api/v1/kafka/consumer/enable
```

### **3. Monitoring and Alerting**

#### **Health Check Script**
```bash
#!/bin/bash
# Check Kafka status
status=$(curl -s http://localhost:10000/api/v1/kafka/status | jq -r '.data.producer_enabled')

if [ "$status" = "false" ]; then
    echo "ALERT: Kafka producer is disabled!"
    exit 1
fi

echo "Kafka producer is enabled"
```

### **4. Testing Scenarios**

#### **Producer-Only Testing**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: true
    CONSUMER_ENABLED: false
```

#### **Consumer-Only Testing**
```yaml
KAFKA:
  CONTROL:
    PRODUCER_ENABLED: false
    CONSUMER_ENABLED: true
```

## 🔧 Testing Commands

### **1. Test Configuration Loading**
```bash
# Start with different configs
KAFKA_CONTROL_PRODUCER_ENABLED=false go run cmd/app/main.go

# Check status
curl http://localhost:10000/api/v1/kafka/producer/status
```

### **2. Test Status Endpoints**
```bash
# Check all status endpoints
curl http://localhost:10000/api/v1/kafka/status
curl http://localhost:10000/api/v1/kafka/producer/status
curl http://localhost:10000/api/v1/kafka/consumer/status
```

### **3. Test Dynamic Control**
```bash
# Disable producer
curl -X POST http://localhost:10000/api/v1/kafka/producer/disable

# Check status
curl http://localhost:10000/api/v1/kafka/producer/status

# Try to send message (should fail)
curl -X POST http://localhost:10000/api/v1/kafka/producer/request \
  -H "Content-Type: application/json" \
  -d '{"topic":"test","key":"test","value":"test"}'

# Re-enable
curl -X POST http://localhost:10000/api/v1/kafka/producer/enable
```

## 🚨 Error Handling

### **Configuration Errors**
```bash
# Missing brokers
Error: kafka brokers are required

# Missing group ID
Error: kafka group ID is required
```

### **Runtime Errors**
```bash
# Producer disabled
{"error": "kafka producer is disabled"}

# Consumer disabled
{"error": "kafka consumer is disabled"}
```

## 📊 Logging

### **Startup Logs**
```
Kafka Control Settings:
  Producer Enabled: true
  Consumer Enabled: false
```

### **Runtime Logs**
```
INFO kafka producer enabled
INFO kafka consumer disabled
WARN kafka consumer is disabled, skipping start all consumers
```

## 🎉 Benefits

1. **Environment Flexibility**: Different configs for different environments
2. **Runtime Control**: Override config settings via API
3. **Status Visibility**: Clear status checking endpoints
4. **Validation**: Configuration validation at startup
5. **Logging**: Comprehensive logging of control settings
6. **CI/CD Ready**: Easy integration with deployment pipelines

Your Kafka configuration is now fully manageable through YAML files with runtime status checking! 🚀 