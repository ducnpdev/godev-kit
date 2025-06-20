package usecase

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/repo"
)

type Kafka interface {
	ProduceMessage(ctx context.Context, topic, key string, value interface{}) error
	ConsumeMessage(ctx context.Context, topic, group string) (string, []byte, error)
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

func (u *kafkaUseCase) ConsumeMessage(ctx context.Context, topic, group string) (string, []byte, error) {
	// TODO: Implement logic to consume one message from Kafka
	return "", nil, nil
}
