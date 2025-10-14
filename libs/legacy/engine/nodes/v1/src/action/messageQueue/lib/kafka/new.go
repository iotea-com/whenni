package kafkaActionNode

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	kafkaActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/kafka/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Kafka Message Queue Publisher"
)

type KafkaActionSubnode struct {
	config        kafkaActionNodeConfig.KafkaActionSubnodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	producer *kafka.Producer

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &KafkaActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "messageQueueInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}
