package main

import (
	"context"
	"testing"

	"github.com/goccy/go-json"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	logActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/log/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("successfully initializes node", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Create input parameters
		ctx := context.Background()
		obsv := helpers.CreateTestObsv()

		// Define the configuration for the Log action node
		config := logActionNodeConfig.LogActionNodeConfig{}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Call Init with the context, config, lifetime, and observability
		initErr := n.Init(node.InitParams{
			Ctx:    ctx,
			Config: configBytes,
			Obsv:   obsv,
		})
		if initErr.Type != node.NoError {
			t.Errorf("Init() returned error: %v", initErr)
		}

		// Additional assertions can be made here to verify internal state if necessary
		deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx})
		if deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node: %v", deinitErr.Reason)
		}
	})
}
