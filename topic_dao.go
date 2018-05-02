package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func writeToTopic(data []byte) {
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					//fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic asynchronously
	topic := "wx.nws.api" // TODO why can't this be a const?
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)

	// Wait for message deliveries
	producer.Flush(15 * 1000)
}
