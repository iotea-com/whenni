package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	httpSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/http/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

// Define a struct representing your data
type TestData struct {
	ID         int    `json:"id"`
	RandomData int    `json:"data"`
	Message    string `json:"message"`
}

func TestExec(t *testing.T) {
	t.Run("successfully executes source node", func(t *testing.T) {
		n := New()

		// Create input parameters with a cancellation context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define the configuration for the MQTT source node
		config := httpSourceNodeConfig.HttpSourceNodeConfig{
			Method:          "GET",
			ResponseTimeout: 1,
			ThingHttpServer: things.HttpServer{
				Host:     "localhost",
				Port:     TestServerPort,
				Protocol: "http",
				Paths:    []string{"/"},
			},
		}

		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Call Init with the context, config, and observability
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Allow some time for the server to be spawn up
		time.Sleep(time.Millisecond * 1000)

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Verify that the subscriber node received the correct number of messages
		numRequests := 2
		receivedRequestsCount := 0

		// Start mock publishing
		go func() {
			// We expect the server to return a 504 because there's no response nodes
			// connected in this test, so the node is expected to timeout
			err := mockRequests(fiber.MethodGet, 504, numRequests)
			if err != nil {
				errChan <- fmt.Errorf("mockRequests() returned error: %v", err)
			}
		}()

	Loop:
		for {
			select {
			case <-ctx.Done():
				// Context canceled, exit the loop
				break Loop
			case <-n.GetOutputChans(ctx)[0].Channel:
				receivedRequestsCount++
				if receivedRequestsCount >= numRequests {
					cancel()
				}
			case notification := <-notifyChan:
				t.Logf("received notification %v", notification)
			case err := <-errChan:
				t.Fatalf("Exec() returned error: %v", err)
			}
		}

		// Verify that the subscriber node received the correct number of messages
		if receivedRequestsCount != numRequests {
			t.Errorf("Number of received messages (%v) does not match expected number of messages (%v)", receivedRequestsCount, numRequests)
		} else {
			t.Logf("Message count sent(%v) and received(%v) MATCH", receivedRequestsCount, numRequests)
		}

		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}

func mockRequests(method string, expectedResponseCode int, numRequests int) error {
	// Create HTTP client and request
	client := http.DefaultClient
	maxRetries := 5
	retryDelay := 500 * time.Millisecond

	for i := 0; i < numRequests; i++ {
		request, err := http.NewRequest(method, fmt.Sprintf("http://127.0.0.1:%v", TestServerPort), nil)
		if err != nil {
			return err
		}

		var response *http.Response
		var lastErr error
		for retry := 0; retry < maxRetries; retry++ {
			response, err = client.Do(request)
			if err == nil {
				break
			}
			lastErr = err
			if retry < maxRetries-1 {
				time.Sleep(retryDelay)
			}
		}
		if lastErr != nil {
			return fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
		}

		if expectedResponseCode != response.StatusCode {
			return fmt.Errorf("unexpected status code received during mock HTTP request: got %d, expected %d", response.StatusCode, expectedResponseCode)
		}
	}

	return nil
}
