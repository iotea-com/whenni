package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	thresholdConditionalNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/conditional/threshold/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

// Define a struct representing your data
type TestData struct {
	Value int `json:"value"`
}

// TestExec validates that the conditional node with short lifetime routes output data appropriately
// based on the input conditions passed
func TestExec(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("successfully executes node and passes the condition", func(t *testing.T) {
		n := New()
		obsv := helpers.CreateTestObsv()

		config := thresholdConditionalNodeConfig.ThresholdConditionalNodeConfig{
			Model: models.Model{
				Attributes: map[string]models.Attribute{
					"f0": {
						Id:   "f0",
						Key:  "num1",
						Type: "number",
					},
				},
			},
			Conditions: []models.Condition{{
				AttributeId: "f0",
				Operator:    "greater-than",
				Value:       0,
			}},
			LogicalOperator: "AND",
		}

		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Create notification channel
		notifyChan := make(chan node.Notification)
		errCh := make(chan error)

		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Send test data
		go func() {
			dataToSend := TestData{Value: 5}
			jsonData, err := json.Marshal(dataToSend)
			if err != nil {
				errCh <- fmt.Errorf("failed to marshal test data: %v", err)
				return
			}

			n.GetInputChans(ctx)[0].Channel <- node.IoData{
				Data: jsonData,
				Type: node.BytesDataType,
				Ctx:  ctx,
			}
		}()

		// Start execution
		go func() {
			execErr := n.Exec(node.ExecParams{Ctx: ctx})
			if execErr.Type != node.NoError {
				errCh <- fmt.Errorf("exec error: %v", execErr)
			}
		}()

		// Wait for notifications and results
		expectedNotifications := map[node.NotificationType]bool{
			node.NotifyDataProcessed: false,
		}

	Loop:
		for {
			select {
			case notification := <-notifyChan:
				expectedNotifications[notification.Type] = true
				if notification.Type == node.NotifyDataInvalid {
					t.Errorf("Unexpected invalid data notification")
				}
			case err := <-errCh:
				t.Fatalf("Error during execution: %v", err)
			case output := <-n.GetOutputChans(ctx)[1].Channel:
				// Verify output data if needed
				t.Logf("Received output on success channel: %v", output)
				break Loop
			case <-time.After(5 * time.Second):
				t.Fatal("Test timed out")
			case <-ctx.Done():
				return
			}

			// Check if we've received all expected notifications
			allReceived := true
			for _, received := range expectedNotifications {
				if !received {
					allReceived = false
					break
				}
			}
			if allReceived {
				break Loop
			}
		}

		// Cleanup
		deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx})
		if deinitErr.Type != node.NoError {
			t.Errorf("Deinit() returned error: %v", deinitErr)
		}
	})

	t.Run("sends notification for invalid data type", func(t *testing.T) {
		n := New()
		obsv := helpers.CreateTestObsv()

		config := thresholdConditionalNodeConfig.ThresholdConditionalNodeConfig{
			Model: models.Model{
				Attributes: map[string]models.Attribute{
					"f0": {
						Id:   "f0",
						Key:  "num1",
						Type: "number",
					},
				},
			},
			Conditions: []models.Condition{{
				AttributeId: "f0",
				Operator:    "greater-than",
				Value:       0,
			}},
			LogicalOperator: "AND",
		}

		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		notifyChan := make(chan node.Notification)
		errCh := make(chan error)

		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Send invalid data
		go func() {
			n.GetInputChans(ctx)[0].Channel <- node.IoData{
				Data: "invalid data",
				Type: node.BytesDataType,
				Ctx:  ctx,
			}
		}()

		// Start execution
		go func() {
			execErr := n.Exec(node.ExecParams{Ctx: ctx})
			if execErr.Type != node.NoError {
				errCh <- fmt.Errorf("exec error: %v", execErr)
			}
		}()

		// Wait for invalid data notification
		select {
		case notification := <-notifyChan:
			if notification.Type != node.NotifyDataInvalid {
				t.Errorf("Expected invalid data notification, got %v", notification.Type)
			}
		case err := <-errCh:
			t.Fatalf("Unexpected error: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for invalid data notification")
		}

		// Cleanup
		deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx})
		if deinitErr.Type != node.NoError {
			t.Errorf("Deinit() returned error: %v", deinitErr.Reason)
		}
	})
}
