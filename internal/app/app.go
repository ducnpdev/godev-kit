// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ducnpdev/godev-kit/config"
	"github.com/ducnpdev/godev-kit/internal/controller/http"
	"github.com/ducnpdev/godev-kit/internal/repo/externalapi"
	vietqrrepo "github.com/ducnpdev/godev-kit/internal/repo/externalapi/vietqr"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent"
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/internal/usecase/billing"
	natuc "github.com/ducnpdev/godev-kit/internal/usecase/nat"
	"github.com/ducnpdev/godev-kit/internal/usecase/payment"
	redisuc "github.com/ducnpdev/godev-kit/internal/usecase/redis"
	"github.com/ducnpdev/godev-kit/internal/usecase/translation"
	"github.com/ducnpdev/godev-kit/internal/usecase/user"
	vietqruc "github.com/ducnpdev/godev-kit/internal/usecase/vietqr"
	"github.com/ducnpdev/godev-kit/pkg/httpserver"
	"github.com/ducnpdev/godev-kit/pkg/kafka"
	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/ducnpdev/godev-kit/pkg/nats"
	"github.com/ducnpdev/godev-kit/pkg/postgres"
	"github.com/ducnpdev/godev-kit/pkg/redis"
	// amqprpc "github.com/ducnpdev/godev-kit/internal/controller/amqp_rpc"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL,
		postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.MinPoolSize(cfg.PG.PoolMin),
		postgres.MaxConnLifetime(cfg.PG.MaxConnLifetime),
		postgres.MaxConnIdleTime(cfg.PG.MaxConnIdleTime),
		postgres.HealthCheckPeriod(cfg.PG.HealthCheckPeriod),
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Redis client
	redisClient, err := redis.New(cfg.Redis.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
	}
	defer redisClient.Close()

	// Kafka Repository
	kafkaRepo := persistent.NewKafkaRepoWithConfig(
		cfg.Kafka.Brokers,
		l.Zerolog(),
		cfg.Kafka.Control.ProducerEnabled,
		cfg.Kafka.Control.ConsumerEnabled,
	)
	defer func() {
		if err := kafkaRepo.Close(); err != nil {
			l.Error(fmt.Errorf("app - Run - kafkaRepo.Close: %w", err))
		}
	}()

	// NATS client
	var (
		natsClient *nats.NatsClient
		errNats    error
	)
	if cfg.NATS.Enable {
		if cfg.NATS.Timeout > 0 {
			natsClient, errNats = nats.New(cfg.NATS.URL, nats.ConnTimeout(cfg.NATS.Timeout))
		} else {
			natsClient, errNats = nats.New(cfg.NATS.URL)
		}
		if errNats != nil {
			l.Fatal(fmt.Errorf("app - Run - nats.New: %w", err))
		}
		defer natsClient.Close()
	}

	// Use-Case
	translationUseCase := translation.New(
		persistent.New(pg),
		externalapi.New(),
	)

	userUseCase := user.New(
		persistent.NewUserRepo(pg),
		cfg.JWT.Secret,
	)
	kafkaUseCase := usecase.NewKafkaUseCase(kafkaRepo)
	redisUseCase := redisuc.NewRedisUseCase(
		persistent.NewRedisRepo(redisClient),
	)
	natsUseCase := natuc.NewNatsUseCase(persistent.NewNatsRepo(natsClient))
	vietqrUseCase := vietqruc.NewVietQRUseCase(
		vietqrrepo.NewVietQRRepo(),
		persistent.NewVietQRRepo(pg),
	)
	billingUseCase := billing.New()

	redisRepo := persistent.NewRedisRepo(redisClient)
	shipperLocationRepo := persistent.NewShipperLocationRepo(pg)
	shipperLocationUsecase := redisuc.NewShipperLocationUseCase(redisRepo, shipperLocationRepo)

	// Payment Use Case
	paymentRepo := persistent.NewPaymentRepo(pg)
	
	// Only create Kafka producer if enabled
	var kafkaProducer *kafka.Producer
	if cfg.Kafka.Control.ProducerEnabled {
		kafkaProducer = kafka.NewProducer(cfg.Kafka.Brokers, l.Zerolog())
	}
	paymentUseCase := payment.NewPaymentUseCase(paymentRepo, kafkaProducer, l.ZerologPtr())

	// Setup context for Kafka operations
	ctx := context.Background()

	// Only create and start payment consumer if Kafka consumer is enabled
	var paymentConsumer *payment.PaymentConsumer
	if cfg.Kafka.Control.ConsumerEnabled {
		paymentConsumer = payment.NewPaymentConsumer(cfg.Kafka.Brokers, "payment-processor", paymentUseCase, l.ZerologPtr())
		
		// Start Payment Consumer
		go func() {
			if err := paymentConsumer.Start(ctx); err != nil {
				l.Error(fmt.Errorf("app - Run - paymentConsumer.Start: %w", err))
			}
		}()
	} else {
		l.Info("Kafka consumer is disabled, skipping payment consumer initialization")
	}

	// Kafka Event Use Case
	// kafkaEventUseCase := usecase.NewKafkaEventUseCase(kafkaRepo, l.Zerolog())

	// Setup Kafka consumers
	// if err := kafkaEventUseCase.ConsumeUserEvents(ctx); err != nil {
	// 	l.Error(fmt.Errorf("app - Run - ConsumeUserEvents: %w", err))
	// }

	// if err := kafkaEventUseCase.ConsumeTranslationEvents(ctx); err != nil {
	// 	l.Error(fmt.Errorf("app - Run - ConsumeTranslationEvents: %w", err))
	// }

	// Start Payment Consumer
	// go func() {
	// 	if err := paymentConsumer.Start(ctx); err != nil {
	// 		l.Error(fmt.Errorf("app - Run - paymentConsumer.Start: %w", err))
	// 	}
	// }()

	// Start Kafka consumers
	// kafkaRepo.StartAllConsumers(ctx)

	// RabbitMQ RPC Server
	// rmqRouter := amqprpc.NewRouter(translationUseCase, l)

	// rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	// if err != nil {
	// 	l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	// }

	// gRPC Server
	// grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	// grpc.NewRouter(grpcServer.App, translationUseCase, l)

	// HTTP Server
	httpServer := httpserver.New(cfg, httpserver.Port(cfg.HTTP.Port))
	http.NewRouter(httpServer.App, cfg, translationUseCase, userUseCase, kafkaUseCase, redisUseCase, natsUseCase, vietqrUseCase, billingUseCase, l, shipperLocationUsecase, paymentUseCase, billingUseCase)

	// Start servers
	// rmqServer.Start()
	// grpcServer.Start()
	httpServer.Start()

	l.Info("app running at port http:%s", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("%s", "app - Run - signal: "+s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
		// case err = <-grpcServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
		// case err = <-rmqServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	// err = grpcServer.Shutdown()
	// if err != nil {
	// 	l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	// }

	// err = rmqServer.Shutdown()
	// if err != nil {
	// 	l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	// }
}
