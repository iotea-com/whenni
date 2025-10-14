package sendgridActionNode

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	sendgridActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/notification/lib/sendgrid/config"
	"github.com/iotea-com/iotea/libs/engine/template"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "Sendgrid Client"
)

type SendgridActionSubnode struct {
	config        sendgridActionNodeConfig.SendgridActionSubnodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	template template.Template

	// Observability
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &SendgridActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "notificationInput",
				Channel: make(chan node.IoData, 100),
				Type:    []node.DataType{node.BytesDataType},
			},
		},
	}
}
