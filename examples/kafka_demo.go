package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent"
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/pkg/logger"
)

func main() {
	// Initialize logger
	l := logger.New("debug")

	// Initialize Kafka repository
	brokers := []string{"localhost:9092"}
	kafkaRepo := persistent.NewKafkaRepo(brokers, l.Zerolog())
	defer kafkaRepo.Close()

	// Initialize Kafka event use case
	kafkaEventUseCase := usecase.NewKafkaEventUseCase(kafkaRepo, l.Zerolog())

	// Setup consumers
	ctx := context.Background()
	if err := kafkaEventUseCase.ConsumeUserEvents(ctx); err != nil {
		log.Fatalf("Failed to setup user events consumer: %v", err)
	}

	if err := kafkaEventUseCase.ConsumeTranslationEvents(ctx); err != nil {
		log.Fatalf("Failed to setup translation events consumer: %v", err)
	}

	// Start consumers in background
	go kafkaRepo.StartAllConsumers(ctx)

	// Wait a bit for consumers to start
	time.Sleep(2 * time.Second)

	// Example 1: Produce user events
	fmt.Println("=== Producing User Events ===")

	// User created event
	err := kafkaEventUseCase.ProduceUserEvent(ctx, entity.UserCreatedEvent, 1, "john@example.com", map[string]interface{}{
		"name": "John Doe",
		"age":  30,
	})
	if err != nil {
		log.Printf("Failed to produce user created event: %v", err)
	}

	// User updated event
	err = kafkaEventUseCase.ProduceUserEvent(ctx, entity.UserUpdatedEvent, 1, "john@example.com", map[string]interface{}{
		"name": "John Smith",
		"age":  31,
	})
	if err != nil {
		log.Printf("Failed to produce user updated event: %v", err)
	}

	// User deleted event
	err = kafkaEventUseCase.ProduceUserEvent(ctx, entity.UserDeletedEvent, 1, "john@example.com", nil)
	if err != nil {
		log.Printf("Failed to produce user deleted event: %v", err)
	}

	// Example 2: Produce translation events
	fmt.Println("=== Producing Translation Events ===")

	// Translation requested event
	err = kafkaEventUseCase.ProduceTranslationEvent(ctx, entity.TranslationRequestEvent, 1, "en", "es", "Hello world", "")
	if err != nil {
		log.Printf("Failed to produce translation requested event: %v", err)
	}

	// Translation completed event
	err = kafkaEventUseCase.ProduceTranslationEvent(ctx, entity.TranslationCompletedEvent, 1, "en", "es", "Hello world", "Hola mundo")
	if err != nil {
		log.Printf("Failed to produce translation completed event: %v", err)
	}

	// Example 3: Direct Kafka repository usage
	fmt.Println("=== Direct Kafka Repository Usage ===")

	// Send a custom message
	customEvent := map[string]interface{}{
		"type":    "custom_event",
		"message": "This is a custom event",
		"time":    time.Now(),
	}

	err = kafkaRepo.SendMessage(ctx, "custom-topic", []byte("custom-key"), customEvent)
	if err != nil {
		log.Printf("Failed to send custom message: %v", err)
	}

	// Wait for some time to see the events being processed
	fmt.Println("Waiting for events to be processed...")
	time.Sleep(5 * time.Second)

	fmt.Println("Demo completed!")
}
