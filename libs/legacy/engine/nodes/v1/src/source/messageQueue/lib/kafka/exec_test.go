package kafkaSourceNode

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	kafkaSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/kafka/config"
	helpers "github.com/iotea-com/iotea/libs/engine/nodes/v1/test"
)

// TestExec ensures that the source node with long lifetime successfully receives and processes messages from a Kafka cluster.
func TestExec(t *testing.T) {
	t.Run("successfully executes source node", func(t *testing.T) {
		n := New()

		// Create input parameters with a cancellation context
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		obsv := helpers.CreateTestObsv()

		// Define the configuration for the Kafka source node
		config := kafkaSourceNodeConfig.KafkaSourceSubnodeConfig{
			ThingKafkaConsumer: things.KafkaConsumer{
				Cluster: things.KafkaCluster{
					BootstrapServers: []string{"127.0.0.1:9094"},
					Topics:           []string{"iotea-test-topic"},
				},
				GroupId:         "some-group",
				AutoOffsetReset: "earliest",
			},
			Topic: "iotea-test-topic",
		}

		// Create an AdminClient using the same config to verify cluster connection
		producerConfig := &kafka.ConfigMap{
			"bootstrap.servers": "127.0.0.1:9094", // Ensure this matches your advertised listener
		}

		adminClient, err := kafka.NewAdminClient(producerConfig)
		if err != nil {
			t.Fatalf("could not create Kafka admin client: %s", err)
		}
		defer adminClient.Close()

		// Verify cluster existence by retrieving the Cluster ID
		_, err = adminClient.ClusterID(ctx)
		if err != nil {
			t.Fatalf("could not verify Kafka cluster existence: %s, is the service running?", err)
		}

		// Create or check the topic exists
		topicSpecification := kafka.TopicSpecification{
			Topic:             config.Topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}

		results, err := adminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpecification})
		if err != nil {
			t.Fatalf("could not create topic %s: %s", config.Topic, err)
		}

		for _, result := range results {
			if result.Error.Code() != kafka.ErrNoError {
				// Check if the error indicates the topic already exists
				if result.Error.Code() == kafka.ErrTopicAlreadyExists {
					log.Printf("topic %s already exists, proceeding...", config.Topic)
					continue
				}
				t.Fatalf("failed to create topic %s: %v", result.Topic, result.Error)
			}
		}

		// Marshal the configuration into JSON
		configBytes, err := json.Marshal(&config)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}

		// Create channels for errors and notifications
		errChan := make(chan error)
		notifyChan := make(chan node.Notification)

		// Call Init with the context, config, and observability
		initErr := n.Init(node.InitParams{
			Ctx:           ctx,
			Config:        configBytes,
			NotifyChannel: notifyChan,
			Obsv:          obsv,
		})
		if initErr.Type != node.NoError {
			t.Fatalf("Init() returned error: %v", initErr.Reason)
		}

		// Start the execution in a goroutine
		go func() {
			execErr := n.Exec(node.ExecParams{
				Ctx: ctx,
			})
			if execErr.Type != node.NoError {
				errChan <- fmt.Errorf("Exec() returned error: %v", execErr)
			}
		}()

		// Allow some time for the notifications channel to be setup
		time.Sleep(time.Second * 1)

		numMessages := 2

		// Create channel to store published messages
		publishedMessagesCh := make(chan map[string]any, numMessages)

		// Start mock publishing
		go func() {
			// Publish input message to Kafka
			producerConfig := &kafka.ConfigMap{
				"bootstrap.servers": "127.0.0.1:9094",
			}

			producer, err := kafka.NewProducer(producerConfig)
			if err != nil {
				fmt.Printf("Error creating Kafka producer: %v\n", err)
			}
			defer producer.Close()

			for i := 0; i < numMessages; i++ {
				message := map[string]any{
					"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
				}

				jsonMessage, _ := json.Marshal(message)

				// Store the published message in the channel for later comparison
				publishedMessagesCh <- message

				kafkaMessage := &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic: &config.Topic,
					},
					Value: jsonMessage,
				}
				err = producer.Produce(kafkaMessage, nil)
				if err != nil {
					log.Printf("could not publish message to %s: %s", config.Topic, err)
				}

				time.Sleep(time.Millisecond * 100) // Sleep to avoid flooding the broker
			}
		}()

		// Verify that the subscriber node received the correct number of messages
		receivedMessagesCount := 0

	Loop:
		for {
			select {
			case <-ctx.Done():
				// Context canceled, exit the loop
				break Loop
			case receivedData := <-n.GetOutputChans(ctx)[0].Channel:
				receivedMessagesCount++
				// Type assertion to convert received data to []byte
				dataBytes, ok := receivedData.Data.([]byte)
				if !ok {
					t.Errorf("Failed to convert received data to []byte")
					continue
				}

				// Unmarshal the received data
				var receivedMsg map[string]any
				if err := json.Unmarshal(dataBytes, &receivedMsg); err != nil {
					t.Errorf("Error unmarshalling received message: %v", err)
					continue
				}

				expectedData, ok := <-publishedMessagesCh
				if !ok {
					t.Error("Expected messages channel closed unexpectedly")
					return
				}

				// Compare sent and received messages
				if !reflect.DeepEqual(receivedMsg, expectedData) {
					t.Errorf("Received message does not match expected message. Received: %v, Expected: %v", receivedMsg, expectedData)
				} else {
					t.Logf("Successful data MATCH. Received: %v, Expected: %v", receivedMsg, expectedData)
				}

				if receivedMessagesCount >= numMessages {
					break Loop
				}
			case notification := <-notifyChan:
				t.Logf("received notification %v", notification)
			case err := <-errChan:
				t.Fatalf("Exec() returned error: %v", err)
			}
		}

		// Verify that the subscriber node received the correct number of messages
		if receivedMessagesCount != numMessages {
			t.Errorf("Number of received messages (%v) does not match expected number of messages (%v)", receivedMessagesCount, numMessages)
		} else {
			t.Logf("Message count sent(%v) and received(%v) MATCH", receivedMessagesCount, numMessages)
		}

		deinitErr := n.Deinit(node.DeinitParams{
			Ctx: ctx,
		})
		if deinitErr.Type != node.NoError {
			t.Fatalf("Could not clean up node: %v", err)
		}
	})
}
