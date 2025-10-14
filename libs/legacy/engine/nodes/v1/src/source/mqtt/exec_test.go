package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/goccy/go-json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	mqttSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/mqtt/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

const (
	TestTopic = "test/source/topic"
)

// Define a struct representing your data
type TestData struct {
	ID         int    `json:"id"`
	RandomData int    `json:"data"`
	Message    string `json:"message"`
}

// TestExec ensures that the source node with long lifetime successfully receives and processes messages from a MQTT broker.
func TestExec(t *testing.T) {

	t.Run("successfully executes source node", func(t *testing.T) {

		n := New()

		// Create input parameters with a cancellation context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define the configuration for the MQTT source node
		config := mqttSourceNodeConfig.MqttSourceNodeConfig{
			ThingMqttClient: things.MqttClient{
				Broker: things.MqttBroker{
					Host:     "127.0.0.1",
					Port:     1883,
					Protocol: "mqtt",
				},
				Username: "test",
				Password: "test",
				ClientId: uuid.NewString(),
			},
			Topic: TestTopic,
			QoS:   0,
		}

		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Call Init with the context, config, lifetime, and observability
		if initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		}); initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr.Reason)
		}

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr.Reason)
			}
		}()

		numMessages := 2

		// Create channel to store published messages
		publishedMessagesCh := make(chan TestData, numMessages)

		// Start mock publishing
		go func() {
			if err := mockPublishData(config.Topic, publishedMessagesCh, numMessages, cancel); err != nil {
				errChan <- err
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
				var receivedMsg TestData
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
				t.Fatalf("Fatal error: %v", err)
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
			t.Errorf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}

func mockPublishData(topic string, publishedMessagesCh chan<- TestData, numMessages int, cancel context.CancelFunc) error {
	// MQTT broker configuration
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("publisher")

	// Create MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("Test client failed to connect to MQTT broker: %v\n", token.Error())
	}

	// Start the publisher goroutine after the client is connected
	if client.IsConnected() {
		// Publish some messages to the topic
		defer cancel()

		for i := 0; i < numMessages; i++ {
			message := fmt.Sprintf("Publishing message %v to %v", i, topic)

			dataToSend := TestData{
				ID:         i,
				RandomData: rand.Intn(100),
				Message:    message,
			}

			jsonData, err := json.Marshal(dataToSend)
			if err != nil {
				return fmt.Errorf("Error marshalling JSON: %v\n", err)
			}

			// Store the published message in the channel for later comparison
			publishedMessagesCh <- dataToSend

			if token := client.Publish(topic, 0, false, jsonData); token.Wait() && token.Error() != nil {
				return fmt.Errorf("Failed to publish message: %v\n", token.Error())
			}

			time.Sleep(time.Millisecond * 100) // Sleep to avoid flooding the broker
		}

		// Close the channel after publishing all messages
		close(publishedMessagesCh)

	}

	return nil
}
