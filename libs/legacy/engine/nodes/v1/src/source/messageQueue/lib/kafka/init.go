package kafkaSourceNode

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

const (
	// The connection to a cluster should not take a very long time
	ClusterConnectTimeout = 1 * time.Second
)

func (n *KafkaSourceSubnode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	// Set the notify channel
	n.notifyChannel = params.NotifyChannel

	// Setup the error channel
	n.errChan = make(chan error)

	// Setup configuration fields
	if err := n.setConfig(params.Config); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: fmt.Sprintf("invalid configuration, %v", err),
		}
	}

	// Connect to the Kafka cluster
	err := n.connectToCluster()
	if err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshals the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *KafkaSourceSubnode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	// Validate config
	v := validator.New()

	if err := v.Struct(n.config); err != nil {
		return err
	}

	return nil
}

// connectToCluster will instantiate a new *kafka.Consumer with options based on the user's
// configuration, and then attempt to connect to the provided cluster.
func (n *KafkaSourceSubnode) connectToCluster() error {
	// Get a things.KafkaCluster from config
	jsonCluster, err := json.Marshal(n.config.ThingKafkaConsumer.Cluster)
	if err != nil {
		return fmt.Errorf("error marshalling cluster: %s", err)
	}

	var cluster things.KafkaCluster
	err = json.Unmarshal(jsonCluster, &cluster)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON into struct: %s", err)
	}

	// Connect a new consumer to the cluster
	consumerConfig := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cluster.BootstrapServers, ","),
		"group.id":          n.config.ThingKafkaConsumer.GroupId,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		return fmt.Errorf("could not connect to the Kafka cluster: %s", err)
	}

	n.consumer = consumer

	// Create an AdminClient using the same config to verify cluster connection
	adminClient, err := kafka.NewAdminClient(consumerConfig)
	if err != nil {
		return fmt.Errorf("could not create Kafka admin client: %s", err)
	}
	defer adminClient.Close()

	// Verify cluster existence by retrieving the Cluster ID
	ctx, cancel := context.WithTimeout(context.Background(), ClusterConnectTimeout)
	defer cancel()

	_, err = adminClient.ClusterID(ctx)
	if err != nil {
		return fmt.Errorf("could not verify Kafka cluster existence: %s, is the service running?", err)
	}

	// Create or check the topic exists
	topicName := n.config.Topic
	topicSpecification := kafka.TopicSpecification{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	results, err := adminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpecification})
	if err != nil {
		return fmt.Errorf("could not create topic %s: %s", topicName, err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			// Check if the error indicates the topic already exists
			if result.Error.Code() == kafka.ErrTopicAlreadyExists {
				n.logger.Debug().Msgf("Topic %s already exists, proceeding...", topicName)
				continue
			}
			return fmt.Errorf("failed to create topic %s: %v", result.Topic, result.Error)
		}
	}

	// Subscribe to topic
	err = consumer.SubscribeTopics([]string{n.config.Topic}, nil)
	if err != nil {
		return fmt.Errorf("error subscribing to topic %s: %v", n.config.Topic, err)
	}

	n.logger.Debug().Msgf("Kafka consumer connected and topic %s is ready.", topicName)

	go func() {
		// Attempt to read a message from the Kafka topic
		for {
			if consumer == nil || consumer.IsClosed() {
				n.errChan <- fmt.Errorf("consumer is closed - cannot read messages")
				break
			}

			msg, err := consumer.ReadMessage(500 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok {
					if kafkaErr.IsFatal() {
						n.errChan <- fmt.Errorf("fatal error while consuming messages: %s", kafkaErr)
						break // Exit the loop on fatal errors
					}

					if kafkaErr.Code() == kafka.ErrTimedOut {
						continue // Continue listening if it's a non-fatal timeout
					}
				}

				n.errChan <- fmt.Errorf("unknown error while consuming messages: %s", err)
				break
			}

			// If a valid message is received, process it
			n.kafkaSubscriptionCallback(msg)
		}

		// Close the consumer when the loop exits (only if there's a fatal error or shutdown)
		n.logger.Debug().Msg("Kafka consumer loop exited, closing consumer")
		err := consumer.Close()
		if err != nil {
			n.errChan <- fmt.Errorf("error closing Kafka consumer: %s", err)
		}
	}()

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *KafkaSourceSubnode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("MessageQueueSourceNodeTracer")
	n.meter = obsv.MeterProvider.Meter("MessageQueueSourceNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
