package natsActionNode

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	natsActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/messageQueue/lib/nats/config"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "NATS Message Queue Publisher"
)

type NatsActionSubnode struct {
	config        natsActionNodeConfig.NatsActionSubnodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	conn *nats.Conn

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &NatsActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "messageQueueInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}
