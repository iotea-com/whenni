package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/goccy/go-json"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	mqttActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/mqtt/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

const (
	TestTopic = "test/action/topic"
)

// Define a struct representing your data
type TestData struct {
	ID         int    `json:"id"`
	RandomData int    `json:"data"`
	Message    string `json:"message"`
}

func compareSentAndReceivedMessages(t *testing.T, sentMessagesCh <-chan TestData, msg MQTT.Message) {
	receivedMsg := TestData{}
	if err := json.Unmarshal(msg.Payload(), &receivedMsg); err != nil {
		t.Fatalf("error unmarshalling received message: %v", err)

	}
	sentMsg := <-sentMessagesCh
	if sentMsg != receivedMsg {
		t.Fatalf("received message %v differs from sent message %v", receivedMsg, sentMsg)
	} else {
		t.Logf("Data match successful")
	}
}

// TestExec_LongLifetime ensures that the publishing node successfully sends multiple messages over MQTT
// with a long lifetime, meaning it remains active for an extended period. It verifies that
// the messages sent by the node match the ones received by the MQTT subscriber.
func TestExec_LongLifetime(t *testing.T) {
	// Define the timeout duration
	timeoutDuration := 10 * time.Second

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numMessages := 5
	sentMessagesCh := make(chan TestData, numMessages)
	receivedMessagesCount := 0 // Track the number of received messages

	// Create channels for signaling and errors
	done := make(chan struct{})
	errCh := make(chan error, 1)
	notifyChannel := make(chan node.Notification)

	t.Run("successfully executes node with long lifetime", func(t *testing.T) {
		n := New()

		// Create input parameters for node init
		obsv := helpers.CreateTestObsv()

		// Define the configuration for the MQTT action node
		config := mqttActionNodeConfig.MqttActionNodeConfig{
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
			t.Fatalf("Init() returned error: %v", initErr.Reason)
		}

		// Set up MQTT subscriber
		opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			t.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
		}
		defer client.Disconnect(250)

		// Subscribe to topic
		topic := config.Topic
		if token := client.Subscribe(topic, 0, func(client MQTT.Client, msg MQTT.Message) {
			t.Logf("Subscriber received message on topic %s:\n %s", msg.Topic(), msg.Payload())
			compareSentAndReceivedMessages(t, sentMessagesCh, msg)
			receivedMessagesCount++

			if receivedMessagesCount == numMessages {
				t.Logf("Number of messages sent and received match")
				close(done)
			}
		}); token.Wait() && token.Error() != nil {
			t.Fatalf("Failed to subscribe to topic %s: %v", topic, token.Error())
		}

		// Producer goroutine
		go func() {
			defer cancel() // Cancel the context when done
			for i := 0; i < numMessages; i++ {
				// Generate random data to send to the action node input
				randomData := rand.Intn(100)
				message := fmt.Sprintf("Publishing message %d to %s", i, topic)

				dataToSend := TestData{
					ID:         i,
					RandomData: randomData,
					Message:    message,
				}

				// Marshal the struct into JSON
				jsonData, err := json.MarshalIndent(dataToSend, "", "  ")
				if err != nil {
					errCh <- fmt.Errorf("Error marshalling JSON: %v", err)
					return
				}

				// Send random data to inputChannels[0].Channel
				n.GetInputChans(ctx)[0].Channel <- node.IoData{
					Data: jsonData,
					Type: node.BytesDataType,
					Ctx:  ctx,
				}

				// Send message to the sent messages channel
				sentMessagesCh <- dataToSend

				// Sleep for a while before sending the next data
				time.Sleep(time.Millisecond * 100)
			}
			close(sentMessagesCh)
		}()

		// Execute node
		go func() {
			if execErr := n.Exec(node.ExecParams{Ctx: ctx}); execErr.Type != node.NoError {
				errCh <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Wait for the subscriber function to print the message or timeout
	Loop:
		for {
			select {
			case <-done:
				t.Logf("Test completed successfully")
				break Loop
			case notification := <-notifyChannel:
				t.Logf("Received notification: %v", notification)
			case err := <-errCh:
				t.Fatalf("Test failed: %v", err)
			case <-time.After(timeoutDuration):
				t.Fatalf("Test timed out after %v", timeoutDuration)
			}
		}

		// Clean up node
		if deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx}); deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}

// TestExec_FailsWithNonByteData confirms that the publishing node fails gracefully
// when provided with non-byte data. It verifies that the node returns an error
// when attempting to send data of an incorrect type over MQTT.
func TestExec_FailsWithNonByteData(t *testing.T) {
	t.Run("successfully fails when non-byte data is passed", func(t *testing.T) {
		// Create a context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		n := New()

		// Create input parameters
		obsv := helpers.CreateTestObsv()

		// Define the configuration for the MQTT action node
		config := mqttActionNodeConfig.MqttActionNodeConfig{
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

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		notifyChannel := make(chan node.Notification)
		errCh := make(chan error, 1)

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

		nonByteData := "not a byte slice"
		ioData := node.IoData{
			Data: nonByteData,
			Type: node.BytesDataType, // This is not bytes, deliberately causing an error
			Ctx:  ctx,
		}

		go func() {
			n.GetInputChans(ctx)[0].Channel <- ioData
		}()

		// Start the execution in a goroutine
		go func() {
			if execErr := n.Exec(node.ExecParams{Ctx: ctx}); execErr.Type != node.NoError {
				errCh <- fmt.Errorf("Exec() returned unexpected error: %v", execErr)
			}
		}()

		// Monitor channels for the dataInvalid notification
	Loop:
		for {
			select {
			case notification := <-notifyChannel:
				if notification.Type == node.NotifyDataInvalid {
					t.Logf("Received expected dataInvalid notification: %v", notification)
					break Loop
				} else {
					t.Errorf("Received unexpected notification: %v", notification)
				}
			case err := <-errCh:
				t.Fatalf("Fatal error: %v", err)
			case <-time.After(5 * time.Second):
				t.Fatalf("Test timed out waiting for dataInvalid notification")
			}
		}

		// Cleanup
		if deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx}); deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}
