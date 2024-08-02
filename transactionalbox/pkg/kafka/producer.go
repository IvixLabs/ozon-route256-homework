package kafka

import (
	"github.com/IBM/sarama"
	"time"
)

func NewAsyncProducer(brokerAddrs []string) (sarama.AsyncProducer, error) {
	pc := NewProducerConfig()

	aProducer, err := sarama.NewAsyncProducer(brokerAddrs, pc)
	if err != nil {
		return nil, err
	}

	return aProducer, nil

}

func NewProducerConfig() *sarama.Config {

	c := sarama.NewConfig()

	c.Producer.Partitioner = sarama.NewHashPartitioner

	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Idempotent = false
	c.Producer.Retry.Max = 100
	c.Producer.Retry.Backoff = 5 * time.Millisecond
	c.Net.MaxOpenRequests = 1
	c.Producer.CompressionLevel = sarama.CompressionLevelDefault
	c.Producer.Compression = sarama.CompressionGZIP
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true

	return c
}
