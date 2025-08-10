package usecase

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/repo"
)

type Kafka interface {
	ProduceMessage(ctx context.Context, topic, key string, value interface{}) error
	ConsumeMessage(ctx context.Context, topic, group string) (string, []byte, error)

	// Control methods
	EnableProducer()
	DisableProducer()
	IsProducerEnabled() bool
	EnableConsumer()
	DisableConsumer()
	IsConsumerEnabled() bool
	GetStatus() map[string]interface{}
}

type kafkaUseCase struct {
	kafkaRepo repo.KafkaRepo
}

func NewKafkaUseCase(kafkaRepo repo.KafkaRepo) Kafka {
	return &kafkaUseCase{kafkaRepo: kafkaRepo}
}

func (u *kafkaUseCase) ProduceMessage(ctx context.Context, topic, key string, value interface{}) error {
	return u.kafkaRepo.SendMessage(ctx, topic, []byte(key), value)
}

// func (u *kafkaUseCase) ConsumeMessage(ctx context.Context, topic, group string) (string, []byte, error) {
// 	msgCh := make(chan struct {
// 		key   string
// 		value []byte
// 	}, 1)
// 	errCh := make(chan error, 1)

// 	// Add a temporary consumer
// 	err := u.kafkaRepo.AddConsumer(topic, group, func(ctx context.Context, key, value []byte) error {
// 		msgCh <- struct {
// 			key   string
// 			value []byte
// 		}{key: string(key), value: value}
// 		return nil
// 	})
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	// Start the consumer in a goroutine
// 	if err := u.kafkaRepo.StartConsumer(ctx, topic); err != nil {
// 		errCh <- err
// 	}

// 	select {
// 	case msg := <-msgCh:
// 		return msg.key, msg.value, nil
// 	case err := <-errCh:
// 		return "", nil, err
// 	case <-ctx.Done():
// 		return "", nil, ctx.Err()
// 	}
// }

func (u *kafkaUseCase) ConsumeMessage(ctx context.Context, topic, group string) (string, []byte, error) {
	msgCh := make(chan struct {
		key   string
		value []byte
	}, 1)
	errCh := make(chan error, 1)

	// Create a cancellable context so we can stop the consumer after one message
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Add a temporary consumer
	err := u.kafkaRepo.AddConsumer(topic, group, func(ctx context.Context, key, value []byte) error {
		select {
		case msgCh <- struct {
			key   string
			value []byte
		}{key: string(key), value: value}:
			cancel() // Stop the consumer after the first message
		default:
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}

	// Start the consumer in a goroutine
	go func() {
		if err := u.kafkaRepo.StartConsumer(ctx, topic); err != nil {
			errCh <- err
		}
	}()

	select {
	case msg := <-msgCh:
		return msg.key, msg.value, nil
	case err := <-errCh:
		return "", nil, err
	case <-ctx.Done():
		return "", nil, ctx.Err()
	}
}

// Control methods implementation

// EnableProducer enables the Kafka producer
func (u *kafkaUseCase) EnableProducer() {
	u.kafkaRepo.EnableProducer()
}

// DisableProducer disables the Kafka producer
func (u *kafkaUseCase) DisableProducer() {
	u.kafkaRepo.DisableProducer()
}

// IsProducerEnabled returns the current producer status
func (u *kafkaUseCase) IsProducerEnabled() bool {
	return u.kafkaRepo.IsProducerEnabled()
}

// EnableConsumer enables the Kafka consumer
func (u *kafkaUseCase) EnableConsumer() {
	u.kafkaRepo.EnableConsumer()
}

// DisableConsumer disables the Kafka consumer
func (u *kafkaUseCase) DisableConsumer() {
	u.kafkaRepo.DisableConsumer()
}

// IsConsumerEnabled returns the current consumer status
func (u *kafkaUseCase) IsConsumerEnabled() bool {
	return u.kafkaRepo.IsConsumerEnabled()
}

// GetStatus returns the current status of both producer and consumer
func (u *kafkaUseCase) GetStatus() map[string]interface{} {
	return u.kafkaRepo.GetStatus()
}
