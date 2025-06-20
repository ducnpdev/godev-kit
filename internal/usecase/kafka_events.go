package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo"
	"github.com/rs/zerolog"
)

// KafkaEventUseCase -.
type KafkaEventUseCase struct {
	kafkaRepo repo.KafkaRepo
	logger    zerolog.Logger
}

// NewKafkaEventUseCase -.
func NewKafkaEventUseCase(kafkaRepo repo.KafkaRepo, logger zerolog.Logger) *KafkaEventUseCase {
	return &KafkaEventUseCase{
		kafkaRepo: kafkaRepo,
		logger:    logger,
	}
}

// ProduceUserEvent -.
func (k *KafkaEventUseCase) ProduceUserEvent(ctx context.Context, eventType string, userID int64, email string, data any) error {
	event := entity.UserEvent{
		ID:        time.Now().UnixNano(),
		EventType: eventType,
		UserID:    userID,
		Email:     email,
		Data:      data,
		Timestamp: time.Now(),
	}

	key := []byte(strconv.FormatInt(userID, 10))

	err := k.kafkaRepo.SendMessage(ctx, "user-events", key, event)
	if err != nil {
		return fmt.Errorf("failed to send user event: %w", err)
	}

	k.logger.Info().
		Str("event_type", eventType).
		Int64("user_id", userID).
		Msg("user event produced")

	return nil
}

// ProduceTranslationEvent -.
func (k *KafkaEventUseCase) ProduceTranslationEvent(ctx context.Context, eventType string, userID int64, source, target, original, translated string) error {
	event := entity.TranslationEvent{
		ID:         time.Now().UnixNano(),
		EventType:  eventType,
		UserID:     userID,
		Source:     source,
		Target:     target,
		Original:   original,
		Translated: translated,
		Timestamp:  time.Now(),
	}

	key := []byte(fmt.Sprintf("%d-%s", userID, eventType))

	err := k.kafkaRepo.SendMessage(ctx, "translation-events", key, event)
	if err != nil {
		return fmt.Errorf("failed to send translation event: %w", err)
	}

	k.logger.Info().
		Str("event_type", eventType).
		Int64("user_id", userID).
		Str("source", source).
		Str("target", target).
		Msg("translation event produced")

	return nil
}

// ConsumeUserEvents -.
func (k *KafkaEventUseCase) ConsumeUserEvents(ctx context.Context) error {
	handler := func(ctx context.Context, key, value []byte) error {
		var event entity.UserEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user event: %w", err)
		}

		k.logger.Info().
			Str("event_type", event.EventType).
			Int64("user_id", event.UserID).
			Str("email", event.Email).
			Msg("user event consumed")

		// Here you can add business logic to handle different event types
		switch event.EventType {
		case entity.UserCreatedEvent:
			k.handleUserCreated(ctx, event)
		case entity.UserUpdatedEvent:
			k.handleUserUpdated(ctx, event)
		case entity.UserDeletedEvent:
			k.handleUserDeleted(ctx, event)
		default:
			k.logger.Warn().Str("event_type", event.EventType).Msg("unknown user event type")
		}

		return nil
	}

	return k.kafkaRepo.AddConsumer("user-events", "user-events-consumer", handler)
}

// ConsumeTranslationEvents -.
func (k *KafkaEventUseCase) ConsumeTranslationEvents(ctx context.Context) error {
	handler := func(ctx context.Context, key, value []byte) error {
		var event entity.TranslationEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal translation event: %w", err)
		}

		k.logger.Info().
			Str("event_type", event.EventType).
			Int64("user_id", event.UserID).
			Str("source", event.Source).
			Str("target", event.Target).
			Msg("translation event consumed")

		// Here you can add business logic to handle different event types
		switch event.EventType {
		case entity.TranslationRequestEvent:
			k.handleTranslationRequested(ctx, event)
		case entity.TranslationCompletedEvent:
			k.handleTranslationCompleted(ctx, event)
		default:
			k.logger.Warn().Str("event_type", event.EventType).Msg("unknown translation event type")
		}

		return nil
	}

	return k.kafkaRepo.AddConsumer("translation-events", "translation-events-consumer", handler)
}

// handleUserCreated -.
func (k *KafkaEventUseCase) handleUserCreated(ctx context.Context, event entity.UserEvent) {
	k.logger.Info().
		Int64("user_id", event.UserID).
		Str("email", event.Email).
		Msg("processing user created event")

	// Add your business logic here
	// For example: send welcome email, create user profile, etc.
}

// handleUserUpdated -.
func (k *KafkaEventUseCase) handleUserUpdated(ctx context.Context, event entity.UserEvent) {
	k.logger.Info().
		Int64("user_id", event.UserID).
		Str("email", event.Email).
		Msg("processing user updated event")

	// Add your business logic here
	// For example: update cache, notify other services, etc.
}

// handleUserDeleted -.
func (k *KafkaEventUseCase) handleUserDeleted(ctx context.Context, event entity.UserEvent) {
	k.logger.Info().
		Int64("user_id", event.UserID).
		Str("email", event.Email).
		Msg("processing user deleted event")

	// Add your business logic here
	// For example: cleanup user data, revoke tokens, etc.
}

// handleTranslationRequested -.
func (k *KafkaEventUseCase) handleTranslationRequested(ctx context.Context, event entity.TranslationEvent) {
	k.logger.Info().
		Int64("user_id", event.UserID).
		Str("source", event.Source).
		Str("target", event.Target).
		Str("original", event.Original).
		Msg("processing translation requested event")

	// Add your business logic here
	// For example: track translation requests, update metrics, etc.
}

// handleTranslationCompleted -.
func (k *KafkaEventUseCase) handleTranslationCompleted(ctx context.Context, event entity.TranslationEvent) {
	k.logger.Info().
		Int64("user_id", event.UserID).
		Str("source", event.Source).
		Str("target", event.Target).
		Str("translated", event.Translated).
		Msg("processing translation completed event")

	// Add your business logic here
	// For example: update translation history, send notifications, etc.
}
