package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	producer *kafka.Producer
	err      error
)

func init() {
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	// Unable to connect to broker is not an error
	if err != nil {
		panic(err)
	}
}

func writeToTopic(data []byte, topic *string) {
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Warn("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)

	// TODO why do I need this? Wait for message deliveries
	producer.Flush(15 * 1000)
}
