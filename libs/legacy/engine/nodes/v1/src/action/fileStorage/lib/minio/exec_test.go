package minio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	minioActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/minio/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	TestBucketName    = "test-bucket"
	TestDataDirectory = "test"
)

func setupMinIO(t *testing.T) *minio.Client {
	t.Helper()

	// Initialize minio client object.
	minioClient, err := minio.New("localhost:9005", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		t.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Check if the bucket exists, if not, create it
	exists, err := minioClient.BucketExists(context.Background(), TestBucketName)
	if err != nil {
		t.Fatalf("Failed to check if bucket exists: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), TestBucketName, minio.MakeBucketOptions{})
		if err != nil {
			t.Fatalf("Failed to create test bucket: %v", err)
		}
	}

	return minioClient
}

func clearTestBucket(t *testing.T, minioClient *minio.Client) {
	t.Helper()

	// Create a new context for listing objects
	listCtx, listCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer listCancel()

	// List all objects in the bucket
	listObjects := minioClient.ListObjects(listCtx, TestBucketName, minio.ListObjectsOptions{})

	// Delete all objects from the bucket
	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), 5*time.Second)
	for object := range listObjects {
		t.Logf("Deleting object %s", object.Key)
		err := minioClient.RemoveObject(deleteCtx, TestBucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			t.Fatalf("Failed to delete object %s from bucket: %v", object.Key, err)
		}
	}
	deleteCancel()
}

func TestExec(t *testing.T) {
	t.Run("uploads JPEG, PNG, and JSON data to Minio bucket", func(t *testing.T) {
		// Set up Minio and clear bucket
		minioClient := setupMinIO(t)
		clearTestBucket(t, minioClient)

		n := New()
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)

		// Set up configuration
		config := minioActionNodeConfig.MinioActionSubnodeConfig{
			ThingMinioBucket: things.MinioBucket{
				MinioEndpoint:        minioEndpoint,
				MinioAccessKeyId:     minioAccessKeyId,
				MinioSecretAccessKey: minioSecretAccessKey,
				BucketName:           TestBucketName,
				UseSSL:               useSSL,
			},
			Key: key,
		}
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		obsv := helpers.CreateTestObsv()

		// Initialize the node
		if initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		}); initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr.Reason)
		}

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr.Reason)
			}
		}()

		// Read test files in a separate goroutine
		go func() {
			files, err := os.ReadDir(TestDataDirectory)
			if err != nil {
				errChan <- fmt.Errorf("Failed to read test data directory: %v", err)
				return
			}

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				fileName := filepath.Join(TestDataDirectory, file.Name())
				fileContent, err := os.ReadFile(fileName)
				if err != nil {
					errChan <- fmt.Errorf("Failed to read test file %s: %v", fileName, err)
					return
				}

				n.GetInputChans(ctx)[0].Channel <- node.IoData{
					Data: fileContent,
					Type: node.BytesDataType,
					Ctx:  ctx,
				}

				// Small delay between files
				time.Sleep(100 * time.Millisecond)
			}

			// Allow time for processing before canceling
			time.Sleep(500 * time.Millisecond)
			cancel()
		}()

		// Monitor channels
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			case notification := <-notifyChan:
				t.Logf("Received notification: %v", notification)
			case err := <-errChan:
				t.Fatalf("Fatal error: %v", err)
			}
		}

		// Verify uploads after main loop exits
		listCtx, listCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer listCancel()

		uploadedFiles := 0
		for obj := range minioClient.ListObjects(listCtx, TestBucketName, minio.ListObjectsOptions{}) {
			if obj.Err != nil {
				t.Fatalf("Error in object listing: %v", obj.Err)
			}
			t.Logf("Object found: %s, Last Modified: %s, Size: %d bytes", obj.Key, obj.LastModified, obj.Size)
			uploadedFiles++
		}

		if uploadedFiles != 3 {
			t.Fatalf("Expected 3 files in bucket, but found %d", uploadedFiles)
		}
		t.Log("Successfully uploaded and verified 3 files in the bucket")

		// Cleanup
		deinitCtx, deinitCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer deinitCancel()

		if deinitErr := n.Deinit(node.DeinitParams{Ctx: deinitCtx}); deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node, %v", deinitErr.Reason)
		}
	})
}
