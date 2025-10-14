package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9094", // Ensure this matches your advertised listener
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		fmt.Printf("Error creating Kafka producer: %v\n", err)
		os.Exit(1)
	}
	defer producer.Close()

	msg := map[string]any{
		"timestamp":    time.Now().UnixMilli(),
		"messageQueue": "Kafka",
	}

	msgJson, _ := json.Marshal(msg)

	topic := "test1"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic,
		},
		Value: msgJson,
	}
	err = producer.Produce(message, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic %s: %v\n", topic, err)
		os.Exit(1)
	}

	fmt.Println("Closing producer...")
}
