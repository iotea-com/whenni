package kafkaActionNode

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

func (n *KafkaActionSubnode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: fmt.Sprintf("error setting up observability: %v", err),
		}
	}

	// Set notify channel
	n.notifyChannel = params.NotifyChannel

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
			Reason: fmt.Sprintf("error connecting to Kafka cluster: %v", err),
		}
	}

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type: node.NoError,
	}
}

// setConfig unmarshals the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *KafkaActionSubnode) setConfig(config []byte) error {
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

// connectToCluster will instantiate a new *kafka.Producer with options based on the user's
// configuration, and then attempt to connect to the provided cluster.
func (n *KafkaActionSubnode) connectToCluster() error {
	// Get a things.KafkaCluster from config
	jsonCluster, err := json.Marshal(n.config.ThingKafkaProducer.Cluster)
	if err != nil {
		return fmt.Errorf("error marshalling cluster: %s", err)
	}

	var cluster things.KafkaCluster
	err = json.Unmarshal(jsonCluster, &cluster)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON into struct: %s", err)
	}

	// Connect a new producer to the cluster
	producerConfig := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cluster.BootstrapServers, ","),
	}

	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		return fmt.Errorf("could not connect to the Kafka cluster: %s", err)
	}

	n.producer = producer

	// Create an AdminClient using the same config to verify cluster connection
	adminClient, err := kafka.NewAdminClient(producerConfig)
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

	n.logger.Debug().Msgf("Kafka producer connected and topic %s is ready.", topicName)

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *KafkaActionSubnode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("MessageQueueActionNodeTracer")
	n.meter = obsv.MeterProvider.Meter("MessageQueueActionNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}
