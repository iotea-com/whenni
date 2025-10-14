package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

const (
	ShortLifetimeExpectedDuration = 100 * time.Millisecond
)

// TestInit tests that the node won't accept a wrong configuration byte array passed
func TestInit(t *testing.T) {
	t.Run("returns error on wrong config", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Set normal input parameters
		ctx, cancel := context.WithTimeout(context.Background(), ShortLifetimeExpectedDuration)
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define the wrong configuration for the node
		type DummyConfig struct {
			DummyVar int `json:"topic"`
		}
		config := DummyConfig{}

		// Make the notification channel
		notifyChannel := make(chan node.Notification)

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Call Init
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			Obsv:          obsv,
			NotifyChannel: notifyChannel,
		})

		fmt.Printf("initErr: %v\n", initErr)

		// Check if error is returned, since we expect it to fail with the dummy configuration
		if initErr.Type != node.ValidationError {
			t.Errorf("Init() did not return error; expected error for wrong config")
		}

		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}
