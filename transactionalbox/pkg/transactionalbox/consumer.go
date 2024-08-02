package transactionalbox

import (
	"context"
	"log"
)

type ConsumerBrokerRequest interface {
	Mark()
	Payload() []byte
}

type ConsumerBroker interface {
	Consume(ctx context.Context, topics []string) <-chan ConsumerBrokerRequest
}

type ConsumerHandler interface {
	Handle(ctx context.Context, payload []byte) error
}

type Consumer struct {
	broker ConsumerBroker
}

func NewConsumer(broker ConsumerBroker) *Consumer {
	return &Consumer{broker: broker}
}

func (c *Consumer) Run(ctx context.Context, topics []string, handler ConsumerHandler) {

	eventCh := c.broker.Consume(ctx, topics)
	for {
		select {
		case req := <-eventCh:
			go func() {
				err := handler.Handle(ctx, req.Payload())
				if err != nil {
					log.Println(err)
					return
				}
				req.Mark()
			}()
		case <-ctx.Done():
			return
		}

	}

}
