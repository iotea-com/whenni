package natsActionNode

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	natsActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/nats/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
	"github.com/nats-io/nats.go"
)

// Define a struct representing your data
type TestData struct {
	ID         int    `json:"id"`
	RandomData int    `json:"data"`
	Message    string `json:"message"`
}

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

	t.Run("successfully executes node using NATS with long lifetime", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		obsv := helpers.CreateTestObsv()

		// Define the configuration for the NATS action node
		server := things.NatsServer{
			Host:   "127.0.0.1",
			Port:   4444,
			Topics: []string{"test-topic"},
		}

		client := things.NatsClient{
			Server: server,
		}

		config := natsActionNodeConfig.NatsActionSubnodeConfig{
			ThingNatsClient: client,
			Topic:           "test-topic",
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
			t.Errorf("Init() returned error: %v", initErr)
		}

		// Set up NATS subscription
		natsConn, err := nats.Connect(fmt.Sprintf("nats://%s:%d", server.Host, server.Port))
		if err != nil {
			t.Fatalf("Failed to connect to NATS server: %v", err)
		}
		defer natsConn.Close()

		_, err = natsConn.Subscribe("test-topic", func(msg *nats.Msg) {
			var receivedData TestData
			if err := json.Unmarshal(msg.Data, &receivedData); err != nil {
				errCh <- fmt.Errorf("error unmarshalling message: %v", err)
				return
			}

			log.Printf("Consumer received message: %s", string(msg.Data))

			receivedMessagesCount++
			if receivedMessagesCount == numMessages {
				close(done)
			}
		})
		if err != nil {
			t.Fatalf("Failed to subscribe to NATS topic: %v", err)
		}

		// Execute node
		go func() {
			if execErr := n.Exec(node.ExecParams{Ctx: ctx}); execErr.Type != node.NoError {
				errCh <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Producer goroutine
		go func() {
			for i := 0; i < numMessages; i++ {
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

				n.GetInputChans(ctx)[0].Channel <- node.IoData{
					Data: jsonData,
					Type: node.BytesDataType,
					Ctx:  ctx,
				}

				time.Sleep(time.Millisecond * 10)
			}
		}()

		// Wait for the consumer to receive all messages or timeout
	Loop:
		for {
			select {
			case <-done:
				log.Printf("Test completed successfully")
				break Loop
			case <-notifyChannel:
				log.Printf("Notification received")
			case err := <-errCh:
				t.Fatalf("Test failed: %v", err)
			case <-time.After(timeoutDuration):
				t.Fatalf("Test timed out after %v", timeoutDuration)
			}
		}

		// Stop Exec() by cancelling the context we called it with
		cancel()

		// Clean up node
		if deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx}); deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}
