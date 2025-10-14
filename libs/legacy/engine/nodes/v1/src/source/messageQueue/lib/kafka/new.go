package kafkaSourceNode

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	kafkaSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/kafka/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Kafka Message Queue Consumer"
)

type KafkaSourceSubnode struct {
	config         kafkaSourceNodeConfig.KafkaSourceSubnodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification

	consumer *kafka.Consumer
	errChan  chan error

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &KafkaSourceSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "passThrough",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
		outputChannels: []node.IoChannel{
			{
				Id:      "trigger",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}
