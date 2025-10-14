package awsS3

import (
	"context"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	s3ActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/awsS3/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

const (
	// Default configuration values. these must match the ones in the
	// configuration for the docker-compose file in the nodes test directry ?
	AwsAccessKeyId     = "testing"
	AwsSecretAccessKey = "garbage-value"
	AwsRegion          = "us-west-2"
	bucketName         = "path-test"
	key                = "some-key"

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

		// Define an invalid configuration for the S3ActionNode
		type DummyConfig struct {
			DummyVar int `json:"dummy_var"`
		}
		config := DummyConfig{}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		notifyChannel := make(chan node.Notification)

		// Call Init with the invalid config
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			Obsv:          obsv,
			NotifyChannel: notifyChannel,
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

func TestInit_ValidConfig(t *testing.T) {
	t.Run("successfully initializes node with valid config", func(t *testing.T) {
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

		// Define a valid configuration for the S3ActionNode
		config := s3ActionNodeConfig.S3ActionSubnodeConfig{
			ThingAwsS3Bucket: things.S3Bucket{
				AwsAccessKeyId:     AwsAccessKeyId,
				AwsSecretAccessKey: AwsSecretAccessKey,
				AwsRegion:          AwsRegion,
				BucketName:         bucketName,
			},
			Key: key,
		}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		notifyChannel := make(chan node.Notification)
		// Call Init with the valid config
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			Obsv:          obsv,
			NotifyChannel: notifyChannel,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned an error: %v", initErr.Reason)
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

func TestInit_MissingConfigurationValues(t *testing.T) {
	t.Run("fails with various invalid configurations", func(t *testing.T) {
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

		// Test cases with invalid configurations
		invalidConfigs := []s3ActionNodeConfig.S3ActionSubnodeConfig{
			{ // Missing AwsAccessKeyId
				ThingAwsS3Bucket: things.S3Bucket{
					AwsSecretAccessKey: AwsSecretAccessKey,
					AwsRegion:          AwsRegion,
					BucketName:         bucketName,
				},
				Key: key,
			},
			{ // Missing BucketName
				ThingAwsS3Bucket: things.S3Bucket{
					AwsAccessKeyId:     AwsAccessKeyId,
					AwsSecretAccessKey: AwsSecretAccessKey,
					AwsRegion:          AwsRegion,
				},
				Key: key,
			},
		}

		for index, config := range invalidConfigs {
			configBytes, err := json.Marshal(&config)
			if err != nil {
				t.Fatalf("Failed to marshal config: %v", err)
			}

			notifyChannel := make(chan node.Notification)
			initErr := n.Init(node.InitParams{
				Ctx:           ctx,
				Config:        configBytes,
				Obsv:          obsv,
				NotifyChannel: notifyChannel,
			})
			if initErr.Type != node.ValidationError {
				t.Fatalf("Init() returned no error for invalid config at index %d", index)
			}
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
