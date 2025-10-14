package kafkaActionNode

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	kafkaActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/kafka/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

// Define a struct representing your data
type TestData struct {
	ID         int    `json:"id"`
	RandomData int    `json:"data"`
	Message    string `json:"message"`
}

// TestExec validates that the producer that is running inside the nodes Exec() function
// is sending the right data and number of messages based on the number of inputs its given.
// We do this check by spinning up a kafka consumer and subscribing to the desired topic
func TestExec(t *testing.T) {
	// Define the timeout duration
	timeoutDuration := 15 * time.Second

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numMessages := 5
	receivedMessagesCount := 0 // Track the number of received messages

	// Create channels for signaling and errors
	done := make(chan struct{})
	errCh := make(chan error, 1)
	notifyChannel := make(chan node.Notification)

	t.Run("successfully executes node with long lifetime using Kafka", func(t *testing.T) {
		n := New()

		// Create input parameters for node init
		obsv := helpers.CreateTestObsv()

		// Define the configuration for the Kafka action node
		cluster := things.KafkaCluster{
			BootstrapServers: []string{"127.0.0.1:9094"},
			Topics:           []string{"test-topic"},
		}

		producer := things.KafkaProducer{
			Cluster: cluster,
		}

		config := kafkaActionNodeConfig.KafkaActionSubnodeConfig{
			Topic:              "test-topic",
			ThingKafkaProducer: producer,
		}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Call Init with the context, config, lifetime, and observability
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChannel,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Set up Kafka consumer
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": strings.Join(cluster.BootstrapServers, ","),
			"group.id":          "test-group",
			"auto.offset.reset": "earliest", // Ensure we read from the start of the topic
		})
		if err != nil {
			t.Fatalf("Failed to create Kafka consumer: %v", err)
		}
		defer consumer.Close()

		err = consumer.Subscribe(config.Topic, nil)
		if err != nil {
			t.Fatalf("Failed to subscribe to topic: %v", err)
		}

		// Consumer goroutine
		go func() {
			for {
				msg, err := consumer.ReadMessage(-1)
				if err != nil {
					errCh <- fmt.Errorf("error reading message: %v", err)
					return
				}

				log.Printf("Consumer received message: %s", string(msg.Value))

				receivedMessagesCount++
				if receivedMessagesCount == numMessages {
					log.Printf("Number of messages sent and received match")
					close(done)
					return
				}
			}
		}()

		// Execute node
		go func() {
			if execErr := n.Exec(node.ExecParams{Ctx: ctx}); execErr.Type != node.NoError {
				errCh <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Producer goroutine
		go func() {
			for i := 0; i < numMessages; i++ {
				select {
				case <-ctx.Done():
					return
				default:
					randomData := rand.Intn(100)
					message := fmt.Sprintf("Publishing message %d to %s", i, config.Topic)

					dataToSend := TestData{
						ID:         i,
						RandomData: randomData,
						Message:    message,
					}

					jsonData, err := json.Marshal(dataToSend)
					if err != nil {
						errCh <- fmt.Errorf("Error marshalling JSON: %v", err)
						return
					}

					log.Printf("Producer sent message: %s", string(jsonData))

					select {
					case n.GetInputChans(ctx)[0].Channel <- node.IoData{
						Data: jsonData,
						Type: node.BytesDataType,
						Ctx:  ctx,
					}:
					case <-ctx.Done():
						return
					}

					time.Sleep(time.Millisecond * 10)
				}
			}
		}()

		// Wait for the consumer to receive all messages or timeout
	Loop:
		for {
			select {
			case <-done:
				break Loop
			case notification := <-notifyChannel:
				log.Printf("Received notification: %v", notification)
			case err := <-errCh:
				t.Fatalf("Test failed: %v", err)
			case <-time.After(timeoutDuration):
				t.Fatalf("Test timed out after %v", timeoutDuration)
			}
		}

		log.Printf("Test completed successfully")

		// Cancel exec
		cancel()

		deinitCtx, deinitCancel := context.WithCancel(context.Background())
		defer deinitCancel()

		// Clean up node
		if deinitErr := n.Deinit(node.DeinitParams{Ctx: deinitCtx}); deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}
