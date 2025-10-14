package main

import (
	"context"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/google/uuid"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	mqttActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/mqtt/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

const (
	ShortLifetimeExpectedDuration = 100 * time.Millisecond
)

// TestInit_WrongConfig tests that the node won't accept a wrong configuration byte array passed
func TestInit_WrongConfig(t *testing.T) {
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

		// Define the configuration for the MQTT action node
		type DummyConfig struct {
			DummyVar int `json:"topic"`
		}
		config := DummyConfig{}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Call Init
		initErr := n.Init(node.InitParams{
			Ctx:    ctx,
			Config: configBytes,
			Obsv:   obsv,
		})

		// Check if error is returned, since we expect it to fail with the dummy configuration
		if initErr.Type != node.ValidationError {
			t.Fatalf("Init() did not return error; expected error for wrong config")
		}

		deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node")
		}
	})
}

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
			Topic: "test/topic",
			QoS:   0,
		}

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
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Additional assertions can be made here to verify internal state if necessary
		deinitErr := n.Deinit(node.DeinitParams{Ctx: ctx})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node")
		}
	})
}
