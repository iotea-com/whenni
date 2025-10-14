package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9094", // Ensure this matches your advertised listener
		"group.id":          "test-group",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %v\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	topic := "test1"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Error subscribing to topic %s: %v\n", topic, err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Consumer is ready, waiting for messages...")

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Received message: %s: %s", msg.TopicPartition, string(msg.Value))
			} else {
				fmt.Printf("Error while consuming message: %v\n", err)
			}

			var m map[string]any
			json.Unmarshal(msg.Value, &m)

			if m != nil {
				log.Printf("m: %#v", m)
			}
		}
	}

	fmt.Println("Closing consumer...")
}
