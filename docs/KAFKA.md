# Kafka Integration

This document describes the Kafka integration in the GoDev Kit project using the [segmentio/kafka-go](https://github.com/segmentio/kafka-go) library.

## Overview

The Kafka integration provides:
- **Producer**: Send messages to Kafka topics
- **Consumer**: Consume messages from Kafka topics
- **Event-driven architecture**: Handle user and translation events
- **Repository pattern**: Clean separation of concerns

## Configuration

### Kafka Configuration

Add Kafka configuration to your `config/config.yaml`:

```yaml
KAFKA:
  BROKERS:
    - localhost:9092
  GROUP_ID: godev-kit-group
  TOPICS:
    USER_EVENTS: user-events
    TRANSLATION_EVENTS: translation-events
```

## Architecture

### Components

1. **Kafka Package** (`pkg/kafka/`)
   - `producer.go`: Kafka producer implementation
   - `consumer.go`: Kafka consumer implementation
   - `manager.go`: Kafka manager for coordinating producer and consumers

2. **Repository Layer** (`internal/repo/`)
   - `contracts.go`: Kafka repository interface
   - `persistent/kafka.go`: Kafka repository implementation

3. **Use Case Layer** (`internal/usecase/`)
   - `kafka_events.go`: Event handling use cases

4. **Entity Layer** (`internal/entity/`)
   - `events.go`: Event models and types

## Usage Examples

### Basic Producer Usage

```go
// Initialize Kafka repository
brokers := []string{"localhost:9092"}
kafkaRepo := persistent.NewKafkaRepo(brokers, logger)

// Send a message
ctx := context.Background()
err := kafkaRepo.SendMessage(ctx, "my-topic", []byte("key"), "Hello Kafka!")
```

### Event-Driven Usage

```go
// Initialize Kafka event use case
kafkaEventUseCase := usecase.NewKafkaEventUseCase(kafkaRepo, logger)

// Produce user events
err := kafkaEventUseCase.ProduceUserEvent(ctx, entity.UserCreatedEvent, 1, "user@example.com", data)

// Setup and start consumers
kafkaEventUseCase.ConsumeUserEvents(ctx)
kafkaRepo.StartAllConsumers(ctx)
```

## Event Types

### User Events
- `user.created`: When a new user is created
- `user.updated`: When a user is updated
- `user.deleted`: When a user is deleted

### Translation Events
- `translation.requested`: When a translation is requested
- `translation.completed`: When a translation is completed

## Testing

Run the Kafka demo to test the integration:

```bash
go run examples/kafka_demo.go
``` 