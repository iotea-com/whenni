package natsSourceNode

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	natsSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/nats/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
	"github.com/nats-io/nats.go"
)

// TestExec_LongLifetime ensures that the source node with long lifetime successfully receives and processes messages from a NATS server.
func TestExec_LongLifetime(t *testing.T) {
	t.Run("successfully executes source node with long lifetime", func(t *testing.T) {
		n := New()

		// Create input parameters with a cancellation context
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define the configuration for the NATS source node
		config := natsSourceNodeConfig.NatsSourceSubnodeConfig{
			ThingNatsClient: things.NatsClient{
				Server: things.NatsServer{
					Host: "127.0.0.1",
					Port: 4444,
				},
			},
			Topic: "test-topic",
		}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}
		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)

		// Call Init with the context, config, lifetime, and observability
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		numMessages := 2

		// Create channel to store published messages
		publishedMessagesCh := make(chan map[string]any, numMessages)

		// Start mock publishing
		go func() {
			time.Sleep(time.Millisecond * 1000) // Sleep to avoid flooding the broker

			for i := 0; i < numMessages; i++ {
				message := map[string]any{
					"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
				}

				jsonData, _ := json.Marshal(message)

				// Store the published message in the channel for later comparison
				publishedMessagesCh <- message

				// Publish input message to NATS
				conn, err := nats.Connect("nats://127.0.0.1:4444")
				if err != nil {
					log.Printf("error connecting to NATS server: %s", err)
					return
				}

				err = conn.Publish("test-topic", jsonData)
				if err != nil {
					log.Printf("could not publish message to test-topic: %s", err)
				}

				time.Sleep(time.Millisecond * 100) // Sleep to avoid flooding the broker
			}
		}()

		// Verify that the subscriber node received the correct number of messages
		receivedMessagesCount := 0

	Loop:
		for {
			select {
			case <-ctx.Done():
				// Context canceled, exit the loop
				break Loop
			case receivedData := <-n.GetOutputChans(ctx)[0].Channel:
				receivedMessagesCount++
				// Type assertion to convert received data to []byte
				dataBytes, ok := receivedData.Data.([]byte)
				if !ok {
					t.Errorf("Failed to convert received data to []byte")
					continue
				}

				// Unmarshal the received data
				var receivedMsg map[string]any
				if err := json.Unmarshal(dataBytes, &receivedMsg); err != nil {
					t.Errorf("Error unmarshalling received message: %v", err)
					continue
				}

				expectedData, ok := <-publishedMessagesCh
				if !ok {
					t.Error("Expected messages channel closed unexpectedly")
					return
				}

				// Compare sent and received messages
				if !reflect.DeepEqual(receivedMsg, expectedData) {
					t.Errorf("Received message does not match expected message. Received: %v, Expected: %v", receivedMsg, expectedData)
				} else {
					t.Logf("Successful data MATCH. Received: %v, Expected: %v", receivedMsg, expectedData)
				}
			case notification := <-notifyChan:
				t.Logf("received notification %v", notification)
			case err := <-errChan:
				t.Fatalf("Exec() returned error: %v", err)
			}

			if receivedMessagesCount >= numMessages {
				break Loop
			}
		}

		// Verify that the subscriber node received the correct number of messages
		if receivedMessagesCount != numMessages {
			t.Errorf("Number of received messages (%v) does not match expected number of messages (%v)", receivedMessagesCount, numMessages)
		} else {
			t.Logf("Message count sent(%v) and received(%v) MATCH", receivedMessagesCount, numMessages)
		}

		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}
