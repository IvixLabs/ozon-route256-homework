package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"route256/transactionalbox/pkg/transactionalbox"
)

type ConsumerBroker struct {
	consumerGroup sarama.ConsumerGroup
}

func NewConsumerBroker(consumerGroup sarama.ConsumerGroup) *ConsumerBroker {
	return &ConsumerBroker{consumerGroup: consumerGroup}
}

func (b *ConsumerBroker) Consume(ctx context.Context, topics []string) <-chan transactionalbox.ConsumerBrokerRequest {

	handler := NewHandler()

	go func() {
		defer handler.Close()
		defer func(cg sarama.ConsumerGroup) {
			closeErr := cg.Close()
			if closeErr != nil {
				log.Panicln(closeErr)
			}
		}(b.consumerGroup)

		for {
			err := b.consumerGroup.Consume(ctx, topics, handler)
			if err != nil {
				log.Printf("consumer error: %v", err)
			}
			if ctx.Err() != nil {
				log.Printf("consumer error: %v", ctx.Err())
				return
			}
		}
	}()

	return handler.requests
}
