package awsSES

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	awsSESActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/awsSES/config"
	"github.com/iotea-com/iotea/libs/engine/template"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "AWS SES"
)

type AwsSESActionSubnode struct {
	config        awsSESActionNodeConfig.AwsSESActionSubnodeConfig
	notifyChannel chan node.Notification
	inputChannels []node.IoChannel

	template template.Template

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &AwsSESActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "notificationInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}
