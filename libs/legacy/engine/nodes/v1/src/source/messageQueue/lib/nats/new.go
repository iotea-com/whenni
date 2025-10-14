package natsSourceNode

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	natsSourceNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/source/messageQueue/lib/nats/config"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "NATS Message Queue Consumer"
)

type NatsSourceSubnode struct {
	config         natsSourceNodeConfig.NatsSourceSubnodeConfig
	inputChannels  []node.IoChannel
	outputChannels []node.IoChannel
	notifyChannel  chan node.Notification

	conn    *nats.Conn
	errChan chan error

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &NatsSourceSubnode{
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
