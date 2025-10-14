package natsSourceNode

import (
	"context"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	natsSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/nats/config"
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

		// Define a fake config
		type DummyConfig struct {
			DummyVar int `json:"topic"`
		}
		config := DummyConfig{}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Create a notify channel
		notifyChan := make(chan node.Notification)

		// Call Init
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})

		// Check if error is returned, since we expect it to fail with the dummy configuration
		if initErr.Type != node.ValidationError {
			t.Errorf("Init() did not return error; expected error for wrong config")
		}

		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node")
		}
	})

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

		// Define the configuration for the NATS source node
		config := natsSourceNodeConfig.NatsSourceSubnodeConfig{
			ThingNatsClient: things.NatsClient{
				Server: things.NatsServer{
					Host:   "127.0.0.1",
					Port:   4444,
					Topics: []string{"test-topic"},
				},
			},
			Topic: "test-topic",
		}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Create a notify channel
		notifyChan := make(chan node.Notification)

		// Call Init with the context, config, lifetime, and observability
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Errorf("Init() returned error: %v", initErr)
		}

		// Additional assertions can be made here to verify internal state if necessary
		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Errorf("Could not clean up node")
		}
	})
}
