package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	transformProcessingNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/processing/transform/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()

	t.Run("returns error on invalid input", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Set normal input parameters
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define a valid configuration for the node
		config := transformProcessingNodeConfig.TransformProcessingNodeConfig{
			ModelInput: models.Model{
				Attributes: map[string]models.Attribute{
					"inputField1": {
						Id:   "inputField1",
						Key:  "myString",
						Type: "string",
					},
					"inputField2": {
						Id:   "inputField2",
						Key:  "myNumber",
						Type: "number",
					},
				},
			},
			ModelOutput: models.Model{
				Attributes: map[string]models.Attribute{
					"outputField1": {
						Id:   "outputField1",
						Key:  "myStringOut",
						Type: "string",
					},
					"outputField2": {
						Id:   "outputField2",
						Key:  "myNumberOut",
						Type: "number",
					},
				},
			},
			Mapping: map[string]string{
				"outputField1": "inputField1",
				"outputField2": "inputField2",
			},
		}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Make a notify channel
		notifyChan := make(chan node.Notification)

		// Call Init
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Prepare the input data that conforms to the input structure
		inputData := map[string]any{
			"myString": "test string",
			"myNumber": 123,
		}

		inputBytes, err := json.Marshal(&inputData)
		if err != nil {
			t.Fatalf("Failed to marshal input data: %v", err)
		}

		// Send the input data to the input channel
		go func() {
			n.(*TransformProcessingNode).inputChannels[0].Channel <- node.IoData{
				Data: inputBytes,
				Type: node.BytesDataType,
				Ctx:  ctx,
			}
		}()

		// Call Exec
		errCh := make(chan error)
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errCh <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Retrieve the output from the output channel
		var outputData node.IoData
	Loop:
		for {
			select {
			case outputData = <-n.(*TransformProcessingNode).outputChannels[0].Channel:
				// success, do nothing
				break Loop
			case <-time.After(30 * time.Second):
				t.Fatal("Timeout waiting for output data")
			case err := <-errCh:
				t.Fatalf("Error on exec %v", err)
			case notification := <-notifyChan:
				t.Logf("received notification %v", notification)
			}
		}

		// Unmarshal the output data
		var outputMap map[string]any
		if outputData.Data != nil {
			err = json.Unmarshal(outputData.Data.([]byte), &outputMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal output data: %v", err)
			}
		} else {
			t.Fatal("Output data is nil")
		}

		t.Logf("outputMap: %#v", outputMap)

		// Validate that the output data matches the expected result
		require.Equal(t, outputMap["myStringOut"], "test string")
		require.Equal(t, outputMap["myNumberOut"], float64(123))

		// Cleanup
		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Errorf("Deinit() returned error: %v", deinitErr)
		}
	})

	t.Run("sends notification on invalid input", func(t *testing.T) {
		// Instantiate a new node
		n := New()
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define a valid configuration
		config := transformProcessingNodeConfig.TransformProcessingNodeConfig{
			ModelInput: models.Model{
				Attributes: map[string]models.Attribute{
					"inputField1": {
						Id:   "inputField1",
						Key:  "myString",
						Type: "string",
					},
					"inputField2": {
						Id:   "inputField2",
						Key:  "myNumber",
						Type: "number",
					},
				},
			},
			ModelOutput: models.Model{
				Attributes: map[string]models.Attribute{
					"outputField1": {
						Id:   "outputField1",
						Key:  "myStringOut",
						Type: "string",
					},
					"outputField2": {
						Id:   "outputField2",
						Key:  "myNumberOut",
						Type: "number",
					},
				},
			},
			Mapping: map[string]string{
				"outputField1": "inputField1",
				"outputField2": "inputField2",
			},
		}

		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Make a notify channel
		notifyChan := make(chan node.Notification)

		// Call Init with notify channel
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Prepare invalid input data
		invalidInputData := map[string]any{
			"myNumber": "invalid", // wrong type
			"myString": "valid",
		}

		invalidInputBytes, err := json.Marshal(&invalidInputData)
		if err != nil {
			t.Fatalf("Failed to marshal input data: %v", err)
		}

		// Send the invalid input data
		go func() {
			n.(*TransformProcessingNode).inputChannels[0].Channel <- node.IoData{
				Data: invalidInputBytes,
				Type: node.BytesDataType,
				Ctx:  ctx,
			}
		}()

		// Call Exec in goroutine
		errCh := make(chan error)
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errCh <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Wait for notification about invalid data
		select {
		case notification := <-notifyChan:
			// Verify the notification indicates invalid data
			require.Equal(t, notification.Type, node.NotifyDataInvalid)
		case err := <-errCh:
			t.Fatalf("Unexpected error during execution: %v", err)
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for invalid data notification")
		}

		// Cleanup
		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Errorf("Deinit() returned error: %v", deinitErr)
		}
	})
}
