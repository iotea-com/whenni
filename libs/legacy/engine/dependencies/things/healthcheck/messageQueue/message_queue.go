package messageQueueHealthcheck

import (
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	"github.com/nats-io/nats.go"
)

func KafkaCluster(attrs *things.KafkaCluster) error {
	// Create a producer to check connection
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(attrs.BootstrapServers, ","),
		"socket.timeout.ms": 3000, // 3 second timeout
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Get metadata to verify connection
	metadata, err := producer.GetMetadata(nil, true, 3000)
	if err != nil {
		return fmt.Errorf("failed to connect a test producer to the Kafka cluster: %v", err)
	}

	if len(metadata.Brokers) == 0 {
		return fmt.Errorf("no Kafka brokers found in the cluster")
	}

	return nil
}

func KafkaProducer(attrs *things.KafkaProducer) error {
	return fmt.Errorf("not implemented - perform a healthcheck for the Kafka cluster instead")
}

func KafkaConsumer(attrs *things.KafkaConsumer) error {
	return fmt.Errorf("not implemented - perform a healthcheck for the Kafka cluster instead")
}

func NatsServer(attrs *things.NatsServer) error {
	// Initialize NATS connection
	opts := []nats.Option{
		nats.Timeout(3 * time.Second),
	}

	conn, err := nats.Connect(fmt.Sprintf("%s:%d", attrs.Host, attrs.Port), opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Check if the connection has been successfully established
	if !conn.IsConnected() {
		return fmt.Errorf("could not connect to the NATS server")
	}

	return nil
}

func NatsClient(attrs *things.NatsClient) error {
	return fmt.Errorf("not implemented - perform a healthcheck for the NATS server instead")
}
