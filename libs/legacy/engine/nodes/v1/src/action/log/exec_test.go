package main

import (
	"context"
	"testing"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	logActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/log/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

func TestExec_LogsDataCorrectly(t *testing.T) {
	t.Run("logs various data types correctly", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		n := New()

		// Create channels for errors and notifications
		errChan := make(chan node.Error)
		notifyChan := make(chan node.Notification)

		// Create input parameters
		obsv := helpers.CreateTestObsv()

		// Define the configuration for the Log action node
		config := logActionNodeConfig.LogActionNodeConfig{}

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

		// Prepare test data
		testData := []node.IoData{
			{Data: []byte(`{"key": "value"}`), Type: node.BytesDataType, Ctx: ctx},
			{Data: map[string]any{"key": "value"}, Type: node.MapDataType, Ctx: ctx},
			{Data: "simple string", Type: node.StringDataType, Ctx: ctx},
		}

		// Send test data
		go func() {
			for _, data := range testData {
				n.GetInputChans(ctx)[0].Channel <- data
			}
			cancel()
		}()

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- execErr
			}
		}()

		// Monitor channels
	Loop:
		for {
			select {
			case <-ctx.Done():
				t.Log("Context cancelled")
				break Loop
			case notification := <-notifyChan:
				t.Logf("Received notification: %v", notification)
			case err := <-errChan:
				t.Fatalf("Fatal error: %v", err.Reason)
			}
		}

		// Cleanup
		if deinitErr := n.Deinit(node.DeinitParams{
			Ctx: context.Background(),
		}); deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}
