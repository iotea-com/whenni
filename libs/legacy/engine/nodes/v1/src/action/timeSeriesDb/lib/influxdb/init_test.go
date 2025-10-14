package influxdbNode

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goccy/go-json"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	influxdbActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/timeSeriesDb/lib/influxdb/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

const (
	// Default configuration values. these must match the ones in the
	// configuration for the docker-compose file in the nodes test directry
	DefaultProtocol = "http"
	DefaultHost     = "127.0.0.1"
	DefaultPort     = 8086
	DefaultToken    = "my-secret-influxdb2-token"
	DefaultOrgName  = "iotea"
	DefaultBucket   = "default"

	InitExpectedDuration = 500 * time.Millisecond
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
		ctx, cancel := context.WithTimeout(context.Background(), InitExpectedDuration)
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

		notifyChannel := make(chan node.Notification)

		// Call Init
		initErr := n.Init(node.InitParams{
			Config:        configBytes,
			NotifyChannel: notifyChannel,
			Ctx:           ctx,
			Obsv:          obsv,
		})

		// Check if error is returned, since we expect it to fail with the dummy configuration
		if initErr.Type != node.ValidationError {
			t.Fatalf("Init() did not return error; expected error for wrong config")
		}

		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}

// TestInit_WrongInputParameters tests that the node won't accept a missing or erroneous configuration
// parameter
func TestInit_WrongInputParameters(t *testing.T) {
	t.Run("successfully initializes node with short lifetime", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Create input parameters with a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), InitExpectedDuration)
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define an array with an entry for each field, with an invalid configuration type
		configs := []influxdbActionNodeConfig.InfluxDBActionNodeConfig{
			{
				ThingInfluxDB: things.InfluxDbDatabase{
					Host:    DefaultHost,
					Port:    DefaultPort,
					Token:   DefaultToken,
					OrgName: DefaultOrgName,
				},
				Bucket: DefaultBucket,
			},
			{
				ThingInfluxDB: things.InfluxDbDatabase{
					Protocol: DefaultProtocol,
					Token:    "",
					OrgName:  DefaultOrgName,
				},
				Bucket: DefaultBucket,
			},
			{
				ThingInfluxDB: things.InfluxDbDatabase{
					Host:    DefaultHost,
					Token:   DefaultToken,
					OrgName: "",
				},
				Bucket: DefaultBucket,
			},
			{
				ThingInfluxDB: things.InfluxDbDatabase{
					Port:    DefaultPort,
					Token:   DefaultToken,
					OrgName: DefaultOrgName,
				},
				Bucket: "",
			},
		}

		notifyChannel := make(chan node.Notification)

		for index, config := range configs {
			// Marshal the configuration into JSON
			configBytes, err := json.Marshal(&config)
			if err != nil {
				t.Fatalf("Failed to marshal config: %v", err)
			}

			// Call Init with the context, config, lifetime, and observability
			initErr := n.Init(node.InitParams{
				Config:        configBytes,
				NotifyChannel: notifyChannel,
				Ctx:           ctx,
				Obsv:          obsv,
			})
			if initErr.Type != node.ValidationError {
				t.Fatalf("Init() returned no error for wrong config field, index %v", index)
			}

			// Additional assertions can be made here to verify internal state if necessary
			deinitErr := n.Deinit(node.DeinitParams{
				Ctx: ctx,
			})
			if deinitErr.Type != node.NoError {
				t.Fatalf("Could not clean up node, %v", deinitErr.Reason)
			}
		}
	})
}

// TestInit_CreateNonExistentBucket tests that the node will succesfully create a new bucket the
// first time the channel is executed
func TestInit_CreateNonExistentBucket(t *testing.T) {
	t.Run("Checks that node creates a new bucket if it doesn't already exist", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Create input parameters with a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), InitExpectedDuration)
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define an array with an entry for each field, with an invalid configuration type
		config := influxdbActionNodeConfig.InfluxDBActionNodeConfig{
			ThingInfluxDB: things.InfluxDbDatabase{
				Protocol: DefaultProtocol,
				Host:     DefaultHost,
				Port:     DefaultPort,
				Token:    DefaultToken,
				OrgName:  DefaultOrgName,
			},
			Bucket: "my_test_bucket",
		}

		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		notifyChannel := make(chan node.Notification)

		// Call Init with the context, config, lifetime, and observability
		initErr := n.Init(node.InitParams{
			Config:        configBytes,
			NotifyChannel: notifyChannel,
			Ctx:           ctx,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Validate bucket creation
		databaseUrl := fmt.Sprintf("%s://%s:%d", config.ThingInfluxDB.Protocol, config.ThingInfluxDB.Host, config.ThingInfluxDB.Port)
		client := influxdb2.NewClient(databaseUrl, config.ThingInfluxDB.Token)
		defer client.Close()
		bucketsAPI := client.BucketsAPI()
		bucket, err := bucketsAPI.FindBucketByName(ctx, config.Bucket)
		if err != nil || bucket == nil {
			t.Fatalf("Bucket was not created as expected: %v", err)
		}

		// Cleanup: Delete the bucket
		if err := bucketsAPI.DeleteBucketWithID(ctx, *bucket.Id); err != nil {
			t.Fatalf("Failed to clean up bucket: %v", err)
		}

		// Cleanup Node
		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}

// TestInit_UseExistentBucket tests that the node will use an existing bucket if it's already present
// in the database
func TestInit_UseExistentBucket(t *testing.T) {
	t.Run("Checks that node uses an existing bucket if it already exists", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Create input parameters with a timeout context
		ctx, cancel := context.WithTimeout(context.Background(), InitExpectedDuration)
		defer cancel()

		obsv := helpers.CreateTestObsv()
		config := influxdbActionNodeConfig.InfluxDBActionNodeConfig{
			ThingInfluxDB: things.InfluxDbDatabase{
				Protocol: DefaultProtocol,
				Host:     DefaultHost,
				Port:     DefaultPort,
				Token:    DefaultToken,
				OrgName:  DefaultOrgName,
			},
			Bucket: DefaultBucket,
		}

		// Serialize the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Initialize the client to use for checking buckets
		databaseUrl := fmt.Sprintf("%s://%s:%d", config.ThingInfluxDB.Protocol, config.ThingInfluxDB.Host, config.ThingInfluxDB.Port)
		client := influxdb2.NewClient(databaseUrl, config.ThingInfluxDB.Token)
		defer client.Close()
		bucketsAPI := client.BucketsAPI()

		// Get the initial count of buckets
		initialBuckets, err := bucketsAPI.GetBuckets(ctx)
		if err != nil {
			t.Fatalf("Failed to list initial buckets: %v", err)
		}
		initialCount := len(*initialBuckets)

		notifyChannel := make(chan node.Notification)
		// Call Init with the context, config, and observability
		initErr := n.Init(node.InitParams{
			Config:        configBytes,
			NotifyChannel: notifyChannel,
			Ctx:           ctx,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		// Get the count of buckets after initialization
		finalBuckets, err := bucketsAPI.GetBuckets(ctx)
		if err != nil {
			t.Fatalf("Failed to list final buckets: %v", err)
		}
		finalCount := len(*finalBuckets)

		// Compare the initial and final counts
		if initialCount != finalCount {
			t.Fatalf("Bucket count changed from %d to %d; expected it to remain the same", initialCount, finalCount)
		}

		// Cleanup Node
		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}
