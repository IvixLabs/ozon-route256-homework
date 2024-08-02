package main

import (
	"context"
	"route256/common/pkg/env"
	"route256/notifier/internal/app/consumer"
	"strings"
)

const EnvKafkaTopics = "APP_KAFKA_TOPICS"

func runConsumer(ctx context.Context) {
	config := getConsumerConfig()
	consumerApp := consumer.New(config)

	consumerApp.Run(ctx)
}

func getConsumerConfig() consumer.Config {
	strKafkaTopics := env.GetEnvVar(EnvKafkaTopics)
	kafkaTopics := strings.Split(strKafkaTopics, ",")
	for i := range kafkaTopics {
		kafkaTopics[i] = strings.Trim(kafkaTopics[i], " ")
	}

	return consumer.Config{
		KafkaTopics: kafkaTopics,
	}
}
