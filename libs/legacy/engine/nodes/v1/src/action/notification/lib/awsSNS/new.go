package awsSNS

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	awsSNSActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/awsSNS/config"
	"github.com/iotea-com/iotea/libs/engine/template"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "AWS SNS"
)

type AwsSNSActionSubnode struct {
	config        awsSNSActionNodeConfig.AwsSNSActionSubnodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	template template.Template

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &AwsSNSActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "notificationInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}
