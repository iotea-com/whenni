package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	conn, err := nats.Connect("nats://127.0.0.1:4444")
	if err != nil {
		fmt.Printf("Error creating NATS connection: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	topic := "test1"

	msg := map[string]any{
		"timestamp":    time.Now().UnixMilli(),
		"messageQueue": "NATS",
	}

	msgJson, _ := json.Marshal(msg)

	err = conn.Publish(topic, msgJson)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v\n", topic, err)
		os.Exit(1)
	}

	fmt.Println("Successfully published message")
}
