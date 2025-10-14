package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	_, err = conn.Subscribe(topic, func(msg *nats.Msg) {
		log.Printf("New message: %s", string(msg.Data))
	})
	if err != nil {
		fmt.Printf("Error subscribing to topic %s: %v\n", topic, err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Subscriber is ready, waiting for messages...")

	for {
		sig := <-sigchan
		fmt.Printf("Caught signal %v: terminating\n", sig)
		break
	}

	fmt.Println("Closing subscriber connection...")
}
