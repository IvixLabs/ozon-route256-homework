package kafka

import (
	"log"
	"route256/common/pkg/env"
	kafka2 "route256/transactionalbox/pkg/kafka"
	"strings"
)

const (
	EnvKafkaGroupID = "APP_KAFKA_GROUP_ID"
)

func NewKafkaConsumerBroker() *ConsumerBroker {
	strKafkaBrokers := env.GetEnvVar(EnvKafkaBrokers)
	kafkaBrokers := strings.Split(strKafkaBrokers, ",")
	for i, _ := range kafkaBrokers {
		kafkaBrokers[i] = strings.Trim(kafkaBrokers[i], " ")
	}

	groupID := env.GetEnvVar(EnvKafkaGroupID)

	cg, err := kafka2.NewConsumerGroup(kafkaBrokers, groupID)
	if err != nil {
		log.Panicln(err)
	}

	return NewConsumerBroker(cg)
}
