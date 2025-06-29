package usecase

import "github.com/ducnpdev/godev-kit/internal/repo"

type Nats interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string, handler func(msg []byte)) (unsubscribe func() error, err error)
}

type natsUseCase struct {
	natsRepo repo.NatsRepo
}

func NewNatsUseCase(natsRepo repo.NatsRepo) Nats {
	return &natsUseCase{natsRepo: natsRepo}
}

func (u *natsUseCase) Publish(subject string, data []byte) error {
	return u.natsRepo.Publish(subject, data)
}

func (u *natsUseCase) Subscribe(subject string, handler func(msg []byte)) (func() error, error) {
	return u.natsRepo.Subscribe(subject, handler)
}
