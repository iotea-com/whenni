package awsS3

import (
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
	s3ActionNodeConfig "github.com/iotea-com/iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/awsS3/config"
	"github.com/iotea-com/iotea/libs/engine/template"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const (
	NodeName = "S3 File Storage"
)

type S3ActionSubnode struct {
	config        s3ActionNodeConfig.S3ActionSubnodeConfig
	inputChannels []node.IoChannel
	notifyChannel chan node.Notification

	keyTemplate template.Template

	// Observability fields
	tracer trace.Tracer
	meter  metric.Meter
	logger zerolog.Logger
}

func New() node.Interface {
	return &S3ActionSubnode{
		inputChannels: []node.IoChannel{
			{
				Id:      "fileStorageInput",
				Channel: make(chan node.IoData, 100),
				Type: []node.DataType{
					node.BytesDataType,
					node.JsonDataType,
					node.FileDataType,
				},
			},
		},
		// No output channels
	}
}

func main() {}
