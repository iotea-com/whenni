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

type TestData struct {
	FieldInteger int     `json:"int"`
	FieldFloat32 float32 `json:"float32"`
	FieldFloat64 float64 `json:"float64"`
}

func TestExec(t *testing.T) {
	t.Run("Checks multiple data points are written and read correctly", func(t *testing.T) {
		// Instantiate a new node
		n := New()

		// Assert that 'n' is not nil
		if n == nil {
			t.Fatal("New() returned nil, expected non-nil node.Interface")
		}

		// Create input parameters with a cancellation context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define configuration with an appropriate bucket name
		config := influxdbActionNodeConfig.InfluxDBActionNodeConfig{
			ThingInfluxDB: things.InfluxDbDatabase{
				Protocol: DefaultProtocol,
				Host:     DefaultHost,
				Port:     DefaultPort,
				Token:    DefaultToken,
				OrgName:  DefaultOrgName,
			},
			Bucket: "exec_long_test_bucket",
		}

		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		notifyChannel := make(chan node.Notification)

		initErr := n.Init(node.InitParams{
			Config:        configBytes,
			NotifyChannel: notifyChannel,
			Ctx:           ctx,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr)
		}

		databaseUrl := fmt.Sprintf("%s://%s:%d", config.ThingInfluxDB.Protocol, config.ThingInfluxDB.Host, config.ThingInfluxDB.Port)
		client := influxdb2.NewClient(databaseUrl, config.ThingInfluxDB.Token)
		defer client.Close()
		bucketsAPI := client.BucketsAPI()
		bucket, err := bucketsAPI.FindBucketByName(ctx, config.Bucket)
		if err != nil || bucket == nil {
			t.Fatalf("Bucket was not created as expected: %v", err)
		}

		// Create channels for errors
		errChan := make(chan error)

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr.Reason)
			}
		}()

		numDataPoints := 10
		sentDataPoints := make([]TestData, 0, numDataPoints)

		// Start sending data points in a separate goroutine
		sendCtx, sendCancel := context.WithCancel(ctx)
		go func() {
			for i := 0; i < numDataPoints; i++ {
				testData := TestData{
					FieldInteger: 10 + i,
					FieldFloat32: 64.43 + float32(i),
					FieldFloat64: 231.98732 + float64(i),
				}
				sentDataPoints = append(sentDataPoints, testData)

				testDataBytes, err := json.Marshal(testData)
				if err != nil {
					errChan <- fmt.Errorf("Failed to marshal data: %v", err)
					return
				}

				n.GetInputChans(context.Background())[0].Channel <- node.IoData{
					Data: testDataBytes,
					Type: node.BytesDataType,
					Ctx:  sendCtx,
				}

				time.Sleep(time.Millisecond * 100) // Avoid flooding
			}

			// Give some time for writes to complete
			time.Sleep(time.Second)
			sendCancel()
		}()

	Loop:
		for {
			select {
			case <-sendCtx.Done():
				break Loop
			case notification := <-notifyChannel:
				t.Logf("received notification %v", notification)
			case err := <-errChan:
				t.Fatalf("Fatal error: %v", err)
			}
		}

		// Query and verify the data
		queryAPI := client.QueryAPI(config.ThingInfluxDB.OrgName)
		query := `from(bucket:"` + config.Bucket + `") |> range(start: -30m) |> filter(fn: (r) => r._measurement == "measurements")`

		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			t.Fatalf("Failed to query data: %v", err)
		}

		foundCount := 0
		for result.Next() {
			if result.Record().Measurement() == "measurements" {
				foundCount++
				t.Logf("Data point %d: %v:%v", foundCount, result.Record().Field(), result.Record().Value())
			}
		}

		if foundCount/3 != numDataPoints {
			t.Errorf("Expected %d data points, but found %d", numDataPoints, foundCount/3)
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
