package main

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	httpActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/http/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

// Define a struct representing your data
type TestData struct {
	ID string `json:"id"`
}

func compareSentAndReceivedMessages(t *testing.T, sentMessagesCh <-chan TestData, msg []byte) {
	receivedMsg := TestData{}
	if err := json.Unmarshal(msg, &receivedMsg); err != nil {
		t.Errorf("error unmarshalling received message: %v", err)
	}
	sentMsg := <-sentMessagesCh
	if sentMsg != receivedMsg {
		t.Errorf("received message %v differs from sent message %v", receivedMsg, sentMsg)
	}

	t.Logf("Data match successful")
}

// TestExec_LongLifetime ensures that the publishing node successfully sends multiple HTTP requests
// with a long lifetime, meaning it remains active for an extended period. It verifies that
// the messages sent by the node match the ones received by the HTTP server.
func TestExec_LongLifetime(t *testing.T) {
	t.Run("successfully executes node with long lifetime", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		numMessages := 5
		sentMessagesCh := make(chan TestData, numMessages)
		receivedMessagesCount := 0
		testPort := 23456

		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)
		done := make(chan struct{})

		n := New()
		obsv := helpers.CreateTestObsv()

		// Define the configuration for the HTTP action node
		config := httpActionNodeConfig.HttpActionNodeConfig{
			ThingHttpServer: things.HttpServer{
				Host:     "localhost",
				Port:     testPort,
				Protocol: "http",
				Paths:    []string{"/"},
			},
			Path:   "/",
			Method: "POST",
		}

		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Initialize the node
		if initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		}); initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr.Reason)
		}

		// Set up a mock HTTP server
		mockServer := fiber.New()
		mockServer.Post("/", func(c *fiber.Ctx) error {
			receivedMessagesCount++
			c.Status(200).JSON(nil)
			compareSentAndReceivedMessages(t, sentMessagesCh, c.Body())

			if receivedMessagesCount == numMessages {
				close(done)
			}
			return nil
		})

		// Start the mock server
		go func() {
			if err := mockServer.Listen(fmt.Sprintf(":%d", testPort)); err != nil {
				errChan <- fmt.Errorf("mock server error: %v", err)
			}
		}()

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr.Reason)
			}
		}()

		// Send requests in a separate goroutine
		go func() {
			for i := 0; i < numMessages; i++ {
				requestBody := TestData{
					ID: uuid.NewString(),
				}

				requestBodyBytes, err := json.Marshal(requestBody)
				if err != nil {
					errChan <- fmt.Errorf("failed to marshal request body: %v", err)
					return
				}

				n.GetInputChans(ctx)[0].Channel <- node.IoData{
					Data: requestBodyBytes,
					Type: node.BytesDataType,
					Ctx:  ctx,
				}

				sentMessagesCh <- requestBody
				time.Sleep(time.Millisecond * 250)
			}
			close(sentMessagesCh)
		}()

		// Monitor channels
	Loop:
		for {
			select {
			case <-ctx.Done():
				t.Log("Context cancelled")
				break Loop
			case <-done:
				t.Logf("Received all %d messages successfully", numMessages)
				cancel()
			case notification := <-notifyChan:
				t.Logf("Received notification: %v", notification)
			case err := <-errChan:
				t.Fatalf("Fatal error: %v", err)
			}
		}

		// Verify message count
		if receivedMessagesCount != numMessages {
			t.Errorf("Expected %d messages, got %d", numMessages, receivedMessagesCount)
		}

		// Cleanup
		if err := mockServer.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown mock server: %v", err)
		}

		if deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		}); deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}

// TestExec_FailsWithNonByteData confirms that the publishing node fails gracefully
// when provided with non-byte data. It verifies that the node returns an error
// when attempting to send data of an incorrect type over HTTP.
func TestExec_FailsWithNonByteData(t *testing.T) {
	t.Run("successfully fails when non-byte data is passed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		n := New()

		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)

		// Create input parameters
		obsv := helpers.CreateTestObsv()

		// Define the configuration for the HTTP action node
		config := httpActionNodeConfig.HttpActionNodeConfig{
			ThingHttpServer: things.HttpServer{
				Host:     "localhost",
				Port:     51111,
				Protocol: "http",
				Paths:    []string{"/"},
			},
			Path:   "/",
			Method: "POST",
		}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Initialize the node
		if initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		}); initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr.Reason)
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
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned unexpected error: %v", execErr)
			}
		}()

		// Monitor channels for the dataInvalid notification
		select {
		case <-ctx.Done():
			t.Log("Context cancelled")
		case notification := <-notifyChan:
			if notification.Type == node.NotifyDataInvalid {
				t.Logf("Received expected dataInvalid notification: %v", notification)
			} else {
				t.Errorf("Received unexpected notification: %v", notification)
			}
		case err := <-errChan:
			t.Fatalf("Fatal error: %v", err)
		}

		// Cleanup
		if deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		}); deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}
