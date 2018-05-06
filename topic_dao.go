package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

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

	// Produce messages to topic asynchronously
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)

	// Wait for message deliveries
	producer.Flush(15 * 1000)
}
