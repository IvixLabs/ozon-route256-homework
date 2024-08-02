package consumer

import (
	"context"
	"route256/logger/pkg/logger"
	"route256/notifier/internal/handler"
	"route256/transactionalbox/pkg/transactionalbox"
	"route256/transactionalbox/pkg/transactionalbox/broker/kafka"
)

type Consumer struct {
	config Config
}

func New(config Config) *Consumer {
	return &Consumer{config: config}
}

func (c *Consumer) Run(ctx context.Context) {

	handlerObj := &handler.Event{}
	consumerCtx := logger.ContextWithLogger(ctx, logger.InfoLevel, "app_notifier")

	logger.Infow(consumerCtx, "Notifier service is started")
	defer logger.Infow(consumerCtx, "Notifier service is stopped")

	broker := kafka.NewKafkaConsumerBroker()
	consumer := transactionalbox.NewConsumer(broker)

	consumer.Run(consumerCtx, c.config.KafkaTopics, handlerObj)
}
