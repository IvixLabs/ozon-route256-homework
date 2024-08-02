package kafka

import (
	"bytes"
	"context"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
	"route256/transactionalbox/pkg/kafka"
	"route256/transactionalbox/pkg/transactionalbox"
	"sync"
)

var (
	HeaderKey = []byte("id")
)

type ProducerBroker struct {
	producer  sarama.AsyncProducer
	successCh chan uuid.UUID
	errorCh   chan transactionalbox.ErrSending
}

var _ transactionalbox.ProducerBroker = (*ProducerBroker)(nil)

func NewProducerBroker(producer sarama.AsyncProducer) *ProducerBroker {
	return &ProducerBroker{producer: producer}
}

func (b *ProducerBroker) Send(ctx context.Context, recordId uuid.UUID, msg transactionalbox.Message) error {

	ctxErr := ctx.Err()
	if ctxErr != nil {
		return nil
	}

	byteRecordId, err := recordId.MarshalBinary()
	if err != nil {
		return err
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic:   msg.Topic,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Body),
		Headers: []sarama.RecordHeader{{HeaderKey, byteRecordId}},
	}
	b.producer.Input() <- kafkaMsg

	return nil
}

func (b *ProducerBroker) Success(ctx context.Context) <-chan uuid.UUID {
	if b.successCh != nil {
		return b.successCh
	}

	sync.OnceFunc(func() {
		b.successCh = make(chan uuid.UUID)
		go func() {
			defer close(b.successCh)

			for {
				select {
				case kafkaMsg, ok := <-b.producer.Successes():
					if !ok {
						return
					}

					var bytesUUID []byte
					for _, header := range kafkaMsg.Headers {
						if bytes.Equal(header.Key, HeaderKey) {
							bytesUUID = header.Value
						}
					}

					if len(bytesUUID) == 0 {
						log.Panicln("recordId is not found")
					}

					key, err := uuid.FromBytes(bytesUUID)
					if err != nil {
						log.Panicln(err)
					}

					b.successCh <- key

				case <-ctx.Done():
					return
				}
			}
		}()
	})()

	return b.successCh

}

func (b *ProducerBroker) Errors(ctx context.Context) <-chan transactionalbox.ErrSending {
	if b.errorCh != nil {
		return b.errorCh
	}

	sync.OnceFunc(func() {
		b.errorCh = make(chan transactionalbox.ErrSending)

		go func() {
			defer close(b.errorCh)

			for {
				select {
				case kafkaErr, ok := <-b.producer.Errors():
					log.Println("ERROR", kafkaErr, ok)

					if !ok {
						return
					}

					key, err := uuid.ParseBytes(kafkaErr.Msg.Key.(sarama.ByteEncoder))
					if err != nil {
						log.Panicln(err)
					}

					b.errorCh <- transactionalbox.ErrSending{Err: kafkaErr.Unwrap(), Key: key}

				case <-ctx.Done():
					return
				}
			}
		}()
	})()

	return b.errorCh
}

func (b *ProducerBroker) ConsumeGroup(ctx context.Context, brokers []string, groupID string, topics []string, handler sarama.ConsumerGroupHandler) {

	cg, err := kafka.NewConsumerGroup(brokers, groupID)
	if err != nil {
		log.Panicln(err)
	}

	defer func(cg sarama.ConsumerGroup) {
		closeErr := cg.Close()
		if closeErr != nil {
			log.Panicln(err)
		}
	}(cg)

	err = cg.Consume(ctx, topics, handler)
	if err != nil {
		log.Panicln(err)
	}
}
