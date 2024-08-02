package kafka

import (
	"log"
	"route256/common/pkg/env"
	"route256/transactionalbox/pkg/kafka"
	"strings"
)

const (
	EnvKafkaBrokers = "APP_KAFKA_BROKERS"
)

func NewKafkaProducerBroker() *ProducerBroker {

	strKafkaBrokers := env.GetEnvVar(EnvKafkaBrokers)
	kafkaBrokers := strings.Split(strKafkaBrokers, ",")

	for i, _ := range kafkaBrokers {
		kafkaBrokers[i] = strings.Trim(kafkaBrokers[i], " ")
	}

	kafkaProducer, err := kafka.NewAsyncProducer(kafkaBrokers)
	if err != nil {
		log.Panicln(err)
	}

	return NewProducerBroker(kafkaProducer)
}
